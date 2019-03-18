package aws

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 handle simple S3 methods
type S3 struct {
	client *s3.S3
}

type s3readcloser struct {
	io.ReadCloser
	i io.ReadCloser
	c []io.Closer
}

type s3ObjectHandler func(*S3ObjectWithOriginal) error

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
func (s *S3) GetReadCloser(o *S3Object) (io.ReadCloser, error) {
	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(o.Bucket),
		Key:    aws.String(o.Key),
	})
	if err != nil {
		return nil, err
	}
	return newS3ReadCloser(output.Body, o.Key)
}

// ListObjects lists objects present on o.Bucket and prefix o.Key
func (s *S3) ListObjects(o *S3Object, oh s3ObjectHandler) (int, error) {
	received := 0
	s3ListObjectsInput := &s3.ListObjectsInput{
		Bucket: aws.String(o.Bucket),
		Prefix: aws.String(o.Key),
	}
	err := s.client.ListObjectsPages(s3ListObjectsInput, func(res *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		received += len(res.Contents)
		for _, r := range res.Contents {
			if err := oh(NewS3ObjectWithOriginal(o.Bucket, r)); err != nil {
				return false
			}
		}
		return true
	})
	return received, err
}

func newS3ReadCloser(i io.ReadCloser, key string) (io.ReadCloser, error) {
	s := &s3readcloser{
		i: i,
		c: []io.Closer{i},
	}
	if strings.HasSuffix(key, ".gz") {
		var err error
		s.i, err = gzip.NewReader(i)
		if err != nil {
			return nil, err
		}
		s.c = append(s.c, s.i)
	}
	return s, nil
}

func (s *s3readcloser) Read(p []byte) (n int, err error) {
	return s.i.Read(p)
}

func (s *s3readcloser) Close() error {
	for i := len(s.c) - 1; i >= 0; i-- {
		s.c[i].Close()
	}
	return nil
}
