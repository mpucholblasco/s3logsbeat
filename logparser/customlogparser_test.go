// +build !integration

package logparser

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/mpucholblasco/s3logsbeat/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	regexTest = regexp.MustCompile(`^(?P<string>[^ ]*) (?P<time>[^ ]*) (?P<int>[-0-9]*) (?P<int8>[-0-9]*) (?P<int16>[-0-9]*) (?P<bool>[^ ]*) (?P<string2>[^ ]*) (?P<float32>[-.0-9]*) (?P<float64>[-.0-9]*)`)
	regexKind = map[string]string{
		"time":    "timeISO8601",
		"int":     "int",
		"int8":    "int8",
		"int16":   "int16",
		"bool":    "bool",
		"string2": "string",
		"float32": "float32",
		"float64": "float64",
	}
)

// isLineIgnored
func TestIsLineIgnored(t *testing.T) {
	reIgnore := regexp.MustCompile(`^\s*#`)
	line := "\n"
	assert.True(t, isLineIgnored(&line, reIgnore))
	line = ""
	assert.True(t, isLineIgnored(&line, reIgnore))
	line = "# comment matching ignore line"
	assert.True(t, isLineIgnored(&line, reIgnore))
	line = "Line not being ignored"
	assert.False(t, isLineIgnored(&line, reIgnore))
}

// KindMapStringToType & MustKindMapStringToType tests
func TestCustomLogParserKindMapStringToTypeCorrect(t *testing.T) {
	m := map[string]string{
		"timeLayout": "time:2006-01-02\t15:04:05",
		"time":       "timeISO8601",
		"int":        "int",
		"int8":       "int8",
		"int16":      "int16",
		"bool":       "bool",
		"string2":    "string",
		"float32":    "float32",
		"float64":    "float64",
		"urlencoded": "urlencoded",
	}
	expected := map[string]kindElement{
		"timeLayout": kindElement{kind: kindTimeLayout, kindExtra: "2006-01-02\t15:04:05", name: "time layout (2006-01-02\t15:04:05)"},
		"time":       kindMap[kindTimeISO8601],
		"int":        kindMap[kindInt],
		"int8":       kindMap[kindInt8],
		"int16":      kindMap[kindInt16],
		"bool":       kindMap[kindBool],
		"string2":    kindMap[kindString],
		"float32":    kindMap[kindFloat32],
		"float64":    kindMap[kindFloat64],
		"urlencoded": kindMap[kindURLEncoded],
	}

	value, err := kindMapStringToType(m)
	assert.NoError(t, err)
	assert.Equal(t, expected, value)
	assert.NotPanics(t, func() {
		mustKindMapStringToType(m)
	})
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

	value, err := kindMapStringToType(m)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.Panics(t, func() {
		mustKindMapStringToType(m)
	})
}

