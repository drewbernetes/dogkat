package resources

import (
	"context"
	"fmt"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1beta1Typed "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type GatewayResource struct {
	Client   v1beta1Typed.GatewayInterface
	Resource *v1beta1.Gateway
	Error    error
}

func (r *GatewayResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *GatewayResource) GetError() error {
	return r.Error
}

func (r *GatewayResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *GatewayResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}
func (r *GatewayResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *GatewayResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.Istio.NetworkingV1beta1().Gateways(namespace)
}

func (r *GatewayResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *GatewayResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *GatewayResource) Update() {
}
func (r *GatewayResource) Delete() {
}
