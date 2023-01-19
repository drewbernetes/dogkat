package web

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
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
