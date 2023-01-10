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

type Secret struct {
	Client   *kubernetes.Clientset
	Resource *v1.Secret
}

// Generate the base Secret.
func (s *Secret) Generate(data map[string]string) {
	s.Resource = &v1.Secret{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
	}
}

// Create creates a Secret on the Kubernetes cluster.
func (s *Secret) Create() error {
	log.Printf("creating Secret:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().Secrets(s.Resource.Namespace)
	_, err := r.Create(context.Background(), s.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Secret:%s created.\n", s.Resource.Name)
	return nil
}

// Validate validates a Secret on the Kubernetes cluster.
func (s *Secret) Validate() error {
	var err error
	log.Printf("confirming Secret:%s...\n", s.Resource.Name)
	r := s.Client.CoreV1().Secrets(s.Resource.Namespace)
	s.Resource, err = r.Get(context.Background(), s.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("Secret: %s exists\n", s.Resource.Name)
	return nil
}

// Delete deletes a Secret from the Kubernetes cluster.
func (s *Secret) Delete() error {
	name := s.Resource.Name
	log.Printf("deleting Secret:%s...\n", name)
	r := s.Client.CoreV1().Secrets(s.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Secret:%s deleted.\n", name)
	return nil
}

func (s *Secret) GetResourceName() string {
	return s.Resource.Name
}

func (s *Secret) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", s.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (s *Secret) IsReady() bool {
	if s.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}