// parseStringToKind tests
func TestCustomLogParserParseStringToKindsWithNoErrors(t *testing.T) {
	type elem struct {
		kind     kindElement
		strValue string
		value    interface{}
	}
	elems := []elem{
		elem{
			kind:     kindElement{kind: kindTimeLayout, kindExtra: "2006-01-02\t15:04:05"},
			strValue: "2014-05-23\t01:15:18",
			value:    time.Date(2014, 5, 23, 1, 15, 18, 0, time.UTC),
		},
		elem{
			kind:     kindMap[kindTimeISO8601],
			strValue: "2016-08-10T22:08:42.945958Z",
			value:    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
		},
		elem{
			kind:     kindMap[kindBool],
			strValue: "true",
			value:    true,
		},
		elem{
			kind:     kindMap[kindInt8],
			strValue: "5",
			value:    int8(5),
		},
		elem{
			kind:     kindMap[kindInt16],
			strValue: "32000",
			value:    int16(32000),
		},
		elem{
			kind:     kindMap[kindInt],
			strValue: "67353",
			value:    int(67353),
		},
		elem{
			kind:     kindMap[kindInt32],
			strValue: "67353",
			value:    int32(67353),
		},
		elem{
			kind:     kindMap[kindInt64],
			strValue: "-35868395685",
			value:    int64(-35868395685),
		},
		elem{
			kind:     kindMap[kindUint8],
			strValue: "250",
			value:    uint8(250),
		},
		elem{
			kind:     kindMap[kindUint16],
			strValue: "32000",
			value:    uint16(32000),
		},
		elem{
			kind:     kindMap[kindUint],
			strValue: "835000",
			value:    uint(835000),
		},
		elem{
			kind:     kindMap[kindUint32],
			strValue: "835000",
			value:    uint32(835000),
		},
		elem{
			kind:     kindMap[kindUint64],
			strValue: "35868395685",
			value:    uint64(35868395685),
		},
		elem{
			kind:     kindMap[kindFloat32],
			strValue: "0.385694",
			value:    float32(0.385694),
		},
		elem{
			kind:     kindMap[kindFloat64],
			strValue: "0.38569355355334",
			value:    0.38569355355334,
		},
		elem{
			kind:     kindMap[kindString],
			strValue: "This is a string",
			value:    "This is a string",
		},
		elem{
			kind:     kindMap[kindURLEncoded],
			strValue: "Mozilla/4.0%20(compatible;%20MSIE%207.0;%20Windows%20NT%205.1)",
			value:    "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
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
		kind     kindElement
		strValue string
	}
	elems := []elem{
		elem{
			kind:     kindElement{kind: kindTimeLayout, kindExtra: "2006-01-02 15:04:05"},
			strValue: "3 Feb 2014 01:14:18",
		},
		elem{
			kind:     kindMap[kindTimeISO8601],
			strValue: "true",
		},
		elem{
			kind:     kindMap[kindBool],
			strValue: "3",
		},
		elem{
			kind:     kindMap[kindInt8],
			strValue: "53535",
		},
		elem{
			kind:     kindMap[kindInt16],
			strValue: "jo",
		},
		elem{
			kind:     kindMap[kindInt],
			strValue: "true",
		},
		elem{
			kind:     kindMap[kindInt32],
			strValue: "false",
		},
		elem{
			kind:     kindMap[kindInt64],
			strValue: "none",
		},
		elem{
			kind:     kindMap[kindUint8],
			strValue: "-35",
		},
		elem{
			kind:     kindMap[kindUint16],
			strValue: "-5",
		},
		elem{
			kind:     kindMap[kindUint],
			strValue: "false",
		},
		elem{
			kind:     kindMap[kindUint32],
			strValue: "true",
		},
		elem{
			kind:     kindMap[kindUint64],
			strValue: "-3235",
		},
		elem{
			kind:     kindMap[kindFloat32],
			strValue: "false",
		},
		elem{
			kind:     kindMap[kindFloat64],
			strValue: "true",
		},
		elem{
			kind:     kindMap[kindURLEncoded],
			strValue: "a%5Z",
		},
	}

	for _, e := range elems {
		_, err := parseStringToKind(e.kind, e.strValue)
		assert.Error(t, err)
	}
}

// Parse tests
func TestCustomLogParserParseSingleLine(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353`
	expected := []common.MapStr{
		common.MapStr{
			"string":  "str1",
			"time":    "2016-08-10T22:08:42.945958Z",
			"int":     "35325",
			"int8":    "120",
			"int16":   "30123",
			"bool":    "true",
			"string2": "str2",
			"float32": "0.325",
			"float64": "0.0318353",
		},
	}

	parser := NewCustomLogParser(regexTest)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseSingleLineWithKindMap(t *testing.T) {
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

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseSingleLineWithEmtpyValues(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z - 120 30123 true str2 - 0.0318353`
	expected := []common.MapStr{
		common.MapStr{
			"string":  "str1",
			"time":    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			"int8":    int8(120),
			"int16":   int16(30123),
			"bool":    true,
			"string2": "str2",
			"float64": 0.0318353,
		},
	}

	emptyValues := map[string]string{
		"float32": "-",
		"int":     "-",
		"int8":    "-",
	}

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind).WithEmptyValues(emptyValues)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLines(t *testing.T) {
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

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLinesWithIgnoredLines(t *testing.T) {
	logs := `# comment, ignored
str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353

    # another ignored line
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

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind).WithReIgnore(regexp.MustCompile(`^\s*#`))
	expectedErrorsPrefix := []string{}
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseErrorLines(t *testing.T) {
	logs := `str1 not-a-valid-date 35325 120 30123 true str2 0.325 0.0318353`
	expected := []common.MapStr{}
	expectedErrorsPrefix := []string{
		`Couldn't parse field (time) to type (timeISO8601). Error: parsing time "not-a-valid-date"`,
	}

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
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

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
	testCustomLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserNothingProcessedOnReaderError(t *testing.T) {
	ok := 0
	ko := 0
	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
	err := parser.Parse(&testReader{}, func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.Error(t, err)
	assert.Equal(t, 0, ok)
	assert.Equal(t, 0, ko)
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
		testutil.AssertEvent(t, expEvent, resultEvent)
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

type testReader struct {
	reader io.Reader
}

func (a *testReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("my custom error")
}
