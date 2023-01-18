package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewDeleteIstioCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Deletes the Istio testing resources",
		Long:  `Deletes the Istio application testing suite.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Istio will be supported soon™")
		},
	}

	return cmd
}
