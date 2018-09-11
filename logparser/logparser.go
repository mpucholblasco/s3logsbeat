package logparser

import (
	"fmt"
	"io"

	"github.com/elastic/beats/libbeat/beat"
)

// LogParser interface to inherit on each type of log parsers
type LogParser interface {
	Parse(io.Reader, func(*beat.Event), func(string, error)) error
}

// GetPredefinedParser gets a predefined parser based on its name
func GetPredefinedParser(n string) (LogParser, error) {
	switch n {
	case "alb":
		return S3ALBLogParser, nil
	case "cloudfront":
		return S3CloudFrontWebLogParser, nil
	}
	return nil, fmt.Errorf("Predefined parser %s not found", n)
}
