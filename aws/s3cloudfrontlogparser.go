package aws

import (
	"regexp"

	"github.com/mpucholblasco/s3logsbeat/logparser"
)

var (

	//   status INT,
	//   referrer STRING,
	//   useragent STRING,
	//   querystring STRING,
	//   cookie STRING,
	//   resulttype STRING,
	//   requestid STRING,
	//   hostheader STRING,
	//   requestprotocol STRING,
	//   requestbytes BIGINT,
	//   timetaken FLOAT,
	//   xforwardedfor STRING,
	//   sslprotocol STRING,
	//   sslcipher STRING,
	//   responseresulttype STRING,
	//   httpversion STRING,
	//   filestatus STRING,
	//   encryptedfields INT
	//

	// #Fields: date time x-edge-location sc-bytes c-ip cs-method cs(Host) cs-uri-stem sc-status cs(Referer) cs(User-Agent) cs-uri-query cs(Cookie) x-edge-result-type x-edge-request-id x-host-header cs-protocol cs-bytes time-taken x-forwarded-for ssl-protocol ssl-cipher x-edge-response-result-type cs-protocol-version fle-status fle-encrypted-fields
	// 2014-05-23 01:13:11 FRA2 182 192.0.2.10 GET d111111abcdef8.cloudfront.net /view/my/file.html 200 www.displaymyfiles.com Mozilla/4.0%20(compatible;%20MSIE%205.0b1;%20Mac_PowerPC) - zip=98101 RefreshHit MRVMF7KydIvxMWfJIglgwHQwZsbG2IhRJ07sn9AkKUFSHS9EXAMPLE== d111111abcdef8.cloudfront.net http - 0.001 - - - RefreshHit HTTP/1.1 Processed 1

	S3CloudFrontWebLogParser = logparser.NewCustomLogParser(regexp.MustCompile(`^(?P<timestamp>[^\t]*\t[^\t]*)\t(?P<x_edge_location>[^\t]*)\t(?P<sc_bytes>[^\t]*)\t(?P<c_ip>[^\t]*)\t(?P<cs_method>[^\t]*)\t(?P<cs_host>[^\t]*)\t(?P<cs_uri_stem>[^\t]*)\t(?P<sc_status>[^\t]*)\t(?P<cs_referer>[^\t]*)\t(?P<cs_user_agent>[^\t]*)\t`)).
		WithKindMap(logparser.MustKindMapStringToType(map[string]string{
			"timestamp":       "time:2006-01-02\t15:04:05",
			"x_edge_location": "urlencoded",
			"sc_bytes":        "uint64",
			"cs_host":         "urlencoded",
			"cs_uri_stem":     "urlencoded",
			"sc_status":       "int16",
			"cs_referer":      "urlencoded",
			"cs_user_agent":   "urlencoded",
		})).
		WithReIgnore(regexp.MustCompile(`^#`))
)
