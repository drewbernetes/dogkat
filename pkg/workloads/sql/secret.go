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

package sql

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/constants"
	workloads "github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
)

// GeneratePostgresqlSecret returns a Secret that will be consumed by the Postgresql service.
func GeneratePostgresqlSecret(namespace string) *workloads.Secret {
	s := &workloads.Secret{}
	s.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlPasswdName, "label": constants.PGSqlName})
	s.Resource.Data = map[string][]byte{
		"passwd": []byte(constants.DBPassword),
	}

	return s
}
