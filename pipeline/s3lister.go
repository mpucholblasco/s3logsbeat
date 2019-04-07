package pipeline

import (
	"fmt"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/sequra/s3logsbeat/aws"
)

const (
	s3ListerWorkers = 2
)

// S3ListerWorker is a worker to list an S3 bucket and pass them to the output (out channel)
type S3ListerWorker struct {
	wg              sync.WaitGroup
	in              <-chan *S3List
	out             chan<- *S3Object
	done            chan struct{}
	wgS3Objects     eventCounter
	keepSQSMessages bool
}

// NewS3ListerWorker creates an S3ListerWorker
func NewS3ListerWorker(in <-chan *S3List, out chan<- *S3Object, wgS3Objects eventCounter) *S3ListerWorker {
	return &S3ListerWorker{
		in:          in,
		out:         out,
		done:        make(chan struct{}),
		wgS3Objects: wgS3Objects,
	}
}

// Start starts the SQS consumer workers
func (w *S3ListerWorker) Start() {
	w.wg.Add(s3ListerWorkers)
	for n := 0; n < s3ListerWorkers; n++ {
		go func(workerID int) {
			defer w.wg.Done()
			logp.Info("S3 lister worker #%d : waiting for input data", workerID)
			for {
				select {
				case <-w.done:
					logp.Info("S3 lister worker #%d finished", workerID)
					return
				case s3, ok := <-w.in:
					if !ok {
						logp.Info("S3 lister worker #%d finished because channel is closed", workerID)
						return
					}
					w.onS3List(workerID, s3)
				}
			}
		}(n)
	}
}

// onS3ListNotification Reads S3 objects from a bucket and prefix
func (w *S3ListerWorker) onS3List(workerID int, s3 *S3List) {
	logp.Debug("s3logsbeat", "Listening S3 objects present on S3 prefix URI %s", s3.s3prefix.String())
	onS3Object := func(o *aws.S3ObjectWithOriginal) error {
		// Case in which we have a lot of filtered S3 files and we want to cancel the process
		select {
		case <-w.done:
			logp.Info("Cancelling ListS3Objects")
			return fmt.Errorf("Cancelling")
		default:
		}
		if o.LastModified.After(s3.since) && o.LastModified.Before(s3.to) {
			// Using a select because w.out could be full
			select {
			case <-w.done:
				logp.Info("Cancelling ListS3Objects")
				return fmt.Errorf("Cancelling")
			case w.out <- NewS3Object(o.S3Object, s3.S3ReaderInformation, NewS3ObjectProcessNotificationsIgnorer()):
				w.wgS3Objects.Add(1)
			}
		} else {
			if logp.IsDebug("s3logsbeat") {
				logp.Debug("s3logsbeat", "Filtering key s3://%s/%s because does not fit with timestamp (%s) filters (since=%s,to=%s)", o.Bucket, o.S3Object.Key, o.LastModified, s3.since.UTC(), s3.to.UTC())
			}
		}
		return nil
	}

	if objectsReceived, err := s3.ListObjects(s3.s3prefix, onS3Object); err != nil {
		logp.Err("Could not list S3 object from S3 prefix URI %s. Error: %v", s3.s3prefix.String(), err)
	} else {
		logp.Debug("s3logsbeat", "Received %d S3 objects from S3 prefix URI %s", objectsReceived, s3.s3prefix.String())
	}
}

// Wait waits until all workers have finished
func (w *S3ListerWorker) Wait() {
	w.wg.Wait()
}

// Stop sends notification to stop to workers and wait untill all workers finish
func (w *S3ListerWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping S3 lister workers")
	close(w.done)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "S3 lister workers stopped")
}
