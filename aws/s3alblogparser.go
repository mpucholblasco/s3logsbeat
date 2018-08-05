package aws

import (
	"bufio"
	"io"
	"regexp"

	"github.com/elastic/beats/libbeat/common"
)

var (
	regexALB      = regexp.MustCompile(`^$`)
	namesRegexALB = regexALB.SubexpNames()
	// Types are important, so we have to convert to this type
)

// S3ObjectSQSMessage contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type S3ALBLogParser struct{}

func NewS3ALBLogParser() *S3ALBLogParser {
	return &S3ALBLogParser{}
}

// ExtractNewS3Objects extracts those new S3 objects present on an SQS message
func (l *S3ALBLogParser) parse(reader io.Reader, mh s3LogParserMessageHandler, eh s3LogParserErrorHandler) error {
	r := bufio.NewReader(reader)
	re := regexALB.Copy()
	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		match := re.FindStringSubmatch(line)
		if match == nil {
			eh(line)
		} else {
			captures := common.MapStr{}
			for i, name := range namesRegexALB {
				// Ignore the whole regexp match and unnamed groups
				if i == 0 || name == "" {
					continue
				}
				captures.Put(name, match[i])
			}
			mh(captures)
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}
