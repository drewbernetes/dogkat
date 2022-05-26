package test_cases

import (
	"context"
	"fmt"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1beta1Typed "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type VirtualServiceResource struct {
	Client   v1beta1Typed.VirtualServiceInterface
	Resource *v1beta1.VirtualService
	Error    error
}

func (r *VirtualServiceResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *VirtualServiceResource) GetError() error {
	return r.Error
}

func (r *VirtualServiceResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *VirtualServiceResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *VirtualServiceResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *VirtualServiceResource) GetClient(namespace string) {
	r.Client = istioClientset.NetworkingV1beta1().VirtualServices(namespace)
}

func (r *VirtualServiceResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *VirtualServiceResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *VirtualServiceResource) Update() {
}
func (r *VirtualServiceResource) Delete() {
}
