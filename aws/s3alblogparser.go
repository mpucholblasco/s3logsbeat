package aws

import (
	"regexp"

	"github.com/mpucholblasco/s3logsbeat/logparser"
)

var (
	regexALB = `^(?P<type>[^ ]*) (?P<time>[^ ]*) (?P<elb>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) (?P<target_ip>[^ ]*)[:-](?P<target_port>[0-9]*) (?P<request_processing_time>[-.0-9]*) (?P<target_processing_time>[-.0-9]*) (?P<response_processing_time>[-.0-9]*) (?P<elb_status_code>|[-0-9]*) (?P<target_status_code>-|[-0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) \"(?P<request_verb>[^ ]*) (?P<request_url>[^ ]*) (?P<request_proto>- |[^ ]*)\" \"(?P<user_agent>[^\"]*)\" (?P<ssl_cipher>[A-Z0-9-]+) (?P<ssl_protocol>[A-Za-z0-9.-]*) (?P<target_group_arn>[^ ]*) \"(?P<trace_id>[^\"]*)\"`

	typesRegexALB, _ = logparser.KindMapKindToType(map[string]logparser.Kind{
		"time":                     logparser.TimeISO8601,
		"client_port":              logparser.Uint16,
		"target_port":              logparser.Uint16,
		"request_processing_time":  logparser.Float64,
		"target_processing_time":   logparser.Float64,
		"response_processing_time": logparser.Float64,
		"received_bytes":           logparser.Int64,
		"sent_bytes":               logparser.Int64,
	})

	S3ALBLogParser = logparser.NewCustomLogParser(regexp.MustCompile(regexALB), typesRegexALB)
)
