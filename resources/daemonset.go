package resources

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	appsv1Typed "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"log"
	"strings"
)

type DaemonSetResource struct {
	Client   appsv1Typed.DaemonSetInterface
	Resource *appsv1.DaemonSet
	Error    error
}

func (r *DaemonSetResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *DaemonSetResource) GetError() error {
	return r.Error
}

func (r *DaemonSetResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *DaemonSetResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *DaemonSetResource) IsReady() bool {
	//TODO could count the nodes and compare to r.Resource.Status.NumberAvailable to ensure they're all deployed
	if r.Resource.Status.NumberUnavailable != 0 || r.Resource.Status.NumberMisscheduled != 0 {
		return false
	}
	return true
}

func (r *DaemonSetResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.AppsV1().DaemonSets(namespace)
}

func (r *DaemonSetResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *DaemonSetResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *DaemonSetResource) Update() {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := r.Client.Update(context.TODO(), r.Resource, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed for %s:%s: %v\n", r.Resource.Kind, r.Resource.Name, retryErr)
	}
}
func (r *DaemonSetResource) Delete() {
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Client.Delete(context.TODO(), "demo-DaemonSet", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
	}
}
