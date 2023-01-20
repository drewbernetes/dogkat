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

package prometheus

import (
	"context"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v12 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// TODO rewrite all of this to match new format

// GenerateServiceMonitorResource returns a ServiceMonitor resource that will be used for testing.
func GenerateServiceMonitorResource(namespace string) *v12.ServiceMonitor {
	sm := &v12.ServiceMonitor{
		ObjectMeta: coreworkloads.GenerateMetadata(namespace, "nginx-e2e", "nginx-e2e"),
		Spec: v12.ServiceMonitorSpec{
			JobLabel: "nginx-e2e",
			Endpoints: []v12.Endpoint{
				{
					Port:     "http-metrics",
					Path:     "/metrics",
					Interval: "1m",
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/instance": "nginx-e2e",
					"app.kubernetes.io/name":     "nginx-e2e",
				},
			},
			NamespaceSelector: v12.NamespaceSelector{
				MatchNames: []string{namespace},
			},
		},
	}
	return sm
}

// CreateServiceMonitor creates a ServiceMonitor on the Kubernetes cluster.
func CreateServiceMonitor(client *promclient.Clientset, w *v12.ServiceMonitor) {
	log.Printf("creating %s:%s...\n", w.Kind, w.Name)
	r := client.MonitoringV1().ServiceMonitors(w.Namespace)
	_, err := r.Create(context.Background(), w, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%s:%s created.\n", w.Kind, w.Name)

	log.Printf("confirming %s:%s...\n", w.Kind, w.Name)
	res, err := r.Get(context.Background(), w.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("* %s\n", res.Name)
}
