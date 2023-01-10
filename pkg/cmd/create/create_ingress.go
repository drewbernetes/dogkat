package create

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newCreateIngressCmd(f util.Factory) *cobra.Command {
	o := &createOptions{}

	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Creates the workloads and resources for testing an ingress controller works.",
		Long: `This will test the ingress setup. This will ensure that cert-manager, 
external-dns and ingress are all working as expected.`,
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

			ing := web.GenerateWebIngressResource(namespace, map[string]string{})
			err = ing.Create()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	cmd.Flags().StringVar(&storageClass, "storage-class", "longhorn", "Used to define the name of the storage class to use for Persistent Volumes.")

	return cmd
}
