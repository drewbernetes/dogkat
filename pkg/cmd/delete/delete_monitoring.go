package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewDeleteMonitoringCmd(f util.Factory) *cobra.Command {
	o := &deleteOptions{}

	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Deletes the Monitoring testing resources",
		Long:  `Deletes the Monitoring application testing suite.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			fmt.Println(namespace)
		},
	}

	return cmd
}
