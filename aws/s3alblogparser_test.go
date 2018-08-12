// +build !integration

package aws

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

const (
	logHTTP  = `http 2016-08-10T22:08:42.945958Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.000 0.001 0.000 200 200 34 366 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.46.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" - -`
	logHTTPS = `https 2016-08-10T23:39:43.065466Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.086 0.048 0.037 200 200 0 57 "GET https://www.example.com:443/ HTTP/1.1" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337281-1d84f3d73c47ec4e58577259" www.example.com arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012`
	logHTTP2 = `h2 2016-08-10T00:10:33.145057Z app/my-loadbalancer/50dc6c495c0c9188 10.0.1.252:48160 10.0.0.66:9000 0.000 0.002 0.000 200 200 5 257 "GET https://10.0.2.105:773/ HTTP/2.0" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337327-72bd00b0343d75b906739c42" - -`
	logWS    = `ws 2016-08-10T00:32:08.923954Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:40914 10.0.1.192:8010 0.001 0.003 0.000 101 101 218 587 "GET http://10.0.0.30:80/ HTTP/1.1" "-" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
	logWSS   = `wss 2016-08-10T00:42:46.423695Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:44244 10.0.0.171:8010 0.000 0.001 0.000 101 101 218 786 "GET https://10.0.0.30:443/ HTTP/1.1" "-" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
)

// Examples present here have been obtained from: https://docs.aws.amazon.com/es_es/elasticloadbalancing/latest/application/load-balancer-access-logs.html
func TestS3ALBLogParseHTTP(t *testing.T) {
	logs := logHTTP
	expected := []common.MapStr{
		common.MapStr{
			"type":                     "http",
			"time":                     time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
			"elb":                      "app/my-loadbalancer/50dc6c495c0c9188",
			"client_ip":                "192.168.131.39",
			"client_port":              uint16(2817),
			"target_ip":                "10.0.0.1",
			"target_port":              uint16(80),
			"request_processing_time":  0.000,
			"target_processing_time":   0.001,
			"response_processing_time": 0.000,
			"elb_status_code":          "200",
			"target_status_code":       "200",
			"received_bytes":           int64(34),
			"sent_bytes":               int64(366),
			"request_verb":             "GET",
			"request_url":              "http://www.example.com:80/",
			"request_proto":            "HTTP/1.1",
			"user_agent":               "curl/7.46.0",
			"ssl_cipher":               "-",
			"ssl_protocol":             "-",
			"target_group_arn":         "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
			"trace_id":                 "Root=1-58337262-36d228ad5d99923122bbe354",
		},
	}

	errorLinesExpected := []string{}
	testALBLogParser(t, &logs, expected, errorLinesExpected)
}

// func TestS3ALBLogParserHTTPS(t *testing.T) {
// 	logs := `https 2016-08-10T23:39:43.065466Z app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.086 0.048 0.037 200 200 0 57 "GET https://www.example.com:443/ HTTP/1.1" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337281-1d84f3d73c47ec4e58577259" www.example.com arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012`
// 	expected := []common.MapStr{
// 		common.MapStr{},
// 	}
// 	errorLinesExpected := []string{}
// 	testALBLogParser(t, &logs, expected, errorLinesExpected)
// }
//
// func TestS3ALBLogParserHTTP2(t *testing.T) {
// 	logs := `h2 2016-08-10T00:10:33.145057Z app/my-loadbalancer/50dc6c495c0c9188 10.0.1.252:48160 10.0.0.66:9000 0.000 0.002 0.000 200 200 5 257 "GET https://10.0.2.105:773/ HTTP/2.0" "curl/7.46.0" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337327-72bd00b0343d75b906739c42" - -`
// 	expected := []common.MapStr{
// 		common.MapStr{},
// 	}
// 	errorLinesExpected := []string{}
// 	testALBLogParser(t, &logs, expected, errorLinesExpected)
// }
//
// func TestS3ALBLogParserWS(t *testing.T) {
// 	logs := `ws 2016-08-10T00:32:08.923954Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:40914 10.0.1.192:8010 0.001 0.003 0.000 101 101 218 587 "GET http://10.0.0.30:80/ HTTP/1.1" "-" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
// 	expected := []common.MapStr{
// 		common.MapStr{},
// 	}
// 	errorLinesExpected := []string{}
// 	testALBLogParser(t, &logs, expected, errorLinesExpected)
// }
//
// func TestS3ALBLogParserWSS(t *testing.T) {
// 	logs := `wss 2016-08-10T00:42:46.423695Z app/my-loadbalancer/50dc6c495c0c9188 10.0.0.140:44244 10.0.0.171:8010 0.000 0.001 0.000 101 101 218 786 "GET https://10.0.0.30:443/ HTTP/1.1" "-" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337364-23a8c76965a2ef7629b185e3" - -`
// 	expected := []common.MapStr{
// 		common.MapStr{},
// 	}
// 	errorLinesExpected := []string{}
// 	testALBLogParser(t, &logs, expected, errorLinesExpected)
// }

func TestS3ALBLogParseMultiline(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(logHTTP)
	b.WriteByte('\n')
	b.WriteString(logHTTPS)
	logs := b.String()

	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(logs), func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, ok)
	assert.Equal(t, 0, ko)
}

func TestS3ALBLogParseMultilineAcceptsEmptyLineAtEOF(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(logHTTP)
	b.WriteByte('\n')
	b.WriteString(logHTTPS)
	b.WriteByte('\n')
	logs := b.String()

	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(logs), func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, ok)
	assert.Equal(t, 0, ko)
}

