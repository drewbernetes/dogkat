package sql

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	workloads "github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
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
