package aws

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

var (
	regexALB      = regexp.MustCompile(`^(?P<type>[^ ]*) (?P<time>[^ ]*) (?P<elb>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) (?P<target_ip>[^ ]*)[:-](?P<target_port>[0-9]*) (?P<request_processing_time>[-.0-9]*) (?P<target_processing_time>[-.0-9]*) (?P<response_processing_time>[-.0-9]*) (?P<elb_status_code>|[-0-9]*) (?P<target_status_code>-|[-0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) \"(?P<request_verb>[^ ]*) (?P<request_url>[^ ]*) (?P<request_proto>- |[^ ]*)\" \"(?P<user_agent>[^\"]*)\" (?P<ssl_cipher>[A-Z0-9-]+) (?P<ssl_protocol>[A-Za-z0-9.-]*) (?P<target_group_arn>[^ ]*) \"(?P<trace_id>[^\"]*)\"`)
	namesRegexALB = regexALB.SubexpNames()

	// TODO: convert types to type.Type (enum)
	typesRegexALB = map[string]string{
		"time":                     "time.RFC3339Nano",
		"client_port":              "uint16",
		"target_port":              "uint16",
		"request_processing_time":  "float64",
		"target_processing_time":   "float64",
		"response_processing_time": "float64",
		"received_bytes":           "int64",
		"sent_bytes":               "int64",
	}
)

// S3ObjectSQSMessage contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type S3ALBLogParser struct{}

func NewS3ALBLogParser() *S3ALBLogParser {
	return &S3ALBLogParser{}
}

// ExtractNewS3Objects extracts those new S3 objects present on an SQS message
func (l *S3ALBLogParser) Parse(reader io.Reader, mh s3LogParserMessageHandler, eh s3LogParserErrorHandler) error {
	r := bufio.NewReader(reader)
	re := regexALB.Copy()
LINE_READER:
	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if line != "" {
			match := re.FindStringSubmatch(line)
			if match == nil {
				eh(line, fmt.Errorf("Line does not match ALB log format"))
			} else {
				captures := common.MapStr{}
				for i, name := range namesRegexALB {
					// Ignore the whole regexp match and unnamed groups
					if i == 0 || name == "" {
						continue
					}

					if kind, ok := typesRegexALB[name]; ok {
						if v, err := parseStringToKind(kind, match[i]); err != nil {
							eh(line, fmt.Errorf("Couldn't parse field (%s) to type (%s). Error: %+v", name, kind, err))
							continue LINE_READER
						} else {
							captures.Put(name, v)
						}
					} else {
						captures.Put(name, match[i])
					}
				}
				mh(captures)
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}

func parseStringToKind(kind string, value string) (interface{}, error) {
	switch kind {
	case "time.RFC3339Nano":
		return time.Parse(time.RFC3339Nano, value)
	case "bool":
		return strconv.ParseBool(value)
	case "int8":
		if v, err := strconv.ParseInt(value, 10, 8); err != nil {
			return nil, err
		} else {
			return int8(v), nil
		}
	case "int16":
		if v, err := strconv.ParseInt(value, 10, 16); err != nil {
			return nil, err
		} else {
			return int16(v), nil
		}
	case "int":
		if v, err := strconv.ParseInt(value, 10, 32); err != nil {
			return nil, err
		} else {
			return int(v), nil
		}
	case "int32":
		if v, err := strconv.ParseInt(value, 10, 32); err != nil {
			return nil, err
		} else {
			return int32(v), nil
		}
	case "int64":
		return strconv.ParseInt(value, 10, 64)
	case "uint8":
		if v, err := strconv.ParseUint(value, 10, 8); err != nil {
			return nil, err
		} else {
			return uint8(v), nil
		}
	case "uint16":
		if v, err := strconv.ParseUint(value, 10, 16); err != nil {
			return nil, err
		} else {
			return uint16(v), nil
		}
	case "uint":
		if v, err := strconv.ParseUint(value, 10, 32); err != nil {
			return nil, err
		} else {
			return uint(v), nil
		}
	case "uint32":
		if v, err := strconv.ParseUint(value, 10, 32); err != nil {
			return nil, err
		} else {
			return uint32(v), nil
		}
	case "uint64":
		return strconv.ParseUint(value, 10, 64)
	case "float32":
		if v, err := strconv.ParseFloat(value, 32); err != nil {
			return nil, err
		} else {
			return float32(v), nil
		}
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "string":
		return value, nil
	default:
		return nil, fmt.Errorf("Can not convert to unsupported kind (%s)", kind)
	}
}
