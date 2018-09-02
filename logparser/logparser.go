package logparser

import (
	"io"

	"github.com/elastic/beats/libbeat/beat"
)

// LogParser interface to inherit on each type of log parsers
type LogParser interface {
	Parse(io.Reader, func(beat.Event), func(string, error)) error
}
