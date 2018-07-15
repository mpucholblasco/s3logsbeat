package worker

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mpucholblasco/s3logsbeat/aws"
)

const (
	SQSConsumerWorkers = 5
)

// NewSQS is a construct function for creating the object
// with session and url of the queue as arguments
func sqsConsumer(in <-chan *aws.SQS) {
	for n := 0; n <= SQSConsumerWorkers; n++ {
		go func() {
			sqs := <-in
			messagesReceived, err := sqs.ReceiveMessages(func(message *aws.SQSMessage) error {
				logp.Debug("input", "Message: %v", message)
				// Generate object to read from S3 and pass to output
				return nil
			})
			if err != nil {
				logp.Err("Could not receive SQS messages: %v", err)
			}
			// TODO: add received to monitor
			logp.Debug("sqsconsumer", "Received %d messages from SQS queue", messagesReceived)
		}()
	}
}
