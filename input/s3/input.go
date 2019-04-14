package s3

import (
	"github.com/sequra/s3logsbeat/aws"
	"github.com/sequra/s3logsbeat/input"
	"github.com/sequra/s3logsbeat/logparser"
	"github.com/sequra/s3logsbeat/pipeline"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

func init() {
	err := input.Register("s3", NewInput)
	if err != nil {
		panic(err)
	}
}

// Input contains the input and its config
type Input struct {
	cfg       *common.Config
	config    config
	done      chan struct{}
	out       chan *pipeline.S3List
	logParser logparser.LogParser
}

// NewInput instantiates a new Log
func NewInput(
	cfg *common.Config,
	context input.Context,
) (input.Input, error) {
	p := &Input{
		config: defaultConfig,
		cfg:    cfg,
		done:   context.Done,
		out:    context.OutS3List,
	}

	if err := cfg.Unpack(&p.config); err != nil {
		return nil, err
	}

	var err error
	p.logParser, err = logparser.GetPredefinedParser(p.config.LogFormat)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Run runs the input
func (p *Input) Run() {
	logp.Debug("s3logsbeat", "Start next scan")
	awsSession := aws.NewSession()

	for _, s3uri := range p.config.Buckets {
		s3prefix, err := aws.NewS3ObjectFromURI(s3uri)
		if err != nil {
			logp.Critical("Couldn't parse S3 URI %s", s3uri)
		}
		ri := pipeline.NewS3ReaderInformation(p.logParser, p.config.KeyRegexFields, p.config.LogFormat)
		s3list := pipeline.NewS3List(awsSession, s3prefix, ri, p.config.Since, p.config.To)

		select {
		case p.out <- s3list:
		case <-p.done:
			return
		}
	}
}

// Wait stops the input
// Once the app is goning to stop, we will not accept more SQS messages, so we can stop
// this input directly
func (p *Input) Wait() {
	p.Stop()
}

// Stop stops the input
func (p *Input) Stop() {
	// Nothing to do, as we don't control done channel and it should already be closed
}
