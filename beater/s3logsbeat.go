package beater

import (
	"fmt"

	"github.com/sequra/s3logsbeat/config"
	"github.com/sequra/s3logsbeat/crawler"
	"github.com/sequra/s3logsbeat/pipeline"
	"github.com/sequra/s3logsbeat/registrar"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/monitoring"
)

// S3logsbeat beater
type S3logsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates beater
func NewS3logsbeat(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
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

// Run runs beater
func (bt *S3logsbeat) Run(b *beat.Beat) error {
	logp.Info("s3logsbeat is running! Hit CTRL-C to stop it.")

	var err error

	waitFinished := newSignalWait()
	waitEvents := newSignalWait()

	// count SQS messages for monitoring purposes
	wgSQSMessages := &eventCounter{
		count: monitoring.NewInt(nil, "s3logsbeat.sqsMessages.active"),
		added: monitoring.NewUint(nil, "s3logsbeat.sqsMessages.added"),
		done:  monitoring.NewUint(nil, "s3logsbeat.sqsMessages.done"),
		err:   monitoring.NewUint(nil, "s3logsbeat.sqsMessages.pullError"),
	}

	// count S3 objects for monitoring purposes
	wgS3Objects := &eventCounter{
		count: monitoring.NewInt(nil, "s3logsbeat.s3objects.active"),
		added: monitoring.NewUint(nil, "s3logsbeat.s3objects.added"),
		done:  monitoring.NewUint(nil, "s3logsbeat.s3objects.done"),
		err:   monitoring.NewUint(nil, "s3logsbeat.s3object.readError"),
	}

	// count active events for waiting on shutdown
	wgEvents := &eventCounter{
		count: monitoring.NewInt(nil, "s3logsbeat.events.active"),
		added: monitoring.NewUint(nil, "s3logsbeat.events.added"),
		done:  monitoring.NewUint(nil, "s3logsbeat.events.done"),
		err:   monitoring.NewUint(nil, "s3logsbeat.events.parserError"),
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

	pipelineChannels := pipeline.NewChannels()

	crawler, err := crawler.New(
		bt.config.Inputs,
		b.Info.Version,
		bt.done,
		*once,
		pipelineChannels.GetSQSChannel(),
		nil,
		[]string{"sqs"},
	)
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

	// Start the pipeline workers
	s3readerWorker := pipeline.NewS3ReaderWorker(pipelineChannels.GetS3Channel(), bt.client, wgEvents, wgS3Objects)
	sqsConsumerWorker := pipeline.NewSQSConsumerWorker(pipelineChannels.GetSQSChannel(), pipelineChannels.GetS3Channel(), wgSQSMessages, wgS3Objects, *keepSQSMessages)
	s3readerWorker.Start()
	sqsConsumerWorker.Start()

	// If run once, add crawler completion check as alternative to done signal
	if *once {
		runOnce := func() {
			logp.Info("Running s3logsbeat once. Waiting for completion ...")
			crawler.WaitForCompletion()
			pipelineChannels.CloseSQSChannel()
			sqsConsumerWorker.Wait()
			pipelineChannels.CloseS3Channel()
			s3readerWorker.Wait()
			wgEvents.Wait()
			logp.Info("All data collection completed. Shutting down.")
		}
		waitFinished.Add(runOnce)
	}

	// Add done channel to wait for shutdown signal
	waitFinished.AddChan(bt.done)
	waitFinished.Wait()

	crawler.Stop()
	sqsConsumerWorker.StopAcceptingMessages()

	timeout := bt.config.ShutdownTimeout
	// Checks if on shutdown it should wait for all events to be published
	waitPublished := timeout > 0 || *once
	if waitPublished {
		// Wait until all will be done + all events published
		waitEvents.Add(withLog(func() {
			sqsConsumerWorker.Wait()
			pipelineChannels.CloseS3Channel()
			s3readerWorker.Wait()
			wgEvents.Wait()
		},
			"Continue shutdown: All enqueued events being published."))
		// Wait for either timeout or all events having been ACKed by outputs.
		if timeout > 0 {
			logp.Info("Shutdown output timer started. Waiting for max %v.", timeout)
			waitEvents.Add(withLog(waitDuration(timeout),
				"Continue shutdown: Time out waiting for events being published."))
		} else {
			waitEvents.AddChan(bt.done)
		}
	}

	// Wait for all events to be processed or timeout
	logp.Debug("s3logsbeat", "Waiting for all events to be processed or timeout")
	waitEvents.Wait()

	sqsConsumerWorker.Stop()
	bt.client.Close() // unlock publish events (if locked)
	s3readerWorker.Stop()

	// Close registrar
	logp.Debug("s3logsbeat", "Stopping registrar")
	registrar.Stop()
	registrarChannel.Close()

	return nil
}

// Stop stops beater
func (bt *S3logsbeat) Stop() {
	close(bt.done)
}
