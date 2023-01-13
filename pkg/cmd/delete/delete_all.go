package delete

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewDeleteAllCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Delete all resources except Istio",
		Long: `Deletes the resources that were deployed to test all the elements of a cluster.
Including Core, GPU, Ingress and Monitoring.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
