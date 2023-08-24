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

package testsuite

import (
	"fmt"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/helpers"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/tracing"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

func CheckReadyForTesting(resource []coreworkloads.Resource) error {
	//Thread the tests to run in parallel
	checksCompleted := make(chan coreworkloads.ResourceReady, len(resource))

	readyCheck := func(obj coreworkloads.Resource) {
		r := coreworkloads.ResourceReady{}
		r.Resource = obj
		statusResults := checkIfResourceIsReady(obj, 0, 5)
		if !statusResults {
			log.Printf("%s:%s is not ready\n", obj.GetResourceKind(), obj.GetResourceName())
			r.Ready = false
			checksCompleted <- r
			return
		}
		log.Printf("%s:%s is ready\n", obj.GetResourceKind(), obj.GetResourceName())
		r.Ready = true
		checksCompleted <- r
	}

	defer close(checksCompleted)
	for _, r := range resource {
		if r == nil {
			continue
		}
		go readyCheck(r)
	}

	for range resource {
		<-checksCompleted
	}

	return nil
}

// checkIfResourceIsReady validates the readiness of the resource.
func checkIfResourceIsReady(r coreworkloads.Resource, counter int, delaySeconds time.Duration) bool {
	delay := time.Second * delaySeconds
	if counter >= 100 {
		return false
	}
	log.Printf("Waiting for resource to be ready: %s/%s\n", r.GetResourceKind(), r.GetResourceName())
	if !r.IsReady() {
		time.Sleep(delay)
		return checkIfResourceIsReady(r, counter+1, delaySeconds)
	}
	return true
}

// ScaleUpStandardNodes scales up the standard nodes that generic workloads will sit on.
func ScaleUpStandardNodes(resource *coreworkloads.Deployment) error {

	tracer := tracing.Tracer{JobName: "e2e_workloads", PushURL: "http://prometheus-push-gateway.prometheus:9091"}
	tracer.NewTimer("scale_workload", "Times the scaling of workloads to determine how long the autoscaler takes")
	timer := tracer.Start()

	replicaSize := int32(5)
	//Get number of nodes
	initialNodeCount := countNodes(resource.Client)
	//Scale up the workload
	initialReplicaSize := *resource.Resource.Spec.Replicas
	log.Println("Testing cluster scaling")
	log.Printf("Node count before Scale %v\n", initialNodeCount)

	resource.Resource.Spec.Replicas = helpers.IntPtr(replicaSize)

	if err := resource.Update(); err != nil {
		return fmt.Errorf("Failed to increase replicas for %s:%s: %v\n", resource.Resource.Kind, resource.Resource.Name, err)
	}

	log.Printf("Waiting for Deployment to scale\n")
	time.Sleep(time.Second * 60)

	isReady := checkIfResourceIsReady(resource, 0, 5)
	if !isReady {
		return fmt.Errorf("There was a problem scaling up the resource - it was not considered ready - you may need to ensure your nodes can support %v of these workloads\n", replicaSize)
	}

	//Get number of nodes
	newNodeAmount := countNodes(resource.Client)
	if newNodeAmount <= initialNodeCount {
		log.Printf("The node count did not increase - either the nodes were not required, cluster-autoscaler didn't kick in or you're running a single node cluster\n")
		//Scale down the workload
		resource.Resource.Spec.Replicas = helpers.IntPtr(initialReplicaSize)
		if err := resource.Update(); err != nil {
			log.Printf("Failed to restore replias for %s:%s: %v\n", resource.Resource.Kind, resource.Resource.Name, err)
		}
		return nil
	}
	log.Printf("Replicas after Scale %v\n", *resource.Resource.Spec.Replicas)
	log.Printf("Nodes after Scale %v\n", newNodeAmount)

	//End timer here as we only care about the scale in time as we need to know how quick the autoscaling is.
	timer.ObserveDuration()

	//Scale down the workload
	resource.Resource.Spec.Replicas = helpers.IntPtr(initialReplicaSize)
	if err := resource.Update(); err != nil {
		log.Printf("Failed to restore replias for %s:%s: %v\n", resource.Resource.Kind, resource.Resource.Name, err)
	}

	return nil
}

// countNodes returns the current number of nodes in the cluster.
func countNodes(client *kubernetes.Clientset) int {
	allNodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return len(allNodes.Items)
}
