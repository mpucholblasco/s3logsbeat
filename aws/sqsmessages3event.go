package aws

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

// SQSMessageS3Event SQS message providing from an S3 event
type SQSMessageS3Event struct {
	sqsMessage *SQSMessage
}

type s3Event struct {
	Records []struct {
		EventSource string `json:"eventSource"`
		AwsRegion   string `json:"awsRegion"`
		EventName   string `json:"eventName"`
		S3          struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key  string `json:"key"`
				Size int    `json:"size"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

// NewSQSMessageS3Event creates a new SQS message from S3 event based on an SQS message
func NewSQSMessageS3Event(sqsMessage *SQSMessage) *SQSMessageS3Event {
	return &SQSMessageS3Event{
		sqsMessage: sqsMessage,
	}
}

// ExtractNewObjects extracts those new S3 objects present on an SQS message
// Returns the number of new S3 objects extracted
func (s *SQSMessageS3Event) ExtractNewObjects(mh func(*S3Object) error) (uint64, error) {
	var s3e s3Event
	if err := json.Unmarshal([]byte(*s.sqsMessage.Body), &s3e); err != nil {
		logp.Warn("Couldn't parse json from S3. Ignoring this SQS mess. Error: %v", err)
		return 0, nil
	}
	var c uint64
	for _, e := range s3e.Records {
		if e.EventSource == "aws:s3" && strings.HasPrefix(e.EventName, "ObjectCreated:") {
			if s3key, err := url.QueryUnescape(e.S3.Object.Key); err != nil {
				logp.Warn("Could not unescape S3 object: %s", e.S3.Object.Key)
			} else {
				c++
				if err := mh(NewS3Object(e.S3.Bucket.Name, s3key)); err != nil {
					// Client want to cancel process, passing as an error to parent
					return 0, err
				}
			}
		}
	}
	return c, nil
}
