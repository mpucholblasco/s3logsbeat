// +build !integration

package aws

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestS3CreateEventCorrect(t *testing.T) {
	body := `
	{"Records":
		[
			{
				"eventVersion":"2.0",
				"eventSource":"aws:s3",
				"awsRegion":"eu-west-1",
				"eventTime":"2018-07-07T09:35:10.990Z",
				"eventName":"ObjectCreated:Put",
				"userIdentity":{
					"principalId":"AWS:MHYPRINCIPAL"
				},
				"requestParameters":{
					"sourceIPAddress":"34.249.104.213"
				},
				"responseElements":{
					"x-amz-request-id":"C6CC46982C978BF5",
					"x-amz-id-2":"myxamzid2"
				},
				"s3":{
					"s3SchemaVersion":"1.0",
					"configurationId":"test-s3-queue",
					"bucket":{
						"name":"mybucket",
						"ownerIdentity":{
							"principalId":"MyPrincipalID"
						},
						"arn":"arn:aws:s3:::mybucket"
					},
					"object":{
						"key":"app-env-3/AWSLogs/123456789012/elasticloadbalancing/eu-west-1/2018/07/07/123456789012_elasticloadbalancing_eu-west-1_app.app-env-3.ad4ceee8a897566c_20180707T0935Z_52.17.184.44_4vsrpn7y.log.gz",
						"size":14313,
						"eTag":"0f0c79b67cf091c2228c16640d75ff3b",
						"sequencer":"005B40894EEA476179"
					}
				}
			}
		]
	}
	`
	h := md5.New()
	io.WriteString(h, body)
	md5body := hex.EncodeToString(h.Sum(nil))
	message := &sqs.Message{
		Body:          &body,
		MD5OfBody:     &md5body,
		MessageId:     aws.String("fakeMessageId"),
		ReceiptHandle: aws.String("fakeReceipt"),
	}
	sqsMessage := NewSQSMessage(nil, message)
	results := 0
	sqsMessage.ExtractNewS3Objects(func(s *S3ObjectSQSMessage) {
		assert.Equal(t, aws.String("eu-west-1"), s.region)
		assert.Equal(t, aws.String("mybucket"), s.s3Bucket)
		assert.Equal(t, aws.String("app-env-3/AWSLogs/123456789012/elasticloadbalancing/eu-west-1/2018/07/07/123456789012_elasticloadbalancing_eu-west-1_app.app-env-3.ad4ceee8a897566c_20180707T0935Z_52.17.184.44_4vsrpn7y.log.gz"), s.s3Key)
		results = results + 1
	})
	assert.Equal(t, 1, results)
}
