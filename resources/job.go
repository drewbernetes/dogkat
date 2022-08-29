package resources

import (
	"context"
	"fmt"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1Typed "k8s.io/client-go/kubernetes/typed/batch/v1"
	"strings"
)

type JobResource struct {
	Client   v1Typed.JobInterface
	Resource *v1.Job
	Error    error
}

func (r *JobResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *JobResource) GetError() error {
	return r.Error
}

func (r *JobResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *JobResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *JobResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *JobResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.BatchV1().Jobs(namespace)
}

func (r *JobResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *JobResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *JobResource) Update() {
}
func (r *JobResource) Delete() {
}
