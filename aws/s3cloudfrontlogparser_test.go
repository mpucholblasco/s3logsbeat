// +build !integration

package aws

import (
	"strings"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/mpucholblasco/s3logsbeat/testutil"
	"github.com/stretchr/testify/assert"
)

// Examples present here have been obtained from: https://docs.aws.amazon.com/es_es/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func TestS3CloudFrontWebLogParse(t *testing.T) {
	logs := `#Version: 1.0
#Fields: date time x-edge-location sc-bytes c-ip cs-method cs(Host) cs-uri-stem sc-status cs(Referer) cs(User-Agent) cs-uri-query cs(Cookie) x-edge-result-type x-edge-request-id x-host-header cs-protocol cs-bytes time-taken x-forwarded-for ssl-protocol ssl-cipher x-edge-response-result-type cs-protocol-version fle-status fle-encrypted-fields
2014-05-23	01:13:11	FRA2	182	192.0.2.10	GET	d111111abcdef8.cloudfront.net	/view/my/file.html	200	www.displaymyfiles.com	Mozilla/4.0%20(compatible;%20MSIE%205.0b1;%20Mac_PowerPC)	-	zip=98101	RefreshHit	MRVMF7KydIvxMWfJIglgwHQwZsbG2IhRJ07sn9AkKUFSHS9EXAMPLE==	d111111abcdef8.cloudfront.net	http	-	0.001	-	-	-	RefreshHit	HTTP/1.1	Processed	1
2014-05-23	01:13:12	LAX1	2390282	192.0.2.202	GET	d111111abcdef8.cloudfront.net	/soundtrack/happy.mp3	304	www.unknownsingers.com	Mozilla/4.0%20(compatible;%20MSIE%207.0;%20Windows%20NT%205.1)	a=b&c=d	zip=50158	Hit	xGN7KWpVEmB9Dp7ctcVFQC4E-nrcOcEKS3QyAez--06dV7TEXAMPLE==	d111111abcdef8.cloudfront.net	http	-	0.002	-	-	-	Hit	HTTP/1.1	-	-`
	expected := []common.MapStr{
		common.MapStr{
			"timestamp":                   time.Date(2014, 5, 23, 1, 13, 11, 0, time.UTC),
			"x_edge_location":             "FRA2",
			"sc_bytes":                    uint64(182),
			"c_ip":                        "192.0.2.10",
			"cs_method":                   "GET",
			"cs_host":                     "d111111abcdef8.cloudfront.net",
			"cs_uri_stem":                 "/view/my/file.html",
			"sc_status":                   int16(200),
			"cs_referer":                  "www.displaymyfiles.com",
			"cs_user_agent":               "Mozilla/4.0 (compatible; MSIE 5.0b1; Mac_PowerPC)",
			"cs_uri_query":                "-",
			"cs_cookie":                   "zip=98101",
			"x_edge_result_type":          "RefreshHit",
			"x_edge_request_id":           "MRVMF7KydIvxMWfJIglgwHQwZsbG2IhRJ07sn9AkKUFSHS9EXAMPLE==",
			"x_host_header":               "d111111abcdef8.cloudfront.net",
			"cs_protocol":                 "http",
			"cs_bytes":                    "-",
			"time_taken":                  0.001,
			"x_forwarded_for":             "-",
			"ssl_protocol":                "-",
			"ssl_cipher":                  "-",
			"x_edge_response_result_type": "RefreshHit",
			"cs_protocol_version":         "HTTP/1.1",
			"fle_status":                  "Processed",
			"fle_encrypted_fields":        "1",
		},
		common.MapStr{
			"timestamp":                   time.Date(2014, 5, 23, 1, 13, 12, 0, time.UTC),
			"x_edge_location":             "LAX1",
			"sc_bytes":                    uint64(2390282),
			"c_ip":                        "192.0.2.202",
			"cs_method":                   "GET",
			"cs_host":                     "d111111abcdef8.cloudfront.net",
			"cs_uri_stem":                 "/soundtrack/happy.mp3",
			"sc_status":                   int16(304),
			"cs_referer":                  "www.unknownsingers.com",
			"cs_user_agent":               "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
			"cs_uri_query":                "a=b&c=d",
			"cs_cookie":                   "zip=50158",
			"x_edge_result_type":          "Hit",
			"x_edge_request_id":           "xGN7KWpVEmB9Dp7ctcVFQC4E-nrcOcEKS3QyAez--06dV7TEXAMPLE==",
			"x_host_header":               "d111111abcdef8.cloudfront.net",
			"cs_protocol":                 "http",
			"cs_bytes":                    "-",
			"time_taken":                  0.002,
			"x_forwarded_for":             "-",
			"ssl_protocol":                "-",
			"ssl_cipher":                  "-",
			"x_edge_response_result_type": "Hit",
			"cs_protocol_version":         "HTTP/1.1",
			"fle_status":                  "-",
			"fle_encrypted_fields":        "-",
		},
	}

	errorLinesExpected := []string{}
	testCloudFrontLogParser(t, &logs, expected, errorLinesExpected)
}

func testCloudFrontLogParser(t *testing.T, logs *string, expected []common.MapStr, expectedErrorLines []string) {
	results := make([]common.MapStr, 0, len(expected))
	errors := make([]string, 0, len(expectedErrorLines))
	err := S3CloudFrontWebLogParser.Parse(strings.NewReader(*logs), func(s common.MapStr) {
		results = append(results, s)
	}, func(errLine string, err error) {
		errors = append(errors, errLine)
	})
	assert.NoError(t, err)
	assert.Len(t, errors, len(expectedErrorLines))
	assert.Len(t, results, len(expected))
	for idx, expEvent := range expected {
		resultEvent := results[idx]
		testutil.AssertEvent(t, expEvent, resultEvent)
	}
	assert.Equal(t, expectedErrorLines, errors)
}
