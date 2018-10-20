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

	"github.com/elastic/beats/libbeat/beat"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
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
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":  "str1",
				"int":     "35325",
				"int8":    "120",
				"int16":   "30123",
				"bool":    "true",
				"string2": "str2",
				"float32": "0.325",
				"float64": "0.0318353",
			},
			Meta: common.MapStr{
				"_id": "2db0da2449bcdbf8f1838844107a490f2773e09e",
			},
		},
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(
		map[string]string{
			"time": "timeISO8601",
		},
	)
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseSingleLineWithKindMap(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":  "str1",
				"int":     int(35325),
				"int8":    int8(120),
				"int16":   int16(30123),
				"bool":    true,
				"string2": "str2",
				"float32": float32(0.325),
				"float64": 0.0318353,
			},
			Meta: common.MapStr{
				"_id": "2db0da2449bcdbf8f1838844107a490f2773e09e",
			},
		},
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind)
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseSingleLineWithEmtpyValues(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z - 120 30123 true str2 - 0.0318353`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":  "str1",
				"int8":    int8(120),
				"int16":   int16(30123),
				"bool":    true,
				"string2": "str2",
				"float64": 0.0318353,
			},
			Meta: common.MapStr{
				"_id": "0c2f4767adbc506b33eb9fd8dc7ad47bf18dfb50",
			},
		},
	}

	emptyValues := map[string]string{
		"float32": "-",
		"int":     "-",
		"int8":    "-",
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind).WithEmptyValues(emptyValues)
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseSingleLineWithConditionals(t *testing.T) {
	r := regexp.MustCompile(`^(?P<string>[^ ]*) (?P<time>[^ ]*) (?P<int>[-0-9]*) ((?P<target_ip>[^ ]+)[:-](?P<target_port>[0-9]+)|-)( (?P<intopt>[-0-9]*))?`)
	km := map[string]string{
		"time":        "timeISO8601",
		"int":         "int",
		"target_port": "uint16",
		"intopt":      "int",
	}

	logs := `str1 2016-08-10T22:08:42.945958Z 57 1.2.3.4:123 4567
str2 2017-12-18T23:09:45.945958Z 58 - 45678
str3 2019-12-18T23:09:45.945958Z 59 -
`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":      "str1",
				"int":         int(57),
				"target_ip":   "1.2.3.4",
				"target_port": uint16(123),
				"intopt":      int(4567),
			},
			Meta: common.MapStr{
				"_id": "6d7a371b5bd42f8240106580733ec63e397db158",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2017, 12, 18, 23, 9, 45, 945958000, time.UTC),
			Fields: common.MapStr{
				"string": "str2",
				"int":    int(58),
				"intopt": int(45678),
			},
			Meta: common.MapStr{
				"_id": "e6657a85ed4f1e9e7c280815af782166c0bbe4bf",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2019, 12, 18, 23, 9, 45, 945958000, time.UTC),
			Fields: common.MapStr{
				"string": "str3",
				"int":    int(59),
			},
			Meta: common.MapStr{
				"_id": "fdd3dd82bce86471e109171626105a34e0f4f798",
			},
		},
	}

	parser := NewCustomLogParser("time", r).WithKindMap(km)
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLines(t *testing.T) {
	logs := `str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353
strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555
strLine3 2006-08-13T02:08:12.544953Z 12345 05 31123 true str2 0.111 0.123456`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":  "str1",
				"int":     int(35325),
				"int8":    int8(120),
				"int16":   int16(30123),
				"bool":    true,
				"string2": "str2",
				"float32": float32(0.325),
				"float64": 0.0318353,
			},
			Meta: common.MapStr{
				"_id": "4fa1020dbfc28237c745009a6e13c87fd8546e91",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			Fields: common.MapStr{
				"string":  "strLine2",
				"int":     int(321345),
				"int8":    int8(25),
				"int16":   int16(27535),
				"bool":    false,
				"string2": "str2Line2",
				"float32": float32(0.312),
				"float64": 0.323454555,
			},
			Meta: common.MapStr{
				"_id": "bf63b2f1b77a4ef54ed707b9d4eca8006e0049d2",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2006, 8, 13, 2, 8, 12, 544953000, time.UTC),
			Fields: common.MapStr{
				"string":  "strLine3",
				"int":     int(12345),
				"int8":    int8(5),
				"int16":   int16(31123),
				"bool":    true,
				"string2": "str2",
				"float32": float32(0.111),
				"float64": 0.123456,
			},
			Meta: common.MapStr{
				"_id": "69feaceba421db7b52f65815422fb5abc55d6e37",
			},
		},
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind)
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseMultipleLinesWithIgnoredLines(t *testing.T) {
	logs := `# comment, ignored
str1 2016-08-10T22:08:42.945958Z 35325 120 30123 true str2 0.325 0.0318353

    # another ignored line
strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555


strLine3 2006-08-13T02:08:12.544953Z 12345 05 31123 true str2 0.111 0.123456

`

	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			Fields: common.MapStr{
				"string":  "str1",
				"int":     int(35325),
				"int8":    int8(120),
				"int16":   int16(30123),
				"bool":    true,
				"string2": "str2",
				"float32": float32(0.325),
				"float64": 0.0318353,
			},
			Meta: common.MapStr{
				"_id": "4fa1020dbfc28237c745009a6e13c87fd8546e91",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			Fields: common.MapStr{
				"string":  "strLine2",
				"int":     int(321345),
				"int8":    int8(25),
				"int16":   int16(27535),
				"bool":    false,
				"string2": "str2Line2",
				"float32": float32(0.312),
				"float64": 0.323454555,
			},
			Meta: common.MapStr{
				"_id": "bf63b2f1b77a4ef54ed707b9d4eca8006e0049d2",
			},
		},
		&beat.Event{
			Timestamp: time.Date(2006, 8, 13, 2, 8, 12, 544953000, time.UTC),
			Fields: common.MapStr{
				"string":  "strLine3",
				"int":     int(12345),
				"int8":    int8(5),
				"int16":   int16(31123),
				"bool":    true,
				"string2": "str2",
				"float32": float32(0.111),
				"float64": 0.123456,
			},
			Meta: common.MapStr{
				"_id": "a9ffb794f6ead3fca7d4a98c58784ade15608430",
			},
		},
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind).WithReIgnore(regexp.MustCompile(`^\s*#`))
	expectedErrorsPrefix := []string{}
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserParseErrorLines(t *testing.T) {
	logs := `str1 not-a-valid-date 35325 120 30123 true str2 0.325 0.0318353`
	expected := []*beat.Event{}
	expectedErrorsPrefix := []string{
		`Couldn't parse field (time) to type (timeISO8601). Error: parsing time "not-a-valid-date"`,
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind)
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserInvalidFormat(t *testing.T) {
	logs := `Incorrect Line
strLine2 2018-07-15T21:18:47.483845Z 321345 25 27535 false str2Line2 0.312 0.323454555
Incorrect line2
`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2018, 7, 15, 21, 18, 47, 483845000, time.UTC),
			Fields: common.MapStr{
				"string":  "strLine2",
				"int":     int(321345),
				"int8":    int8(25),
				"int16":   int16(27535),
				"bool":    false,
				"string2": "str2Line2",
				"float32": float32(0.312),
				"float64": 0.323454555,
			},
		},
	}
	expectedErrorsPrefix := []string{
		"Line does not match expected format",
		"Line does not match expected format",
	}

	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind)
	assertLogParser(t, parser, &logs, expected, expectedErrorsPrefix)
}

