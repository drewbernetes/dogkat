/*
Copyright 2024 EscherCloud.
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

package options

import (
	"github.com/spf13/viper"
)

type IngressOptions struct {
	Enabled       bool
	IngressClass  string
	Host          string
	Annotations   map[string]string
	EnableTLS     bool
	TLSHost       string
	TLSSecretName string
}

func newIngressOptions() IngressOptions {
	return IngressOptions{
		Enabled:       viper.GetBool("ingress.enabled"),
		IngressClass:  viper.GetString("ingress.ingressClassName"),
		Host:          viper.GetString("ingress.host"),
		Annotations:   viper.GetStringMapString("ingress.annotations"),
		EnableTLS:     viper.GetBool("ingress.tls.enable"),
		TLSHost:       viper.GetString("ingress.tls.host"),
		TLSSecretName: viper.GetString("ingress.tls.secretName"),
	}
}
