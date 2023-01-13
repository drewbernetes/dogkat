package coreworkloads

import (
	"context"
	"fmt"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type Ingress struct {
	Client   *kubernetes.Clientset
	Resource *v1.Ingress
}

// Generate a base Ingress.
func (i *Ingress) Generate(data map[string]string) {
	pathType := v1.PathTypePrefix
	className := data["className"]

	i.Resource = &v1.Ingress{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["name"]),
		Spec: v1.IngressSpec{
			IngressClassName: &className,
			Rules: []v1.IngressRule{
				{
					Host: data["host"],
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: data["name"],
											Port: v1.ServiceBackendPort{
												Name:   "http",
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create creates an Ingress on the Kubernetes cluster.
func (i *Ingress) Create() error {
	log.Printf("creating Ingress:%s...\n", i.Resource.Name)
	r := i.Client.NetworkingV1().Ingresses(i.Resource.Namespace)
	res, err := r.Create(context.Background(), i.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	i.Resource = res
	log.Printf("Ingress:%s created.\n", i.Resource.Name)
	return nil
}

// Validate validates an Ingress on the Kubernetes cluster.
func (i *Ingress) Validate() error {
	var err error
	r := i.Client.NetworkingV1().Ingresses(i.Resource.Namespace)
	i.Resource, err = r.Get(context.Background(), i.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies an Ingress in the Kubernetes cluster.
func (i *Ingress) Update() error {
	return nil
}

// Delete deletes an Ingress from the Kubernetes cluster.
func (i *Ingress) Delete() error {
	name := i.Resource.Name
	log.Printf("deleting Ingress:%s...\n", name)
	r := i.Client.NetworkingV1().Ingresses(i.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Ingress:%s deleted.\n", name)
	return nil
}

func (i *Ingress) GetResourceName() string {
	return i.Resource.Name
}

func (i *Ingress) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", i.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (i *Ingress) IsReady() bool {
	if err := i.Validate(); err != nil {
		log.Println(err)
		return false
	}
	for _, v := range i.Resource.Status.LoadBalancer.Ingress {
		if v.Hostname == "" && v.IP == "" {
			return false
		}
	}
	return true
}
