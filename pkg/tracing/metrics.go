package tracing

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Tracer struct {
	Timer   prometheus.Histogram
	PushURL string
	JobName string
}

func (t *Tracer) NewTimer(name, help string) prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: help,
		//TODO: check if this needs extrapolating for different workloads and tests.
		Buckets: []float64{
			1, 5, 10, 20, 30, 45, 60, 90, 120, 240, 600, 1200, 2400,
		},
	})
}

func (t *Tracer) RegisterTimer() {
	prometheus.MustRegister(t.Timer)
}

func (t *Tracer) PushToGateway() {
	pusher := push.New(t.PushURL, t.JobName)

	if err := pusher.Push(); err != nil {
		panic(err)
	}
}

func (t *Tracer) Start() *prometheus.Timer {
	return prometheus.NewTimer(t.Timer)
}
