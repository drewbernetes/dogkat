package resources

import (
	promscheme "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/scheme"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"
	// _ "K8S.io/client-go/plugin/pkg/client/auth/oidc"
)

// TODO: Template this to prevent repetitive code.
var manifests = []string{
	`---
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: v1
kind: Secret
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: v1
kind: Service
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: batch/v1
kind: Job
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: test
  namespace: test
`,
	`---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: test
  namespace: test
`,
}

// TestParseResourceKind simply tests the ParseResourceKind by passing a POD resource in.
func TestParseResourceKind(t *testing.T) {
	promscheme.AddToScheme(scheme.Scheme)
	istioscheme.AddToScheme(scheme.Scheme)
	for _, manifest := range manifests {
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode([]byte(manifest), nil, nil)
		if err != nil {
			t.Errorf("There was an error decoding: %s, %s\n", manifest, err)
		}

		r := ParseResourceKind(obj)
		if r == nil {
			t.Errorf("There was no resources returned suggesting an error with the parser or missing API Resource.\n")
		}
	}
}
