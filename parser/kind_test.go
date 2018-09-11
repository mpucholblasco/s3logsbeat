// +build !integration

package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
