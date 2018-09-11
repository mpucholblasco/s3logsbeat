package aws

import "fmt"

// S3Object represents an object on S3
type S3Object struct {
	Bucket string
	Key    string
}

// NewS3Object creates a new S3 object
func NewS3Object(bucket, key string) *S3Object {
	return &S3Object{
		Bucket: bucket,
		Key:    key,
	}
}

// String converts current object into string
func (s *S3Object) String() string {
	return fmt.Sprintf("S3Object{Bucket:%s, Key: %s}", s.Bucket, s.Key)
}
