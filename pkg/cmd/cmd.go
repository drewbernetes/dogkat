package cmd

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/cmd/create"
	"github.com/drew-viles/k8s-e2e-tester/pkg/cmd/delete"
	"github.com/drew-viles/k8s-e2e-tester/pkg/cmd/validate"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
)

var (
	valuesFile    string
	namespaceName = "e2e-testing"
)

// init is auto run by cobra - all commands should be added here.
func newRootCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)
	f := util.NewFactory(configFlags)

	cmd := &cobra.Command{
		Use:   "k8s-e2e-tester",
		Short: "K8S End-2-End tester is an end-2-end tester which can deploy workloads to a provided cluster.",
		Long: `An End-2-End tester that can be used to test all elements of a cluster rollout.
Documentation is available here: https://github.com/drew-viles/k8s-e2e-tester/blob/main/README.md`,
	}

	commands := []*cobra.Command{
		create.NewCreateCommand(f),
		validate.NewValidateCommand(f),
		delete.NewDeleteCommand(f),
		NewVersionCmd(),
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(commands...)
	//TODO: Move this where it's required.
	//cmd.Flags().StringVarP(&valuesFile, "values", "v", "", "The Helm values file to use - required")
	//
	//err := cmd.MarkFlagRequired("values")
	//if err != nil {
	//	fmt.Printf("there was an error with a required flag: %s\n", err.Error())
	//}
	return cmd
}

func Generate() *cobra.Command {
	return newRootCommand()
}

// determineTestCase will parse the flags and run the appropriate test
//func determineTestCase() {
//	if testAllFlag {
//		runCoreTests(valuesFile)
//		//test_cases.CoreWorkloadChecks(valuesFile, namespaceName, clientsets)
//		//TODO: Add more here as more are added.
//	} else {
//		if testWorkloadsFlag {
//			runCoreTests(valuesFile)
//			//test_cases.CoreWorkloadChecks(valuesFile, namespaceName, clientsets)
//		}
//	}
//}

// parseResource will read through the supplied manifest file and work out what kind of API resource they are.
//func parseResource(manifest string) resources.ApiResource {
//	decode := scheme.Codecs.UniversalDeserializer().Decode
//	obj, _, err := decode([]byte(manifest), nil, nil)
//	if err != nil {
//		log.Printf("There was an error decoding: %s, %s\n", manifest, err)
//		return nil
//	}
//
//	r := resources.ParseResourceKind(obj)
//	if r == nil {
//		return nil
//	}
//
//	r.GetClient(namespaceName, cmd.clientsets)
//	return r
//}