func TestCustomLogParserNothingProcessedOnReaderError(t *testing.T) {
	ok := 0
	ko := 0
	parser := NewCustomLogParser("time", regexTest).WithKindMap(regexKind)
	err := parser.Parse(&testReader{}, func(event *beat.Event) {
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

func assertLogParser(t *testing.T, p LogParser, logs *string, expectedEvents []*beat.Event, expectedErrorsPrefix []string) {
	results := make([]*beat.Event, 0, len(expectedEvents))
	errors := make([]error, 0, len(expectedErrorsPrefix))
	err := p.Parse(strings.NewReader(*logs), func(event *beat.Event) {
		results = append(results, event)
	}, func(errLine string, err error) {
		errors = append(errors, err)
	})
	assert.NoError(t, err)
	assert.Len(t, errors, len(expectedErrorsPrefix))
	assert.Len(t, results, len(expectedEvents))
	for idx, expEvent := range expectedEvents {
		resultEvent := results[idx]
		assertEventFields(t, expEvent.Fields, resultEvent.Fields)
		assertEventFields(t, expEvent.Meta, resultEvent.Meta)
		assert.Equal(t, expEvent.Timestamp, resultEvent.Timestamp)
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

func assertEventFields(t *testing.T, expected, event common.MapStr) {
	for field, exp := range expected {
		val, found := event[field]
		if !found {
			t.Errorf("Missing field: %v", field)
			continue
		}

		if sub, ok := exp.(common.MapStr); ok {
			assertEventFields(t, sub, val.(common.MapStr))
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
