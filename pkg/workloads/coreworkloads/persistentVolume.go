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

type PersistentVolume struct {
	Client   *kubernetes.Clientset
	Resource *v1.PersistentVolume
}

// Generate the base Secret.
func (p *PersistentVolume) Generate(data map[string]string) {
	p.Resource = &v1.PersistentVolume{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a Secret on the Kubernetes cluster.
func (p *PersistentVolume) Create() error {
	log.Printf("creating PersistentVolume:%s...\n", p.Resource.Name)
	r := p.Client.CoreV1().PersistentVolumes()
	res, err := r.Create(context.Background(), p.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	p.Resource = res
	log.Printf("PersistentVolume:%s created.\n", p.Resource.Name)
	return nil
}

// Validate validates a PersistentVolume on the Kubernetes cluster.
func (p *PersistentVolume) Validate() error {
	var err error
	r := p.Client.CoreV1().PersistentVolumes()
	p.Resource, err = r.Get(context.Background(), p.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a PersistentVolume in the Kubernetes cluster.
func (p *PersistentVolume) Update() error {
	return nil
}

// Delete deletes a PersistentVolume from the Kubernetes cluster.
func (p *PersistentVolume) Delete() error {
	name := p.Resource.Name
	log.Printf("deleting PersistentVolume:%s...\n", name)
	r := p.Client.CoreV1().PersistentVolumes()
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("PersistentVolume:%s deleted.\n", name)
	return nil
}

func (p *PersistentVolume) GetResourceName() string {
	return p.Resource.Name
}

func (p *PersistentVolume) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", p.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (p *PersistentVolume) IsReady() bool {
	if err := p.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if p.Resource.Status.Phase != "Bound" {
		return false
	}
	return true
}
