package pipeline

// S3ObjectProcessNotifications interface implemented by elements passed on event.Private
type S3ObjectProcessNotifications interface {
	// EventACKed executed when an event is ACKed
	EventACKed()

	// EventSent executed when an event is sent to beats pipeline to be published
	EventSent()

	// S3ObjectProcessed executed when an S3 object is completely processed
	S3ObjectProcessed()
}

type s3ObjectProcessNotificationsIgnorer struct{}

// NewS3ObjectProcessNotificationsIgnorer creates an S3 object process notifications which ignores all events
func NewS3ObjectProcessNotificationsIgnorer() S3ObjectProcessNotifications {
	return &s3ObjectProcessNotificationsIgnorer{}
}

// S3ObjectProcessNotifications implementations

// S3ObjectProcessed reduces the number of pending S3 objects to process and executed DeleteOnJobCompleted
func (s *s3ObjectProcessNotificationsIgnorer) S3ObjectProcessed() {}

// EventSent adds the number of events to the counter (to know the number of events pending to ACK)
func (s *s3ObjectProcessNotificationsIgnorer) EventSent() {}

// EventACKed reduces the number of events to the counter (to know the number of events pending to ACK).
// If all events have been processed, the SQS message is deleted.
func (s *s3ObjectProcessNotificationsIgnorer) EventACKed() {}
