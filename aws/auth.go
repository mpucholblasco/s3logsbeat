package aws

import (
	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Auth is a common authenticator for AWS services
type Auth struct {
	Config  *aws.Config
	Session *session.Session
}

// NewSession tries to create AWS authentication in this order:
// 1. Environment
// 2. EC2 ROle
// 3. Shared credentials located in ~/.aws/credentials
func NewSession(region string) *session.Session {
	usr, _ := user.Current()
	credentialsFile := filepath.Join(usr.HomeDir, ".aws/credentials")

	sessionRole := session.Must(session.NewSession())

	// Try first environment variables, then with IAM roles and lastly with user configuration
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sessionRole),
			},
			&credentials.SharedCredentialsProvider{
				Filename: "", // leave empty to use default values
				Profile:  "",
			},
		})

	return session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
}
