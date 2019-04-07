package pipeline

import (
	"fmt"
	"regexp"

	"github.com/elastic/beats/libbeat/common"
	"github.com/sequra/s3logsbeat/logparser"
)

// S3ReaderInformation information present on inputs needed at S3 reader stage
type S3ReaderInformation struct {
	logParser      logparser.LogParser
	keyRegexFields *regexp.Regexp
	metadataType   string
}

// NewS3ReaderInformation creates a new S3 reader information
func NewS3ReaderInformation(logParser logparser.LogParser, keyRegexFields *regexp.Regexp, metadataType string) *S3ReaderInformation {
	return &S3ReaderInformation{
		logParser:      logParser,
		keyRegexFields: keyRegexFields,
		metadataType:   metadataType,
	}
}

// GetLogParser obtains the log parser
func (ri *S3ReaderInformation) GetLogParser() logparser.LogParser {
	return ri.logParser
}

// GetMetadataType obtains metadata type
func (ri *S3ReaderInformation) GetMetadataType() string {
	return ri.metadataType
}

// GetKeyFields extract fields from key if input set it
func (ri S3ReaderInformation) GetKeyFields(key string) (*common.MapStr, error) {
	f := &common.MapStr{}

	if ri.keyRegexFields != nil {
		re := ri.keyRegexFields.Copy()
		match := re.FindStringSubmatch(key)
		if match == nil {
			return nil, fmt.Errorf("Couldn't match key regex fields %s with S3 key %s", re.String(), key)
		}
		for i, name := range re.SubexpNames() {
			// Ignore the whole regexp match, unnamed groups, and empty values
			if i == 0 || name == "" || match[i] == "" {
				continue
			}
			f.Put(name, match[i])
		}
	}
	return f, nil
}
