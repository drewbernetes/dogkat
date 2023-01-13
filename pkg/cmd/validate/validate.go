package validate

import (
	"github.com/docker/distribution/context"
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

var (
	storageClassFlag string
)

type validateOptions struct {
	client     *kubernetes.Clientset
	istio      *istioclient.Clientset
	prometheus *promclient.Clientset
}

func NewValidateCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Create and validates test resources",
		Long: `Creates a selection of resources as per the subcommand provided. Once all resources are deployed and confirmed as ready,
the required test suite will run against the resources to ensure everything is working as expected within a cluster.`,
	}

	commands := []*cobra.Command{
		newValidateAllCmd(f),
		newValidateCoreCmd(f),
		newValidateGpuCmd(f),
		newValidateIngressCmd(f),
		newValidateIstioCmd(f),
		newValidateMonitoringCmd(f),
	}

	cmd.PersistentFlags().StringVar(&storageClassFlag, "storage-class", "longhorn", "Used to define the name of the storage class to use for Persistent Volumes.")

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
