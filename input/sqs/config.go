package sqs

import (
	"fmt"

	"github.com/mpucholblasco/s3logsbeat/input"
)

var (
	defaultConfig = config{}
)

type config struct {
	input.GlobalConfig `config:",inline"`
	QueuesURL          []string `config:"queues_url"`
}

func (c *config) Validate() error {
	if err := c.GlobalConfig.Validate(); err != nil {
		return err
	}

	if len(c.QueuesURL) == 0 {
		return fmt.Errorf("No queues_url defined for sqs input")
	}
	return nil
}
