package input

import (
	"fmt"

	"github.com/sequra/s3logsbeat/pipeline"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Context input context
type Context struct {
	Done      chan struct{}
	BeatDone  chan struct{}
	OutSQS    chan *pipeline.SQS
	OutS3List chan *pipeline.S3List
}

// Factory is used to register functions creating new Input instances.
type Factory = func(config *common.Config, context Context) (Input, error)

var registry = make(map[string]Factory)

// Register registers an input
func Register(name string, factory Factory) error {
	logp.Info("Registering input factory")
	if name == "" {
		return fmt.Errorf("Error registering input: name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("Error registering input '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering input '%v': already registered", name)
	}

	registry[name] = factory
	logp.Info("Successfully registered input")

	return nil
}

// GetFactory gets a factory from a name
func GetFactory(name string) (Factory, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating input. No such input type exist: '%v'", name)
	}
	return registry[name], nil
}
