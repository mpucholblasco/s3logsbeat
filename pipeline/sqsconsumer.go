package pipeline

import (
	"fmt"
	"sync"

	"github.com/mpucholblasco/s3logsbeat/aws"

	"github.com/elastic/beats/libbeat/logp"
)

const (
	sqsConsumerWorkers = 2
)

// SQSConsumerWorker is a worker to read SQS notifications for reading messages from AWS (present on
// in channel), extract new S3 objects present on messages and pass to the output (out channel)
type SQSConsumerWorker struct {
	wg              sync.WaitGroup
	in              <-chan *SQS
	out             chan<- *S3Object
	done            chan struct{}
	doneForced      chan struct{}
	wgSQSMessages   eventCounter
	wgS3Objects     eventCounter
	keepSQSMessages bool
}

// NewSQSConsumerWorker creates an SQSConsumerWorker
func NewSQSConsumerWorker(in <-chan *SQS, out chan<- *S3Object, wgSQSMessages eventCounter, wgS3Objects eventCounter, keepSQSMessages bool) *SQSConsumerWorker {
	return &SQSConsumerWorker{
		in:              in,
		out:             out,
		done:            make(chan struct{}),
		doneForced:      make(chan struct{}),
		wgSQSMessages:   wgSQSMessages,
		wgS3Objects:     wgS3Objects,
		keepSQSMessages: keepSQSMessages,
	}
}

// Start starts the SQS consumer workers
func (w *SQSConsumerWorker) Start() {
	w.wg.Add(sqsConsumerWorkers)
	for n := 0; n < sqsConsumerWorkers; n++ {
		go func(workerID int) {
			defer w.wg.Done()
			logp.Info("SQS consumer worker #%d : waiting for input data", workerID)
			for {
				select {
				case <-w.done:
					logp.Info("SQS consumer worker #%d finished", workerID)
					return
				case sqs, ok := <-w.in:
					if !ok {
						logp.Info("SQS consumer worker #%d finished because channel is closed", workerID)
						return
					}
					w.onSQSNotification(workerID, sqs)
				}
			}
		}(n)
	}
}

// Reads SQS messages from SQS queue until empty or AWS returns less than
// maximum
func (w *SQSConsumerWorker) onSQSNotification(workerID int, sqs *SQS) {
	logp.Debug("s3logsbeat", "Reading SQS messages from queue %s", sqs.String())
	var messagesReceived int
	var err error
	more := true

	onNewS3Object := func(s3object *S3Object) error {
		// Using a select because w.out could be full
		select {
		case <-w.doneForced:
			logp.Info("Cancelling ExtractNewS3Objects")
			return fmt.Errorf("Cancelling")
		case w.out <- s3object:
			w.wgS3Objects.Add(1)
		}
		return nil
	}

	onSQSMessage := func(message *aws.SQSMessage) error {
		// Monitoring
		w.wgSQSMessages.Add(1)
		m := NewSQSMessage(sqs, message, w.keepSQSMessages)
		m.OnDelete(func() { w.wgSQSMessages.Done() })

		// Extract new S3 objects from SQS message
		return m.ExtractNewS3Objects(onNewS3Object)
	}

	for more {
		// Avoid reading more SQS messages on stop
		select {
		case <-w.done:
			return
		default:
			if messagesReceived, more, err = sqs.ReceiveMessages(onSQSMessage); err != nil {
				w.wgSQSMessages.Error(1)
				logp.Err("Could not receive SQS messages from queue %s. Error: %v", sqs.String(), err)
				// more is false when err != nil -> exiting from loop
			} else {
				logp.Debug("s3logsbeat", "Received %d messages from SQS queue %s", messagesReceived, sqs.String())
			}
		}
	}
}

// StopAcceptingMessages sends notification to stop to workers and wait until all workers finish
func (w *SQSConsumerWorker) StopAcceptingMessages() {
	logp.Debug("s3logsbeat", "SQS consumers not accepting more messages")
	w.in = nil
	close(w.done)
}

// Wait waits until all workers have finished
func (w *SQSConsumerWorker) Wait() {
	w.wg.Wait()
}

// Stop sends notification to stop to workers and wait untill all workers finish
func (w *SQSConsumerWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping SQS consumer workers")
	close(w.doneForced)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "SQS consumer workers stopped")
}
