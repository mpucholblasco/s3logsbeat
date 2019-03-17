package aws

import (
	"fmt"
	"regexp"
)

var (
	s3uriRE = regexp.MustCompile(`^s3://(?P<bucket>[^/]+)/(?P<key>.*)$`)
)

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

// NewS3ObjectFromURI creates a new S3 object from a URI with format:
// s3://bucket/path
func NewS3ObjectFromURI(uri string) (*S3Object, error) {
	re := s3uriRE.Copy()
	match := re.FindStringSubmatch(uri)
	if match == nil {
		return nil, fmt.Errorf("Incorrect S3 URI %s", uri)
	}

	return &S3Object{
		Bucket: match[1],
		Key:    match[2],
	}, nil
}

// String converts current object into string
func (s *S3Object) String() string {
	return fmt.Sprintf("S3Object{Bucket:%s, Key: %s}", s.Bucket, s.Key)
}
