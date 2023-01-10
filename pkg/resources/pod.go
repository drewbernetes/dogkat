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

// Deprecated: No longer in use
type PodResource struct {
	Client   v1Typed.PodInterface
	Resource *v1.Pod
	Error    error
}

func (r *PodResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *PodResource) GetError() error {
	return r.Error
}

func (r *PodResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *PodResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *PodResource) IsReady() bool {
	r.Get()
	for _, pod := range r.Resource.Status.ContainerStatuses {
		if pod.Ready != true {
			return false
		}
	}
	return true
}

func (r *PodResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.CoreV1().Pods(namespace)
}

func (r *PodResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}

func (r *PodResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}

func (r *PodResource) Update() {

}

func (r *PodResource) Delete() {
}

func rawPod() {

}
