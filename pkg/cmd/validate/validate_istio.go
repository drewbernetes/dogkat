package validate

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func newValidateIstioCmd(f util.Factory) *cobra.Command {
	//o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Creates and tests Istio",
		Long: `Creates the core workloads and add Istio resources on top to allow testing to be done.
This will ensure Istio is deployed and working on a basic level by testing things such as the VirtualService and Gateway.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Istio will be supported soon™")
			//var err error

			// Connect to cluster
			//if o.client, err = f.KubernetesClientSet(); err != nil {
			//	log.Fatalln(err)
			//}

			// Configure namespace
			//namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())
			//
			//_, _ = workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)
		},
	}
	addCoreFlags(cmd)

	return cmd
}
