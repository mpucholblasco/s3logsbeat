package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQS handle simple SQS queue functions used by a consumer
type SQS struct {
	client *sqs.SQS
	url    *string
}

// NewSQS is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQS(session.Session *session, queueUrl string) (*SQS, error) {

	client := sqs.New(session)

	sqs := &SQS{
		client: client,
		url:    queueUrl,
	}

	return sqs, nil
}

// ReceiveMessage from queue
func (q *SQS) ReceiveMessage() (*sqs.Message, error) {
	messageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            q.url,
		MaxNumberOfMessages: aws.Int64(1), // TODO: change by config
	}
	resp, err := q.client.ReceiveMessage(messageInput)
	if err != nil {
		return nil, err
	}
	// SQS messages should not be more than one
	if len(resp.Messages) > 1 {
		return nil, fmt.Errorf("Too many SQS messages")
	} else if len(resp.Messages) == 1 {
		return resp.Messages[0], nil
	} else {
		return &sqs.Message{}, nil
	}
}

// DeleteMessage from queue
func (q *SQS) DeleteMessage(receiptHandle *string) error {
	var err error
	if receiptHandle != nil {
		messageInput := &sqs.DeleteMessageInput{
			QueueUrl:      q.url,
			ReceiptHandle: receiptHandle,
		}

		_, err = q.client.DeleteMessage(messageInput)
	}

	return err
}
