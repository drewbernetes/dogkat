package validate

import (
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/spf13/cobra"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
)

type validateOptions struct {
	client     *kubernetes.Clientset
	istio      *istioclient.Clientset
	prometheus *promclient.Clientset
}

func NewValidateCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "",
		Long:  "",
	}

	commands := []*cobra.Command{
		NewValidateAllCmd(f),
		NewValidateCoreCmd(f),
		NewValidateGpuCmd(f),
		NewValidateIngressCmd(f),
		NewValidateIstioCmd(f),
	}

	cmd.AddCommand(commands...)

	return cmd
}
