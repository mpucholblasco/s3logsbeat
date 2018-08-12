// +build !integration

package logparser

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

var (
	regexTest = regexp.MustCompile(`^(?P<string>[^ ]*) (?P<time>[^ ]*) (?P<int>[-0-9]*) (?P<int8>[-0-9]*) (?P<int16>[-0-9]*) (?P<bool>[^ ]*) (?P<string2>[^ ]*) (?P<float32>[-.0-9]*) (?P<float64>[-.0-9]*)`)
	regexKind = map[string]Kind{
		"time":    TimeISO8601,
		"int":     Int,
		"int8":    Int8,
		"int16":   Int16,
		"bool":    Bool,
		"string2": String,
		"float32": Float32,
		"float64": Float64,
	}
)

// KindMapStringToType tests
func TestCustomLogParserKindMapStringToTypeCorrect(t *testing.T) {
	m := map[string]string{
		"time":    "timeISO8601",
		"int":     "int",
		"int8":    "int8",
		"int16":   "int16",
		"bool":    "bool",
		"string2": "string",
		"float32": "float32",
		"float64": "float64",
	}
	expected := map[string]Kind{
		"time":    TimeISO8601,
		"int":     Int,
		"int8":    Int8,
		"int16":   Int16,
		"bool":    Bool,
		"string2": String,
		"float32": Float32,
		"float64": Float64,
	}

	value, err := KindMapStringToType(m)
	assert.NoError(t, err)
	for idx, e := range value {
		assert.Equal(t, expected[idx], e.kind)
	}
}

func TestCustomLogParserKindMapStringToTypeUnsupportedType(t *testing.T) {
	m := map[string]string{
		"time":    "timeISO8601",
		"int8":    "int8",
		"int16":   "int16",
		"bool":    "bool",
		"string2": "string",
		"int":     "unsupportedType",
		"float32": "float32",
		"float64": "float64",
	}

	value, err := KindMapStringToType(m)
	assert.Nil(t, value)
	assert.Error(t, err)
}

// KindMapKindToType tests
func TestCustomLogParserKindMapKindToTypeCorrect(t *testing.T) {
	m := map[string]Kind{
		"time":    TimeISO8601,
		"int":     Int,
		"int8":    Int8,
		"int16":   Int16,
		"bool":    Bool,
		"string2": String,
		"float32": Float32,
		"float64": Float64,
	}

	expected := map[string]string{
		"time":    "timeISO8601",
		"int":     "int",
		"int8":    "int8",
		"int16":   "int16",
		"bool":    "bool",
		"string2": "string",
		"float32": "float32",
		"float64": "float64",
	}

	value, err := KindMapKindToType(m)
	assert.NoError(t, err)
	for idx, e := range value {
		assert.Equal(t, expected[idx], e.name)
	}
}

