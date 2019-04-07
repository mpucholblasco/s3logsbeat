package pipeline

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sequra/s3logsbeat/aws"
)

// SQS SQS element to send thru pipeline
type SQS struct {
	*aws.SQS
	*S3ReaderInformation
}

// NewSQS creates a new SQS to be sent thru pipeline
func NewSQS(session *session.Session, queueURL *string, ri *S3ReaderInformation) *SQS {
	return &SQS{
		SQS:                 aws.NewSQS(session, queueURL),
		S3ReaderInformation: ri,
	}
}
