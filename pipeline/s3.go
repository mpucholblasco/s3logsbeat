package pipeline

import (
	"fmt"

	"github.com/elastic/beats/libbeat/common"
	"github.com/mpucholblasco/s3logsbeat/aws"
)

// S3Object S3 object element to send thru pipeline
type S3Object struct {
	*aws.S3Object
	sqsMessage *SQSMessage
}

// NewS3Object creates a new S3 object to be sent thru pipeline
func NewS3Object(awsS3Object *aws.S3Object, sqsMessage *SQSMessage) *S3Object {
	return &S3Object{
		S3Object:   awsS3Object,
		sqsMessage: sqsMessage,
	}
}

// GetKeyFields extract fields from key if input set it
func (s *S3Object) GetKeyFields() (*common.MapStr, error) {
	f := &common.MapStr{}

	if s.sqsMessage.sqs.keyRegexFields != nil {
		re := s.sqsMessage.sqs.keyRegexFields.Copy()
		match := re.FindStringSubmatch(s.Key)
		if match == nil {
			return nil, fmt.Errorf("Couldn't match key regex fields %s with S3 key %s", re.String(), s.Key)
		}
		for i, name := range re.SubexpNames() {
			// Ignore the whole regexp match, unnamed groups, and empty values
			if i == 0 || name == "" || match[i] == "" {
				continue
			}
			f.Put(name, match[i])
		}
	}
	return f, nil
}
