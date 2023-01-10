package validate

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewValidateIngressCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
