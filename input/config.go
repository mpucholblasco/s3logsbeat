// Package input configuration
package input

import (
	"regexp"
	"time"

	"github.com/elastic/beats/libbeat/common"
	cfg "github.com/mpucholblasco/s3logsbeat/config"
)

// GlobalConfig global config for all kind of inputs
type GlobalConfig struct {
	Type             string            `config:"type" validate:"required"`
	PollFrequency    time.Duration     `config:"poll_frequency" validate:"min=0,nonzero"`
	LogFormat        string            `config:"log_format" validate:"required"`
	LogFormatOptions *common.Config    `config:"log_format_options"`
	KeyRegexFields   *regexp.Regexp    `config:"key_regex_fields"`
	Fields           map[string]string `config:"fields"`
}

var (
	defaultConfig = GlobalConfig{
		Type: cfg.DefaultType,
	}
)

// Validate validates global config logic
func (c *GlobalConfig) Validate() error {
	return nil
}
