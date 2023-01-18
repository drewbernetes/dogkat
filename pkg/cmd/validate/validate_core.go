package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newValidateCoreCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "core",
		Short: "Create and test core resources",
		Long: `Creates a selection of standard resources. Once all resources are deployed and confirmed as ready,
the test suite will then run against them to ensure everything is working as it should be within a standard cluster.

The following is tested:
* Deployments, StatefulSets & Services to validate Cluster DNS
* Pod Disruption Budgets to confirm HA
* ConfigMaps & Secrets toc onverm ENVs and Volume Mounting of ConfigMaps
* PV creation & Mounting via PVC in StatefulSet
`,
		Run: func(cmd *cobra.Command, args []string) {

			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			nginxWorkload, _ := workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)

			err = testsuite.ScaleUpStandardNodes(nginxWorkload.Workload)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	addCoreFlags(cmd)
	return cmd
}
