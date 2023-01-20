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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	hostFlag         string
	ingressClassFlag string
	enableTLSFlag    bool
	annotationsFlag  string
)

func newValidateIngressCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Creates and tests an ingress",
		Long: `Creates the core workload resources and corresponding ingress resource. Once the ingress is validated, 
testing of the ingress setup will occur. This will ensure that cert-manager, external-dns and ingress are all working as expected.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			_, _ = workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)

			ing := web.CreateIngressResource(o.client, namespace.Name, annotationsFlag, hostFlag, ingressClassFlag, enableTLSFlag)
			web.ValidateIngressResource(ing)

			err = testsuite.TestIngress(hostFlag)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	addCoreFlags(cmd)
	addIngressFlags(cmd)

	return cmd
}
