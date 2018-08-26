package logparser

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type kind int

const (
	kindBool kind = iota
	kindInt
	kindInt8
	kindInt16
	kindInt32
	kindInt64
	kindUint
	kindUint8
	kindUint16
	kindUint32
	kindUint64
	kindFloat32
	kindFloat64
	kindString
	kindURLEncoded

	kindTimeISO8601
	kindTimeLayout // based on https://golang.org/pkg/time/#Parse

	// aliases
	kindByte = kindUint8
	kindRune = kindInt32
)

type kindElement struct {
	kind      kind
	kindExtra interface{}
	name      string
}

var (
	kindElements = []kindElement{
		kindElement{
			kind: kindBool,
			name: "bool",
		},
		kindElement{
			kind: kindInt,
			name: "int",
		},
		kindElement{
			kind: kindInt8,
			name: "int8",
		},
		kindElement{
			kind: kindInt16,
			name: "int16",
		},
		kindElement{
			kind: kindInt32,
			name: "int32",
		},
		kindElement{
			kind: kindInt64,
			name: "int64",
		},
		kindElement{
			kind: kindUint,
			name: "uint",
		},
		kindElement{
			kind: kindUint8,
			name: "uint8",
		},
		kindElement{
			kind: kindUint16,
			name: "uint16",
		},
		kindElement{
			kind: kindUint32,
			name: "uint32",
		},
		kindElement{
			kind: kindUint64,
			name: "uint64",
		},
		kindElement{
			kind: kindFloat32,
			name: "float32",
		},
		kindElement{
			kind: kindFloat64,
			name: "float64",
		},
		kindElement{
			kind: kindString,
			name: "string",
		},
		kindElement{
			kind: kindURLEncoded,
			name: "urlencoded",
		},
		kindElement{
			kind: kindTimeISO8601,
			name: "timeISO8601",
		},
		// aliases
		kindElement{
			kind: kindByte,
			name: "byte",
		},
		kindElement{
			kind: kindRune,
			name: "rune",
		},
	}

	kindStringMap = func() map[string]kindElement {
		r := make(map[string]kindElement)
		for _, e := range kindElements {
			r[e.name] = e
		}
		return r
	}()

	kindMap = func() map[kind]kindElement {
		r := make(map[kind]kindElement)
		for _, e := range kindElements {
			r[e.kind] = e
		}
		return r
	}()
)

// CustomLogParser contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type CustomLogParser struct {
	re          *regexp.Regexp
	reIgnore    *regexp.Regexp
	reNames     []string
	reKindMap   map[string]kindElement
	emptyValues map[string]string
}

// NewCustomLogParser creates a new custom log parser based on regular expression
// to detect fields in a log line (re)
func NewCustomLogParser(re *regexp.Regexp) *CustomLogParser {
	return &CustomLogParser{
		re:      re,
		reNames: re.SubexpNames(),
	}
}

// WithKindMap configures current log parser to map types passed on reKindMap
func (c *CustomLogParser) WithKindMap(reKindMap map[string]string) *CustomLogParser {
	c.reKindMap = mustKindMapStringToType(reKindMap)
	return c
}

// WithReIgnore configures current log parser to ignore lines that match reIgnore
func (c *CustomLogParser) WithReIgnore(reIgnore *regexp.Regexp) *CustomLogParser {
	c.reIgnore = reIgnore
	return c
}

// WithEmptyValues configures current log parser to take into account emptyValues
func (c *CustomLogParser) WithEmptyValues(emptyValues map[string]string) *CustomLogParser {
	c.emptyValues = emptyValues
	return c
}

// Parse parses a reader and sends errors and parsed elements to handlers
func (c *CustomLogParser) Parse(reader io.Reader, mh logParserMessageHandler, eh logParserErrorHandler) error {
	r := bufio.NewReader(reader)
	re := c.re.Copy()
	var reIgnore *regexp.Regexp
	if c.reIgnore != nil {
		reIgnore = c.reIgnore.Copy()
	}
LINE_READER:
	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if !isLineIgnored(&line, reIgnore) {
			match := re.FindStringSubmatch(line)
			if match == nil {
				eh(line, fmt.Errorf("Line does not match expected format"))
			} else {
				captures := common.MapStr{}
				for i, name := range c.reNames {
					// Ignore the whole regexp match and unnamed groups
					if i == 0 || name == "" {
						continue
					}

					if k, ok := c.reKindMap[name]; ok {
						if v, err := parseStringToKind(k, match[i]); err != nil {
							eh(line, fmt.Errorf("Couldn't parse field (%s) to type (%s). Error: %+v", name, k.name, err))
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

func isLineIgnored(line *string, reIgnore *regexp.Regexp) bool {
	if *line == "" || *line == "\n" {
		return true
	}
	if reIgnore != nil {
		return reIgnore.MatchString(*line)
	}
	return false
}

func mustKindMapStringToType(o map[string]string) map[string]kindElement {
	r, err := kindMapStringToType(o)
	if err != nil {
		panic(`logparser: KindMapStringToType error: ` + err.Error())
	}
	return r
}

// KindMapStringToType obtains a map[string]kindElement from a
// map[string]string or an error if kind is not supported
func kindMapStringToType(o map[string]string) (map[string]kindElement, error) {
	r := make(map[string]kindElement)
	for k, v := range o {
		if kind, ok := kindStringMap[v]; ok {
			r[k] = kind
		} else if strings.HasPrefix(v, "time:") {
			timeLayout := strings.TrimPrefix(v, "time:")
			r[k] = kindElement{
				kind:      kindTimeLayout,
				kindExtra: timeLayout,
				name:      fmt.Sprintf("time layout (%s)", timeLayout),
			}
		} else {
			return nil, fmt.Errorf("Unsupported kind (%s)", k)
		}
	}
	return r, nil
}

func parseStringToKind(e kindElement, value string) (interface{}, error) {
	switch e.kind {
	case kindTimeLayout:
		return time.Parse(e.kindExtra.(string), value)
	case kindTimeISO8601:
		return time.Parse(time.RFC3339Nano, value)
	case kindBool:
		return strconv.ParseBool(value)
	case kindInt8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(v), nil
	case kindInt16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(v), nil
	case kindInt:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return int(v), nil
	case kindInt32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(v), nil
	case kindInt64:
		return strconv.ParseInt(value, 10, 64)
	case kindUint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, err
		}
		return uint8(v), nil
	case kindUint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, err
		}
		return uint16(v), nil
	case kindUint:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint(v), nil
	case kindUint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(v), nil
	case kindUint64:
		return strconv.ParseUint(value, 10, 64)
	case kindFloat32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(v), nil
	case kindFloat64:
		return strconv.ParseFloat(value, 64)
	case kindURLEncoded:
		return url.QueryUnescape(value)
	}
	return value, nil
}
