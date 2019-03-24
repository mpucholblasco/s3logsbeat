package logparser

import (
	"fmt"
	"net/url"
	"reflect"
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
	kindDeepURLEncoded

	kindTimeISO8601
	kindTimeUnixMilliseconds
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
			kind: kindDeepURLEncoded,
			name: "deepurlencoded",
		},
		kindElement{
			kind: kindTimeISO8601,
			name: "timeISO8601",
		},
		kindElement{
			kind: kindTimeUnixMilliseconds,
			name: "timeUnixMilliseconds",
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
	var err error
	for k, v := range o {
		r[k], err = kindFromString(v)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func mustKindFromString(v string) kindElement {
	r, err := kindFromString(v)
	if err != nil {
		panic(`parser: mustKindFromString error: ` + err.Error())
	}
	return r
}

func kindFromString(v string) (kindElement, error) {
	if kind, ok := kindStringMap[v]; ok {
		return kind, nil
	} else if strings.HasPrefix(v, "time:") {
		timeLayout := strings.TrimPrefix(v, "time:")
		return kindElement{
			kind:      kindTimeLayout,
			kindExtra: timeLayout,
			name:      fmt.Sprintf("time layout (%s)", timeLayout),
		}, nil
	} else {
		return kindElement{}, fmt.Errorf("Unsupported kind (%s)", v)
	}
}

// parseToKind parses a value to convert it into the kind passed as argument
// NOTE: tried to impre performance (obtained ~46.5ns/op) by using functions inside kindElement
// but it did it slower (~90ns/op)
func parseToKind(e kindElement, value interface{}) (interface{}, error) {
	switch e.kind {
	case kindTimeLayout:
		switch s := value.(type) {
		case string:
			return time.Parse(e.kindExtra.(string), s)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindTimeISO8601:
		switch s := value.(type) {
		case string:
			return time.Parse(time.RFC3339Nano, s)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindTimeUnixMilliseconds:
		var milliseconds int64
		var err error
		switch s := value.(type) {
		case string:
			milliseconds, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
		case int:
			milliseconds = int64(s)
		case int32:
			milliseconds = int64(s)
		case int64:
			milliseconds = s
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
		return time.Unix(milliseconds/1000, milliseconds%1000*1000000).UTC(), nil
	case kindBool:
		switch s := value.(type) {
		case string:
			return strconv.ParseBool(s)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindInt8:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return nil, err
			}
			return int8(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindInt16:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return nil, err
			}
			return int16(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindInt:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return nil, err
			}
			return int(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindInt32:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return nil, err
			}
			return int32(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindInt64:
		switch s := value.(type) {
		case string:
			return strconv.ParseInt(s, 10, 64)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindUint8:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return nil, err
			}
			return uint8(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindUint16:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return nil, err
			}
			return uint16(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindUint:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil, err
			}
			return uint(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindUint32:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil, err
			}
			return uint32(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindUint64:
		switch s := value.(type) {
		case string:
			return strconv.ParseUint(s, 10, 64)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindFloat32:
		switch s := value.(type) {
		case string:
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return nil, err
			}
			return float32(v), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindFloat64:
		switch s := value.(type) {
		case string:
			return strconv.ParseFloat(s, 64)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindURLEncoded:
		switch s := value.(type) {
		case string:
			return url.QueryUnescape(s)
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindDeepURLEncoded:
		switch s := value.(type) {
		case string:
			return deepURLDecode(s), nil
		default:
			return nil, fmt.Errorf("Couldn't convert %s to %s", reflect.TypeOf(value), e.name)
		}
	case kindString:
		switch s := value.(type) {
		case int8:
			return strconv.FormatInt(int64(s), 10), nil
		case int16:
			return strconv.FormatInt(int64(s), 10), nil
		case int32:
			return strconv.FormatInt(int64(s), 10), nil
		case int:
			return strconv.FormatInt(int64(s), 10), nil
		case int64:
			return strconv.FormatInt(s, 10), nil
		case uint8:
			return strconv.FormatUint(uint64(s), 10), nil
		case uint16:
			return strconv.FormatUint(uint64(s), 10), nil
		case uint32:
			return strconv.FormatUint(uint64(s), 10), nil
		case uint:
			return strconv.FormatUint(uint64(s), 10), nil
		case uint64:
			return strconv.FormatUint(s, 10), nil
		case float32:
			return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
		case float64:
			return strconv.FormatFloat(s, 'f', -1, 64), nil
		case bool:
			return strconv.FormatBool(s), nil
		}
	}

	return value, nil
}

func deepURLDecode(u string) string {
	p := u
	for {
		n, err := url.QueryUnescape(p)
		if err != nil || n == p {
			return p
		}
		p = n
	}
}
