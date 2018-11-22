package sqs

import (
	"fmt"

	"github.com/mpucholblasco/s3logsbeat/input"
)

var (
	defaultConfig = config{}
)

type config struct {
	input.GlobalConfig
	Buckets []string `config:"buckets"`
}

func (c *config) Validate() error {
	if err := c.GlobalConfig.Validate(); err != nil {
		return err
	}

	if len(c.Buckets) == 0 {
		return fmt.Errorf("No bucket defined for s3 input")
	}
	return nil
}
