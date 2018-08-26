package aws

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage SQS message
type SQSMessage struct {
	sqs     *SQS
	message *sqs.Message
}

// NewSQSMessage is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQSMessage(sqs *SQS, message *sqs.Message) *SQSMessage {
	sqsMessage := &SQSMessage{
		sqs:     sqs,
		message: message,
	}

	return sqsMessage
}

// GetID get message ID
func (sm *SQSMessage) GetID() *string {
	return sm.message.MessageId
}

// GetBody get message body
func (sm *SQSMessage) GetBody() *string {
	return sm.message.Body
}

// Delete deletes message
func (sm *SQSMessage) Delete() error {
	return sm.sqs.DeleteMessage(sm.message.ReceiptHandle)
}

// String converts to String
func (sm *SQSMessage) String() string {
	return fmt.Sprintf("%+v", sm.message)
}

// VerifyMD5Sum returns true if MD5 passed on message corresponds with the one
// obtained from body.
func (sm *SQSMessage) VerifyMD5Sum() bool {
	h := md5.New()
	io.WriteString(h, *sm.message.Body)
	md5body := hex.EncodeToString(h.Sum(nil))
	return md5body == *sm.message.MD5OfBody
}
