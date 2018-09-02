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
	SQS     *SQS
	Message *sqs.Message
}

// NewSQSMessage is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQSMessage(sqs *SQS, message *sqs.Message) *SQSMessage {
	sqsMessage := &SQSMessage{
		SQS:     sqs,
		Message: message,
	}

	return sqsMessage
}

// GetID get message ID
func (sm *SQSMessage) GetID() *string {
	return sm.Message.MessageId
}

// GetBody get message body
func (sm *SQSMessage) GetBody() *string {
	return sm.Message.Body
}

// Delete deletes message
func (sm *SQSMessage) Delete() error {
	return sm.SQS.DeleteMessage(sm.Message.ReceiptHandle)
}

// String converts to String
func (sm *SQSMessage) String() string {
	return fmt.Sprintf("%+v", sm.Message)
}

// VerifyMD5Sum returns true if MD5 passed on message corresponds with the one
// obtained from body.
func (sm *SQSMessage) VerifyMD5Sum() bool {
	h := md5.New()
	io.WriteString(h, *sm.Message.Body)
	md5body := hex.EncodeToString(h.Sum(nil))
	return md5body == *sm.Message.MD5OfBody
}