func TestS3ALBLogParseWithParserErrorLines(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(`http not-a-valid-date app/my-loadbalancer/50dc6c495c0c9188 192.168.131.39:2817 10.0.0.1:80 0.000 0.001 0.000 200 200 34 366 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.46.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" - -
`)
	b.WriteString(logHTTP)
	b.WriteByte('\n')
	b.WriteString(logHTTPS)
	b.WriteByte('\n')
	logs := b.String()

	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(logs), func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
		assert.True(t, strings.HasPrefix(err.Error(), `Couldn't parse field (time) to type (time.RFC3339Nano). Error: parsing time "not-a-valid-date"`))
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, ok)
	assert.Equal(t, 1, ko)
}

func TestS3ALBLogParseWithErrorLines(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("Incorrect line\n")
	b.WriteString(logHTTP)
	b.WriteByte('\n')
	b.WriteString("Incorrect line2\n")
	b.WriteString(logHTTPS)
	b.WriteByte('\n')
	logs := b.String()

	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(logs), func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, ok)
	assert.Equal(t, 2, ko)
}

func TestS3ALBLogParseMultilineWithIncorrectKindDoesNotProcess(t *testing.T) {
	var b bytes.Buffer
	b.WriteString(logHTTP)
	b.WriteByte('\n')
	b.WriteString(logHTTPS)
	logs := b.String()

	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(logs), func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, ok)
	assert.Equal(t, 0, ko)
}

func TestS3ALBLogParseReaderErrorProcessNothing(t *testing.T) {
	ok := 0
	ko := 0
	parser := NewS3ALBLogParser()
	err := parser.Parse(&testReader{}, func(s common.MapStr) {
		ok++
	}, func(errLine string, err error) {
		ko++
	})
	assert.Error(t, err)
	assert.Equal(t, 0, ok)
	assert.Equal(t, 0, ko)
}

func TestS3ALBLogParseStringToKindsWithNoErrors(t *testing.T) {
	type elem struct {
		kind     string
		strValue string
		value    interface{}
	}
	elems := []elem{
		elem{
			kind:     "time.RFC3339Nano",
			strValue: "2016-08-10T22:08:42.945958Z",
			value:    time.Date(2016, 8, 10, 22, 8, 42, 945958000, time.UTC),
		},
		elem{
			kind:     "bool",
			strValue: "true",
			value:    true,
		},
		elem{
			kind:     "int8",
			strValue: "5",
			value:    int8(5),
		},
		elem{
			kind:     "int16",
			strValue: "32000",
			value:    int16(32000),
		},
		elem{
			kind:     "int",
			strValue: "67353",
			value:    int(67353),
		},
		elem{
			kind:     "int32",
			strValue: "67353",
			value:    int32(67353),
		},
		elem{
			kind:     "int64",
			strValue: "-35868395685",
			value:    int64(-35868395685),
		},
		elem{
			kind:     "uint8",
			strValue: "250",
			value:    uint8(250),
		},
		elem{
			kind:     "uint16",
			strValue: "32000",
			value:    uint16(32000),
		},
		elem{
			kind:     "uint",
			strValue: "835000",
			value:    uint(835000),
		},
		elem{
			kind:     "uint32",
			strValue: "835000",
			value:    uint32(835000),
		},
		elem{
			kind:     "uint64",
			strValue: "35868395685",
			value:    uint64(35868395685),
		},
		elem{
			kind:     "float32",
			strValue: "0.385694",
			value:    float32(0.385694),
		},
		elem{
			kind:     "float64",
			strValue: "0.38569355355334",
			value:    0.38569355355334,
		},
		elem{
			kind:     "string",
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

func TestS3ALBLogParseStringToKindsWithParseErrors(t *testing.T) {
	type elem struct {
		kind     string
		strValue string
	}
	elems := []elem{
		elem{
			kind:     "time.RFC3339Nano",
			strValue: "true",
		},
		elem{
			kind:     "bool",
			strValue: "3",
		},
		elem{
			kind:     "int8",
			strValue: "53535",
		},
		elem{
			kind:     "int16",
			strValue: "jo",
		},
		elem{
			kind:     "int",
			strValue: "true",
		},
		elem{
			kind:     "int32",
			strValue: "false",
		},
		elem{
			kind:     "int64",
			strValue: "none",
		},
		elem{
			kind:     "uint8",
			strValue: "-35",
		},
		elem{
			kind:     "uint16",
			strValue: "-5",
		},
		elem{
			kind:     "uint",
			strValue: "false",
		},
		elem{
			kind:     "uint32",
			strValue: "true",
		},
		elem{
			kind:     "uint64",
			strValue: "-3235",
		},
		elem{
			kind:     "float32",
			strValue: "false",
		},
		elem{
			kind:     "float64",
			strValue: "true",
		},
	}

	for _, e := range elems {
		_, err := parseStringToKind(e.kind, e.strValue)
		assert.Error(t, err)
	}
}

func TestS3ALBLogParseStringToKindsWithError(t *testing.T) {
	v, err := parseStringToKind("NotExisting", "353")
	assert.Nil(t, v)
	assert.EqualError(t, err, "Can not convert to unsupported kind (NotExisting)")
}

func testALBLogParser(t *testing.T, logs *string, expected []common.MapStr, expectedErrorLines []string) {
	results := make([]common.MapStr, 0, len(expected))
	errors := make([]string, 0, len(expectedErrorLines))
	parser := NewS3ALBLogParser()
	err := parser.Parse(strings.NewReader(*logs), func(s common.MapStr) {
		results = append(results, s)
	}, func(errLine string, err error) {
		errors = append(errors, errLine)
	})
	assert.NoError(t, err)
	assert.Len(t, errors, 0)
	assert.Len(t, results, 1)
	for idx, expEvent := range expected {
		resultEvent := results[idx]
		assertEvent(t, expEvent, resultEvent)
	}
	assert.Equal(t, expectedErrorLines, errors)
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
