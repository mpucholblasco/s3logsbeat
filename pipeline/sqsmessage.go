package pipeline

import (
	"fmt"
	"sync"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/mpucholblasco/s3logsbeat/aws"
)

// SQSMessage SQS message to be passed thru pipeline.
// We have to keep how much S3 objects and how much events
// are generated from this message in order to delete it
// from SQS once it finishes
type SQSMessage struct {
	*aws.SQSMessage

	sqs *SQS

	// Control S3 objects to be processed and events to be acked
	mutex     *sync.Mutex
	s3objects uint64
	events    uint64

	// Events
	onDeleteCallbacks []func()
}

// NewSQSMessage is a construct function for creating the object
// with session and url of the queue as arguments
func NewSQSMessage(sqs *SQS, sqsMessage *aws.SQSMessage) *SQSMessage {
	return &SQSMessage{
		SQSMessage: sqsMessage,
		sqs:        sqs,
		mutex:      &sync.Mutex{},
	}
}

// Events

// OnDelete adds callback for OnDelete event
func (s *SQSMessage) OnDelete(f func()) {
	s.onDeleteCallbacks = append(s.onDeleteCallbacks, f)
}

// S3ObjectProcessed reduces the number of pending S3 objects to process and executed DeleteOnJobCompleted
func (s *SQSMessage) S3ObjectProcessed() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.s3objects--
	s.deleteOnJobCompleted()
}

// AddEvents adds the number of events to the counter (to know the number of events pending to ACK)
func (s *SQSMessage) AddEvents(c uint64) {
	s.mutex.Lock()
	s.events += c
	s.mutex.Unlock()
}

// EventsProcessed reduces the number of events to the counter (to know the number of events pending to ACK).
// If all events have been processed, the SQS message is deleted.
func (s *SQSMessage) EventsProcessed(c uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.events -= c
	if s.events < 0 {
		panic(fmt.Sprintf("Acked %d more events than added", -s.events))
	}
	s.deleteOnJobCompleted()
}

func (s *SQSMessage) deleteOnJobCompleted() {
	if s.s3objects == 0 && s.events == 0 {
		s.delete()
	}
}

func (s *SQSMessage) delete() {
	//sm.sqs.DeleteMessage(sm.ReceiptHandle) TODO uncomment on final version
}

// ExtractNewS3Objects extracts those new S3 objects present on an SQS message
// This function is executed on a mutex to avoid the following case:
// Time 0 -> Goroutine A (GA) : executes ExtractNewS3Objects with first S3 element and keeps on the loop
// Time 1 -> Goroutine B (GB) : downloads S3 object and is empty. It executes DeleteOnJobCompleted and deletes SQS message
// Time 2 -> app crashes
// Problem: as SQS message has already been deleted, it can not be processed again
func (s *SQSMessage) ExtractNewS3Objects(mh func(s3object *S3Object) error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s3event := aws.NewSQSMessageS3Event(s.SQSMessage)
	c, err := s3event.ExtractNewObjects(func(awsS3Object *aws.S3Object) error {
		return mh(NewS3Object(awsS3Object, s))
	})

	if err != nil {
		return err
	}

	if c == 0 {
		logp.Warn("No S3 objects extracted from SQS message with ID %s", s.MessageId)
		s.delete()
	} else {
		s.s3objects += c
	}
	return nil
}
