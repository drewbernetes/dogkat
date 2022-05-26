package test_cases

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1beta12 "istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"strings"
	"time"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type Resource interface {
	GetObject() runtime.Object
	GetClient(namespace string)
	GetError() error
	GetResourceName() string
	GetResourceKind() string
	IsReady() bool
	Get()
	Create()
	Update()
	Delete()
}

//CoreWorkloadTests will run the basic tests. Deployments, Ingress, Cluster scaling, Cluster DNS validation
func CoreWorkloadTests(valuesFile, namespaceName string) {
	resourcesReady := []bool{}
	values := parseValues(valuesFile)
	if values == nil {
		log.Fatalln("no values provided or couldn't parse them")
	}
	actionCfg, rel, err := deployChart(namespaceName, values)
	if err != nil {
		log.Println(err.Error())
		return
	}

	allManifests := strings.Split(rel.Manifest, "---")
	res := []Resource{}

	for _, v := range allManifests {
		if v == "" {
			continue
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode([]byte(v), nil, nil)
		if err != nil {
			log.Printf("There was an error decoding: %#v, %s", v, err)
			continue
		}
		r := parseResourceKind(obj)
		if r == nil {
			continue
		}
		r.GetClient(namespaceName)

		statusResults := checkIfResourceIsReady(r, 0, 20)
		if !statusResults {
			log.Printf("%s:%s is not ready\n", r.GetResourceKind(), r.GetResourceName())
			continue
		}
		log.Printf("%s:%s is ready\n", r.GetResourceKind(), r.GetResourceName())
		resourcesReady = append(resourcesReady, true)
		res = append(res, r)
	}

	if len(resourcesReady) != len(res) {
		log.Println(res, resourcesReady)
		log.Fatalln("one of the resources was not created - cannot continue testing")
	}
	log.Println("**ALL RESOURCES ARE DEPLOYED**")

	scalingTested := false
	for _, v := range res {
		runCoreTests(v)
		if v.GetResourceKind() == "Deployment" {
			if !scalingTested {
				scalingTested = runScalingTest(v)
			} else {
				break
			}
		}
	}

	resp, err := uninstallChart(actionCfg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if resp != nil {
		log.Println(resp.Info)
	}

	log.Println("the helm release has been removed - please remove the namespace e2e-testing manually - (this will be automated in a future release)")
}

func checkIfResourceIsReady(v Resource, counter int, delaySeconds time.Duration) bool {
	delay := time.Second * delaySeconds
	if counter >= 10 {
		return false
	}
	v.Get()
	if !v.IsReady() {
		time.Sleep(delay)
		return checkIfResourceIsReady(v, counter+1, delaySeconds)
	}
	return true
}

//parseResourceKind will check the kind and return a valid K8S object so that it can be validated
func parseResourceKind(obj runtime.Object) (r Resource) {
	switch obj.GetObjectKind().GroupVersionKind().Kind {
	case "ConfigMap":
		r = &ConfigMapResource{
			Resource: obj.(*v1.ConfigMap),
		}
	case "Secret":
		r = &SecretResource{
			Resource: obj.(*v1.Secret),
		}
	case "Deployment":
		r = &DeploymentResource{
			Resource: obj.(*appsv1.Deployment),
		}
	case "ServiceAccount":
		r = &ServiceAccountResource{
			Resource: obj.(*v1.ServiceAccount),
		}
	case "Service":
		r = &ServiceResource{
			Resource: obj.(*v1.Service),
		}
	case "Job":
		r = &JobResource{
			Resource: obj.(*v1batch.Job),
		}
	case "PersistentVolumeClaim":
		r = &PersistentVolumeClaimResource{
			Resource: obj.(*v1.PersistentVolumeClaim),
		}
	case "Ingress":
		r = &IngressResource{
			Resource: obj.(*networkingv1.Ingress),
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
		log.Printf("couldn't ascertain the Kind of the resource, skipping %+v\n", r)
		return nil
	}

	return r
}
