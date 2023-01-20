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

type ConfigMap struct {
	Client   *kubernetes.Clientset
	Resource *v1.ConfigMap
}

// Generate the base ConfigMap.
func (c *ConfigMap) Generate(data map[string]string) {
	c.Resource = &v1.ConfigMap{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a ConfigMap on the Kubernetes cluster.
func (c *ConfigMap) Create() error {
	log.Printf("creating ConfigMap:%s...\n", c.Resource.Name)
	r := c.Client.CoreV1().ConfigMaps(c.Resource.Namespace)
	res, err := r.Create(context.Background(), c.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	c.Resource = res
	log.Printf("ConfigMap:%s created.\n", c.Resource.Name)
	return nil
}

// Validate validates a ConfigMap exists on the Kubernetes cluster.
func (c *ConfigMap) Validate() error {
	var err error
	r := c.Client.CoreV1().ConfigMaps(c.Resource.Namespace)
	c.Resource, err = r.Get(context.Background(), c.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a ConfigMap in the Kubernetes cluster.
func (c *ConfigMap) Update() error {
	return nil
}

// Delete deletes a ConfigMap from the Kubernetes cluster.
func (c *ConfigMap) Delete() error {
	name := c.Resource.Name

	log.Printf("deleting ConfigMap:%s...\n", name)
	r := c.Client.CoreV1().ConfigMaps(c.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("ConfigMap:%s deleted.\n", name)
	return nil
}

func (c *ConfigMap) GetResourceName() string {
	return c.Resource.Name
}

func (c *ConfigMap) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", c.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (c *ConfigMap) IsReady() bool {
	if err := c.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if c.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}
