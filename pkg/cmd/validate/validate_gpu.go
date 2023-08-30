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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	numberOfGPUsFlag string
)

func newValidateGpuCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Create and test a GPU workload",
		Long: `Creates the resources for GPU testing and then 
runs the validation and test suite against the GPU workload to ensure it is working as expected.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			//TODO: This repeats - let's clean it up!
			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String(), pushGatewayURLFlag)

			pod, err := workloads.DeployGPUWorkloads(o.client, namespace.Name, numberOfGPUsFlag, pushGatewayURLFlag)
			if err != nil {
				log.Fatalln(err)
			}

			err = testsuite.TestGPU(pod, pushGatewayURLFlag)
			if err != nil {
				log.Fatalln(err)
			}

			//err = testsuite.ScaleUpGPUNodes(pod)
			//if err != nil {
			//	log.Fatalln(err)
			//}
		},
	}
	cmd.Flags().StringVar(&numberOfGPUsFlag, "number-of-gpus", "1", "Sets the number of GPUS in resources.limits")

	return cmd
}
