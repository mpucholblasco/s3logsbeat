package pipeline

import "sync"

// Channels Pipeline channels
type Channels struct {
	sqsChannel chan *SQS
	s3Channel  chan *S3Object
	mutex      sync.Mutex
}

// NewChannels creates a new pipeline channels object
func NewChannels() *Channels {
	return &Channels{
		sqsChannel: make(chan *SQS, 5),
		s3Channel:  make(chan *S3Object, 10),
	}
}

// GetSQSChannel gets SQS channel
func (c *Channels) GetSQSChannel() chan *SQS {
	return c.sqsChannel
}

// GetS3Channel gets S3 channel
func (c *Channels) GetS3Channel() chan *S3Object {
	return c.s3Channel
}

// Close methods
// Set into a lock because we can be executing a once execution, which waits for completion
// and closes channels and press CTL+C, which also closes same channels. This could produce
// error on already closed channel (mixed with conditions to know if channel is closed or
// not)

// CloseSQSChannel closes SQS channel
func (c *Channels) CloseSQSChannel() {
	c.mutex.Lock()
	if c.sqsChannel != nil {
		close(c.sqsChannel)
		c.sqsChannel = nil
	}
	c.mutex.Unlock()
}

// CloseS3Channel closes S3 channel
func (c *Channels) CloseS3Channel() {
	c.mutex.Lock()
	if c.s3Channel != nil {
		close(c.s3Channel)
		c.s3Channel = nil
	}
	c.mutex.Unlock()
}
