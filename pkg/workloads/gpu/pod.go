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

package gpu

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/constants"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// GenerateGPUPod Generates the GPU workload that will be deployed into the cluster for testing.
func GenerateGPUPod(namespace string, amount string) *coreworkloads.Pod {

	container := coreworkloads.GenerateContainer("cuda-vectoradd", "nvidia/samples", "vectoradd-cuda11.2.1")
	container.Resources = apiv1.ResourceRequirements{
		Limits: apiv1.ResourceList{
			"nvidia.com/gpu": resource.MustParse(amount),
		},
	}
	containers := []apiv1.Container{
		container,
	}

	pod := &coreworkloads.Pod{}
	pod.Generate(map[string]string{"namespace": namespace, "name": constants.GPUName})
	pod.Resource.Spec.Containers = containers

	return pod
}
