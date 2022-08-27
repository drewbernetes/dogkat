package resources

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
)

type ServiceResource struct {
	Client   v1Typed.ServiceInterface
	Resource *v1.Service
	Error    error
}

func (r *ServiceResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *ServiceResource) GetError() error {
	return r.Error
}

func (r *ServiceResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *ServiceResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *ServiceResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *ServiceResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.CoreV1().Services(namespace)
}

func (r *ServiceResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *ServiceResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *ServiceResource) Update() {
}
func (r *ServiceResource) Delete() {
}
