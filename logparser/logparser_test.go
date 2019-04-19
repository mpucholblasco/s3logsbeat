// +build !integration

package logparser

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

func TestCreateEvent(t *testing.T) {
	line := `2019-03-23T17:04:53 mytest`
	timestamp := time.Date(2019, 3, 23, 17, 4, 53, 208000000, time.UTC)
	fields := common.MapStr{
		"field": "mytest",
	}
	event := CreateEvent(&line, timestamp, fields)

	expectedFields := common.MapStr{
		"field": "mytest",
	}
	expectedMeta := common.MapStr{
		"_id": "5e5aaa8b1837066efeb3b048f9c7048e2b8261ec",
	}
	assertEventFields(t, expectedFields, event.Fields)
	assertEventFields(t, expectedMeta, event.Meta)
}
