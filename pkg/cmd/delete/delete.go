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
	"context"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

type deleteOptions struct {
	client *kubernetes.Clientset
	//TODO: Enable these once They are implemented
	//istio      *istioclient.Clientset
	//prometheus *promclient.Clientset
}

func NewDeleteCommand(f util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes test resources",
		Long:  `Deletes the resources deployed for testing and any associated namespaces as per the subcommand provided.`,
	}

	commands := []*cobra.Command{
		NewDeleteAllCmd(f),
		NewDeleteCoreCmd(f),
		NewDeleteGpuCmd(f),
		NewDeleteIngressCmd(f),
		NewDeleteIstioCmd(f),
		NewDeleteMonitoringCmd(f),
	}

	cmd.AddCommand(commands...)

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
