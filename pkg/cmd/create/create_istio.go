package create

import (
	"github.com/spf13/cobra"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newCreateIstioCmd(f util.Factory) *cobra.Command {
	o := &createOptions{}

	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Creates the Istio testing resource from Nvidia.",
		Long:  `Creates an Istio application to allow testing to be done to ensure Istio is deployed and working.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			config, err := f.ToRESTConfig()
			if err != nil {
				log.Fatalln(err)
			}

			addIstioToScheme()

			o.istio, err = istioclient.NewForConfig(config)
			if err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			log.Println("checking for namespace, will create if doesn't exist")
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			createNamespaceIfNotExists(o.client, namespace)

		},
	}

	return cmd
}
