package pipeline

import "sync"

const (
	maxSQSChanCapacity    = 5
	maxS3ListChanCapacity = 5
	maxS3ObjectsCapacity  = 10
)

type baseChannels struct {
	s3Channel chan *S3Object
	m         sync.Mutex
}

// GetS3Channel gets S3 channel
func (c *baseChannels) GetS3Channel() chan *S3Object {
	return c.s3Channel
}

// CloseS3Channel closes S3 channel
func (c *baseChannels) CloseS3Channel() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.s3Channel != nil {
		close(c.s3Channel)
		c.s3Channel = nil
	}
}

// Channels Pipeline channels
type Channels struct {
	baseChannels
	sqsChannel chan *SQS
}

// NewChannels creates a new pipeline channels object
func NewChannels() *Channels {
	return &Channels{
		baseChannels: baseChannels{
			s3Channel: make(chan *S3Object, maxS3ObjectsCapacity),
		},
		sqsChannel: make(chan *SQS, maxSQSChanCapacity),
	}
}

// GetSQSChannel gets SQS channel
func (c *Channels) GetSQSChannel() chan *SQS {
	return c.sqsChannel
}

// CloseSQSChannel closes SQS channel
func (c *Channels) CloseSQSChannel() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.sqsChannel != nil {
		close(c.sqsChannel)
		c.sqsChannel = nil
	}
}

// S3ImportsChannels Pipeline channels used on s3imports command
type S3ImportsChannels struct {
	baseChannels
	s3ListChannel chan *S3List
}

// NewS3ImportsChannels creates a new s3 list pipeline channels object
func NewS3ImportsChannels() *S3ImportsChannels {
	return &S3ImportsChannels{
		baseChannels: baseChannels{
			s3Channel: make(chan *S3Object, maxS3ObjectsCapacity),
		},
		s3ListChannel: make(chan *S3List, maxS3ListChanCapacity),
	}
}

// GetS3ListChannel gets S3 list channel
func (c *S3ImportsChannels) GetS3ListChannel() chan *S3List {
	return c.s3ListChannel
}

// CloseS3ListChannel closes S3 list channel
func (c *S3ImportsChannels) CloseS3ListChannel() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.s3ListChannel != nil {
		close(c.s3ListChannel)
		c.s3ListChannel = nil
	}
}
