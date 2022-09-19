package resources

import (
	"context"
	"fmt"
	v1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1_typed "k8s.io/client-go/kubernetes/typed/policy/v1"
	"strings"
)

type PDBResource struct {
	Client   v1_typed.PodDisruptionBudgetInterface
	Resource *v1.PodDisruptionBudget
	Error    error
}

func (r *PDBResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.ApiResource)
	return r.Resource
}

func (r *PDBResource) GetError() error {
	return r.Error
}

func (r *PDBResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *PDBResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *PDBResource) IsReady() bool {
	if r.Resource.Status.CurrentHealthy < r.Resource.Status.DesiredHealthy || r.Resource.Status.DisruptionsAllowed == 0 {
		return false
	}
	return true
}

func (r *PDBResource) GetClient(namespace string, clientset *ClientSets) {
	r.Client = clientset.K8S.PolicyV1().PodDisruptionBudgets(namespace)
}

func (r *PDBResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *PDBResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *PDBResource) Update() {
}
func (r *PDBResource) Delete() {
}
