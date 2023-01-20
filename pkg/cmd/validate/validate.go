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
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	storageClassFlag  string
	requestCPUFlag    string
	requestMemoryFlag string
)

type validateOptions struct {
	client *kubernetes.Clientset
	//TODO: Enable these once They are implemented
	//istio      *istioclient.Clientset
	//prometheus *promclient.Clientset
}

func NewValidateCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Create and validates test resources",
		Long: `Creates a selection of resources as per the subcommand provided. Once all resources are deployed and confirmed as ready,
the required test suite will run against the resources to ensure everything is working as expected within a cluster.`,
	}

	commands := []*cobra.Command{
		newValidateAllCmd(f),
		newValidateCoreCmd(f),
		newValidateGpuCmd(f),
		newValidateIngressCmd(f),
		newValidateIstioCmd(f),
		newValidateMonitoringCmd(f),
	}

	cmd.AddCommand(commands...)

	return cmd
}

// addCoreFlags adds the flags required for the core workload. It can't be persistent as some just won't be required for all tests.
func addCoreFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&storageClassFlag, "storage-class", "", "The name of the storage class to use for Persistent Volumes")
	cmd.Flags().StringVar(&requestCPUFlag, "cpu", "1", "The request CPU value to ensure scaling happens")
	cmd.Flags().StringVar(&requestMemoryFlag, "memory", "1Gi", "The request memory value to ensure scaling happens")
}

// addIngressFlags adds the flags required for the ingress workload. It can't be persistent as some just won't be required for all tests.
func addIngressFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&ingressClassFlag, "ingress-class", "nginx", "The IngressClass name")
	cmd.Flags().BoolVar(&enableTLSFlag, "enable-tls", false, "Whether to enable TLS on the Ingress endpoint. You must have cert-manager enabled and configured for this test to succeed")
	cmd.Flags().StringVar(&annotationsFlag, "annotations", "", "Any additional annotations to add to the ingress resource in the format 'a=1,b=2'")
	cmd.Flags().StringVar(&hostFlag, "host", "", "The fqdn for the ingress resource")
	err := cmd.MarkFlagRequired("host")
	if err != nil {
		log.Fatalln(err)
	}
}

//TODO: Enable these once They are implemented

// addPrometheusToScheme adds the Prometheus scheme to the scheme so that the clientset can use it.
//func addPrometheusToScheme() {
//	err := promscheme.AddToScheme(scheme.Scheme)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

// addIstioScheme adds the Istio scheme to the scheme so that the clientset can use it.
//func addIstioToScheme() {
//	err := istioscheme.AddToScheme(scheme.Scheme)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
