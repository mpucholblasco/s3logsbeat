package logparser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/jsontransform"
)

// JSONLogParserConfig JSONLogParser configuration
type JSONLogParserConfig struct {
	TimestampField string `config:"timestamp_field" validate:"required"`
	TimestampKing  string `config:"timestamp_format" validate:"required"`
}

// JSONLogParser JSON log parser
type JSONLogParser struct {
	timestampField string
	timestampKind  kindElement
}

// NewJSONLogParserConfig creates a new JSON log parser based on a map os strins
func NewJSONLogParserConfig(cfg *common.Config) (*JSONLogParser, error) {
	var config JSONLogParserConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, err
	}

	timestampKind, err := kindFromString(config.TimestampKing)
	if err != nil {
		return nil, err
	}

	return NewJSONLogParser(config.TimestampField, timestampKind), nil
}

// NewJSONLogParser creates a new JSON log parser
func NewJSONLogParser(timestampField string, timestampKind kindElement) *JSONLogParser {
	return &JSONLogParser{
		timestampField: timestampField,
		timestampKind:  timestampKind,
	}
}

// Parse parses a reader and sends errors and parsed elements to handlers
func (j *JSONLogParser) Parse(reader io.Reader, mh func(*beat.Event), eh func(string, error)) error {
	r := bufio.NewReader(reader)
LINE_READER:
	for {
		line, errReadString := r.ReadString('\n')
		if errReadString != nil && errReadString != io.EOF {
			return errReadString
		}

		if line != "" && line != "\n" {
			var fields map[string]interface{}
			if err := unmarshal([]byte(line), &fields); err != nil {
				eh(line, fmt.Errorf("Couldn't parse json line (%s). Error: %+v", line, err))
				continue LINE_READER
			}

			timestamp, err := j.getTimestamp(&fields)
			if err != nil {
				eh(line, err)
				continue LINE_READER
			}
			delete(fields, j.timestampField)

			event := CreateEvent(&line, timestamp, fields)
			mh(event)
		}

		if errReadString == io.EOF {
			break
		}
	}
	return nil
}

func (j *JSONLogParser) getTimestamp(fields *map[string]interface{}) (time.Time, error) {
	timestampValue, found := (*fields)[j.timestampField]
	if !found {
		return time.Time{}, fmt.Errorf("Couldn't find timestamp field %s", j.timestampField)
	}

	v, err := parseToKind(j.timestampKind, timestampValue)
	if err != nil {
		return time.Time{}, err
	}

	timestamp, ok := v.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("Field %s set as timestamp, but it's kind is not time", j.timestampField)
	}
	return timestamp, nil
}

// unmarshal is equivalent with json.Unmarshal but it converts numbers
// to int64 where possible, instead of using always float64.
func unmarshal(text []byte, fields *map[string]interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(text))
	dec.UseNumber()
	err := dec.Decode(fields)
	if err != nil {
		return err
	}
	jsontransform.TransformNumbers(*fields)
	return nil
}
