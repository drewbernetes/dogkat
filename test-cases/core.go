package test_cases

import (
	"context"
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/hack"
	"github.com/drew-viles/k8s-e2e-tester/resources"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

// CoreWorkloadChecks will run the basic tests. Deployments, Ingress, Cluster scaling, Cluster DNS validation
func CoreWorkloadChecks(obj resources.ApiResource, res chan resources.ResourceReady) {
	r := resources.ResourceReady{}
	r.Resource = obj
	statusResults := checkIfResourceIsReady(r.Resource, 0, 5)
	if !statusResults {
		log.Printf("%s:%s is not ready\n", obj.GetResourceKind(), obj.GetResourceName())
		r.Ready = false
		res <- r
		return
	}
	log.Printf("%s:%s is ready\n", obj.GetResourceKind(), obj.GetResourceName())
	r.Ready = true
	res <- r
}

func checkIfResourceIsReady(r resources.ApiResource, counter int, delaySeconds time.Duration) bool {
	delay := time.Second * delaySeconds
	if counter >= 100 {
		return false
	}
	r.Get()
	log.Printf("Waiting for resource to be ready: %s/%s\n", r.GetResourceKind(), r.GetResourceName())
	if !r.IsReady() {
		time.Sleep(delay)
		return checkIfResourceIsReady(r, counter+1, delaySeconds)
	}
	return true
}

// RunScalingTest will
func RunScalingTest(r resources.ApiResource, clientsets *resources.ClientSets) bool {
	replicaSize := int32(20)
	resource := r.(*resources.DeploymentResource)
	//Get number of nodes
	_, initialNodeCount := countNodes(clientsets)
	//Scale up the workload
	initialReplicaSize := *resource.Resource.Spec.Replicas
	log.Println("Testing cluster scaling")
	log.Printf("Replicas before Scale %v\n", initialReplicaSize)
	log.Printf("Nodes before Scale %v\n", initialNodeCount)

	resource.Resource.Spec.Replicas = hack.IntPtr(replicaSize)

	resource.Update()

	log.Printf("Waiting for Deployment/StatefulSet to scale\n")
	time.Sleep(time.Second * 60)

	isReady := checkIfResourceIsReady(r, 0, 5)
	if !isReady {
		log.Fatalf("there was a problem scaling up the resource - it was not considered ready - you may need to ensure your nodes can support %v of these workloads\n", replicaSize)
		return false
	}

	//Get number of nodes
	_, newNodeAmount := countNodes(clientsets)
	if newNodeAmount <= initialNodeCount {
		log.Printf("The node count did not increase - either the nodes were not required, cluster-autoscaler didn't kick in or you're running a single node cluster\n")
		//Scale down the workload
		log.Printf("Replicas after Scale %v\n", *resource.Resource.Spec.Replicas)
		resource.Resource.Spec.Replicas = hack.IntPtr(initialReplicaSize)
		resource.Update()
		return true
	}
	log.Printf("Replicas after Scale %v\n", *resource.Resource.Spec.Replicas)
	log.Printf("Nodes after Scale %v\n", newNodeAmount)

	//Scale down the workload
	resource.Resource.Spec.Replicas = hack.IntPtr(initialReplicaSize)
	resource.Update()

	return true
}

func ScalingValidation(resource resources.ApiResource) {
	switch resource.GetResourceKind() {
	case "Deployment":
		if !resource.IsReady() {
			log.Printf("%s:%s - all replicas up and running\n", resource.GetResourceKind(), resource.GetResourceName())
		}
	case "DaemonSet":
		if !resource.IsReady() {
			log.Printf("%s:%s - all pods up and running\n", resource.GetResourceKind(), resource.GetResourceName())
		}
	case "StatefulSet":
		if !resource.IsReady() {
			log.Printf("%s:%s - all replicas up and running\n", resource.GetResourceKind(), resource.GetResourceName())
		}
	case "Ingress":
		ing := resource.GetObject().(*networkingv1.Ingress)
		//Check ingress hostnames are responding
		err := testIngress(ing.Spec.TLS)
		if err != nil {
			log.Fatalf(fmt.Sprintf("ingress request error: %s\n", err.Error()))
		}
		log.Printf("%s:%s - responding as expected\n", resource.GetResourceKind(), resource.GetResourceName())
	//case "Gateway":
	case "VirtualService":
		vs := resource.GetObject().(*v1beta1.VirtualService)
		//Check VirtualService hostnames are responding
		for _, host := range vs.Spec.Hosts {
			err := testHostEndpoints(host, 0)
			if err != nil {
				log.Fatalf(fmt.Sprintf("ingress request error: %s\n", err.Error()))
			}
			log.Printf("%s:%s - responding as expected\n", resource.GetResourceKind(), resource.GetResourceName())
		}
	}
}

func countNodes(clientsets *resources.ClientSets) (*v1.NodeList, int) {
	allNodes, err := clientsets.K8S.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return nil, 0
	}
	return allNodes, len(allNodes.Items)
}
