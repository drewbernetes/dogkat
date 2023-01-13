package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
	"strings"
)

func newValidateCoreCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "core",
		Short: "Create and test core resources",
		Long: `Creates a selection of standard resources. Once all resources are deployed and confirmed as ready,
the test suite will then run against them to ensure everything is working as it should be within a standard cluster.

The following is tested:
* Deployments, StatefulSets & Services to validate Cluster DNS
* Pod Disruption Budgets to confirm HA
* ConfigMaps & Secrets toc onverm ENVs and Volume Mounting of ConfigMaps
* PV creation & Mounting via PVC in StatefulSet
`,
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

			err = testsuite.ScaleUpStandardNodes(nginxWorkload.Workload)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	return cmd
}

func createAndConfirmResources(client *kubernetes.Clientset, namespace string) (*workloads.NginxWorkloads, *workloads.PostgresWorkloads) {
	nginxWorkload := workloads.CreateNginxWorkloadItems(client, namespace)
	sqlWorkload := workloads.CreateSQLWorkloadItems(client, namespace, storageClassFlag)

	err := nginxWorkload.ValidateNginxWorkloadItems()
	if err != nil {
		log.Fatalln(err)
	}
	err = sqlWorkload.ValidateSQLWorkloadItems()
	if err != nil {
		log.Fatalln(err)
	}
	return nginxWorkload, sqlWorkload
}

func parseVolumesFromStatefulSet(client *kubernetes.Clientset, sqlWorkload *workloads.PostgresWorkloads, namespace string) []coreworkloads.Resource {
	volumeResource := []coreworkloads.Resource{}
	volNames := []string{strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "0"}, "-"),
		strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "1"}, "-"),
		strings.Join([]string{"data", sqlWorkload.Workload.GetResourceName(), "2"}, "-"),
	}

	for _, name := range volNames {
		pvc := &coreworkloads.PersistentVolumeClaim{
			Client:   client,
			Resource: &v1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}},
		}
		err := pvc.Validate()
		if err != nil {
			log.Println(err)
		}
		volumeResource = append(volumeResource, pvc)

		pvName := pvc.Resource.Spec.VolumeName
		pv := &coreworkloads.PersistentVolume{
			Client:   client,
			Resource: &v1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: pvName}},
		}
		err = pv.Validate()
		if err != nil {
			log.Println(err)
		}
		volumeResource = append(volumeResource, pv)
	}

	return volumeResource
}
