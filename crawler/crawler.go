package crawler

import (
	"fmt"
	"sync"

	"github.com/sequra/s3logsbeat/input"
	"github.com/sequra/s3logsbeat/pipeline"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	// Used to load all available inputs
	_ "github.com/sequra/s3logsbeat/include"
)

// Crawler object
type Crawler struct {
	inputs       map[uint64]*input.Runner
	inputConfigs []*common.Config
	wg           sync.WaitGroup
	once         bool
	beatVersion  string
	beatDone     chan struct{}
	outSQS       chan *pipeline.SQS
	outS3List    chan *pipeline.S3List
	allowedTypes []string
}

// New creates a new crawler
func New(inputConfigs []*common.Config, beatVersion string, beatDone chan struct{}, once bool, outSQS chan *pipeline.SQS, outS3List chan *pipeline.S3List, allowedTypes []string) (*Crawler, error) {
	return &Crawler{
		inputs:       map[uint64]*input.Runner{},
		inputConfigs: inputConfigs,
		once:         once,
		beatVersion:  beatVersion,
		beatDone:     beatDone,
		outSQS:       outSQS,
		outS3List:    outS3List,
		allowedTypes: allowedTypes,
	}, nil
}

// Start starts the crawler with all inputs
func (c *Crawler) Start() error {
	logp.Info("Loading Inputs: %v", len(c.inputConfigs))

	for _, inputConfig := range c.inputConfigs {
		err := c.startInput(inputConfig)
		if err != nil {
			return err
		}
	}

	logp.Info("Loading and starting Inputs completed. Enabled inputs: %v", len(c.inputs))

	return nil
}

func (c *Crawler) startInput(
	config *common.Config,
) error {
	if !config.Enabled() {
		return nil
	}

	p, err := input.New(config, c.beatDone, c.outSQS, c.outS3List)
	if err != nil {
		return fmt.Errorf("Error in initing input: %s", err)
	}

	if c.isTypeAllowed(p.Type()) {
		p.Once = c.once

		if _, ok := c.inputs[p.ID]; ok {
			return fmt.Errorf("Input with same ID already exists: %d", p.ID)
		}

		c.inputs[p.ID] = p

		p.Start()
	} else {
		logp.Info("Ignoring not allowed type %s", p.Type())
	}

	return nil
}

// Stop stops all inputs in parallel and waits until all them will be stopped
func (c *Crawler) Stop() {
	logp.Info("Stopping Crawler")

	logp.Info("Stopping %v inputs", len(c.inputs))
	for _, i := range c.inputs {
		// Stop inputs in parallel
		c.wg.Add(1)
		go func(i *input.Runner) {
			defer c.wg.Done()
			i.Stop()
		}(i)
	}
	c.wg.Wait()

	logp.Info("Crawler stopped")
}

// WaitForCompletion waits untill all inputs will be stopped
func (c *Crawler) WaitForCompletion() {
	c.wg.Wait()
}

func (c *Crawler) isTypeAllowed(t string) bool {
	for _, e := range c.allowedTypes {
		if e == t {
			return true
		}
	}
	return false
}
