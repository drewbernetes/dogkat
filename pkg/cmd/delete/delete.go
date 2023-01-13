package delete

import (
	"context"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/spf13/cobra"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

type deleteOptions struct {
	client     *kubernetes.Clientset
	istio      *istioclient.Clientset
	prometheus *promclient.Clientset
}

func NewDeleteCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes test resources",
		Long:  `Deletes the resources deployed for testing and any associated namespaces as per the subcommand provided.`,
	}

	commands := []*cobra.Command{
		NewDeleteAllCmd(f),
		NewDeleteCoreCmd(f),
		NewDeleteGpuCmd(f),
		NewDeleteIngressCmd(f),
		NewDeleteIstioCmd(f),
	}

	cmd.AddCommand(commands...)

	return cmd
}

func deleteNamespaceOnSuccess(client *kubernetes.Clientset, namespace string) error {
	if namespace == "default" {
		return nil
	}

	_, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	log.Printf("deleting namespace: %s\n", namespace)

	err = client.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})

	if err != nil {
		return err
	}
	return nil
}
