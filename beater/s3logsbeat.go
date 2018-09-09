package beater

import (
	"flag"
	"fmt"

	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/config"
	"github.com/mpucholblasco/s3logsbeat/crawler"
	"github.com/mpucholblasco/s3logsbeat/registrar"
	"github.com/mpucholblasco/s3logsbeat/worker"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/monitoring"
)

var (
	once = flag.Bool("once", false, "Run s3logsbeat only once until all inputs will be read")
)

type S3logsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &S3logsbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *S3logsbeat) Run(b *beat.Beat) error {
	logp.Info("s3logsbeat is running! Hit CTRL-C to stop it.")

	var err error

	waitFinished := newSignalWait()
	waitEvents := newSignalWait()

	// count active events for waiting on shutdown
	wgEvents := &eventCounter{
		count: monitoring.NewInt(nil, "s3logsbeat.events.active"),
		added: monitoring.NewUint(nil, "s3logsbeat.events.added"),
		done:  monitoring.NewUint(nil, "s3logsbeat.events.done"),
	}
	finishedLogger := newFinishedLogger(wgEvents)

	// Setup registrar to persist state
	registrar := registrar.New(finishedLogger)

	// Make sure all events that were published in
	registrarChannel := newRegistrarLogger(registrar)

	err = b.Publisher.SetACKHandler(beat.PipelineACKHandler{
		ACKEvents: newEventACKer(registrarChannel).ackEvents,
	})
	if err != nil {
		logp.Err("Failed to install the registry with the publisher pipeline: %v", err)
		return err
	}

	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	chanSQS := make(chan *aws.SQS, 100)

	crawler, err := crawler.New(
		bt.config.Inputs,
		b.Info.Version,
		bt.done,
		*once,
		chanSQS)
	if err != nil {
		logp.Err("Could not init crawler: %v", err)
		return err
	}

	// Start the registrar
	registrar.Start()

	err = crawler.Start()
	if err != nil {
		crawler.Stop()
		return err
	}

	chanS3 := make(chan *aws.S3ObjectSQSMessage, 100)

	s3readerWorker := worker.NewS3ReaderWorker(chanS3, bt.client, wgEvents)
	s3readerWorker.Start()

	sqsConsumerWorker := worker.NewSQSConsumerWorker(chanSQS, chanS3)
	sqsConsumerWorker.Start()

	// If run once, add crawler completion check as alternative to done signal
	if *once {
		runOnce := func() {
			logp.Info("Running s3logsbeat once. Waiting for completion ...")
			crawler.WaitForCompletion()
			logp.Info("All data collection completed. Shutting down.")
		}
		waitFinished.Add(runOnce)
	}

	// Add done channel to wait for shutdown signal
	waitFinished.AddChan(bt.done)
	waitFinished.Wait()

	crawler.Stop()
	sqsConsumerWorker.Stop() // Do not read more SQS messages because we are closing

	timeout := bt.config.ShutdownTimeout
	// Checks if on shutdown it should wait for all events to be published
	waitPublished := timeout > 0 || *once
	if waitPublished {
		logp.Debug("s3logsbeat", "AAAA")
		// Wait for registrar to finish writing registry
		waitEvents.Add(withLog(wgEvents.Wait,
			"Continue shutdown: All enqueued events being published."))
		// Wait for either timeout or all events having been ACKed by outputs.
		if timeout > 0 {
			logp.Info("Shutdown output timer started. Waiting for max %v.", timeout)
			waitEvents.Add(withLog(waitDuration(timeout),
				"Continue shutdown: Time out waiting for events being published."))
		} else {
			logp.Debug("s3logsbeat", "BBBB")
			waitEvents.AddChan(bt.done)
		}
	}

	// Wait for all events to be processed or timeout
	logp.Debug("s3logsbeat", "Waiting for all events to be processed")
	waitEvents.Wait()

	logp.Debug("s3logsbeat", "Stopping S3 reader workers")
	s3readerWorker.Stop()

	// Close publisher
	bt.client.Close()

	// Close registrar
	logp.Debug("s3logsbeat", "Stopping registrar")
	registrar.Stop()
	registrarChannel.Close()

	return nil
}

func (bt *S3logsbeat) Stop() {
	close(bt.done)
}
