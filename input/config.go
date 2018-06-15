// Input configuration
package input

import (
	"regexp"
	"time"

	cfg "github.com/mpucholblasco/s3logsbeat/config"
)

type InputConfig struct {
	QueuesURL      []string          `config:"queues_url"`
	LogFormat      string            `config:"log_format" validate:"required"`
	Type           string            `config:"type" validate:"required"`
	KeyRegexFields *regexp.Regexp    `config:"key_regex_fields"`
	PollFrequency  time.Duration     `config:"poll_frequency" validate:"required,min=0,nonzero"`
	Fields         map[string]string `config:"fields"`
}

var (
	defaultConfig = InputConfig{
		Type: cfg.DefaultType,
	}
)
