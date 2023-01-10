package delete

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewDeleteIngressCmd(f util.Factory) *cobra.Command {
	o := &deleteOptions{}

	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Deletes the workloads and resources for testing an ingress controller works.",
		Long:  `Removes all elements of the ingress testing suite.`,
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

			workloads.DeleteNginxWorkloadItems(o.client, namespace)
			workloads.DeleteSQLWorkloadItems(o.client, namespace)
		},
	}

	return cmd
}