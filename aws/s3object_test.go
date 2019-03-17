// +build !integration

package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3ObjectInvalidURI(t *testing.T) {
	_, err := NewS3ObjectFromURI(`s3://no`)
	assert.Error(t, err)
}

func TestS3ObjectURIWithoutPath(t *testing.T) {
	s3object, err := NewS3ObjectFromURI(`s3://valid.com/`)
	assert.NoError(t, err)
	assert.Equal(t, "valid.com", s3object.Bucket)
	assert.Equal(t, "", s3object.Key)
}

func TestS3ObjectURIComplete(t *testing.T) {
	s3object, err := NewS3ObjectFromURI(`s3://valid.com/a/b/c`)
	assert.NoError(t, err)
	assert.Equal(t, "valid.com", s3object.Bucket)
	assert.Equal(t, "a/b/c", s3object.Key)
}
