package parser

import (
	"regexp"
)

var (
	// S3CloudFrontWebLogParser parser for CloudFront Web logs
	S3CloudFrontWebLogParser = NewCustomLogParser("timestamp", regexp.MustCompile(`^(?P<timestamp>[^\t]*\t[^\t]*)\t(?P<x_edge_location>[^\t]*)\t(?P<sc_bytes>[^\t]*)\t(?P<c_ip>[^\t]*)\t(?P<cs_method>[^\t]*)\t(?P<cs_host>[^\t]*)\t(?P<cs_uri_stem>[^\t]*)\t(?P<sc_status>[^\t]*)\t(?P<cs_referer>[^\t]*)\t(?P<cs_user_agent>[^\t]*)\t(?P<cs_uri_query>[^\t]*)\t(?P<cs_cookie>[^\t]*)\t(?P<x_edge_result_type>[^\t]*)\t(?P<x_edge_request_id>[^\t]*)\t(?P<x_host_header>[^\t]*)\t(?P<cs_protocol>[^\t]*)\t(?P<cs_bytes>[^\t]*)\t(?P<time_taken>[^\t]*)\t(?P<x_forwarded_for>[^\t]*)\t(?P<ssl_protocol>[^\t]*)\t(?P<ssl_cipher>[^\t]*)\t(?P<x_edge_response_result_type>[^\t]*)\t(?P<cs_protocol_version>[^\t]*)\t(?P<fle_status>[^\t]*)\t(?P<fle_encrypted_fields>[^\s]*)`)).
		WithKindMap(map[string]string{
			"timestamp":       "time:2006-01-02\t15:04:05",
			"x_edge_location": "urlencoded",
			"sc_bytes":        "uint64",
			"cs_host":         "urlencoded",
			"cs_uri_stem":     "urlencoded",
			"sc_status":       "int16",
			"cs_referer":      "urlencoded",
			"cs_user_agent":   "urlencoded",
			"cs_uri_query":    "urlencoded",
			"cs_cookie":       "urlencoded",
			"time_taken":      "float64",
		}).
		WithReIgnore(regexp.MustCompile(`^#`)).
		WithEmptyValues(map[string]string{
			"cs_uri_query":         "-",
			"cs_bytes":             "-",
			"x_forwarded_for":      "-",
			"ssl_protocol":         "-",
			"ssl_cipher":           "-",
			"fle_status":           "-",
			"fle_encrypted_fields": "-",
		})
)
