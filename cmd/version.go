package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of EKS E2E-Tester",
	Long:  `All software has versions. This is E2E-Tester's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("E2E-Tester v0.0.1 -- HEAD")
	},
}
