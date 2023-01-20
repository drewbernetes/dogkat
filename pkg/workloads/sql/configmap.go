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

// GeneratePostgresqlConfigMap returns a configmap that will be consumed by the Postgresql service.
func GeneratePostgresqlConfigMap(namespace string) *workloads.ConfigMap {
	cm := &workloads.ConfigMap{}
	cm.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlConfName, "label": constants.PGSqlName})
	initUserDB := `#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE TABLE IF NOT EXISTS web (id SERIAL PRIMARY KEY, value VARCHAR(10) NOT NULL );
INSERT INTO web (value) VALUES('ok');
EOSQL`
	cm.Resource.Data = map[string]string{
		"init-user-db.sh": initUserDB,
	}

	return cm
}
