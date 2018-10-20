package logparser

import (
	"regexp"
)

var (
	// S3ELBLogParser S3 ELB logs parser
	S3ELBLogParser = NewCustomLogParser("timestamp", regexp.MustCompile(`^(?P<timestamp>[^ ]*) (?P<elb>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) ((?P<backend_ip>[^ ]+)[:-](?P<backend_port>[0-9]+)|-) (?P<request_processing_time>[-.0-9]*) (?P<backend_processing_time>[-.0-9]*) (?P<response_processing_time>[-.0-9]*) (?P<elb_status_code>|[-0-9]*) (?P<backend_status_code>-|[-0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) \"(?P<request_verb>[^ ]*) (?P<request_url>[^ ]*) (?P<request_proto>- |[^ ]*)\" \"(?P<user_agent>[^\"]*)\" (?P<ssl_cipher>[A-Z0-9-]+) (?P<ssl_protocol>[A-Za-z0-9.-]*)`)).
		WithKindMap(map[string]string{
			"timestamp":                "timeISO8601",
			"client_port":              "uint16",
			"backend_port":             "uint16",
			"request_processing_time":  "float64",
			"backend_processing_time":  "float64",
			"response_processing_time": "float64",
			"request_url":              "urlencoded",
			"received_bytes":           "int64",
			"sent_bytes":               "int64",
			"elb_status_code":          "int16",
			"backend_status_code":      "int16",
		}).
		WithEmptyValues(map[string]string{
			"user_agent":               "-",
			"ssl_cipher":               "-",
			"ssl_protocol":             "-",
			"elb_status_code":          "-",
			"request_processing_time":  "-1",
			"backend_processing_time":  "-1",
			"response_processing_time": "-1",
			"backend_status_code":      "-",
		})
)
