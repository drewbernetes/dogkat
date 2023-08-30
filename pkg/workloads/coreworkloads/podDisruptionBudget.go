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

package coreworkloads

import (
	"context"
	"fmt"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type PodDisruptionBudget struct {
	Client   *kubernetes.Clientset
	Resource *policyv1.PodDisruptionBudget
}

// Generate the base PodDisruptionBudget.
func (p *PodDisruptionBudget) Generate(data map[string]string) {
	p.Resource = &policyv1.PodDisruptionBudget{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["name"]),
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: &intstr.IntOrString{
				Type:   0,
				IntVal: 1,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": data["label"]},
			},
		},
	}
}

// Create creates a PodDisruptionBudget on the Kubernetes cluster.
func (p *PodDisruptionBudget) Create() error {
	log.Printf("creating PodDisruptionBudget:%s...\n", p.Resource.Name)
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	res, err := r.Create(context.Background(), p.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	p.Resource = res
	log.Printf("PodDisruptionBudget:%s created.\n", p.Resource.Name)
	return nil
}

// Validate validates a PodDisruptionBudget on the Kubernetes cluster.
func (p *PodDisruptionBudget) Validate() error {
	var err error
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	p.Resource, err = r.Get(context.Background(), p.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a PodDisruptionBudget in the Kubernetes cluster.
func (p *PodDisruptionBudget) Update() error {
	return nil
}

// Delete deletes a PodDisruptionBudget from the Kubernetes cluster.
func (p *PodDisruptionBudget) Delete() error {
	name := p.Resource.Name
	log.Printf("deleting PodDisruptionBudget:%s...\n", name)
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("PodDisruptionBudget:%s deleted.\n", name)
	return nil
}

func (p *PodDisruptionBudget) GetResourceName() string {
	return p.Resource.Name
}

func (p *PodDisruptionBudget) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", p.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (p *PodDisruptionBudget) IsReady() bool {
	if err := p.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if p.Resource.Status.CurrentHealthy < p.Resource.Status.DesiredHealthy {
		return false
	}
	if p.Resource.Status.DisruptionsAllowed == 0 && p.Resource.Status.CurrentHealthy > 1 {
		return false
	}
	return true
}
