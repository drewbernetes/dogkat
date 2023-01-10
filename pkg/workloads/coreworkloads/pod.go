package coreworkloads

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type Pod struct {
	Client   *kubernetes.Clientset
	Resource *v1.Pod
}

// Generate a base Pod definition.
func (p *Pod) Generate(data map[string]string) *v1.Pod {
	p.Resource = &v1.Pod{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["name"]),
		Spec: v1.PodSpec{
			RestartPolicy: "OnFailure",
		},
	}
	return nil
}

// Create creates a Pod on the Kubernetes cluster.
func (p *Pod) Create() error {
	log.Printf("creating Pod:%s...\n", p.Resource.Name)
	r := p.Client.CoreV1().Pods(p.Resource.Namespace)
	_, err := r.Create(context.Background(), p.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Pod:%s created.\n", p.Resource.Name)
	return nil
}

// Validate validates a Pod on the Kubernetes cluster.
func (p *Pod) Validate() error {
	var err error
	log.Printf("confirming Pod:%s...\n", p.Resource.Name)
	r := p.Client.CoreV1().Pods(p.Resource.Namespace)
	p.Resource, err = r.Get(context.Background(), p.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("Pod: %s exists\n", p.Resource.Name)
	return nil
}

// Delete deletes a Pod from the Kubernetes cluster.
func (p *Pod) Delete() error {
	name := p.Resource.Name
	log.Printf("deleting Pod:%s...\n", name)
	r := p.Client.CoreV1().Pods(p.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Pod:%s deleted.\n", name)
	return nil
}

func (p *Pod) GetResourceName() string {
	return p.Resource.Name
}

func (p *Pod) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", p.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (p *Pod) IsReady() bool {
	for _, pod := range p.Resource.Status.ContainerStatuses {
		if pod.Ready != true {
			return false
		}
	}
	return true
}
