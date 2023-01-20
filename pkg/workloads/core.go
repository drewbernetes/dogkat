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

package workloads

import (
	"github.com/docker/distribution/context"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/testsuite"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/sql"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/web"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// CreateNamespaceIfNotExists create the provided namespace if it doesn't exist.
// It will return the namespace once created or if it is located.
func CreateNamespaceIfNotExists(client *kubernetes.Clientset, ns string) *v1.Namespace {
	log.Println("checking for namespace, will create if doesn't exist")
	namespace := "default"
	if len(ns) != 0 {
		namespace = ns
	}

	n, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err == nil {
		return n
	}

	log.Println(err)
	log.Printf("creating namespace: %s\n", namespace)

	n, err = client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		log.Fatalln(err)
	}
	return n
}

// DeployBaseWorkloads deploys the basic applications required for the majority of testing.
func DeployBaseWorkloads(client *kubernetes.Clientset, namespace, storageClass, cpuRequest, memoryRequest string) (*web.NginxWorkloads, *sql.PostgresWorkloads) {
	var err error

	//TODO: Check storage class is available or that a CNI is available
	//TODO: Check that cluster-autoscaler is available

	// Generate and create workloads
	nginxWorkload := web.CreateNginxWorkloadItems(client, namespace, cpuRequest, memoryRequest)
	sqlWorkload := sql.CreateSQLWorkloadItems(client, namespace, storageClass)

	//confirm the workloads were deployed
	err = nginxWorkload.ValidateNginxWorkloadItems()
	if err != nil {
		log.Fatalln(err)
	}
	err = sqlWorkload.ValidateSQLWorkloadItems()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("all resources are deployed, running tests...")

	coreResource := []coreworkloads.Resource{
		nginxWorkload.Configuration,
		nginxWorkload.WebPages,
		nginxWorkload.Workload,
		nginxWorkload.ServiceAccount,
		nginxWorkload.Service,
		nginxWorkload.PodDisruptionBudget,
		sqlWorkload.Workload,
		sqlWorkload.ServiceAccount,
		sqlWorkload.InitConf,
		sqlWorkload.Secret,
		sqlWorkload.Service,
		sqlWorkload.PodDisruptionBudget,
	}

	err = testsuite.CheckReadyForTesting(coreResource)
	if err != nil {
		log.Fatalln(err)
	}

	//Check volumes after STS confirmation to ensure volumes have been created.
	volumeResource := parseVolumesFromStatefulSet(client, sqlWorkload, namespace)
	err = testsuite.CheckReadyForTesting(volumeResource)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("** ALL RESOURCES ARE DEPLOYED AND READY **")

	return nginxWorkload, sqlWorkload
}

// parseVolumesFromStatefulSet
func parseVolumesFromStatefulSet(client *kubernetes.Clientset, sqlWorkload *sql.PostgresWorkloads, namespace string) []coreworkloads.Resource {
	volumeResource := []coreworkloads.Resource{}
	volNames := []string{strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "0"}, "-"),
		strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "1"}, "-"),
		strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "2"}, "-"),
	}

	for _, name := range volNames {
		pvc := &coreworkloads.PersistentVolumeClaim{
			Client:   client,
			Resource: &v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}},
		}
		err := pvc.Validate()
		if err != nil {
			log.Println(err)
		}
		volumeResource = append(volumeResource, pvc)

		pvName := pvc.Resource.Spec.VolumeName
		pv := &coreworkloads.PersistentVolume{
			Client:   client,
			Resource: &v1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: pvName}},
		}
		err = pv.Validate()
		if err != nil {
			log.Println(err)
		}
		volumeResource = append(volumeResource, pv)
	}

	return volumeResource
}
