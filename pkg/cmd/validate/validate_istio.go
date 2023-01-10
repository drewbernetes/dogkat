package validate

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewValidateIstioCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "istio",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
