package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/gpu"
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

			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

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

			err = testsuite.TestGPU(pod)
			if err != nil {
				log.Fatalln(err)
			}

			//err = testsuite.ScaleUpGPUNodes(pod)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//err = testsuite.TestGPU(pod)
			//if err != nil {
			//	log.Fatalln(err)
			//}
		},
	}
	cmd.Flags().StringVar(&numberOfGPUsFlag, "number-of-gpus", "1", "Sets the number of GPUS in resources.limits")

	return cmd
}
