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

type SecretResource struct {
	Client   v1Typed.SecretInterface
	Resource *v1.Secret
	Error    error
}

func (r *SecretResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *SecretResource) GetError() error {
	return r.Error
}

func (r *SecretResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *SecretResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *SecretResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *SecretResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.CoreV1().Secrets(namespace)
}

func (r *SecretResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *SecretResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *SecretResource) Update() {
}
func (r *SecretResource) Delete() {
}
