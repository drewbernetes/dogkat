package delete

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewDeleteGpuCmd(f util.Factory) *cobra.Command {
	o := &deleteOptions{}

	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Delete a GPU workload.",
		Long:  "Deletes the GPU workload that has been deployed for e2e-testing",
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

			pod := &coreworkloads.Pod{
				Client: o.client,
				Resource: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: constants.GPUName, Namespace: namespace},
				},
			}

			err = pod.Delete()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
