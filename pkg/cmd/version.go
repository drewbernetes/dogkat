package cmd

import (
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns a version command that prints out application
// and versioning information.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print this command's version",
		Long:  "Print this command's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(constants.Version)
		},
	}
}
