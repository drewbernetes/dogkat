package sql

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GeneratePostgresServiceResource returns a Service resource that will be used for postgres testing.
func GeneratePostgresServiceResource(namespace string) *coreworkloads.Service {
	selectors := map[string]string{
		"app": constants.PGSqlName,
	}

	svc := &coreworkloads.Service{}
	svc.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlName, "label": constants.PGSqlName})
	svc.Resource.Spec.Selector = selectors
	svc.Resource.Spec.ClusterIP = v1.ClusterIPNone
	svc.Resource.Spec.Ports = []v1.ServicePort{
		{
			Name:     "sql",
			Protocol: v1.ProtocolTCP,
			Port:     5432,
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: 5432,
			},
		},
	}

	return svc
}
