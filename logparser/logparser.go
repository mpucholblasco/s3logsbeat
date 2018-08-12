package logparser

import (
	"io"

	"github.com/elastic/beats/libbeat/common"
)

type logParserMessageHandler func(common.MapStr)
type logParserErrorHandler func(string, error)

// LogParser interface to inherit on each type of log parsers
type LogParser interface {
	Parse(io.Reader, logParserMessageHandler, logParserErrorHandler) error
}
