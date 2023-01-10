package create

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	storageClass string
)

func newCreateCoreCmd(f util.Factory) *cobra.Command {
	o := &createOptions{}

	cmd := &cobra.Command{
		Use:   "core",
		Short: "Creates the workloads and resources for testing basic elements of the cluster.",
		Long: `This will create the resources for testing of cluster elements such as dns, 
storage, cluster-autoscalers and more.`,
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
			workloads.CreateNginxWorkloadItems(o.client, namespace)
			workloads.CreateSQLWorkloadItems(o.client, namespace, storageClass)
		},
	}

	cmd.Flags().StringVar(&storageClass, "storage-class", "longhorn", "Used to define the name of the storage class to use for Persistent Volumes.")

	return cmd
}
