package worker

import (
	"fmt"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/mpucholblasco/s3logsbeat/aws"
)

const (
	sqsConsumerWorkers = 2
)

// SQSConsumerWorker is a worker to read SQS notifications for reading messages from AWS (present on
// in channel), extract new S3 objects present on messages and pass to the output (out channel)
type SQSConsumerWorker struct {
	wg         sync.WaitGroup
	in         <-chan *aws.SQS
	out        chan<- *aws.S3ObjectSQSMessage
	done       chan struct{}
	doneForced chan struct{}
}

// NewSQSConsumerWorker creates a ne SQSConsumerWorke
func NewSQSConsumerWorker(in <-chan *aws.SQS, out chan<- *aws.S3ObjectSQSMessage) *SQSConsumerWorker {
	return &SQSConsumerWorker{
		in:         in,
		out:        out,
		done:       make(chan struct{}),
		doneForced: make(chan struct{}),
	}
}

// Start starts the SQSConsumerWorker
func (w *SQSConsumerWorker) Start() {
	for n := 0; n < sqsConsumerWorkers; n++ {
		w.wg.Add(1)
		go func(workerId int) {
			defer w.wg.Done()
			logp.Info("SQS consumer worker #%d : waiting for input data", workerId)
		INPUT_LOOP:
			for {
				select {
				case <-w.done:
					logp.Info("SQS consumer worker #%d finished", workerId)
					return
				case sqs := <-w.in:
					logp.Debug("s3logsbeat", "Reading SQS messages from queue: %s", sqs.String())
					for {
						select {
						case <-w.done:
							logp.Info("SQS consumer worker #%d finished", workerId)
							return
						default:
							messagesReceived, more, err := sqs.ReceiveMessages(func(message *aws.SQSMessage) error {
								return message.ExtractNewS3Objects(
									func(s3object *aws.S3ObjectSQSMessage) error {
										// Using a select because w.out could be full
										select {
										case <-w.doneForced:
											logp.Info("Cancelling ExtractNewS3Objects")
											return fmt.Errorf("Cancelling")
										case w.out <- s3object:
										}
										return nil
									},
								)
							})
							if err != nil {
								logp.Err("Could not receive SQS messages: %v", err)
								continue INPUT_LOOP
							}
							if !more {
								continue INPUT_LOOP
							}
							logp.Debug("s3logsbeat", "Received %d messages from SQS queue", messagesReceived)
						}
					}
				}
			}
		}(n)
	}
}

// StopAcceptingMessages sends notification to stop to workers and wait untill all workers finish
func (w *SQSConsumerWorker) StopAcceptingMessages() {
	logp.Debug("s3logsbeat", "SQS consumers not accepting more messages")
	w.in = nil
	close(w.done)
}

// Stop sends notification to stop to workers and wait untill all workers finish
func (w *SQSConsumerWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping SQS consumer workers")
	close(w.doneForced)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "SQS consumer workers stopped")
}
