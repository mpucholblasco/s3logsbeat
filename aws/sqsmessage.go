package aws

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage SQS message
type SQSMessage struct {
	SQS     *SQS
	Message *sqs.Message

	// Control S3 objects to be processed and events to be acked
	mutex     sync.Mutex
	s3objects uint64
	events    uint64

	// Events
	onDeleteCallbacks []func()
}

// NewSQSMessage is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQSMessage(sqs *SQS, message *sqs.Message) *SQSMessage {
	sqsMessage := &SQSMessage{
		SQS:     sqs,
		Message: message,
	}
	logp.Info("Generated new SQS message with ID %s", *message.MessageId)

	return sqsMessage
}

// Events

// OnDelete adds callback for OnDelete event
func (sm *SQSMessage) OnDelete(f func()) {
	sm.onDeleteCallbacks = append(sm.onDeleteCallbacks, f)
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
	if sm.SQS != nil {
		for _, f := range sm.onDeleteCallbacks {
			f()
		}
		logp.Info("Deleting message with ID %s from SQS", *sm.Message.MessageId)
		//return sm.SQS.DeleteMessage(sm.Message.ReceiptHandle) TODO uncomment
	}
	return nil
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

// S3ObjectProcessed reduces the number of pending S3 objects to process and executed DeleteOnJobCompleted
func (sm *SQSMessage) S3ObjectProcessed() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.s3objects--
	sm.deleteOnJobCompleted()
}

// AddEvents adds the number of events to the counter (to know the number of events pending to ACK)
func (sm *SQSMessage) AddEvents(c uint64) {
	sm.mutex.Lock()
	sm.events += c
	sm.mutex.Unlock()
}

// EventsProcessed reduces the number of events to the counter (to know the number of events pending to ACK).
// If all events have been processed, the SQS message is deleted.
func (sm *SQSMessage) EventsProcessed(c uint64) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.events -= c
	if sm.events < 0 {
		panic(fmt.Sprintf("Acked %d more events than added", -sm.events))
	}
	sm.deleteOnJobCompleted()
}

func (sm *SQSMessage) deleteOnJobCompleted() {
	if sm.s3objects == 0 && sm.events == 0 {
		sm.Delete()
	}
}
