package create

import (
	"context"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	promscheme "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/scheme"
	"github.com/spf13/cobra"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

type createOptions struct {
	client     *kubernetes.Clientset
	istio      *istioclient.Clientset
	prometheus *promclient.Clientset
}

func NewCreateCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the resources required for the testing.",
		Long: `Create enables the creation of specific resources for the test cases to be deployed.
This should be run first to ensure that the resources are in place to allow testing to continue and succeed correctly.`,
	}

	commands := []*cobra.Command{
		NewCreateAllCmd(f),
		newCreateCoreCmd(f),
		newCreateIngressCmd(f),
		newCreateGpuCmd(f),
		newCreateIstioCmd(f),
		newCreateMonitoringCmd(f),
	}

	cmd.AddCommand(commands...)

	return cmd
}

func createNamespaceIfNotExists(client *kubernetes.Clientset, namespace string) {
	if namespace == "default" {
		return
	}
	_, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err == nil {
		return
	}

	log.Println(err)
	log.Printf("creating namespace: %s\n", namespace)

	_, err = client.CoreV1().Namespaces().Create(context.Background(), &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		log.Fatalln(err)
	}
}

// addPrometheusToScheme adds the Prometheus scheme to the scheme so that the clientset can use it.
func addPrometheusToScheme() {
	err := promscheme.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalln(err)
	}
}

// addIstioScheme adds the Istio scheme to the scheme so that the clientset can use it.
func addIstioToScheme() {
	err := istioscheme.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalln(err)
	}
}
