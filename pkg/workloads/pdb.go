/*
Copyright 2024 Drewbernetes.

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

package workloads

import (
	"github.com/drewbernetes/dogkat/pkg/helm"
	"golang.org/x/net/context"
	policyv1 "k8s.io/api/policy/v1"
	optsv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/policy/v1"
)

type PodDisruptionBudget struct {
	PDBName                      string
	PodDisruptionBudget          *policyv1.PodDisruptionBudget
	PodDisruptionBudgetInterface v1.PodDisruptionBudgetInterface
}

func NewPodDisruptionBudget(client *helm.Client, name string) (*PodDisruptionBudget, error) {
	var err error
	p := &PodDisruptionBudget{}

	p.PDBName = name
	p.PodDisruptionBudgetInterface = client.KubeClient.PolicyV1().PodDisruptionBudgets(client.Settings.Namespace())
	p.PodDisruptionBudget, err = p.get()

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PodDisruptionBudget) get() (*policyv1.PodDisruptionBudget, error) {
	return p.PodDisruptionBudgetInterface.Get(context.Background(), p.PDBName, optsv1.GetOptions{})
}

func (p *PodDisruptionBudget) Name() string {
	return p.PodDisruptionBudget.ObjectMeta.Name
}

func (p *PodDisruptionBudget) Kind() string {
	return "PodDisruptionBudget"
}

func (p *PodDisruptionBudget) IsReady() bool {
	var err error
	p.PodDisruptionBudget, err = p.get()
	if err != nil {
		return false
	}

	if p.PodDisruptionBudget.Status.CurrentHealthy < p.PodDisruptionBudget.Status.DesiredHealthy {
		return false
	}
	if p.PodDisruptionBudget.Status.DisruptionsAllowed == 0 && p.PodDisruptionBudget.Status.CurrentHealthy > 1 {
		return false
	}
	return true
}
