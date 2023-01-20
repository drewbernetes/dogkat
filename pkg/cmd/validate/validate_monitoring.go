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
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newValidateMonitoringCmd(f util.Factory) *cobra.Command {
	//TODO: Enable these once They are implemented
	//o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Creates and validates the core resources with additional monitoring",
		Long: `Creates an application with additional grafana dashboards and service monitors. 
This requires Prometheus to be installed. Once deployed the monitoring test suite will run to confirm 
resources are deployed and working as expected.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("monitoring will be supported soon™")
			//TODO: Enable these once They are implemented
			//var err error
			//
			//// Connect to cluster
			//if o.client, err = f.KubernetesClientSet(); err != nil {
			//	log.Fatalln(err)
			//}
			//
			//addPrometheusToScheme()
			//
			//// Configure namespace
			//namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())
			//
			//_, _ = workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)
			//
			////TODO Check for service monitor resource
			//sm := prometheus.GenerateServiceMonitorResource(namespace.Name)
			//prometheus.CreateServiceMonitor(o.prometheus, sm)
			//prometheus.GenerateGrafanaDashboardConfigMap(namespace.Name)
		},
	}
	addCoreFlags(cmd)
	return cmd
}
