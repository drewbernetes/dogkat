package test_cases

import (
	"context"
	"fmt"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

func runCoreTests(v Resource) {
	switch v.GetResourceKind() {
	//case "serviceAccount":
	//case "ConfigMap":
	//case "Secret":
	case "Deployment":
		deploy := v.GetObject().(*appsv1.Deployment)
		if deploy.Status.UnavailableReplicas > 0 {
			log.Printf("there are unavailable replicas in Deployment: %s\n", deploy.Name)
			break
		}
		log.Printf("%s:%s - all replicas up and running\n", v.GetResourceKind(), v.GetResourceName())
	//case "Service":
	//case "Job":
	//case "PersistentVolumeClaim":
	case "Ingress":
		ing := v.GetObject().(*networkingv1.Ingress)
		//Check ingress hostnames are responding
		err := testIngress(ing.Spec.TLS)
		if err != nil {
			log.Printf(fmt.Sprintf("ingress request error: %s\n", err.Error()))
			break
		}
		log.Printf("%s:%s - responding as expected\n", v.GetResourceKind(), v.GetResourceName())
	//case "ServiceMonitor":
	//case "Gateway":
	case "VirtualService":
		vs := v.GetObject().(*v1beta1.VirtualService)
		//Check VirtualService hostnames are responding
		for _, host := range vs.Spec.Hosts {
			testHostEndpoints(host, 0)
		}
	//case "PeerAuthentication":
	//case "DestinationRule":
	default:
		log.Printf("%s:%s - no tests defined for this resource, skipping\n", v.GetResourceKind(), v.GetResourceName())
	}
}

func runScalingTest(r Resource) bool {
	replicaSize := int32(20)
	resource := r.(*DeploymentResource)
	//Get number of nodes
	_, initialNodeCount := countNodes()
	//Scale up the workload
	initialReplicaSize := *resource.Resource.Spec.Replicas
	log.Println("Testing cluster scaling")
	log.Printf("Replicas before Scale %v\n", initialReplicaSize)
	log.Printf("Nodes before Scale %v\n", initialNodeCount)

	resource.Resource.Spec.Replicas = intPtr(replicaSize)

	resource.Update()

	time.Sleep(time.Second * 60)

	isReady := checkIfResourceIsReady(r, 0, 60)
	if !isReady {
		log.Printf("there was a problem scaling up the resource - it was not considered ready - you may need to ensure your nodes can support %v of these workloads\n", replicaSize)
		return false
	}

	//Get number of nodes
	_, newNodeAmount := countNodes()
	if newNodeAmount <= initialNodeCount {
		log.Printf("The node count did not increase - either the nodes were not required, cluster-autoscaler didn't kick in or you're running a single node cluster\n")
		//Scale down the workload
		log.Printf("Replicas after Scale %v\n", *resource.Resource.Spec.Replicas)
		resource.Resource.Spec.Replicas = intPtr(initialReplicaSize)
		resource.Update()
		return true
	}
	log.Printf("Replicas after Scale %v\n", *resource.Resource.Spec.Replicas)
	log.Printf("Nodes after Scale %v\n", newNodeAmount)

	//Scale down the workload
	resource.Resource.Spec.Replicas = intPtr(initialReplicaSize)
	resource.Update()

	return true
}

func countNodes() (*v1.NodeList, int) {
	allNodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return nil, 0
	}
	return allNodes, len(allNodes.Items)
}
