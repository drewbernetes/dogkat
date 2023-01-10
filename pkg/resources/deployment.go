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

// Deprecated: No longer in use
type DeploymentResource struct {
	Client   appsv1Typed.DeploymentInterface
	Resource *appsv1.Deployment
	Error    error
}

func (r *DeploymentResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *DeploymentResource) GetError() error {
	return r.Error
}

func (r *DeploymentResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *DeploymentResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *DeploymentResource) IsReady() bool {
	if r.Resource.Status.UnavailableReplicas != 0 {
		return false
	}
	return true
}

func (r *DeploymentResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.AppsV1().Deployments(namespace)
}

func (r *DeploymentResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *DeploymentResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *DeploymentResource) Update() {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := r.Client.Update(context.TODO(), r.Resource, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed for %s:%s: %v\n", r.Resource.Kind, r.Resource.Name, retryErr)
	}
}
func (r *DeploymentResource) Delete() {
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Client.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
	}
}
