package delete

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewDeleteAllCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Delete all generates all resources except Istio resources.",
		Long: `Deletes the application that was deployed to test all the elements of a cluster.
from Core workloads, to Ingress, to monitoring.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
