package create

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/gpu"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	numberOfGPUs string
)

func newCreateGpuCmd(f util.Factory) *cobra.Command {
	o := &createOptions{}

	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Creates the GPU testing resource from Nvidia.",
		Long: `Creates the GPU workload to allow testing to be done against a GPU 
to ensure it is enabled and working as expected.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			log.Println("checking for namespace, will create if doesn't exist")
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			createNamespaceIfNotExists(o.client, namespace)

			// Generate and create workloads
			pod := gpu.GenerateGPUPod(namespace, numberOfGPUs)
			helpers.HandleCreateError(pod.Create())
		},
	}
	cmd.Flags().StringVar(&numberOfGPUs, "number-of-gpus", "1", "Sets the number of GPUS in resources.limits")

	return cmd
}
