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
	wg  sync.WaitGroup
	in  <-chan *aws.SQS
	out chan<- *aws.S3ObjectSQSMessage
}

// NewSQSConsumerWorker creates a ne SQSConsumerWorker
func NewSQSConsumerWorker(in <-chan *aws.SQS, out chan<- *aws.S3ObjectSQSMessage) (*SQSConsumerWorker, error) {
	return &SQSConsumerWorker{
		in:  in,
		out: out,
	}, nil
}

// Start starts the SQSConsumerWorker
func (w *SQSConsumerWorker) Start() {
	for n := 0; n < sqsConsumerWorkers; n++ {
		w.wg.Add(1)
		go func(workerId int) {
			defer w.wg.Done()
			logp.Info("SQS consumer worker #%d : waiting for input data", workerId)
			for sqs := range w.in {
				logp.Debug("Reading SQS messages from queue: %s", sqs.String())
				messagesReceived, err := sqs.ReceiveMessages(func(message *aws.SQSMessage) error {
					logp.Debug("input", "Message: %v", message)
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
				}
				// TODO: add received to monitor
				logp.Debug("sqsconsumer", "Received %d messages from SQS queue", messagesReceived)
			}
			logp.Info("SQS consumer worker #%d finished", workerId)
		}(n)
	}
}

// Stop stops SQSConsumerWorker and closes output channel
func (w *SQSConsumerWorker) Stop() {
	logp.Info("Stopping SQS consumer workers")
	w.wg.Wait()
	close(w.out)
	logp.Info("SQS consumer workers stopped")
}
