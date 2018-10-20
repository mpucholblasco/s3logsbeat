// +build !integration

package logparser

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

// Examples present here have been obtained from: https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/access-log-collection.html
func TestS3ELBLogParser(t *testing.T) {
	logs := `2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000073 0.001048 0.000057 200 200 0 29 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.38.0" - -
2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000086 0.001048 0.001337 200 200 0 57 "GET https://www.example.com:443/ HTTP/1.1" "curl/7.38.0" DHE-RSA-AES128-SHA TLSv1.2
2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.001069 0.000028 0.000041 - - 82 305 "- - - " "-" - -
2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.001065 0.000015 0.000023 - - 57 502 "- - - " "-" ECDHE-ECDSA-AES128-GCM-SHA256 TLSv1.2`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2015, 5, 13, 23, 39, 43, 945958000, time.UTC),
			Fields: common.MapStr{
				"elb":                      "my-loadbalancer",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"backend_ip":               "10.0.0.1",
				"backend_port":             uint16(80),
				"request_processing_time":  0.000073,
				"backend_processing_time":  0.001048,
				"response_processing_time": 0.000057,
				"elb_status_code":          int16(200),
				"backend_status_code":      int16(200),
				"received_bytes":           int64(0),
				"sent_bytes":               int64(29),
				"request_verb":             "GET",
				"request_url":              "http://www.example.com:80/",
				"request_proto":            "HTTP/1.1",
				"user_agent":               "curl/7.38.0",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2015, 5, 13, 23, 39, 43, 945958000, time.UTC),
			Fields: common.MapStr{
				"elb":                      "my-loadbalancer",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"backend_ip":               "10.0.0.1",
				"backend_port":             uint16(80),
				"request_processing_time":  0.000086,
				"backend_processing_time":  0.001048,
				"response_processing_time": 0.001337,
				"elb_status_code":          int16(200),
				"backend_status_code":      int16(200),
				"received_bytes":           int64(0),
				"sent_bytes":               int64(57),
				"request_verb":             "GET",
				"request_url":              "https://www.example.com:443/",
				"request_proto":            "HTTP/1.1",
				"user_agent":               "curl/7.38.0",
				"ssl_cipher":               "DHE-RSA-AES128-SHA",
				"ssl_protocol":             "TLSv1.2",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2015, 5, 13, 23, 39, 43, 945958000, time.UTC),
			Fields: common.MapStr{
				"elb":                      "my-loadbalancer",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"backend_ip":               "10.0.0.1",
				"backend_port":             uint16(80),
				"request_processing_time":  0.001069,
				"backend_processing_time":  0.000028,
				"response_processing_time": 0.000041,
				"received_bytes":           int64(82),
				"sent_bytes":               int64(305),
			},
		},
		&beat.Event{
			Timestamp: time.Date(2015, 5, 13, 23, 39, 43, 945958000, time.UTC),
			Fields: common.MapStr{
				"elb":                      "my-loadbalancer",
				"client_ip":                "192.168.131.39",
				"client_port":              uint16(2817),
				"backend_ip":               "10.0.0.1",
				"backend_port":             uint16(80),
				"request_processing_time":  0.001065,
				"backend_processing_time":  0.000015,
				"response_processing_time": 0.000023,
				"received_bytes":           int64(57),
				"sent_bytes":               int64(502),
				"ssl_cipher":               "ECDHE-ECDSA-AES128-GCM-SHA256",
				"ssl_protocol":             "TLSv1.2",
			},
		},
	}
	errorLinesExpected := []string{}
	assertLogParser(t, S3ELBLogParser, &logs, expected, errorLinesExpected)
}
