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

package coreworkloads

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type PersistentVolumeClaim struct {
	Client   *kubernetes.Clientset
	Resource *v1.PersistentVolumeClaim
}

// Generate the base Secret.
func (p *PersistentVolumeClaim) Generate(data map[string]string) {
	p.Resource = &v1.PersistentVolumeClaim{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a Secret on the Kubernetes cluster.
func (p *PersistentVolumeClaim) Create() error {
	log.Printf("creating PersistentVolumeClaim:%s...\n", p.Resource.Name)
	r := p.Client.CoreV1().PersistentVolumeClaims(p.Resource.Namespace)
	res, err := r.Create(context.Background(), p.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	p.Resource = res
	log.Printf("PersistentVolumeClaim:%s created.\n", p.Resource.Name)
	return nil
}

// Validate validates a persistentVolumeClaim on the Kubernetes cluster.
func (p *PersistentVolumeClaim) Validate() error {
	var err error
	r := p.Client.CoreV1().PersistentVolumeClaims(p.Resource.Namespace)
	p.Resource, err = r.Get(context.Background(), p.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a PersistentVolumeClaim in the Kubernetes cluster.
func (p *PersistentVolumeClaim) Update() error {
	return nil
}

// Delete deletes a persistentVolumeClaim from the Kubernetes cluster.
func (p *PersistentVolumeClaim) Delete() error {
	name := p.Resource.Name
	log.Printf("deleting PersistentVolumeClaim:%s...\n", name)
	r := p.Client.CoreV1().PersistentVolumeClaims(p.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("PersistentVolumeClaim:%s deleted.\n", name)
	return nil
}

func (p *PersistentVolumeClaim) GetResourceName() string {
	return p.Resource.Name
}

func (p *PersistentVolumeClaim) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", p.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (p *PersistentVolumeClaim) IsReady() bool {
	if err := p.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if p.Resource.Status.Phase != "Bound" {
		return false
	}
	return true
}
