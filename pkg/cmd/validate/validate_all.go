package validate

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	"github.com/drew-viles/k8s-e2e-tester/pkg/testsuite"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/gpu"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/prometheus"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/sql"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/web"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func newValidateAllCmd(f util.Factory) *cobra.Command {
	o := &validateOptions{}

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Create and validates all resources",
		Long: `Creates and validates all resources from core, gpu, ingress and monitoring.
Istio test suites will not be affected by this.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := workloads.CreateNamespaceIfNotExists(o.client, cmd.Flag("namespace").Value.String())

			// Generate and create workloads
			nginxWorkload, _ := workloads.DeployBaseWorkloads(o.client, namespace.Name, storageClassFlag, requestCPUFlag, requestMemoryFlag)
			err = testsuite.ScaleUpStandardNodes(nginxWorkload.Workload)
			if err != nil {
				log.Fatalln(err)
			}

			ing := web.CreateIngressResource(o.client, namespace.Name, annotationsFlag, hostFlag, ingressClassFlag, enableTLSFlag)
			web.ValidateIngressResource(ing)
			err = testsuite.TestIngress(hostFlag)
			if err != nil {
				log.Fatalln(err)
			}

			//TODO Check for service monitor resource
			sm := prometheus.GenerateServiceMonitorResource(namespace.Name)
			prometheus.CreateServiceMonitor(o.prometheus, sm)

			web.DeleteNginxWorkloadItems(o.client, namespace.Name)
			web.DeleteIngressWorkloadItems(o.client, namespace.Name)
			sql.DeleteSQLWorkloadItems(o.client, namespace.Name)

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
		},
	}
	addCoreFlags(cmd)
	addIngressFlags(cmd)

	return cmd
}
