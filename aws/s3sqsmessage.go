package aws

import (
	"encoding/json"
	"net/url"

	"github.com/elastic/beats/libbeat/logp"
)

// S3SQSMessage interface to extract new S3 objects from SQS messages
type S3SQSMessage interface {
	ExtractNewS3Objects()
}

// S3ObjectSQSMessage contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type S3ObjectSQSMessage struct {
	SQSMessage *SQSMessage
	Region     string
	S3Bucket   string
	S3Key      string
}

type s3messageHandler func(*S3ObjectSQSMessage)

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

// ExtractNewS3Objects extracts those new S3 objects present on an SQS message
func (sm *SQSMessage) ExtractNewS3Objects(mh s3messageHandler) error {
	var s3e s3Event
	if err := json.Unmarshal([]byte(*sm.Message.Body), &s3e); err != nil {
		return err
	}
	for _, e := range s3e.Records {
		if e.EventSource == "aws:s3" && e.EventName == "ObjectCreated:Put" {
			if s3key, err := url.QueryUnescape(e.S3.Object.Key); err != nil {
				logp.Warn("Could not unescape S3 object: %s", e.S3.Object.Key)
			} else {
				mh(&S3ObjectSQSMessage{
					SQSMessage: sm,
					Region:     e.AwsRegion,
					S3Bucket:   e.S3.Bucket.Name,
					S3Key:      s3key,
				})
			}
		}
	}
	return nil
}
