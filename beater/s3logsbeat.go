package beater

import (
	"flag"
	"fmt"

	"github.com/mpucholblasco/s3logsbeat/aws"
	"github.com/mpucholblasco/s3logsbeat/crawler"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/mpucholblasco/s3logsbeat/config"
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

	waitFinished := newSignalWait()

	var err error
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

	err = crawler.Start()
	if err != nil {
		crawler.Stop()
		return err
	}

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

	//timeout := fb.config.ShutdownTimeout
	// Checks if on shutdown it should wait for all events to be published
	// waitPublished := timeout > 0 || *once
	// if waitPublished {
	// 	// Wait for registrar to finish writing registry
	// 	waitEvents.Add(withLog(wgEvents.Wait,
	// 		"Continue shutdown: All enqueued events being published."))
	// 	// Wait for either timeout or all events having been ACKed by outputs.
	// 	if fb.config.ShutdownTimeout > 0 {
	// 		logp.Info("Shutdown output timer started. Waiting for max %v.", timeout)
	// 		waitEvents.Add(withLog(waitDuration(timeout),
	// 			"Continue shutdown: Time out waiting for events being published."))
	// 	} else {
	// 		waitEvents.AddChan(fb.done)
	// 	}
	// }

	return nil
}

func (bt *S3logsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
