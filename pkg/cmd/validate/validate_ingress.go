package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
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
			log.Println("checking for namespace, will create if doesn't exist")
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			createNamespaceIfNotExists(o.client, namespace)

			// Generate and create workloads
			nginxWorkload, sqlWorkload := createAndConfirmResources(o.client, namespace)

			log.Println("all resources are deployed, running tests...")

			coreResource := []coreworkloads.Resource{
				nginxWorkload.Configuration,
				nginxWorkload.WebPages,
				nginxWorkload.Workload,
				nginxWorkload.ServiceAccount,
				nginxWorkload.Service,
				nginxWorkload.PodDisruptionBudget,
				sqlWorkload.Workload,
				sqlWorkload.ServiceAccount,
				sqlWorkload.InitConf,
				sqlWorkload.Secret,
				sqlWorkload.Service,
				sqlWorkload.PodDisruptionBudget,
			}

			err = testsuite.CheckReadyForTesting(coreResource)
			if err != nil {
				log.Fatalln(err)
			}

			//Check volumes after STS confirmation to ensure volumes have been created.
			volumeResource := parseVolumesFromStatefulSet(o.client, sqlWorkload, namespace)
			err = testsuite.CheckReadyForTesting(volumeResource)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("** ALL RESOURCES ARE DEPLOYED AND READY **")

			ing := web.GenerateWebIngressResource(namespace, map[string]string{})
			err = ing.Create()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
