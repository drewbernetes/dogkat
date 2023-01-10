package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewValidateCoreCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "core",
		Short: "Validates all core workloads.",
		Long: `Validates the by running test suites against to confirm the following are functional:
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
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			nginxWorkload, err := workloads.ValidateNginxWorkloadItems(o.client, namespace)
			if err != nil {
				log.Fatalln(err)
			}
			sqlWorkload, err := workloads.ValidateSQLWorkloadItems(o.client, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("all resources are deployed, running tests...")

			err = testsuite.TestReady(nginxWorkload, sqlWorkload)
			if err != nil {
				log.Fatalln(err)
			}

			err = testsuite.TestCore(nginxWorkload, sqlWorkload)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
