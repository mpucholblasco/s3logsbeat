package log

import (
	"fmt"
	"regexp"
	"time"
)

var (
	defaultConfig = config{}
)

type config struct {
	QueuesURL      []string          `config:"queues_url"`
	LogFormat      string            `config:"log_format" validate:"required"`
	Type           string            `config:"type" validate:"required"`
	KeyRegexFields *regexp.Regexp    `config:"key_regex_fields"`
	PollFrequency  time.Duration     `config:"poll_frequency" validate:"required,min=0,nonzero"`
	Fields         map[string]string `config:"fields"`
}

func (c *config) Validate() error {
	if len(c.QueuesURL) == 0 {
		return fmt.Errorf("No QueuesURL were defined for input")
	}
	return nil
}
