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
	awsSession := aws.NewSession()
	s3 := aws.NewS3(awsSession)

	for n := 0; n < s3ReaderWorkers; n++ {
		w.wg.Add(1)
		go func(workerId int) {
			defer w.wg.Done()
			logp.Info("S3 reader worker #%d : waiting for input data", workerId)
			for {
				select {
				case <-w.done:
					logp.Info("S3 reader worker #%d finished", workerId)
					return
				case s3object, ok := <-w.in:
					if !ok {
						logp.Info("S3 reader worker #%d finished because channel is closed", workerId)
						return
					}
					logp.Debug("s3logsbeat", "Reading S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
					readCloser, err := s3.GetReadCloser(s3object.S3Bucket, s3object.S3Key)
					if err != nil {
						logp.Err("Could not download S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
					} else {
						defer readCloser.Close()
						s3object.SQSMessage.SQS.Parser.Parse(readCloser, func(event beat.Event) {
							event.Private = s3object.SQSMessage // store to reduce on ACK function
							s3object.SQSMessage.AddEvents(1)
							w.wgEvents.Add(1)
							w.out.Publish(event)
						}, func(errLine string, err error) {
							logp.Warn("Could not parse line: %s, reason: %+v", errLine, err)
						})
					}
					w.wgS3Objects.Done()
					s3object.SQSMessage.S3ObjectProcessed()
				}
			}
		}(n)
	}
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
