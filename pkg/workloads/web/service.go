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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateNginxServiceResource returns a Service resource for the Nginx workload testing.
func GenerateNginxServiceResource(namespace string) *coreworkloads.Service {
	selectors := map[string]string{
		"app": constants.NginxName,
	}
	svc := &coreworkloads.Service{}
	svc.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "label": constants.NginxName})
	svc.Resource.Spec.Selector = selectors
	svc.Resource.Spec.Type = v1.ServiceTypeClusterIP
	svc.Resource.Spec.Ports = []v1.ServicePort{
		{
			Name:     "http",
			Protocol: v1.ProtocolTCP,
			Port:     80,
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: 80,
			},
		},
		{
			Name:     "http-metrics",
			Protocol: v1.ProtocolTCP,
			Port:     9113,
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: 9113,
			},
		},
	}

	return svc
}
