package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/prometheus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newValidateMonitoringCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Creates and validates the core resources with additional monitoring",
		Long: `Creates an application with additional grafana dashboards and service monitors. 
This requires Prometheus to be installed. Once deployed the monitoring test suite will run to confirm 
resources are deployed and working as expected.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			_, _ = workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)

			//TODO Check for service monitor resource
			sm := prometheus.GenerateServiceMonitorResource(namespace.Name)
			prometheus.CreateServiceMonitor(o.prometheus, sm)
		},
	}
	addCoreFlags(cmd)
	return cmd
}
