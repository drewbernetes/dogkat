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

type Service struct {
	Client   *kubernetes.Clientset
	Resource *v1.Service
}

// Generate a Service resource that will be used for testing.
func (s *Service) Generate(data map[string]string) {
	s.Resource = &v1.Service{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a Service on the Kubernetes cluster.
func (s *Service) Create() error {
	log.Printf("creating Service:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().Services(s.Resource.Namespace)
	res, err := r.Create(context.Background(), s.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	s.Resource = res
	log.Printf("Service:%s created.\n", s.Resource.Name)
	return nil
}

// Validate validates a Service on the Kubernetes cluster.
func (s *Service) Validate() error {
	var err error
	r := s.Client.CoreV1().Services(s.Resource.Namespace)
	s.Resource, err = r.Get(context.Background(), s.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a Service in the Kubernetes cluster.
func (s *Service) Update() error {
	return nil
}

// Delete deletes a Service from the Kubernetes cluster.
func (s *Service) Delete() error {
	name := s.Resource.Name
	log.Printf("deleting Service:%s...\n", name)
	r := s.Client.CoreV1().Services(s.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Service:%s deleted.\n", name)
	return nil
}

func (s *Service) GetResourceName() string {
	return s.Resource.Name
}

func (s *Service) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", s.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (s *Service) IsReady() bool {
	if err := s.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if s.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}
