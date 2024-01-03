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

package workloads

import (
	"github.com/eschercloudai/dogkat/pkg/constants"
	"github.com/eschercloudai/dogkat/pkg/helm"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	optsv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type StatefulSet struct {
	*appsv1.StatefulSet
	v1.StatefulSetInterface
}

func NewStatefulSet(client *helm.Client) (*StatefulSet, error) {
	var err error
	s := &StatefulSet{}

	s.StatefulSetInterface = client.KubeClient.AppsV1().StatefulSets(client.Settings.Namespace())
	s.StatefulSet, err = s.get()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *StatefulSet) get() (*appsv1.StatefulSet, error) {
	return s.StatefulSetInterface.Get(context.Background(), constants.PGSqlName, optsv1.GetOptions{})
}

func (s *StatefulSet) Name() string {
	return s.StatefulSet.ObjectMeta.Name
}

func (s *StatefulSet) Kind() string {
	return "StatefulSet"
}

// IsReady Will report if the stateful set has the correct number of replicas.
// If this is true then we also know the PVC/PV has been created and bound as this would not become ready until that is done.
func (s *StatefulSet) IsReady() bool {
	var err error
	s.StatefulSet, err = s.get()
	if err != nil {
		return false
	}

	if s.Status.AvailableReplicas != s.Status.Replicas {
		return false
	}
	return true
}
