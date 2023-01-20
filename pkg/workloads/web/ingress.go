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

package web

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/constants"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v1 "k8s.io/api/networking/v1"
)

// GenerateWebIngressResource returns an Ingress resource that will be used for Ingress testing.
func GenerateWebIngressResource(namespace, host, ingressClass string, annotations map[string]string, enableTLS bool) *coreworkloads.Ingress {
	secret := "e2e-test-secret"

	tls := []v1.IngressTLS{}
	if enableTLS {
		tls = append(tls, v1.IngressTLS{
			Hosts: []string{
				host,
			},
			SecretName: secret,
		})
	}

	ing := &coreworkloads.Ingress{}
	ing.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "className": ingressClass, "host": host})
	ing.Resource.Spec.TLS = tls
	ing.Resource.Annotations = map[string]string{}

	for k, v := range annotations {
		ing.Resource.Annotations[k] = v
	}
	return ing
}
