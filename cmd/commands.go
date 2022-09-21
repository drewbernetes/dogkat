package cmd

import (
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/resources"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

var (
	//cfgFile           string
	kubeconfig        string
	valuesFile        string
	testAllFlag       bool
	testWorkloadsFlag bool
	rootCmd           = &cobra.Command{
		Use:   "k8s-e2e-tester",
		Short: "K8S End-2-End tester is an end-2-end tester which can locate an available sandbox and deploy workloads to it.",
		Long: `A End-2-End tester that can be used to spin up a sandbox cluster in EKS, 
			test all elements of a cluster rollout,
			and then spin it down again.
			Documentation is available here: https://github.com/drew-viles/k8s-e2e-tester/blob/main/README.md`,
		Run: func(cmd *cobra.Command, args []string) {
			connectToKubernetes(kubeconfig)
			determineTestCase()
		},
	}
	namespaceName = "e2e-testing"
)

// Execute is run by main and starts the program
func Execute() error {
	return rootCmd.Execute()
}

// init is auto run by cobra - all commands should be added here.
func init() {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")

	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", kubeConfigPath, fmt.Sprintf("kubeconfig to use defaults to: %s", kubeConfigPath))
	rootCmd.Flags().StringVarP(&valuesFile, "values", "v", "", "The Helm values file to use - required")
	rootCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "e2e-testing", "The Namespace to deploy the tests to")
	rootCmd.Flags().BoolVarP(&testAllFlag, "test-all", "a", false, "Simply tests everything it can - invokes all test commands - won't test Istio")
	rootCmd.Flags().BoolVarP(&testWorkloadsFlag, "test-standard-workload", "w", false, "Test that a workload can be deployed - this also tests Ingress, Cluster DNS, Storage and Scaling")

	//TODO: Next stage of functionality
	//rootCmd.Flags().BoolVarP(&testAWSConnectivityFlag, "test-oidc", "o", false, "Test that the AWS connectivity works via OIDC")
	//rootCmd.Flags().BoolVarP(&testIstioFlag, "test-istio", "m", false, "Test that the Istio service mesh is working at a basic level")

	err := rootCmd.MarkFlagRequired("values")
	if err != nil {
		fmt.Printf("there was an error with a required flag: %s\n", err.Error())
	}
	//TODO: Enable with Istio flag
	//rootCmd.MarkFlagsMutuallyExclusive("test-all", "test-istio")
	//rootCmd.MarkFlagsMutuallyExclusive("test-standard-workload", "test-istio")
	//rootCmd.MarkFlagsMutuallyExclusive("test-oidc", "test-istio")
}

// determineTestCase will parse the flags and run the appropriate test
func determineTestCase() {
	if testAllFlag {
		runCoreTests(valuesFile)
		//test_cases.CoreWorkloadChecks(valuesFile, namespaceName, clientsets)
		//TODO: Add more here as more are added.
	} else {
		if testWorkloadsFlag {
			runCoreTests(valuesFile)
			//test_cases.CoreWorkloadChecks(valuesFile, namespaceName, clientsets)
		}
	}
}

// parseResource will read through the supplied manifest file and work out what kind of API resource they are.
func parseResource(manifest string) resources.ApiResource {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(manifest), nil, nil)
	if err != nil {
		log.Printf("There was an error decoding: %s, %s\n", manifest, err)
		return nil
	}

	r := resources.ParseResourceKind(obj)
	if r == nil {
		return nil
	}

	r.GetClient(namespaceName, clientsets)
	return r
}
