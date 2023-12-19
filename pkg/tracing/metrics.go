package tracing

import "github.com/spf13/viper"

var gatherer *MetricsGather

type MetricsGather struct {
	Enabled     bool
	PushGateway string
}

func NewGatherer() *MetricsGather {
	gatherer = &MetricsGather{
		Enabled:     viper.GetBool("metrics.enabled"),
		PushGateway: viper.GetString("metrics.pushGatewayURI"),
	}

	return gatherer
}

func Gatherer() *MetricsGather {
	return gatherer
}
