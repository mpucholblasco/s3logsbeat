// +build !integration

package testutil

import (
	"reflect"
	"testing"

	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func AssertEvent(t *testing.T, expected, event common.MapStr) {
	for field, exp := range expected {
		val, found := event[field]
		if !found {
			t.Errorf("Missing field: %v", field)
			continue
		}

		if sub, ok := exp.(common.MapStr); ok {
			AssertEvent(t, sub, val.(common.MapStr))
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
