package test_cases

import (
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	promscheme "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/scheme"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"strings"
)

var (
	ClusterName    string
	clientset      *kubernetes.Clientset
	promClientset  *promclient.Clientset
	istioClientset *istioclient.Clientset
)

func ConnectToKubernetes(kubeconfig string) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientsets
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	promClientset, err = promclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	istioClientset, err = istioclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	prepareClient()
}

func prepareClient() {
	promscheme.AddToScheme(scheme.Scheme)
	istioscheme.AddToScheme(scheme.Scheme)
}

//getHandler simply handles error when running a client.Get function
func getHandler(resource, name string, err error) bool {
	if err != nil && strings.Contains(err.Error(), strings.ToLower(field.ErrorTypeNotFound.String())) {
		log.Printf("couldn't locate %s: %s - it may need creating.\n", resource, err.Error())
		return false
	}
	return true
}
