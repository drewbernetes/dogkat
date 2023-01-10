package k8s

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/resources"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	promscheme "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/scheme"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

// ConnectToKubernetes connects to kubernetes and setups any additional clientsets.
//
// Deprecated: This is no longer used in the codebase
func ConnectToKubernetes(kubeconfig *string) *resources.ClientSets {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientsets := &resources.ClientSets{}
	// create the clientsets
	clientsets.K8S, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientsets.Prometheus, err = promclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	clientsets.Istio, err = istioclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	addSchemesToClientset()

	return clientsets
}

// addSchemesToClientset adds any additional schemes so that the clientset can can interact with custom types.
//
// Deprecated: This is no longer used in the codebase
func addSchemesToClientset() {
	err := promscheme.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalln(err)
	}
	err = istioscheme.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalln(err)
	}
}
