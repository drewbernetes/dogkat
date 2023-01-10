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
	_, err := r.Create(context.Background(), s.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Service:%s created.\n", s.Resource.Name)
	return nil
}

// Validate validates a Service on the Kubernetes cluster.
func (s *Service) Validate() error {
	var err error
	log.Printf("confirming Service:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().Services(s.Resource.Namespace)
	s.Resource, err = r.Get(context.Background(), s.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("Service: %s exists\n", s.Resource.Name)
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
	if s.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}
