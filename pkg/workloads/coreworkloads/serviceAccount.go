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

type ServiceAccount struct {
	Client   *kubernetes.Clientset
	Resource *v1.ServiceAccount
}

// Generate a ServiceAccount resource that will be used for testing.
func (s *ServiceAccount) Generate(data map[string]string) {
	s.Resource = &v1.ServiceAccount{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a ServiceAccount on the Kubernetes cluster.
func (s *ServiceAccount) Create() error {
	log.Printf("creating ServiceAccount:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().ServiceAccounts(s.Resource.Namespace)
	_, err := r.Create(context.Background(), s.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("ServiceAccount:%s created.\n", s.Resource.Name)
	return nil
}

// Validate validates a ServiceAccount on the Kubernetes cluster.
func (s *ServiceAccount) Validate() error {
	var err error
	log.Printf("confirming ServiceAccount:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().ServiceAccounts(s.Resource.Namespace)
	s.Resource, err = r.Get(context.Background(), s.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("ServiceAccount: %s exists\n", s.Resource.Name)
	return nil
}

// Delete deletes a ServiceAccount from the Kubernetes cluster.
func (s *ServiceAccount) Delete() error {
	name := s.Resource.Name
	log.Printf("deleting ServiceAccount:%s...\n", name)
	r := s.Client.CoreV1().ServiceAccounts(s.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("ServiceAccount:%s deleted.\n", name)
	return nil
}

func (s *ServiceAccount) GetResourceName() string {
	return s.Resource.Name
}

func (s *ServiceAccount) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", s.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (s *ServiceAccount) IsReady() bool {
	if s.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}
