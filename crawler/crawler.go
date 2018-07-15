package crawler

import (
	"fmt"
	"sync"

	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/input"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	_ "github.com/mpucholblasco/s3logsbeat/include"
)

type Crawler struct {
	inputs       map[uint64]*input.Runner
	inputConfigs []*common.Config
	//out             channel.Factory
	wg sync.WaitGroup
	//InputsFactory   cfgfile.RunnerFactory
	once        bool
	beatVersion string
	beatDone    chan struct{}
	chanSQS     chan *aws.SQS
}

func New(inputConfigs []*common.Config, beatVersion string, beatDone chan struct{}, once bool, chanSQS chan *aws.SQS) (*Crawler, error) {
	return &Crawler{
		//out:          out,
		inputs:       map[uint64]*input.Runner{},
		inputConfigs: inputConfigs,
		once:         once,
		beatVersion:  beatVersion,
		beatDone:     beatDone,
		chanSQS:      chanSQS,
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

	/*	c.InputsFactory = input.NewRunnerFactory(c.out, r, c.beatDone)
		if configInputs.Enabled() {
			c.inputReloader = cfgfile.NewReloader(pipeline, configInputs)
			if err := c.inputReloader.Check(c.InputsFactory); err != nil {
				return err
			}

			go func() {
				c.inputReloader.Run(c.InputsFactory)
			}()
		}*/

	logp.Info("Loading and starting Inputs completed. Enabled inputs: %v", len(c.inputs))

	return nil
}

func (c *Crawler) startInput(
	config *common.Config,
) error {
	if !config.Enabled() {
		return nil
	}

	//connector := channel.ConnectTo(pipeline, c.out)
	p, err := input.New(config, c.beatDone, c.chanSQS)
	if err != nil {
		return fmt.Errorf("Error in initing input: %s", err)
	}
	p.Once = c.once

	if _, ok := c.inputs[p.ID]; ok {
		return fmt.Errorf("Input with same ID already exists: %d", p.ID)
	}

	c.inputs[p.ID] = p

	p.Start()

	return nil
}

func (c *Crawler) Stop() {
	logp.Info("Stopping Crawler")

	asyncWaitStop := func(stop func()) {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			stop()
		}()
	}

	logp.Info("Stopping %v inputs", len(c.inputs))
	for _, p := range c.inputs {
		// Stop inputs in parallel
		asyncWaitStop(p.Stop)
	}

	c.WaitForCompletion()

	logp.Info("Crawler stopped")
}

func (c *Crawler) WaitForCompletion() {
	c.wg.Wait()
}
