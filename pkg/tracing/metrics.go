/*
Copyright 2022 EscherCloud.
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log"
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

func (t *Duration) SetupMetricsGatherer(name, help string) {
	t.Timer = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})

	t.Registry = prometheus.NewRegistry()
	t.Registry.MustRegister(t.Timer)
	t.Pusher = push.New(t.PushURL, t.JobName).Gatherer(t.Registry)
}

func (t *Duration) Start() {
	t.Begin = time.Now()
}

func (t *Duration) CompleteGathering() {
	t.Timer.Set(time.Since(t.Begin).Seconds())
	if err := t.Pusher.Add(); err != nil {
		log.Println("Could not push metrics to PushGateway:", err)
	}
}
