package resources

import (
	"context"
	"fmt"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1Typed "k8s.io/client-go/kubernetes/typed/networking/v1"
	"strings"
)

type IngressResource struct {
	Client   v1Typed.IngressInterface
	Resource *v1.Ingress
	Error    error
}

func (r *IngressResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *IngressResource) GetError() error {
	return r.Error
}

func (r *IngressResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *IngressResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *IngressResource) IsReady() bool {
	for _, v := range r.Resource.Status.LoadBalancer.Ingress {
		if v.Hostname == "" && v.IP == "" {
			return false
		}
	}
	return true
}

func (r *IngressResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.NetworkingV1().Ingresses(namespace)
}

func (r *IngressResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *IngressResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *IngressResource) Update() {
}
func (r *IngressResource) Delete() {
}
