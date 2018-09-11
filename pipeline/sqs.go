package pipeline

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/parser"
)

// SQS SQS element to send thru pipeline
type SQS struct {
	*aws.SQS
	logParser      parser.LogParser
	keyRegexFields *regexp.Regexp
}

// NewSQS creates a new SQS to be sent thru pipeline
func NewSQS(session *session.Session, queueURL *string, logParser parser.LogParser, keyRegexFields *regexp.Regexp) *SQS {
	return &SQS{
		SQS:            aws.NewSQS(session, queueURL),
		logParser:      logParser,
		keyRegexFields: keyRegexFields,
	}
}
