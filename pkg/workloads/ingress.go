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
	networkingv1 "k8s.io/api/networking/v1"
	optsv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

type Ingress struct {
	*networkingv1.Ingress
	v1.IngressInterface
}

func NewIngress(client *helm.Client) (*Ingress, error) {
	var err error
	i := &Ingress{}

	i.IngressInterface = client.KubeClient.NetworkingV1().Ingresses(client.Settings.Namespace())
	i.Ingress, err = i.get()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Ingress) get() (*networkingv1.Ingress, error) {
	return i.IngressInterface.Get(context.Background(), constants.NginxName, optsv1.GetOptions{})
}

func (i *Ingress) Name() string {
	return i.Ingress.ObjectMeta.Name
}

func (i *Ingress) Kind() string {
	return "Ingress"
}

func (i *Ingress) IsReady() bool {
	var err error
	i.Ingress, err = i.get()
	if err != nil {
		return false
	}

	for _, v := range i.Status.LoadBalancer.Ingress {
		if v.Hostname == "" && v.IP == "" {
			return false
		}
	}
	return true
}
