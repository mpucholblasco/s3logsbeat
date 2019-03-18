package pipeline

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mpucholblasco/s3logsbeat/aws"
)

// S3List S3 list object to send thru pipeline
type S3List struct {
	*aws.S3
	*S3ReaderInformation
	s3prefix  *aws.S3Object
	since, to time.Time
}

// NewS3List creates a new S3 to be sent thru pipeline
func NewS3List(session *session.Session, s3prefix *aws.S3Object, ri *S3ReaderInformation, since, to time.Time) *S3List {
	return &S3List{
		S3:                  aws.NewS3(session),
		S3ReaderInformation: ri,
		s3prefix:            s3prefix,
		since:               since,
		to:                  to,
	}
}

// S3Object S3 object element to send thru pipeline
type S3Object struct {
	*aws.S3Object
	*S3ReaderInformation
	s3ObjectProcessNotifications S3ObjectProcessNotifications
}

// NewS3Object creates a new S3 object to be sent thru pipeline
func NewS3Object(awsS3Object *aws.S3Object, ri *S3ReaderInformation, s3ObjectProcessNotifications S3ObjectProcessNotifications) *S3Object {
	return &S3Object{
		S3Object:                     awsS3Object,
		S3ReaderInformation:          ri,
		s3ObjectProcessNotifications: s3ObjectProcessNotifications,
	}
}
