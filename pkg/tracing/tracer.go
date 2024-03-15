/*
Copyright 2024 Drewbernetes.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tracing

import (
	"fmt"
	"github.com/drewbernetes/dogkat/pkg/constants"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
)

type Duration struct {
	Timer    prometheus.Gauge
	Registry *prometheus.Registry
	Pusher   *push.Pusher
	PushURL  string
	JobName  string
	Begin    time.Time
}

func NewCollector(pushURL, metricsName, description string) *Duration {
	d := &Duration{JobName: constants.MetricsJobName, PushURL: pushURL}

	d.Timer = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_%s", constants.MetricsPrefix, metricsName),
		Help: description,
	})

	d.Registry = prometheus.NewRegistry()
	d.Registry.MustRegister(d.Timer)
	d.Pusher = push.New(d.PushURL, d.JobName).Gatherer(d.Registry)

	return d
}

func (t *Duration) Start() {
	t.Begin = time.Now()
}

func (t *Duration) CompleteGathering() error {
	t.Timer.Set(time.Since(t.Begin).Seconds())
	if err := t.Pusher.Add(); err != nil {
		return fmt.Errorf("could not push metrics to PushGateway: %s", err.Error())
	}
	return nil
}
