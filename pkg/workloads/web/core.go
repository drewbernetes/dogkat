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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/helpers"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type NginxWorkloads struct {
	Workload            *coreworkloads.Deployment
	Configuration       *coreworkloads.ConfigMap
	WebPages            *coreworkloads.ConfigMap
	ServiceAccount      *coreworkloads.ServiceAccount
	Service             *coreworkloads.Service
	PodDisruptionBudget *coreworkloads.PodDisruptionBudget
}

// CreateNginxWorkloadItems creates the Nginx workloads required for the core tests to run.
func CreateNginxWorkloadItems(client *kubernetes.Clientset, namespace, cpuRequest, memoryRequest string) *NginxWorkloads {
	w := &NginxWorkloads{}

	w.Configuration = GenerateNginxConfigMap(namespace)
	w.Configuration.Client = client

	w.WebPages = GenerateWebpageConfigMap(namespace)
	w.WebPages.Client = client

	w.ServiceAccount = &coreworkloads.ServiceAccount{}
	w.ServiceAccount.Generate(map[string]string{"namespace": namespace, "name": constants.NginxSAName, "label": constants.NginxName})
	w.ServiceAccount.Client = client

	w.Workload = GenerateNginxDeploy(namespace, cpuRequest, memoryRequest)
	w.Workload.Client = client

	w.Service = GenerateNginxServiceResource(namespace)
	w.Service.Client = client

	w.PodDisruptionBudget = &coreworkloads.PodDisruptionBudget{}
	w.PodDisruptionBudget.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "label": constants.NginxName})
	w.PodDisruptionBudget.Client = client

	helpers.HandleCreateError(w.Configuration.Create())
	helpers.HandleCreateError(w.WebPages.Create())
	helpers.HandleCreateError(w.ServiceAccount.Create())
	helpers.HandleCreateError(w.Workload.Create())
	helpers.HandleCreateError(w.Service.Create())
	helpers.HandleCreateError(w.PodDisruptionBudget.Create())

	return w
}

// ValidateNginxWorkloadItems validates the Nginx workloads required for the core tests to run.
func (w *NginxWorkloads) ValidateNginxWorkloadItems() error {
	var err error

	err = w.Configuration.Validate()
	if err != nil {
		return err
	}
	err = w.WebPages.Validate()
	if err != nil {
		return err
	}
	err = w.ServiceAccount.Validate()
	if err != nil {
		return err
	}
	err = w.Workload.Validate()
	if err != nil {
		return err
	}
	err = w.Service.Validate()
	if err != nil {
		return err
	}
	err = w.PodDisruptionBudget.Validate()
	if err != nil {
		return err
	}

	return nil
}

// DeleteNginxWorkloadItems deletes the Nginx workloads required for the core tests to run.
func DeleteNginxWorkloadItems(client *kubernetes.Clientset, namespace string) {
	nginxCM := &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxConfName, Namespace: namespace}},
	}

	nginxWebPages := &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxPagesName, Namespace: namespace}},
	}

	nginxSa := &coreworkloads.ServiceAccount{
		Client:   client,
		Resource: &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxSAName, Namespace: namespace}},
	}

	nginxDeploy := &coreworkloads.Deployment{
		Client:   client,
		Resource: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	nginxSvc := &coreworkloads.Service{
		Client:   client,
		Resource: &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	nginxPdb := &coreworkloads.PodDisruptionBudget{
		Client:   client,
		Resource: &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	err := nginxCM.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxWebPages.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxSa.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxDeploy.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxSvc.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxPdb.Delete()
	if err != nil {
		log.Println(err)
	}
}

// CreateIngressResource creates the ingress resource for testing.
func CreateIngressResource(client *kubernetes.Clientset, namespace, annotations, host, ingressClass string, enableTLS bool) *coreworkloads.Ingress {
	a := map[string]string{}
	if len(annotations) != 0 {
		p := strings.Split(annotations, ",")
		for _, v := range p {
			s := strings.Split(v, "=")
			if len(s) == 2 {
				a[s[0]] = s[1]
			}
		}
	}

	ing := GenerateWebIngressResource(namespace, host, ingressClass, a, enableTLS)
	ing.Client = client

	helpers.HandleCreateError(ing.Create())
	return ing
}

// ValidateIngressResource validates the ingress resource for testing.
func ValidateIngressResource(ing *coreworkloads.Ingress) {
	err := ing.Validate()
	if err != nil {
		log.Fatalln(err)
	}
}

// DeleteIngressWorkloadItems deletes the Ingress resource.
func DeleteIngressWorkloadItems(client *kubernetes.Clientset, namespace string) {
	var ing = coreworkloads.Ingress{
		Client: client,
		Resource: &networkv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: constants.NginxName},
		},
	}
	err := ing.Delete()
	if err != nil {
		log.Println(err)
	}
}
