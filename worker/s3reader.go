package worker

import (
	"sync"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mpucholblasco/s3logsbeat/aws"
)

const (
	s3ReaderWorkers = 5
)

type eventCounter interface {
	Add(n int)
	Done()
	Error(n uint64)
}

// S3ReaderWorker is a worker to read objects from S3, parse their content, and send events to output
type S3ReaderWorker struct {
	wg          sync.WaitGroup
	in          <-chan *aws.S3ObjectSQSMessage
	out         beat.Client
	done        chan struct{}
	wgEvents    eventCounter
	wgS3Objects eventCounter
}

// NewS3ReaderWorker creates a new S3ReaderWorker
func NewS3ReaderWorker(in <-chan *aws.S3ObjectSQSMessage, out beat.Client, wgEvents eventCounter, wgS3Objects eventCounter) *S3ReaderWorker {
	return &S3ReaderWorker{
		in:          in,
		out:         out,
		done:        make(chan struct{}),
		wgEvents:    wgEvents,
		wgS3Objects: wgS3Objects,
	}
}

// Start starts the worker
func (w *S3ReaderWorker) Start() {
	s3 := aws.NewS3(aws.NewSession())

	w.wg.Add(s3ReaderWorkers)
	for n := 0; n < s3ReaderWorkers; n++ {
		go func(workerID int) {
			defer w.wg.Done()
			logp.Info("S3 reader worker #%d : waiting for input data", workerID)
			for {
				select {
				case <-w.done:
					logp.Info("S3 reader worker #%d finished", workerID)
					return
				case s3object, ok := <-w.in:
					if !ok {
						logp.Info("S3 reader worker #%d finished because channel is closed", workerID)
						return
					}
					w.onS3ObjectFromSQSMessage(s3, s3object)
				}
			}
		}(n)
	}
}

func (w *S3ReaderWorker) onS3ObjectFromSQSMessage(s3 *aws.S3, s3object *aws.S3ObjectSQSMessage) {
	onLogParserSucceed := func(event beat.Event) {
		event.Private = s3object.SQSMessage // store to reduce on ACK function
		s3object.SQSMessage.AddEvents(1)
		w.wgEvents.Add(1)
		w.out.Publish(event)
	}

	onLogParserError := func(errLine string, err error) {
		w.wgEvents.Error(1)
		logp.Warn("Could not parse line: %s, reason: %+v", errLine, err)
	}

	logp.Debug("s3logsbeat", "Reading S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
	if readCloser, err := s3.GetReadCloser(s3object.S3Bucket, s3object.S3Key); err != nil {
		w.wgS3Objects.Error(1)
		logp.Err("Could not download S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
	} else {
		defer readCloser.Close()
		s3object.SQSMessage.SQS.Parser.Parse(readCloser, onLogParserSucceed, onLogParserError)
	}

	// Monitoring
	w.wgS3Objects.Done()

	// Counting how much remaining events are on this SQS message to delete it when all will be processed
	s3object.SQSMessage.S3ObjectProcessed()
}

// Wait waits until all workers have finished
func (w *S3ReaderWorker) Wait() {
	w.wg.Wait()
}

// Stop sends notification to stop to workers and wait untill all workers finish.
// We will not accept more S3 objects
func (w *S3ReaderWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping S3 reader workers")
	w.in = nil
	close(w.done)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "S3 reader workers stopped")
}
