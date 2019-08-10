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

func TestS3CreateEventIncorrect(t *testing.T) {
	body := `
	{"Records":
		[
			{
				"eventVersion":"2.0",
				"eventSource":"aws:s3",
				"awsRegion":"eu-west-1",
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
	s := NewSQSMessageS3Event(newSQSMessage(message))
	c, err := s.ExtractNewObjects(func(o *S3Object) error {
		return nil
	})
	assert.NoError(t, err) // only generates an error on log
	assert.Equal(t, uint64(0), c)
}

func TestS3CreateEventCorrectSimple(t *testing.T) {
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
	s := NewSQSMessageS3Event(newSQSMessage(message))
	c, err := s.ExtractNewObjects(func(s *S3Object) error {
		assert.Equal(t, "mybucket", s.Bucket)
		assert.Equal(t, "app-env-3/AWSLogs/123456789012/elasticloadbalancing/eu-west-1/2018/07/07/123456789012_elasticloadbalancing_eu-west-1_app.app-env-3.ad4ceee8a897566c_20180707T0935Z_52.17.184.44_4vsrpn7y.log.gz", s.Key)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), c)
}

func TestS3CreateEventCorrectEncoded(t *testing.T) {
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
						"key":"My+simple+%5Bkey%5D",
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
	s := NewSQSMessageS3Event(newSQSMessage(message))
	c, err := s.ExtractNewObjects(func(s *S3Object) error {
		assert.Equal(t, "mybucket", s.Bucket)
		assert.Equal(t, "My simple [key]", s.Key)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), c)
}

func TestS3CreateEventIncorrectlyEncoded(t *testing.T) {
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
						"key":"My+simple+%5key%5D",
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
	s := NewSQSMessageS3Event(newSQSMessage(message))
	c, err := s.ExtractNewObjects(func(s *S3Object) error {
		// Not called, only shown a warn on log
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), c)
}

func TestS3CreateEventMultiPartUpload(t *testing.T) {
	body := `
	{
    "Records": [
			{
					"eventVersion": "2.1",
					"eventSource": "aws:s3",
					"awsRegion": "eu-west-1",
					"eventTime": "2019-08-09T17:40:00.623Z",
					"eventName": "ObjectCreated:CompleteMultipartUpload",
					"userIdentity": {
							"principalId": "AWS:AIDAUUSOPYXONJ6WWRD54"
					},
					"requestParameters": {
							"sourceIPAddress": "199.27.72.20"
					},
					"responseElements": {
							"x-amz-request-id": "F53AE4D787EE7FBD",
							"x-amz-id-2": "GeYFbT+p3UpK/4EslyHGfnA/n95Ie7eUHWpDficcI21vLzrjpiETX6M1Ea/ORXq3LlWZWvvaWt8="
					},
					"s3": {
							"s3SchemaVersion": "1.0",
							"configurationId": "my-logs",
							"bucket": {
									"name": "mybucket",
									"ownerIdentity": {
											"principalId": "ABC1EFGHIJKLM"
									},
									"arn": "arn:aws:s3:::mybucket"
							},
							"object": {
									"key": "2019-08-09T17-35-00.000-4EOusK9ws_69koeNqTBf.log",
									"size": 28897,
									"eTag": "1f40ff64a9136f43f8915ed4f6640339-1",
									"sequencer": "005D4DAF014D3FE184"
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
	s := NewSQSMessageS3Event(newSQSMessage(message))
	c, err := s.ExtractNewObjects(func(s *S3Object) error {
		assert.Equal(t, "mybucket", s.Bucket)
		assert.Equal(t, "2019-08-09T17-35-00.000-4EOusK9ws_69koeNqTBf.log", s.Key)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), c)
}
