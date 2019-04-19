// +build !integration

package logparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// KindMapStringToType & MustKindMapStringToType tests
func TestCustomLogParserKindMapStringToTypeCorrect(t *testing.T) {
	m := map[string]string{
		"time":           "timeISO8601",
		"int":            "int",
		"int8":           "int8",
		"int16":          "int16",
		"bool":           "bool",
		"string2":        "string",
		"float32":        "float32",
		"float64":        "float64",
		"urlencoded":     "urlencoded",
		"deepurlencoded": "deepurlencoded",
	}
	expected := map[string]kindElement{
		"time":           kindMap[kindTimeISO8601],
		"int":            kindMap[kindInt],
		"int8":           kindMap[kindInt8],
		"int16":          kindMap[kindInt16],
		"bool":           kindMap[kindBool],
		"string2":        kindMap[kindString],
		"float32":        kindMap[kindFloat32],
		"float64":        kindMap[kindFloat64],
		"urlencoded":     kindMap[kindURLEncoded],
		"deepurlencoded": kindMap[kindDeepURLEncoded],
	}

	value, err := kindMapStringToType(m)
	assert.NoError(t, err)
	for k, v := range value {
		assert.Equal(t, v.kind, expected[k].kind)
	}
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
func TestCustomLogParserParseToKindsWithNoErrors(t *testing.T) {
	type elem struct {
		kind    kindElement
		inValue interface{}
		value   interface{}
	}
	elems := []elem{
		elem{
			kind:    kindMap[kindTimeISO8601],
			inValue: "2016-08-10T22:08:42.945958Z",
			value:   time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
		},
		elem{
			kind:    kindMap[kindTimeUnixMilliseconds],
			inValue: "1553360693208",
			value:   time.Date(2019, 3, 23, 17, 4, 53, 208000000, time.UTC),
		},
		elem{
			kind:    kindMap[kindTimeUnixMilliseconds],
			inValue: 1553360693208,
			value:   time.Date(2019, 3, 23, 17, 4, 53, 208000000, time.UTC),
		},
		elem{
			kind:    kindMap[kindTimeUnixMilliseconds],
			inValue: int64(1553360693208),
			value:   time.Date(2019, 3, 23, 17, 4, 53, 208000000, time.UTC),
		},
		elem{
			kind:    kindMap[kindBool],
			inValue: "true",
			value:   true,
		},
		elem{
			kind:    kindMap[kindInt8],
			inValue: "5",
			value:   int8(5),
		},
		elem{
			kind:    kindMap[kindInt16],
			inValue: "32000",
			value:   int16(32000),
		},
		elem{
			kind:    kindMap[kindInt],
			inValue: "67353",
			value:   int(67353),
		},
		elem{
			kind:    kindMap[kindInt32],
			inValue: "67353",
			value:   int32(67353),
		},
		elem{
			kind:    kindMap[kindInt64],
			inValue: "-35868395685",
			value:   int64(-35868395685),
		},
		elem{
			kind:    kindMap[kindUint8],
			inValue: "250",
			value:   uint8(250),
		},
		elem{
			kind:    kindMap[kindUint16],
			inValue: "32000",
			value:   uint16(32000),
		},
		elem{
			kind:    kindMap[kindUint],
			inValue: "835000",
			value:   uint(835000),
		},
		elem{
			kind:    kindMap[kindUint32],
			inValue: "835000",
			value:   uint32(835000),
		},
		elem{
			kind:    kindMap[kindUint64],
			inValue: "35868395685",
			value:   uint64(35868395685),
		},
		elem{
			kind:    kindMap[kindFloat32],
			inValue: "0.385694",
			value:   float32(0.385694),
		},
		elem{
			kind:    kindMap[kindFloat64],
			inValue: "0.38569355355334",
			value:   0.38569355355334,
		},
		elem{
			kind:    kindMap[kindString],
			inValue: "This is a string",
			value:   "This is a string",
		},
		elem{
			kind:    kindMap[kindURLEncoded],
			inValue: "Mozilla/4.0%20(compatible;%20MSIE%207.0;%20Windows%20NT%205.1)",
			value:   "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)",
		},
		elem{
			kind:    kindMap[kindDeepURLEncoded],
			inValue: `8%2525%2520fresa & maracuyá`,
			value:   "8% fresa & maracuyá",
		},

		// To string
		elem{
			kind:    kindMap[kindString],
			inValue: 3,
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: int8(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: int16(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: int32(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: int64(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: 3,
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: uint8(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: uint16(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: uint32(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: uint64(3),
			value:   "3",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: 3.56463655,
			value:   "3.56463655",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: float32(3.6436653),
			value:   "3.6436653",
		},
		elem{
			kind:    kindMap[kindString],
			inValue: false,
			value:   "false",
		},
	}

	for _, e := range elems {
		v, err := parseToKind(e.kind, e.inValue)
		assert.NoError(t, err)
		assert.Equal(t, e.value, v)
	}
}

func TestTimePattern(t *testing.T) {
	k, e := kindFromString("time:2006-01-02\t15:04:05")
	assert.NoError(t, e)

	in := "2014-05-23\t01:15:18"
	expected := time.Date(2014, 5, 23, 1, 15, 18, 0, time.UTC)
	result, e := parseToKind(k, in)
	assert.NoError(t, e)
	assert.Equal(t, expected, result)
}

func TestTimePatternError(t *testing.T) {
	k, e := kindFromString("time:2006-01-02 15:04:05")
	assert.NoError(t, e)

	in := "3 Feb 2014 01:14:18"
	_, e = parseToKind(k, in)
	assert.Error(t, e)
}

func TestTimePatternInvalidType(t *testing.T) {
	k, e := kindFromString("time:2006-01-02 15:04:05")
	assert.NoError(t, e)

	in := 123456
	_, e = parseToKind(k, in)
	assert.Error(t, e)
}

func TestCustomLogParserParseToKindsWithParseErrors(t *testing.T) {
	type elem struct {
		kind    kindElement
		inValue interface{}
	}
	elems := []elem{
		elem{
			kind:    kindMap[kindTimeISO8601],
			inValue: "true",
		},
		elem{
			kind:    kindMap[kindBool],
			inValue: "3",
		},
		elem{
			kind:    kindMap[kindInt8],
			inValue: "53535",
		},
		elem{
			kind:    kindMap[kindInt16],
			inValue: "jo",
		},
		elem{
			kind:    kindMap[kindInt],
			inValue: "true",
		},
		elem{
			kind:    kindMap[kindInt32],
			inValue: "false",
		},
		elem{
			kind:    kindMap[kindInt64],
			inValue: "none",
		},
		elem{
			kind:    kindMap[kindUint8],
			inValue: "-35",
		},
		elem{
			kind:    kindMap[kindUint16],
			inValue: "-5",
		},
		elem{
			kind:    kindMap[kindUint],
			inValue: "false",
		},
		elem{
			kind:    kindMap[kindUint32],
			inValue: "true",
		},
		elem{
			kind:    kindMap[kindUint64],
			inValue: "-3235",
		},
		elem{
			kind:    kindMap[kindFloat32],
			inValue: "false",
		},
		elem{
			kind:    kindMap[kindFloat64],
			inValue: "true",
		},
		elem{
			kind:    kindMap[kindURLEncoded],
			inValue: "a%5Z",
		},
	}

	for _, e := range elems {
		_, err := parseToKind(e.kind, e.inValue)
		assert.Error(t, err)
	}
}

func TestCustomLogParserParseToKindsWithInvalidType(t *testing.T) {
	type elem struct {
		kind    kindElement
		inValue interface{}
	}
	elems := []elem{
		elem{
			kind:    kindMap[kindTimeISO8601],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindBool],
			inValue: 3.6536,
		},
		elem{
			kind:    kindMap[kindInt8],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindInt16],
			inValue: false,
		},
		elem{
			kind:    kindMap[kindInt],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindInt32],
			inValue: false,
		},
		elem{
			kind:    kindMap[kindInt64],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindUint8],
			inValue: -563,
		},
		elem{
			kind:    kindMap[kindUint16],
			inValue: -5,
		},
		elem{
			kind:    kindMap[kindUint],
			inValue: false,
		},
		elem{
			kind:    kindMap[kindUint32],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindUint64],
			inValue: -3235,
		},
		elem{
			kind:    kindMap[kindFloat32],
			inValue: false,
		},
		elem{
			kind:    kindMap[kindFloat64],
			inValue: true,
		},
		elem{
			kind:    kindMap[kindURLEncoded],
			inValue: 56465,
		},
	}

	for _, e := range elems {
		_, err := parseToKind(e.kind, e.inValue)
		assert.Error(t, err)
	}
}

func BenchmarkParseToKind(b *testing.B) {
	// old : 45.9 ns/op, 45.7...
	for n := 0; n < b.N; n++ {
		parseToKind(kindElements[b.N%len(kindElements)], "in")
	}
}
