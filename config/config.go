// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations
package config

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
)

// Defaults for config variables which are not set
const (
	DefaultType = "sqs"
)

type Config struct {
	Inputs          []*common.Config `config:"inputs" validate:"required"`
	ShutdownTimeout time.Duration    `config:"shutdown_timeout"`
}

var DefaultConfig = Config{
	ShutdownTimeout: 0,
}
