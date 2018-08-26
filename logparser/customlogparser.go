package logparser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/elastic/beats/libbeat/common"
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

// Copy generates a new CustomLogParser from current one
func (c *CustomLogParser) Copy() *CustomLogParser {
	r := &CustomLogParser{
		re:          c.re.Copy(),
		reIgnore:    c.re.Copy(),
		reKindMap:   make(map[string]kindElement),
		emptyValues: make(map[string]string),
	}
	copy(r.reNames, c.reNames)
	for k, v := range c.reKindMap {
		r.reKindMap[k] = v
	}
	for k, v := range c.emptyValues {
		r.emptyValues[k] = v
	}
	return r
}

// WithKindMap configures current log parser to map types passed on reKindMap
func (c *CustomLogParser) WithKindMap(reKindMap map[string]string) *CustomLogParser {
	c.reKindMap = mustKindMapStringToType(reKindMap)
	return c
}

// SetKindMap configures current log parser to map types passed on reKindMap and
// returns error (if present)
func (c *CustomLogParser) SetKindMap(reKindMap map[string]string) error {
	var err error
	c.reKindMap, err = kindMapStringToType(reKindMap)
	return err
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
func (c *CustomLogParser) Parse(reader io.Reader, mh func(common.MapStr), eh func(string, error)) error {
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

					if emptyValue, ok := c.emptyValues[name]; !ok || emptyValue != match[i] {
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
