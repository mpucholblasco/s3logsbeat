package registrar

import (
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/sequra/s3logsbeat/pipeline"
)

// Registrar registers events published on pipeline
type Registrar struct {
	Channel chan []pipeline.S3ObjectProcessNotifications
	out     successLogger
	done    chan struct{}
	wg      sync.WaitGroup
}

type successLogger interface {
	Published(n int) bool
}

// New new registrar
func New(out successLogger) *Registrar {
	return &Registrar{
		done:    make(chan struct{}),
		Channel: make(chan []pipeline.S3ObjectProcessNotifications, 1),
		out:     out,
		wg:      sync.WaitGroup{},
	}
}

// Start starts the registrar
func (r *Registrar) Start() {
	r.wg.Add(1)
	go r.Run()
}

// Run runs the registrar
func (r *Registrar) Run() {
	logp.Info("Starting Registrar")
	// Writes registry on shutdown
	defer func() {
		r.wg.Done()
	}()

	for {
		select {
		case <-r.done:
			logp.Info("Ending Registrar")
			return
		case states := <-r.Channel:
			r.onEvents(states)
		}
	}
}

// onEvents processes events received from the publisher pipeline
func (r *Registrar) onEvents(states []pipeline.S3ObjectProcessNotifications) {
	logp.Debug("registrar", "Processing %d events", len(states))

	for _, s := range states {
		s.EventACKed()
	}

	r.out.Published(len(states))
}

// Stop stops the registry. It waits until Run function finished.
func (r *Registrar) Stop() {
	logp.Info("Stopping Registrar")
	close(r.done)
	r.wg.Wait()
	logp.Info("Registrar stopped")
}
