package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	maxNumberOfMessages = 10
)

// SQS handle simple SQS queue functions used by a consumer
type SQS struct {
	client *sqs.SQS
	url    *string
}

type messageHandler func(*SQSMessage) error

// NewSQS is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQS(session *session.Session, queueURL *string) *SQS {
	client := sqs.New(session)

	sqs := &SQS{
		client: client,
		url:    queueURL,
	}

	return sqs
}

// ReceiveMessages receives messages from queue and executes message handler for each message
// Returns the number of messages received and error (if any)
// Fields present per message:
//   Body: "{jsonbody}"
//   MD5OfBody: "1212f7afeed9f2bff8e8ee2b4f81020a"
// MessageId: "b872e5af-be32-4a67-82d5-87f062937c8a"
// ReceiptHandle: "base64encodedstring"
func (s *SQS) ReceiveMessages(mh messageHandler) (int, error) {
	received := 0
	for {
		logp.Debug("sqsconsumer", "Waiting for messages")
		receiveMessageInput := &sqs.ReceiveMessageInput{
			QueueUrl:            s.url,
			MaxNumberOfMessages: aws.Int64(maxNumberOfMessages), // 1 to 10 (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html)
		}
		resp, err := s.client.ReceiveMessage(receiveMessageInput)

		if err != nil {
			return 0, err
		}

		logp.Debug("sqsconsumer", "Received %d messages", len(resp.Messages))
		received += len(resp.Messages)
		for i := range resp.Messages {
			mh(NewSQSMessage(s, resp.Messages[i]))
		}
		if len(resp.Messages) < maxNumberOfMessages {
			logp.Debug("sqsconsumer", "Received all messages (%d)", received)
			return received, nil
		}
	}
}

// DeleteMessage deletes a message from queue
func (s *SQS) DeleteMessage(receiptHandle *string) error {
	var err error
	if receiptHandle != nil {
		messageInput := &sqs.DeleteMessageInput{
			QueueUrl:      s.url,
			ReceiptHandle: receiptHandle,
		}

		_, err = s.client.DeleteMessage(messageInput)
	}

	return err
}

func (s *SQS) String() string {
	return fmt.Sprintf("%s", *s.url)
}
