package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

// AssertLogParser asserts that expectedEvents and expectedErrorsPrefix are equal to the obtained one when
// parsing logs
func AssertLogParser(t *testing.T, p LogParser, logs *string, expectedEvents []*beat.Event, expectedErrorsPrefix []string) {
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
		AssertEventFields(t, expEvent.Fields, resultEvent.Fields)
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

// AssertEventFields asserts that expected and event maps are equal
func AssertEventFields(t *testing.T, expected, event common.MapStr) {
	for field, exp := range expected {
		val, found := event[field]
		if !found {
			t.Errorf("Missing field: %v", field)
			continue
		}

		if sub, ok := exp.(common.MapStr); ok {
			AssertEventFields(t, sub, val.(common.MapStr))
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
