package sqs

import (
	"fmt"
	"time"

	"github.com/sequra/s3logsbeat/input"
)

var (
	defaultConfig = config{}
)

type config struct {
	input.GlobalConfig `config:",inline"`
	Buckets            []string  `config:"buckets"`
	SinceStr           string    `config:"since"`
	ToStr              string    `config:"to"`
	Since              time.Time `config:",ignore"`
	To                 time.Time `config:",ignore"`
}

func (c *config) Validate() error {
	var err error
	if err = c.GlobalConfig.Validate(); err != nil {
		return err
	}

	if len(c.Buckets) == 0 {
		return fmt.Errorf("No bucket defined for s3 input")
	}

	if c.SinceStr != "" {
		c.Since, err = time.Parse(time.RFC3339Nano, c.SinceStr)
		if err != nil {
			return err
		}
	}

	if c.ToStr == "" {
		c.To = time.Unix(1<<63-62135596801, 999999999)
	} else {
		c.To, err = time.Parse(time.RFC3339Nano, c.ToStr)
		if err != nil {
			return err
		}
	}
	return nil
}
