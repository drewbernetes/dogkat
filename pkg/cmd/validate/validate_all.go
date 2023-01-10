package validate

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewValidateAllCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Validates all resources.",
		Long:  "Validates all resources from core, gpu and ingress. Istio test suites will not be affected by this.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
