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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/helpers"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

type PostgresWorkloads struct {
	Workload            *coreworkloads.StatefulSet
	InitConf            *coreworkloads.ConfigMap
	Secret              *coreworkloads.Secret
	ServiceAccount      *coreworkloads.ServiceAccount
	Service             *coreworkloads.Service
	PodDisruptionBudget *coreworkloads.PodDisruptionBudget
}

// CreateSQLWorkloadItems creates the Postgres workloads required for the core tests to run.
func CreateSQLWorkloadItems(client *kubernetes.Clientset, namespace, storageClass string) *PostgresWorkloads {
	w := &PostgresWorkloads{}
	w.InitConf = GeneratePostgresqlConfigMap(namespace)
	w.InitConf.Client = client

	w.Secret = GeneratePostgresqlSecret(namespace)
	w.Secret.Client = client

	w.ServiceAccount = &coreworkloads.ServiceAccount{}
	w.ServiceAccount.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlSAName, "label": constants.PGSqlName})
	w.ServiceAccount.Client = client

	w.Workload = GeneratePostgresStatefulSet(namespace, storageClass)
	w.Workload.Client = client

	w.Service = GeneratePostgresServiceResource(namespace)
	w.Service.Client = client

	w.PodDisruptionBudget = &coreworkloads.PodDisruptionBudget{}
	w.PodDisruptionBudget.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlName, "label": constants.PGSqlName})
	w.PodDisruptionBudget.Client = client

	helpers.HandleCreateError(w.InitConf.Create())
	helpers.HandleCreateError(w.Secret.Create())
	helpers.HandleCreateError(w.ServiceAccount.Create())
	helpers.HandleCreateError(w.Workload.Create())
	helpers.HandleCreateError(w.Service.Create())
	helpers.HandleCreateError(w.PodDisruptionBudget.Create())

	return w
}

// ValidateSQLWorkloadItems validates the Postgres workloads required for the core tests to run.
func (w *PostgresWorkloads) ValidateSQLWorkloadItems() error {
	var err error

	err = w.InitConf.Validate()
	if err != nil {
		return err
	}
	err = w.Secret.Validate()
	if err != nil {
		return err
	}
	err = w.ServiceAccount.Validate()
	if err != nil {
		return err
	}
	err = w.Workload.Validate()
	if err != nil {
		return err
	}
	err = w.Service.Validate()
	if err != nil {
		return err
	}
	err = w.PodDisruptionBudget.Validate()
	if err != nil {
		return err
	}

	return nil
}

// DeleteSQLWorkloadItems deletes the Postgres workloads required for the core tests to run.
func DeleteSQLWorkloadItems(client *kubernetes.Clientset, namespace string) {
	sqlCM := &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlConfName, Namespace: namespace}},
	}

	sqlSa := &coreworkloads.ServiceAccount{
		Client:   client,
		Resource: &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlSAName, Namespace: namespace}},
	}

	sqlSts := &coreworkloads.StatefulSet{
		Client:   client,
		Resource: &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	sqlSecret := &coreworkloads.Secret{
		Client:   client,
		Resource: &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlPasswdName, Namespace: namespace}},
	}

	sqlSvc := &coreworkloads.Service{
		Client:   client,
		Resource: &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	sqlPdb := &coreworkloads.PodDisruptionBudget{
		Client:   client,
		Resource: &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	err := sqlCM.Delete()
	if err != nil {
		log.Println(err)
	}
	err = sqlSa.Delete()
	if err != nil {
		log.Println(err)
	}
	err = sqlSts.Delete()
	if err != nil {
		log.Println(err)
	}
	err = sqlSecret.Delete()
	if err != nil {
		log.Println(err)
	}
	err = sqlSvc.Delete()
	if err != nil {
		log.Println(err)
	}
	err = sqlPdb.Delete()
	if err != nil {
		log.Println(err)
	}
}
