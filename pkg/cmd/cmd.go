/*
Copyright 2022 EscherCloud.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/cmd/delete"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/cmd/metrics"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/cmd/validate"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
)

// init is auto run by cobra - all commands should be added here.
func newRootCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)
	f := util.NewFactory(configFlags)

	cmd := &cobra.Command{
		Use: "k8s-e2e-tester",
		Long: `Deploys resources to allow End-2-End testing to be conducted. 
It can be used to test most elements of a cluster to ensure consistent stability and functionality.
Documentation is available here: https://github.com/eschercloudai/k8s-e2e-tester/blob/main/README.md`,
	}

	commands := []*cobra.Command{
		metrics.NewMetricsCommand(f),
		validate.NewValidateCommand(f),
		delete.NewDeleteCommand(f),
		NewVersionCmd(),
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(commands...)
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
