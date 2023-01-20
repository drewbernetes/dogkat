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

package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"
)

func NewDeleteIstioCmd(util.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Deletes the Istio testing resources",
		Long:  `Deletes the Istio application testing suite.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Istio will be supported soon™")
			//if namespace != "default" {
			//	err = deleteNamespaceOnSuccess(o.client, namespace)
			//	log.Fatalln(err)
			//}
		},
	}

	return cmd
}
