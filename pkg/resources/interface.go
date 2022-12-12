package resources

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1beta12 "istio.io/client-go/pkg/apis/security/v1beta1"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	v1batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	v1_policy "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// ClientSets contians all available clientset types.
type ClientSets struct {
	K8S        *kubernetes.Clientset
	Istio      *istioclient.Clientset
	Prometheus *promclient.Clientset
}

// ResourceReady is used to determine if a resource can be marked as ready during the threaded checks.
type ResourceReady struct {
	Ready    bool
	Resource ApiResource
}

// ApiResource is used to allow collation of resources into a single one to allow for similar tests to be run across the board without having to do repetitive code.

type ApiResource interface {
	GetObject() runtime.Object
	GetClient(namespace string, clientset *ClientSets)
	GetError() error
	GetResourceName() string
	GetResourceKind() string
	IsReady() bool
	Get()
	Create()
	Update()
	Delete()
}

// ParseResourceKind will check the resource kind and return a valid K8S object so that it can be validated
// To add more checks, the resource must be added here so that it can be parsed.
func ParseResourceKind(obj runtime.Object) (r ApiResource) {
	if obj.GetObjectKind().GroupVersionKind().Kind == "" {
		return nil
	}
	kind := obj.GetObjectKind().GroupVersionKind().Kind

	switch kind {
	case "Pod":
		r = &PodResource{
			Resource: obj.(*v1.Pod),
		}
	case "ConfigMap":
		r = &ConfigMapResource{
			Resource: obj.(*v1.ConfigMap),
		}
	case "Secret":
		r = &SecretResource{
			Resource: obj.(*v1.Secret),
		}
	case "ServiceAccount":
		r = &ServiceAccountResource{
			Resource: obj.(*v1.ServiceAccount),
		}
	case "Service":
		r = &ServiceResource{
			Resource: obj.(*v1.Service),
		}
	case "PersistentVolumeClaim":
		r = &PersistentVolumeClaimResource{
			Resource: obj.(*v1.PersistentVolumeClaim),
		}
	case "Deployment":
		r = &DeploymentResource{
			Resource: obj.(*appsv1.Deployment),
		}
	case "DaemonSet":
		r = &DaemonSetResource{
			Resource: obj.(*appsv1.DaemonSet),
		}
	case "StatefulSet":
		r = &StatefulSetResource{
			Resource: obj.(*appsv1.StatefulSet),
		}
	case "Job":
		r = &JobResource{
			Resource: obj.(*v1batch.Job),
		}
	case "Ingress":
		r = &IngressResource{
			Resource: obj.(*networkingv1.Ingress),
		}
	case "PodDisruptionBudget":
		r = &PDBResource{
			Resource: obj.(*v1_policy.PodDisruptionBudget),
		}
	case "ServiceMonitor":
		r = &ServiceMonitorResource{
			Resource: obj.(*monitoringv1.ServiceMonitor),
		}
	case "Gateway":
		r = &GatewayResource{
			Resource: obj.(*v1beta1.Gateway),
		}
	case "VirtualService":
		r = &VirtualServiceResource{
			Resource: obj.(*v1beta1.VirtualService),
		}
	case "PeerAuthentication":
		r = &PeerAuthenticationResource{
			Resource: obj.(*v1beta12.PeerAuthentication),
		}
	case "DestinationRule":
		r = &DestinationRuleResource{
			Resource: obj.(*v1beta1.DestinationRule),
		}
	default:
		log.Printf("Couldn't ascertain the Kind of the resource, skipping %s\n", kind)
		return nil
	}

	return r
}

// getHandler simply handles error when running a client.Get function
func getHandler(resource, name string, err error) bool {
	if err != nil && strings.Contains(err.Error(), strings.ToLower(field.ErrorTypeNotFound.String())) {
		log.Printf("%s - The Helm release may still exist with no resources deployed - please remove the Helm release and then re-run.\n", err.Error())
		return false
	}
	return true
}
