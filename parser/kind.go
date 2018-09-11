package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
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

func mustKindMapStringToType(o map[string]string) map[string]kindElement {
	r, err := kindMapStringToType(o)
	if err != nil {
		panic(`parser: KindMapStringToType error: ` + err.Error())
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
