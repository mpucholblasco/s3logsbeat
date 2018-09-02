package aws

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 handle simple S3 methods
type S3 struct {
	client *s3.S3
}

// NewS3 is a construct function for creating the object
// with session
func NewS3(session *session.Session) *S3 {
	client := s3.New(session)

	s3 := &S3{
		client: client,
	}

	return s3
}

// GetReadCloser returns a io.ReadCloser to be readed (and then closed) by another method.
func (s *S3) GetReadCloser(bucket string, key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return output.Body, err
}
