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

// S3ReaderWorker is a worker to read objects from S3, parse their content, and send events to output
type S3ReaderWorker struct {
	wg   sync.WaitGroup
	in   <-chan *aws.S3ObjectSQSMessage
	out  beat.Client
	done chan struct{}
}

// NewS3ReaderWorker creates a new S3ReaderWorker
func NewS3ReaderWorker(in <-chan *aws.S3ObjectSQSMessage, out beat.Client) *S3ReaderWorker {
	return &S3ReaderWorker{
		in:   in,
		out:  out,
		done: make(chan struct{}),
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
				case s3object := <-w.in:
					logp.Debug("s3logsbeat", "Reading S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
					readCloser, err := s3.GetReadCloser(s3object.S3Bucket, s3object.S3Key)
					if err != nil {
						logp.Err("Could not download S3 object from region=%s, bucket=%s, key=%s", s3object.Region, s3object.S3Bucket, s3object.S3Key)
						continue
					}
					defer readCloser.Close()
					s3object.SQSMessage.SQS.Parser.Parse(readCloser, func(event beat.Event) {
						w.out.Publish(event)
					}, func(errLine string, err error) {
						logp.Warn("Could not parse line: %s, reason: %+v", errLine, err)
					})
				}
			}
		}(n)
	}
}

// Stop sends notification to stop to workers and wait untill all workers finish
func (w *S3ReaderWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping S3 reader workers")
	close(w.done)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "S3 reader workers stopped")
}
