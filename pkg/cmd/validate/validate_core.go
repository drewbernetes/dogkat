/*
Copyright 2022 EscherCloud.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validate

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/testsuite"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/tracing"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads"
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

			fullTracer := tracing.Duration{JobName: "e2e_workloads", PushURL: pushGatewayURLFlag}
			fullTracer.SetupMetricsGatherer("full_e2e_test_all_duration_seconds", "Times the entire e2e workload testing for a full run")
			fullTracer.Start()
			//TODO: This repeats - let's clean it up!
			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String(), pushGatewayURLFlag)

			nginxWorkload, _ := workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag, pushGatewayURLFlag)

			err = testsuite.ScaleUpStandardNodes(nginxWorkload.Workload, pushGatewayURLFlag)
			if err != nil {
				log.Fatalln(err)
			}

			fullTracer.CompleteGathering()
		},
	}
	addCoreFlags(cmd)
	return cmd
}
