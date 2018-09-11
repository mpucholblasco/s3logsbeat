package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	sqsMaxNumberOfMessages = 10
)

// SQS handle simple SQS queue functions used by a consumer
type SQS struct {
	client *sqs.SQS
	url    *string
}

type sqsMessageHandler func(*SQSMessage) error

// NewSQS is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQS(session *session.Session, queueURL *string) *SQS {
	return &SQS{
		client: sqs.New(session),
		url:    queueURL,
	}
}

// ReceiveMessages receives messages from queue and executes message handler for each message
// Returns the number of messages received and error (if any)
// Fields present per message:
//   Body: "{jsonbody}"
//   MD5OfBody: "1212f7afeed9f2bff8e8ee2b4f81020a"
// MessageId: "b872e5af-be32-4a67-82d5-87f062937c8a"
// ReceiptHandle: "base64encodedstring"
// Returns an integer with the number of messages received, a boolean indicating that more possible
// available messages are present on the queue, and the error (if any)
func (s *SQS) ReceiveMessages(mh sqsMessageHandler) (int, bool, error) {
	received := 0
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            s.url,
		MaxNumberOfMessages: aws.Int64(sqsMaxNumberOfMessages), // 1 to 10 (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ReceiveMessage.html)
	}
	resp, err := s.client.ReceiveMessage(receiveMessageInput)

	if err != nil {
		return 0, false, err
	}

	received += len(resp.Messages)
	for _, m := range resp.Messages {
		if err := mh(newSQSMessage(m)); err != nil {
			return 0, false, err
		}
	}
	return received, len(resp.Messages) == sqsMaxNumberOfMessages, nil
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
