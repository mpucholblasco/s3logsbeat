// +build !integration

package logparser

import (
	"strings"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"

	"github.com/stretchr/testify/assert"
)

func TestJSONLogParser(t *testing.T) {
	logs := `{"timestamp":1553360693208,"formatVersion":1,"webaclId":"2668f4a5-da32-4d63-bea8-1ed02607da4f","terminatingRuleId":"Default_Action","terminatingRuleType":"REGULAR","action":"BLOCK","httpSourceName":"ALB","httpSourceId":"12345678901-app/myalb/70e9bf09a00ca695","ruleGroupList":[],"rateBasedRuleList":[{"rateBasedRuleId":"56d04362-6fab-4cb7-a314-819a1bd40cc6","limitKey":"IP","maxRateAllowed":2000}],"nonTerminatingMatchingRules":[],"httpRequest":{"clientIp":"37.133.193.245","country":"ES","headers":[{"name":"Host","value":"www.example.com"},{"name":"Content-Length","value":"0"},{"name":"user-agent","value":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"},{"name":"accept","value":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},{"name":"accept-language","value":"en-GB,en;q=0.5"},{"name":"accept-encoding","value":"REDACTED"},{"name":"dnt","value":"1"},{"name":"upgrade-insecure-requests","value":"1"},{"name":"pragma","value":"no-cache"},{"name":"cache-control","value":"REDACTED"}],"uri":"REDACTED","args":"REDACTED","httpVersion":"HTTP/2.0","httpMethod":"REDACTED","requestId":null}}
					 {"timestamp":1553360693035,"formatVersion":1,"webaclId":"2668f4a5-da32-4d63-bea8-1ed02607da4f","terminatingRuleId":"Default_Action","terminatingRuleType":"REGULAR","action":"BLOCK","httpSourceName":"ALB","httpSourceId":"12345678901-app/myalb/70e9bf09a00ca695","ruleGroupList":[],"rateBasedRuleList":[{"rateBasedRuleId":"56d04362-6fab-4cb7-a314-819a1bd40cc6","limitKey":"IP","maxRateAllowed":2000}],"nonTerminatingMatchingRules":[],"httpRequest":{"clientIp":"37.133.193.245","country":"ES","headers":[{"name":"Host","value":"www.example.com"},{"name":"Content-Length","value":"0"},{"name":"user-agent","value":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"},{"name":"accept","value":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},{"name":"accept-language","value":"en-GB,en;q=0.5"},{"name":"accept-encoding","value":"REDACTED"},{"name":"dnt","value":"1"},{"name":"pragma","value":"no-cache"},{"name":"cache-control","value":"REDACTED"}],"uri":"REDACTED","args":"REDACTED","httpVersion":"HTTP/2.0","httpMethod":"REDACTED","requestId":null}}`
	expected := []*beat.Event{
		&beat.Event{
			Timestamp: time.Date(2019, 3, 23, 17, 4, 53, 208000000, time.UTC),
			Fields: common.MapStr{
				"formatVersion":       int64(1),
				"webaclId":            "2668f4a5-da32-4d63-bea8-1ed02607da4f",
				"terminatingRuleId":   "Default_Action",
				"terminatingRuleType": "REGULAR",
				"action":              "BLOCK",
				"httpSourceName":      "ALB",
				"httpSourceId":        "12345678901-app/myalb/70e9bf09a00ca695",
				"rateBasedRuleList": []interface{}{
					map[string]interface{}{
						"rateBasedRuleId": "56d04362-6fab-4cb7-a314-819a1bd40cc6",
						"limitKey":        "IP",
						"maxRateAllowed":  int64(2000),
					},
				},
				"httpRequest": map[string]interface{}{
					"clientIp": "37.133.193.245",
					"country":  "ES",
					"headers": []interface{}{
						map[string]interface{}{
							"name":  "Host",
							"value": "www.example.com",
						},
						map[string]interface{}{
							"name":  "Content-Length",
							"value": "0",
						},
						map[string]interface{}{
							"name":  "user-agent",
							"value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0",
						},
						map[string]interface{}{
							"name":  "accept",
							"value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
						},
						map[string]interface{}{
							"name":  "accept-language",
							"value": "en-GB,en;q=0.5",
						},
						map[string]interface{}{
							"name":  "accept-encoding",
							"value": "REDACTED",
						},
						map[string]interface{}{
							"name":  "dnt",
							"value": "1",
						},
						map[string]interface{}{
							"name":  "upgrade-insecure-requests",
							"value": "1",
						},
						map[string]interface{}{
							"name":  "pragma",
							"value": "no-cache",
						},
						map[string]interface{}{
							"name":  "cache-control",
							"value": "REDACTED",
						},
					},
					"uri":         "REDACTED",
					"args":        "REDACTED",
					"httpVersion": "HTTP/2.0",
					"httpMethod":  "REDACTED",
					"requestId":   nil,
				},
			},
		},
		&beat.Event{
			Timestamp: time.Date(2019, 3, 23, 17, 4, 53, 35000000, time.UTC),
			Fields: common.MapStr{
				"formatVersion":       int64(1),
				"webaclId":            "2668f4a5-da32-4d63-bea8-1ed02607da4f",
				"terminatingRuleId":   "Default_Action",
				"terminatingRuleType": "REGULAR",
				"action":              "BLOCK",
				"httpSourceName":      "ALB",
				"httpSourceId":        "12345678901-app/myalb/70e9bf09a00ca695",
				"rateBasedRuleList": []interface{}{
					map[string]interface{}{
						"rateBasedRuleId": "56d04362-6fab-4cb7-a314-819a1bd40cc6",
						"limitKey":        "IP",
						"maxRateAllowed":  int64(2000),
					},
				},
				"httpRequest": map[string]interface{}{
					"clientIp": "37.133.193.245",
					"country":  "ES",
					"headers": []interface{}{
						map[string]interface{}{
							"name":  "Host",
							"value": "www.example.com",
						},
						map[string]interface{}{
							"name":  "Content-Length",
							"value": "0",
						},
						map[string]interface{}{
							"name":  "user-agent",
							"value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0",
						},
						map[string]interface{}{
							"name":  "accept",
							"value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
						},
						map[string]interface{}{
							"name":  "accept-language",
							"value": "en-GB,en;q=0.5",
						},
						map[string]interface{}{
							"name":  "accept-encoding",
							"value": "REDACTED",
						},
						map[string]interface{}{
							"name":  "dnt",
							"value": "1",
						},
						map[string]interface{}{
							"name":  "pragma",
							"value": "no-cache",
						},
						map[string]interface{}{
							"name":  "cache-control",
							"value": "REDACTED",
						},
					},
					"uri":         "REDACTED",
					"args":        "REDACTED",
					"httpVersion": "HTTP/2.0",
					"httpMethod":  "REDACTED",
					"requestId":   nil,
				},
			},
		},
	}
	errorLinesExpected := []string{}
	logParser := NewJSONLogParser("timestamp", mustKindFromString("timeUnixMilliseconds"))
	assertLogParser(t, logParser, &logs, expected, errorLinesExpected)
}

