package sql

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	workloads "github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
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
