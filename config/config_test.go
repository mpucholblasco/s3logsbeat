// +build !integration

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/cfgfile"
)

func TestReadConfig(t *testing.T) {
	absPath, err := filepath.Abs("../tests/files/")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	tmpConfig := struct {
		S3logsbeat Config
	}{}

	// Reads config file
	err = cfgfile.Read(&tmpConfig, absPath+"/config.yml")
	assert.Nil(t, err)
}

func TestReadNoinputsConfig(t *testing.T) {
	absPath, err := filepath.Abs("../tests/files/")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	tmpConfig := struct {
		S3logsbeat Config
	}{}

	// Reads config file
	err = cfgfile.Read(&tmpConfig, absPath+"/config_no_inputs.yml")
	assert.Contains(t, err.Error(), "empty field accessing 's3logsbeat.inputs'")
}

func TestContentConfig(t *testing.T) {
	absPath, err := filepath.Abs("../tests/files/")
	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	tmpConfig := struct {
		S3logsbeat Config
	}{}

	// Reads config file
	err = cfgfile.Read(&tmpConfig, absPath+"/config.yml")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(tmpConfig.S3logsbeat.Inputs))
}
