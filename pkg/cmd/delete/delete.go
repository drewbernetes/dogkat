/*
Copyright 2024 EscherCloud.

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
	"context"
	"errors"
	"github.com/eschercloudai/dogkat/pkg/helm"
	"github.com/eschercloudai/dogkat/pkg/util/options"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/storage/driver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"log"
)

func NewDeleteCommand(cf *genericclioptions.ConfigFlags) *cobra.Command {
	configFlags := cf
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes test resources",
		Long:  `Deletes the resources deployed for testing and any associated namespaces as per the subcommand provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			o := options.NewOptions(configFlags)
			//Create a new client
			client, err := helm.NewClient(*configFlags.Namespace)
			client.KubeClient = o.Client

			if err != nil {
				return err
			}

			_, err = client.ChartDeployed()
			if err != nil {
				if errors.Is(err, driver.ErrReleaseNotFound) {
					return nil
				}
				return err
			}

			if err = client.Uninstall(); err != nil {
				return err
			}

			if err = deleteNamespaceOnSuccess(client.KubeClient, *configFlags.Namespace); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func deleteNamespaceOnSuccess(client *kubernetes.Clientset, namespace string) error {
	if namespace == "default" {
		return nil
	}

	_, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	log.Printf("deleting namespace: %s\n", namespace)

	err = client.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})

	if err != nil {
		return err
	}
	return nil
}
