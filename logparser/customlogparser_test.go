// +build !integration

package logparser

import (
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseErrorLines(t *testing.T) {
	logs := `str1 not-a-valid-date 35325 120 30123 true str2 0.325 0.0318353`
	expected := []common.MapStr{}
	expectedErrorsPrefix := []string{
		`Couldn't parse field (time) to type (timeISO8601). Error: parsing time "not-a-valid-date"`,
	}

	parser := NewCustomLogParser(regexTest).WithKindMap(regexKind)
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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
	AssertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
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

type testReader struct {
	reader io.Reader
}

func (a *testReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("my custom error")
}
