package pipeline

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/logparser"
)

// SQS SQS element to send thru pipeline
type SQS struct {
	*aws.SQS
	logParser      logparser.LogParser
	keyRegexFields *regexp.Regexp
}

// NewSQS creates a new SQS to be sent thru pipeline
func NewSQS(session *session.Session, queueURL *string, logParser logparser.LogParser, keyRegexFields *regexp.Regexp) *SQS {
	return &SQS{
		SQS:            aws.NewSQS(session, queueURL),
		logParser:      logParser,
		keyRegexFields: keyRegexFields,
	}
}
