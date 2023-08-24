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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/sql"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newValidateAllCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Create and validates all resources",
		Long: `Creates and validates all resources from core, gpu, ingress and monitoring.
Istio test suites will not be affected by this.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}
			//TODO: This repeats - let's clean it up!
			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			// Generate and create workloads
			nginxWorkload, _ := workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)
			err = testsuite.ScaleUpStandardNodes(nginxWorkload.Workload)
			if err != nil {
				log.Fatalln(err)
			}

			web.CreateIngressResource(o.client, namespace.Name, annotationsFlag, hostFlag, ingressClassFlag, enableTLSFlag)
			err = testsuite.TestIngress(hostFlag)
			if err != nil {
				log.Fatalln(err)
			}

			//TODO Check for service monitor resource and implement creation
			//sm := prometheus.GenerateServiceMonitorResource(namespace.Name)
			//prometheus.CreateServiceMonitor(o.prometheus, sm)

			web.DeleteNginxWorkloadItems(o.client, namespace.Name)
			web.DeleteIngressWorkloadItems(o.client, namespace.Name)
			sql.DeleteSQLWorkloadItems(o.client, namespace.Name)

			// Generate and create GPU workloads

			tracer := tracing.Tracer{JobName: "e2e_workloads", PushURL: "http://prometheus-push-gateway.prometheus:9091"}
			tracer.NewTimer("deploy_gpu", "Times the deployment of a gpu resource")
			timer := tracer.Start()

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
			timer.ObserveDuration()

			log.Println("** ALL RESOURCES ARE DEPLOYED AND READY **")

			err = testsuite.TestGPU(pod)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	addCoreFlags(cmd)
	addIngressFlags(cmd)

	return cmd
}
