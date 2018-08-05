// +build !integration

package aws

import (
	"strings"
	"testing"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

// Examples present here have been obtained from: https://docs.aws.amazon.com/es_es/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func TestS3ALBLogParseHTTP(t *testing.T) {
	logs := `http 2016-08-10T22:08:42.945958Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.000 0.001 0.000 200 200 34 366 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.46.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" - -`
	expected := []common.MapStr{
		common.MapStr{},
	}
	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

func TestS3ALBLogParserHTTPS(t *testing.T) {
	logs := `https 2016-08-10T23:39:43.065466Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.086 0.048 0.037 200 200 0 57 "GET https://www.example.com:443/ HTTP/1.1" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337281-1d84f3d73c47ec4e58577259" www.example.com arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012`
	expected := []common.MapStr{
		common.MapStr{},
	}
	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

func TestS3ALBLogParserHTTP2(t *testing.T) {
	logs := `h2 2016-08-10T00:10:33.145057Z app/my-loadbalancer/50dc6c495c0c9188 10.0.1.252:48160 10.0.0.66:9000 0.000 0.002 0.000 200 200 5 257 "GET https://10.0.2.105:773/ HTTP/2.0" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337327-72bd00b0343d75b906739c42" - -`
	expected := []common.MapStr{
		common.MapStr{},
	}
	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

func TestS3ALBLogParserWS(t *testing.T) {
	logs := `ws 2016-08-10T00:32:08.923954Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:40914 10.0.1.192:8010 0.001 0.003 0.000 101 101 218 587 "GET http://10.0.0.30:80/ HTTP/1.1" "-" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
	expected := []common.MapStr{
		common.MapStr{},
	}
	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

func TestS3ALBLogParserWSS(t *testing.T) {
	logs := `wss 2016-08-10T00:42:46.423695Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:44244 10.0.0.171:8010 0.000 0.001 0.000 101 101 218 786 "GET https://10.0.0.30:443/ HTTP/1.1" "-" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
	expected := []common.MapStr{
		common.MapStr{},
	}
	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

func testALBLogParser(t *testing.T, logs *string, expected []common.MapStr, expectedErrorLines []string) {
	results := make([]common.MapStr, len(expected))
	errors := make([]string, len(expectedErrorLines))
	parser := NewS3ALBLogParser()
	parser.parse(strings.NewReader(*logs), func(s common.MapStr) {
		results = append(results, s)
	}, func(errLine string) {
		errors = append(errors, errLine)
	})
	assert.Equal(t, expectedErrorLines, errors)
	assert.Equal(t, expected, results)
}