func TestJSONLogParserEmptyLineEOF(t *testing.T) {
	logs := `{"timestamp":1553360693208,"formatVersion":1,"webaclId":"2668f4a5-da32-4d63-bea8-1ed02607da4f","terminatingRuleId":"Default_Action","terminatingRuleType":"REGULAR","action":"BLOCK","httpSourceName":"ALB","httpSourceId":"12345678901-app/myalb/70e9bf09a00ca695","ruleGroupList":[],"rateBasedRuleList":[{"rateBasedRuleId":"56d04362-6fab-4cb7-a314-819a1bd40cc6","limitKey":"IP","maxRateAllowed":2000}],"nonTerminatingMatchingRules":[],"httpRequest":{"clientIp":"37.133.193.245","country":"ES","headers":[{"name":"Host","value":"www.example.com"},{"name":"Content-Length","value":"0"},{"name":"user-agent","value":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"},{"name":"accept","value":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},{"name":"accept-language","value":"en-GB,en;q=0.5"},{"name":"accept-encoding","value":"REDACTED"},{"name":"dnt","value":"1"},{"name":"upgrade-insecure-requests","value":"1"},{"name":"pragma","value":"no-cache"},{"name":"cache-control","value":"REDACTED"}],"uri":"REDACTED","args":"REDACTED","httpVersion":"HTTP/2.0","httpMethod":"REDACTED","requestId":null}}
					 {"timestamp":1553360693035,"formatVersion":1,"webaclId":"2668f4a5-da32-4d63-bea8-1ed02607da4f","terminatingRuleId":"Default_Action","terminatingRuleType":"REGULAR","action":"BLOCK","httpSourceName":"ALB","httpSourceId":"12345678901-app/myalb/70e9bf09a00ca695","ruleGroupList":[],"rateBasedRuleList":[{"rateBasedRuleId":"56d04362-6fab-4cb7-a314-819a1bd40cc6","limitKey":"IP","maxRateAllowed":2000}],"nonTerminatingMatchingRules":[],"httpRequest":{"clientIp":"37.133.193.245","country":"ES","headers":[{"name":"Host","value":"www.example.com"},{"name":"Content-Length","value":"0"},{"name":"user-agent","value":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:65.0) Gecko/20100101 Firefox/65.0"},{"name":"accept","value":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},{"name":"accept-language","value":"en-GB,en;q=0.5"},{"name":"accept-encoding","value":"REDACTED"},{"name":"dnt","value":"1"},{"name":"pragma","value":"no-cache"},{"name":"cache-control","value":"REDACTED"}],"uri":"REDACTED","args":"REDACTED","httpVersion":"HTTP/2.0","httpMethod":"REDACTED","requestId":null}}
`
	logParser := NewJSONLogParser("timestamp", mustKindFromString("timeUnixMilliseconds"))
	ack := make(chan interface{}, 1)
	go func() {
		err := logParser.Parse(strings.NewReader(logs), func(event *beat.Event) {
		}, func(errLine string, err error) {
		})
		assert.NoError(t, err)
		close(ack)
	}()
	select {
	case <-ack:
	case <-time.After(100 * time.Millisecond):
		t.Error("Parser stuck processing input")
	}
}