// Parse tests
func TestCustomLogParserParseSingleLine(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353`
	expected := []common.MapStr{
		common.MapStr{
			"string":  "str1",
			"time":    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			"int":     int(35325),
			"int8":    int8(120),
			"int16":   int16(30123),
			"bool":    true,
			"string2": "str2",
			"float32": float32(0.325),
			"float64": 0.0318353,
		},
	}

	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLine(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353
strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555
strLine3 2006-08-13T02:08:12.544953Z 12345 05 31123 true str2 0.111 0.123456`

	expected := []common.MapStr{
		common.MapStr{
			"string":  "str1",
			"time":    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			"int":     int(35325),
			"int8":    int8(120),
			"int16":   int16(30123),
			"bool":    true,
			"string2": "str2",
			"float32": float32(0.325),
			"float64": 0.0318353,
		},
		common.MapStr{
			"string":  "strLine2",
			"time":    time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			"int":     int(321345),
			"int8":    int8(25),
			"int16":   int16(27535),
			"bool":    false,
			"string2": "str2Line2",
			"float32": float32(0.312),
			"float64": 0.323454555,
		},
		common.MapStr{
			"string":  "strLine3",
			"time":    time.Date(2006, 8, 13, 2, 8, 12, 544953000, time.UTC),
			"int":     int(12345),
			"int8":    int8(5),
			"int16":   int16(31123),
			"bool":    true,
			"string2": "str2",
			"float32": float32(0.111),
			"float64": 0.123456,
		},
	}

	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLineWithEmptyLines(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353


strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555


strLine3 2006-08-13T02:08:12.544953Z 12345 05 31123 true str2 0.111 0.123456

`

	expected := []common.MapStr{
		common.MapStr{
			"string":  "str1",
			"time":    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			"int":     int(35325),
			"int8":    int8(120),
			"int16":   int16(30123),
			"bool":    true,
			"string2": "str2",
			"float32": float32(0.325),
			"float64": 0.0318353,
		},
		common.MapStr{
			"string":  "strLine2",
			"time":    time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			"int":     int(321345),
			"int8":    int8(25),
			"int16":   int16(27535),
			"bool":    false,
			"string2": "str2Line2",
			"float32": float32(0.312),
			"float64": 0.323454555,
		},
		common.MapStr{
			"string":  "strLine3",
			"time":    time.Date(2006, 8, 13, 2, 8, 12, 544953000, time.UTC),
			"int":     int(12345),
			"int8":    int8(5),
			"int16":   int16(31123),
			"bool":    true,
			"string2": "str2",
			"float32": float32(0.111),
			"float64": 0.123456,
		},
	}

	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseErrorLines(t *testing.T) {
	logs := `str1 not-a-valid-date 35325 120 30123 true str2 0.325 0.0318353`
	expected := []common.MapStr{}
	expectedErrorsPrefix := []string{
		`Couldn't parse field (time) to type (timeISO8601). Error: parsing time "not-a-valid-date"`,
	}

	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserInvalidFormat(t *testing.T) {
	logs := `Incorrect Line
strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555
Incorrect line2
`
	expected := []common.MapStr{
		common.MapStr{
			"string":  "strLine2",
			"time":    time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			"int":     int(321345),
			"int8":    int8(25),
			"int16":   int16(27535),
			"bool":    false,
			"string2": "str2Line2",
			"float32": float32(0.312),
			"float64": 0.323454555,
		},
	}
	expectedErrorsPrefix := []string{
		"Line does not match expected format",
		"Line does not match expected format",
	}

	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserNothingProcessedOnReaderError(t *testing.T) {
	ok := 0
	ko := 0
	k, _ := KindMapKindToType(regexKind)
	parser := NewCustomLogParser(regexTest, k)
	err := parser.Parse(&testReader{}, func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.Error(t, err)
	assert.Equal(t, 0, ok)
	assert.Equal(t, 0, ko)
}

// parseStringToKind tests
func TestCustomLogParserParseStringToKindsWithNoErrors(t *testing.T) {
	type elem struct {
		kind     Kind
		strValue string
		value    interface{}
	}
	elems := []elem{
		elem{
			kind:     TimeISO8601,
			strValue: "2016-08-10T22:08:42.945958Z",
			value:    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
		},
		elem{
			kind:     Bool,
			strValue: "true",
			value:    true,
		},
		elem{
			kind:     Int8,
			strValue: "5",
			value:    int8(5),
		},
		elem{
			kind:     Int16,
			strValue: "32000",
			value:    int16(32000),
		},
		elem{
			kind:     Int,
			strValue: "67353",
			value:    int(67353),
		},
		elem{
			kind:     Int32,
			strValue: "67353",
			value:    int32(67353),
		},
		elem{
			kind:     Int64,
			strValue: "-35868395685",
			value:    int64(-35868395685),
		},
		elem{
			kind:     Uint8,
			strValue: "250",
			value:    uint8(250),
		},
		elem{
			kind:     Uint16,
			strValue: "32000",
			value:    uint16(32000),
		},
		elem{
			kind:     Uint,
			strValue: "835000",
			value:    uint(835000),
		},
		elem{
			kind:     Uint32,
			strValue: "835000",
			value:    uint32(835000),
		},
		elem{
			kind:     Uint64,
			strValue: "35868395685",
			value:    uint64(35868395685),
		},
		elem{
			kind:     Float32,
			strValue: "0.385694",
			value:    float32(0.385694),
		},
		elem{
			kind:     Float64,
			strValue: "0.38569355355334",
			value:    0.38569355355334,
		},
		elem{
			kind:     String,
			strValue: "This is a string",
			value:    "This is a string",
		},
	}

	for _, e := range elems {
		v, err := parseStringToKind(e.kind, e.strValue)
		assert.NoError(t, err)
		assert.Equal(t, e.value, v)
	}
}

func TestCustomLogParserParseStringToKindsWithParseErrors(t *testing.T) {
	type elem struct {
		kind     Kind
		strValue string
	}
	elems := []elem{
		elem{
			kind:     TimeISO8601,
			strValue: "true",
		},
		elem{
			kind:     Bool,
			strValue: "3",
		},
		elem{
			kind:     Int8,
			strValue: "53535",
		},
		elem{
			kind:     Int16,
			strValue: "jo",
		},
		elem{
			kind:     Int,
			strValue: "true",
		},
		elem{
			kind:     Int32,
			strValue: "false",
		},
		elem{
			kind:     Int64,
			strValue: "none",
		},
		elem{
			kind:     Uint8,
			strValue: "-35",
		},
		elem{
			kind:     Uint16,
			strValue: "-5",
		},
		elem{
			kind:     Uint,
			strValue: "false",
		},
		elem{
			kind:     Uint32,
			strValue: "true",
		},
		elem{
			kind:     Uint64,
			strValue: "-3235",
		},
		elem{
			kind:     Float32,
			strValue: "false",
		},
		elem{
			kind:     Float64,
			strValue: "true",
		},
	}

	for _, e := range elems {
		_, err := parseStringToKind(e.kind, e.strValue)
		assert.Error(t, err)
	}
}

func testCustomLogParser(t *testing.T, p LogParser, logs *string, expected []common.MapStr, expectedErrorsPrefix []string) {
	results := make([]common.MapStr, 0, len(expected))
	errors := make([]error, 0, len(expectedErrorsPrefix))
	err := p.Parse(strings.NewReader(*logs), func(s common.MapStr) {
		results = append(results, s)
	}, func(errLine string, err error) {
		errors = append(errors, err)
	})
	assert.NoError(t, err)
	assert.Len(t, errors, len(expectedErrorsPrefix))
	assert.Len(t, results, len(expected))
	for idx, expEvent := range expected {
		resultEvent := results[idx]
		assertEvent(t, expEvent, resultEvent)
	}
	for idx, expErr := range expectedErrorsPrefix {
		err := errors[idx]
		if !assert.True(t, strings.HasPrefix(err.Error(), expErr)) {
			t.Logf("expected error prefix: %s", expErr)
			t.Logf("      but found error: %s", err.Error())
			t.Logf("------------------------------")
		}
	}
}

func assertEvent(t *testing.T, expected, event common.MapStr) {
	for field, exp := range expected {
		val, found := event[field]
		if !found {
			t.Errorf("Missing field: %v", field)
			continue
		}

		if sub, ok := exp.(common.MapStr); ok {
			assertEvent(t, sub, val.(common.MapStr))
		} else {
			if !assert.Equal(t, exp, val) {
				t.Logf("failed in field: %v", field)
				t.Logf("type expected: %v", reflect.TypeOf(exp))
				t.Logf("type event: %v", reflect.TypeOf(val))
				t.Logf("------------------------------")
			}
		}
	}
}

type testReader struct {
	reader io.Reader
}

func (a *testReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("my custom error")
}
