package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
	"os"
)

func NewValidateGpuCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Validate and test the GPU workload",
		Long:  "Runs the validation and test suite against the GPU workload to ensure it is working as expected.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			//Check the pod exists
			pod := &coreworkloads.Pod{
				Client: o.client,
				Resource: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: constants.GPUName, Namespace: namespace},
				},
			}
			err = pod.Validate()
			if err != nil {
				log.Fatalln(err)
			}

			err = testsuite.TestGPU(pod)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
