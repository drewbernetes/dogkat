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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/helpers"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/testsuite"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/tracing"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/gpu"
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
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			tracer := tracing.Tracer{JobName: "e2e_workloads", PushURL: "http://prometheus-push-gateway.prometheus:9091"}
			tracer.NewTimer("deploy_gpu", "Times the deployment of a gpu resource")
			timer := tracer.Start()

			// Generate and create workloads
			pod := gpu.GenerateGPUPod(namespace.Name, numberOfGPUsFlag)
			pod.Client = o.client
			helpers.HandleCreateError(pod.Create())

			//Check the pod exists
			err = pod.Validate()
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("all resources are deployed, running tests...")

			coreResource := []coreworkloads.Resource{
				pod,
			}

			err = testsuite.CheckReadyForTesting(coreResource)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("** ALL RESOURCES ARE DEPLOYED AND READY **")

			timer.ObserveDuration()
			err = testsuite.TestGPU(pod)
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
