package aws

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage SQS message
type SQSMessage struct {
	*sqs.Message
}

// NewSQSMessage is a construct function for creating the object
// with session and url of the queue as arguments
func newSQSMessage(message *sqs.Message) *SQSMessage {
	sqsMessage := &SQSMessage{
		message,
	}
	logp.Info("Generated new SQS message with ID %s", *message.MessageId)

	return sqsMessage
}

// VerifyMD5Sum returns true if MD5 passed on message corresponds with the one
// obtained from body.
func (sm *SQSMessage) VerifyMD5Sum() bool {
	h := md5.New()
	io.WriteString(h, *sm.Body)
	md5body := hex.EncodeToString(h.Sum(nil))
	return md5body == *sm.MD5OfBody
}
