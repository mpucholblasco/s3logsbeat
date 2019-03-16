// +build !integration

package pipeline

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
)

func TestGetKeyFields(t *testing.T) {
	key := "myenvironment-myapp/myawsregion/myfile.gz"

	ri := &S3ReaderInformation{
		keyRegexFields: regexp.MustCompile(`^(?P<environment>[^\-]+)-(?P<application>[^/]+)/`),
	}
	keyFields, err := ri.GetKeyFields(key)
	expectedKeyFields := common.MapStr{
		"environment": "myenvironment",
		"application": "myapp",
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedKeyFields, *keyFields)
}
