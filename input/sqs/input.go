package log

import (
	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/input"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

func init() {
	err := input.Register("sqs", NewInput)
	if err != nil {
		panic(err)
	}
}

// Input contains the input and its config
type Input struct {
	cfg     *common.Config
	config  config
	done    chan struct{}
	chanSQS chan *aws.SQS
}

// NewInput instantiates a new Log
func NewInput(
	cfg *common.Config,
	context input.Context,
) (input.Input, error) {
	p := &Input{
		config:  defaultConfig,
		cfg:     cfg,
		done:    context.Done,
		chanSQS: context.ChanSQS,
	}

	if err := cfg.Unpack(&p.config); err != nil {
		return nil, err
	}

	return p, nil
}

// Run runs the input
func (p *Input) Run() {
	logp.Debug("input-sqs", "Start next scan")
	awsSession := aws.NewSession()
	for _, queue := range p.config.QueuesURL {
		sqs := aws.NewSQS(awsSession, &queue)
		//p.chanSQS <- sqs
		_, err := sqs.ReceiveMessages(func(message *aws.SQSMessage) error {
			logp.Debug("input-sqs", "Message: %v", message)
			// Generate object to read from S3 and pass to output
			return nil
		})
		if err != nil {
			logp.Err("Could not receive SQS messages: %v", err)
		}

	}
}

// Wait waits for the all harvesters to complete and only then call stop
func (p *Input) Wait() {
	//p.harvesters.WaitForCompletion()
	//p.Stop()
}

// Stop stops all harvesters and then stops the input
func (p *Input) Stop() {
	// Stop all harvesters
	// In case the beatDone channel is closed, this will not wait for completion
	// Otherwise Stop will wait until output is complete
	//p.harvesters.Stop()

	// close state updater
	//p.stateOutlet.Close()

	// stop all communication between harvesters and publisher pipeline
	//p.outlet.Close()
}
