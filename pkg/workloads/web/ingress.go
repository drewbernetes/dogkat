package web

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v1 "k8s.io/api/networking/v1"
)

// GenerateWebIngressResource returns an Ingress resource that will be used for Ingress testing.
func GenerateWebIngressResource(namespace string, annotations map[string]string) *coreworkloads.Ingress {
	secret := "e2e-test-secret"
	host := "e2e-test.nl1.eschercloud.dev"

	tls := []v1.IngressTLS{
		{
			Hosts: []string{
				host,
			},
			SecretName: secret,
		},
	}

	ing := &coreworkloads.Ingress{}
	ing.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "className": "nginx", "host": host})
	ing.Resource.Spec.TLS = tls

	for k, v := range annotations {
		ing.Resource.Annotations[k] = v
	}
	return ing
}
