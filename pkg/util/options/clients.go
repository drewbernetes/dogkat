/*
Copyright 2024 Drewbernetes.

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

package options

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"log"
)

type K8SClients struct {
	Client *kubernetes.Clientset
	//TODO: Enable these once they are implemented
	//Istio      *istioclient.Clientset
	//Prometheus *promclient.Clientset
}

func newClients(configFlags *genericclioptions.ConfigFlags) K8SClients {
	var err error
	var client *kubernetes.Clientset

	f := util.NewFactory(configFlags)
	if client, err = f.KubernetesClientSet(); err != nil {
		log.Fatalln(err)
	}
	return K8SClients{
		Client: client,
	}
}
