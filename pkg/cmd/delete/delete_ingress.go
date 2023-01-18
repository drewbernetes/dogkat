package delete

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/sql"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewDeleteIngressCmd(f util.Factory) *cobra.Command {
	o := &deleteOptions{}

	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Deletes the Ingress testing resources.",
		Long:  `Deletes all elements of the ingress testing suite.`,
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

			web.DeleteNginxWorkloadItems(o.client, namespace)
			web.DeleteIngressWorkloadItems(o.client, namespace)
			sql.DeleteSQLWorkloadItems(o.client, namespace)

		},
	}

	return cmd
}
