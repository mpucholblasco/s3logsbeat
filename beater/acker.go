package beater

import "github.com/mpucholblasco/s3logsbeat/pipeline"

// eventAcker handles publisher pipeline ACKs and forwards
// them to the registrar.
type eventACKer struct {
	out successLogger
}

type successLogger interface {
	Published(privateElements []pipeline.S3ObjectProcessNotifications)
}

func newEventACKer(out successLogger) *eventACKer {
	return &eventACKer{out: out}
}

func (a *eventACKer) ackEvents(data []interface{}) {
	states := make([]pipeline.S3ObjectProcessNotifications, 0, len(data))
	for _, datum := range data {
		if datum == nil {
			continue
		}

		st, ok := datum.(pipeline.S3ObjectProcessNotifications)
		if !ok {
			continue
		}

		states = append(states, st)
	}

	if len(states) > 0 {
		a.out.Published(states)
	}
}
