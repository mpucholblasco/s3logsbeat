package aws

import (
	"regexp"

	"github.com/mpucholblasco/s3logsbeat/logparser"
)

var (
	// S3ALBLogParser S3 ALB logs parser
	S3ALBLogParser = logparser.
		NewCustomLogParser("timestamp", regexp.MustCompile(`^(?P<type>[^ ]*) (?P<timestamp>[^ ]*) (?P<elb>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) (?P<target_ip>[^ ]*)[:-](?P<target_port>[0-9]*) (?P<request_processing_time>[-.0-9]*) (?P<target_processing_time>[-.0-9]*) (?P<response_processing_time>[-.0-9]*) (?P<elb_status_code>|[-0-9]*) (?P<target_status_code>-|[-0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) \"(?P<request_verb>[^ ]*) (?P<request_url>[^ ]*) (?P<request_proto>- |[^ ]*)\" \"(?P<user_agent>[^\"]*)\" (?P<ssl_cipher>[A-Z0-9-]+) (?P<ssl_protocol>[A-Za-z0-9.-]*) (?P<target_group_arn>[^ ]*) \"(?P<trace_id>[^\"]*)\"`)).
		WithKindMap(map[string]string{
			"timestamp":                "timeISO8601",
			"client_port":              "uint16",
			"target_port":              "uint16",
			"request_processing_time":  "float64",
			"target_processing_time":   "float64",
			"response_processing_time": "float64",
			"received_bytes":           "int64",
			"sent_bytes":               "int64",
			"elb_status_code":          "int16",
			"target_status_code":       "int16",
		}).
		WithEmptyValues(map[string]string{
			"user_agent":   "-",
			"ssl_cipher":   "-",
			"ssl_protocol": "-",
		})
)
