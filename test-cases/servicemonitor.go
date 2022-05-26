package test_cases

import (
	"context"
	"fmt"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1Typed "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

type ServiceMonitorResource struct {
	Client   v1Typed.ServiceMonitorInterface
	Resource *v1.ServiceMonitor
	Error    error
}

func (r *ServiceMonitorResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *ServiceMonitorResource) GetError() error {
	return r.Error
}

func (r *ServiceMonitorResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *ServiceMonitorResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *ServiceMonitorResource) IsReady() bool {
	if r.Resource.CreationTimestamp.IsZero() {
		return false
	}
	return true
}

func (r *ServiceMonitorResource) GetClient(namespace string) {
	r.Client = promClientset.MonitoringV1().ServiceMonitors(namespace)
}

func (r *ServiceMonitorResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *ServiceMonitorResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *ServiceMonitorResource) Update() {
}
func (r *ServiceMonitorResource) Delete() {
}
