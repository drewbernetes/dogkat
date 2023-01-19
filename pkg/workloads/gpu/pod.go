package gpu

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
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
