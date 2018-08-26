package worker

import (
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
	wg   sync.WaitGroup
	in   <-chan *aws.SQS
	out  chan<- *aws.S3ObjectSQSMessage
	done chan struct{}
}

// NewSQSConsumerWorker creates a ne SQSConsumerWorke
func NewSQSConsumerWorker(in <-chan *aws.SQS, out chan<- *aws.S3ObjectSQSMessage) *SQSConsumerWorker {
	return &SQSConsumerWorker{
		in:   in,
		out:  out,
		done: make(chan struct{}),
	}
}

// Start starts the SQSConsumerWorker
func (w *SQSConsumerWorker) Start() {
	for n := 0; n < sqsConsumerWorkers; n++ {
		w.wg.Add(1)
		go func(workerId int) {
			defer w.wg.Done()
			logp.Info("SQS consumer worker #%d : waiting for input data", workerId)
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
								if err := message.ExtractNewS3Objects(
									func(s3object *aws.S3ObjectSQSMessage) {
										w.out <- s3object
									},
								); err != nil {
									logp.Err("Error extracting S3 objects from event: %v", err)
								}
								return nil
							})
							if err != nil {
								logp.Err("Could not receive SQS messages: %v", err)
								break
							}
							if !more {
								break
							}
							// TODO: add messagesReceived to monitor
							logp.Debug("s3logsbeat", "Received %d messages from SQS queue", messagesReceived)
						}
					}
				}
			}
		}(n)
	}
}

// Stop stops SQSConsumerWorker and closes output channel
func (w *SQSConsumerWorker) Stop() {
	logp.Debug("s3logsbeat", "Stopping SQS consumer workers")
	close(w.done)
	w.wg.Wait()
	logp.Debug("s3logsbeat", "SQS consumer workers stopped")
}
