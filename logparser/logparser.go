package logparser

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

// LogParser interface to inherit on each type of log parsers
type LogParser interface {
	Parse(io.Reader, func(*beat.Event), func(string, error)) error
}

// GetPredefinedParser gets a predefined parser based on its name
func GetPredefinedParser(n string, config *common.Config) (LogParser, error) {
	switch n {
	case "elb":
		return S3ELBLogParser, nil
	case "alb":
		return S3ALBLogParser, nil
	case "cloudfront":
		return S3CloudFrontWebLogParser, nil
	case "waf":
		return S3WAFLogParser, nil
	case "json":
		return NewJSONLogParserConfig(config)
	}
	return nil, fmt.Errorf("Predefined parser %s not found", n)
}

// CreateEvent creates an event to be passed to elastic output
func CreateEvent(line *string, timestamp time.Time, fields common.MapStr) *beat.Event {
	h := sha1.New()
	io.WriteString(h, *line)
	return &beat.Event{
		Timestamp: timestamp,
		Fields:    fields,
		Meta: common.MapStr{
			"_id": hex.EncodeToString(h.Sum(nil)),
		},
	}
}
