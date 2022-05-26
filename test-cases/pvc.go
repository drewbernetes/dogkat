package test_cases

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
)

type PersistentVolumeClaimResource struct {
	Client   v1Typed.PersistentVolumeClaimInterface
	Resource *v1.PersistentVolumeClaim
	Error    error
}

func (r *PersistentVolumeClaimResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *PersistentVolumeClaimResource) GetError() error {
	return r.Error
}

func (r *PersistentVolumeClaimResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *PersistentVolumeClaimResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *PersistentVolumeClaimResource) IsReady() bool {
	if r.Resource.Status.Phase != v1.ClaimBound {
		return false
	}
	return true
}

func (r *PersistentVolumeClaimResource) GetClient(namespace string) {
	r.Client = clientset.CoreV1().PersistentVolumeClaims(namespace)
}

func (r *PersistentVolumeClaimResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *PersistentVolumeClaimResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *PersistentVolumeClaimResource) Update() {
}
func (r *PersistentVolumeClaimResource) Delete() {
}
