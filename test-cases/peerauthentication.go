package test_cases

import (
	"context"
	"fmt"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	v1beta1Typed "istio.io/client-go/pkg/clientset/versioned/typed/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type PeerAuthenticationResource struct {
	Client   v1beta1Typed.PeerAuthenticationInterface
	Resource *v1beta1.PeerAuthentication
	Error    error
}

func (r *PeerAuthenticationResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *PeerAuthenticationResource) GetError() error {
	return r.Error
}

func (r *PeerAuthenticationResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *PeerAuthenticationResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *PeerAuthenticationResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *PeerAuthenticationResource) GetClient(namespace string) {
	r.Client = istioClientset.SecurityV1beta1().PeerAuthentications(namespace)
}

func (r *PeerAuthenticationResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *PeerAuthenticationResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *PeerAuthenticationResource) Update() {
}
func (r *PeerAuthenticationResource) Delete() {
}
