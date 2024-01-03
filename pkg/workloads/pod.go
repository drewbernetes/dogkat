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
	"github.com/eschercloudai/dogkat/pkg/constants"
	"github.com/eschercloudai/dogkat/pkg/helm"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	optsv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Pod struct {
	*v1.Pod
	v1typed.PodInterface
}

func NewPod(client *helm.Client) (*Pod, error) {
	var err error
	p := &Pod{}

	p.PodInterface = client.KubeClient.CoreV1().Pods(client.Settings.Namespace())
	p.Pod, err = p.get()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Pod) get() (*v1.Pod, error) {
	return p.PodInterface.Get(context.Background(), constants.GPUName, optsv1.GetOptions{})
}

func (p *Pod) Name() string {
	return p.Pod.ObjectMeta.Name
}

func (p *Pod) Kind() string {
	return "Pod"
}

func (p *Pod) IsReady() bool {
	var err error
	p.Pod, err = p.get()
	if err != nil {
		return false
	}

	ready := false
	for _, container := range p.Status.ContainerStatuses {
		if !container.Ready {
			ready = false
		}
	}
	if p.Status.Phase == v1.PodSucceeded {
		ready = true
	}
	return ready
}
