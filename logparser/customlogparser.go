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

type Kind int

const (
	Bool Kind = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String
	UrlEncoded

	TimeISO8601
	TimeLayout // based on https://golang.org/pkg/time/#Parse

	// aliases
	Byte = Uint8
	Rune = Int32
)

type KindElement struct {
	kind      Kind
	kindExtra interface{}
	name      string
}

var (
	kindElements = []KindElement{
		KindElement{
			kind: Bool,
			name: "bool",
		},
		KindElement{
			kind: Int,
			name: "int",
		},
		KindElement{
			kind: Int8,
			name: "int8",
		},
		KindElement{
			kind: Int16,
			name: "int16",
		},
		KindElement{
			kind: Int32,
			name: "int32",
		},
		KindElement{
			kind: Int64,
			name: "int64",
		},
		KindElement{
			kind: Uint,
			name: "uint",
		},
		KindElement{
			kind: Uint8,
			name: "uint8",
		},
		KindElement{
			kind: Uint16,
			name: "uint16",
		},
		KindElement{
			kind: Uint32,
			name: "uint32",
		},
		KindElement{
			kind: Uint64,
			name: "uint64",
		},
		KindElement{
			kind: Float32,
			name: "float32",
		},
		KindElement{
			kind: Float64,
			name: "float64",
		},
		KindElement{
			kind: String,
			name: "string",
		},
		KindElement{
			kind: UrlEncoded,
			name: "urlencoded",
		},
		KindElement{
			kind: TimeISO8601,
			name: "timeISO8601",
		},
		// aliases
		KindElement{
			kind: Byte,
			name: "byte",
		},
		KindElement{
			kind: Rune,
			name: "rune",
		},
	}

	kindStringMap = func() map[string]KindElement {
		r := make(map[string]KindElement)
		for _, e := range kindElements {
			r[e.name] = e
		}
		return r
	}()

	kindMap = func() map[Kind]KindElement {
		r := make(map[Kind]KindElement)
		for _, e := range kindElements {
			r[e.kind] = e
		}
		return r
	}()
)

// CustomLogParser contains information of S3 objects (sqsMessage not
// null implies that this object is extracted from an SQS message)
type CustomLogParser struct {
	re        *regexp.Regexp
	reIgnore  *regexp.Regexp
	reNames   []string
	reKindMap map[string]KindElement
}

// NewCustomLogParser creates a new custom log parser based on regular expression
// to detect fields in a log line (re)
func NewCustomLogParser(re *regexp.Regexp) *CustomLogParser {
	return &CustomLogParser{
		re:      re,
		reNames: re.SubexpNames(),
	}
}

func (c *CustomLogParser) WithKindMap(reKindMap map[string]KindElement) *CustomLogParser {
	c.reKindMap = reKindMap
	return c
}

func (c *CustomLogParser) WithReIgnore(reIgnore *regexp.Regexp) *CustomLogParser {
	c.reIgnore = reIgnore
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

func MustKindMapStringToType(o map[string]string) map[string]KindElement {
	r, err := KindMapStringToType(o)
	if err != nil {
		panic(`logparser: KindMapStringToType error: ` + err.Error())
	}
	return r
}

func KindMapStringToType(o map[string]string) (map[string]KindElement, error) {
	r := make(map[string]KindElement)
	for k, v := range o {
		if kind, ok := kindStringMap[v]; ok {
			r[k] = kind
		} else if strings.HasPrefix(v, "time:") {
			timeLayout := strings.TrimPrefix(v, "time:")
			r[k] = KindElement{
				kind:      TimeLayout,
				kindExtra: timeLayout,
				name:      fmt.Sprintf("time layout (%s)", timeLayout),
			}
		} else {
			return nil, fmt.Errorf("Unsupported kind (%s)", k)
		}
	}
	return r, nil
}

func KindMapKindToType(o map[string]Kind) map[string]KindElement {
	r := make(map[string]KindElement)
	for k, v := range o {
		r[k], _ = kindMap[v]
	}
	return r
}

func parseStringToKind(e KindElement, value string) (interface{}, error) {
	switch e.kind {
	case TimeLayout:
		return time.Parse(e.kindExtra.(string), value)
	case TimeISO8601:
		return time.Parse(time.RFC3339Nano, value)
	case Bool:
		return strconv.ParseBool(value)
	case Int8:
		if v, err := strconv.ParseInt(value, 10, 8); err != nil {
			return nil, err
		} else {
			return int8(v), nil
		}
	case Int16:
		if v, err := strconv.ParseInt(value, 10, 16); err != nil {
			return nil, err
		} else {
			return int16(v), nil
		}
	case Int:
		if v, err := strconv.ParseInt(value, 10, 32); err != nil {
			return nil, err
		} else {
			return int(v), nil
		}
	case Int32:
		if v, err := strconv.ParseInt(value, 10, 32); err != nil {
			return nil, err
		} else {
			return int32(v), nil
		}
	case Int64:
		return strconv.ParseInt(value, 10, 64)
	case Uint8:
		if v, err := strconv.ParseUint(value, 10, 8); err != nil {
			return nil, err
		} else {
			return uint8(v), nil
		}
	case Uint16:
		if v, err := strconv.ParseUint(value, 10, 16); err != nil {
			return nil, err
		} else {
			return uint16(v), nil
		}
	case Uint:
		if v, err := strconv.ParseUint(value, 10, 32); err != nil {
			return nil, err
		} else {
			return uint(v), nil
		}
	case Uint32:
		if v, err := strconv.ParseUint(value, 10, 32); err != nil {
			return nil, err
		} else {
			return uint32(v), nil
		}
	case Uint64:
		return strconv.ParseUint(value, 10, 64)
	case Float32:
		if v, err := strconv.ParseFloat(value, 32); err != nil {
			return nil, err
		} else {
			return float32(v), nil
		}
	case Float64:
		return strconv.ParseFloat(value, 64)
	case UrlEncoded:
		return url.QueryUnescape(value)
	}
	return value, nil
}
