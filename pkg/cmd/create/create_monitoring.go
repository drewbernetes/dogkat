package create

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/prometheus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newCreateMonitoringCmd(f util.Factory) *cobra.Command {
	o := &createOptions{}

	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Creates the core resources with additional monitoring.",
		Long: `Creates an application with additional grafana dashboards and service monitors. 
This requires Prometheus to be installed.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			log.Println("checking for namespace, will create if doesn't exist")
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			createNamespaceIfNotExists(o.client, namespace)

			// Generate and create workloads
			workloads.CreateNginxWorkloadItems(o.client, namespace)
			workloads.CreateSQLWorkloadItems(o.client, namespace, storageClass)

			sm := prometheus.GenerateServiceMonitorResource(namespace)
			prometheus.CreateServiceMonitor(o.prometheus, sm)
		},
	}
	cmd.Flags().StringVar(&storageClass, "storage-class", "longhorn", "Used to define the name of the storage class to use for Persistent Volumes.")

	return cmd
}
