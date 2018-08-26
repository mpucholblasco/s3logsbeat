package logparser

import (
	"io"

	"github.com/elastic/beats/libbeat/common"
)

// LogParser interface to inherit on each type of log parsers
type LogParser interface {
	Parse(io.Reader, func(common.MapStr), func(string, error)) error
}
