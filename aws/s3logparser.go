package aws

import (
	"io"

	"github.com/elastic/beats/libbeat/common"
)

type s3LogParserMessageHandler func(common.MapStr)
type s3LogParserErrorHandler func(string, error)

// S3LogParser interface to inherit on each type of S3 log parsers
type S3LogParser interface {
	Parse(io.Reader, s3LogParserMessageHandler, s3LogParserErrorHandler)
}
