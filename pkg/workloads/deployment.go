/*
Copyright 2024 EscherCloud.

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
	"github.com/drewbernetes/dogkat/pkg/constants"
	"github.com/drewbernetes/dogkat/pkg/helm"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	optsv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type Deployment struct {
	v1.DeploymentInterface
	*appsv1.Deployment
}

func NewDeployment(client *helm.Client) (*Deployment, error) {
	var err error
	d := &Deployment{}

	d.DeploymentInterface = client.KubeClient.AppsV1().Deployments(client.Settings.Namespace())
	d.Deployment, err = d.get()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Deployment) get() (*appsv1.Deployment, error) {
	return d.DeploymentInterface.Get(context.Background(), constants.NginxName, optsv1.GetOptions{})
}

func (d *Deployment) Get() (*appsv1.Deployment, error) {
	return d.get()
}

func (d *Deployment) Name() string {
	return d.Deployment.ObjectMeta.Name
}

func (d *Deployment) Kind() string {
	return "Deployment"
}

func (d *Deployment) IsReady() bool {
	var err error
	d.Deployment, err = d.get()
	if err != nil {
		return false
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		return false
	}
	return true
}
