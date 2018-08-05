// +build !integration

package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3ALBLogParse(t *testing.T) {
	logs := `
	https 2018-05-16T14:44:13.509969Z app/balancer-name/1234567890abcdef 52.214.210.139:46134 172.31.10.249:80 0.000 0.005 0.000 200 200 337 562 "GET https://www.example.com:443/p1/page?a1=v1&a2=v2 HTTP/1.1" "http.rb/0.9.9" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:eu-west-1:123456789012:targetgroup/balancer-default/1234567890abcdef "Root=1-5afc43bd-34d550b3832942ea29b68552" "www.example.com" "arn:aws:iam::123456789012:server-certificate/www.example.com" 0 2018-05-16T14:44:13.503000Z "forward"
	`
	results := 0
	parser := NewS3ALBLogParser()
	parser.parse(&logs, func(s *struct{}) {
		results = results + 1
	})
	assert.Equal(t, 1, results)
}
