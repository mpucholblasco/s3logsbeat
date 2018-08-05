package aws

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
//maxNumberOfMessages = 10
)

// SQS handle simple SQS queue functions used by a consumer
type S3 struct {
	client *s3.S3
	url    *string
}

// NewSQS is a construct function for creating the object
// with session and url of the queue as arguments
func News3(session *session.Session, queueURL *string) *S3 {
	client := s3.New(session)

	s3 := &S3{
		client: client,
	}

	return s3
}

// ReceiveMessages receives messages from queue and executes message handler for each message
// Returns the number of messages received and error (if any)
// Fields present per message:
//   Body: "{jsonbody}"
//   MD5OfBody: "1212f7afeed9f2bff8e8ee2b4f81020a"
// MessageId: "b872e5af-be32-4a67-82d5-87f062937c8a"
// ReceiptHandle: "base64encodedstring"
func (s *S3) GetReader(bucket string, key string) (*bytes.Reader, error) {
	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	// TODO: maybe not needed -> see https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
	// Convert S3 to a buffer
	buff, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	// Converts to a read seeker
	reader := bytes.NewReader(buff)
	return reader, nil
}
