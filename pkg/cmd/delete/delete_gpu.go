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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/constants"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

func NewDeleteGpuCmd(f util.Factory) *cobra.Command {
	o := &deleteOptions{}

	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Delete a GPU workload",
		Long:  "Deletes the GPU application testing suite.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			// Connect to cluster
			if o.client, err = f.KubernetesClientSet(); err != nil {
				log.Fatalln(err)
			}

			// Configure namespace
			namespace := "default"
			if cmd.Flag("namespace").Value.String() != "" {
				namespace = cmd.Flag("namespace").Value.String()
			}

			pod := &coreworkloads.Pod{
				Client: o.client,
				Resource: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: constants.GPUName, Namespace: namespace},
				},
			}

			err = pod.Delete()
			if err != nil {
				log.Fatalln(err)
			}
			if namespace != "default" {
				err = deleteNamespaceOnSuccess(o.client, namespace)
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
