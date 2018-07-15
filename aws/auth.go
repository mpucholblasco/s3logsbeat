package aws

import (
	"sync"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	sess *session.Session
	once sync.Once
)

// NewSession creates an AWS session. Credentials loaded from SDK default credential chain:
// 1. Environment
// 2. Shared crednetials (~/.aws/credentials)
// 3. EC2 instance role
func NewSession() *session.Session {
	once.Do(func() {
		sess = session.Must(session.NewSession())
		if aws.StringValue(sess.Config.Region) == "" {
			ec2m := ec2metadata.New(sess)
			regionFound, err := ec2m.Region()
			if err != nil {
				logp.Err("Region not found in environment variables, shared credentials, or instance role. Error: %v", err)
			}
			sess.Config.Region = aws.String(regionFound)
		}
	})
	return sess

}
