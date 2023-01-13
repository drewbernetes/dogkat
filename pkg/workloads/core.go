package workloads

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/sql"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/web"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

type NginxWorkloads struct {
	Workload            *coreworkloads.Deployment
	Configuration       *coreworkloads.ConfigMap
	WebPages            *coreworkloads.ConfigMap
	ServiceAccount      *coreworkloads.ServiceAccount
	Service             *coreworkloads.Service
	PodDisruptionBudget *coreworkloads.PodDisruptionBudget
}

// CreateNginxWorkloadItems creates the Nginx workloads required for the core tests to run.
func CreateNginxWorkloadItems(client *kubernetes.Clientset, namespace string) *NginxWorkloads {
	w := &NginxWorkloads{}

	w.Configuration = web.GenerateNginxConfigMap(namespace)
	w.Configuration.Client = client

	w.WebPages = web.GenerateWebpageConfigMap(namespace)
	w.WebPages.Client = client

	w.ServiceAccount = &coreworkloads.ServiceAccount{}
	w.ServiceAccount.Generate(map[string]string{"namespace": namespace, "name": constants.NginxSAName, "label": constants.NginxName})
	w.ServiceAccount.Client = client

	w.Workload = web.GenerateNginxDeploy(namespace)
	w.Workload.Client = client

	w.Service = web.GenerateNginxServiceResource(namespace)
	w.Service.Client = client

	w.PodDisruptionBudget = &coreworkloads.PodDisruptionBudget{}
	w.PodDisruptionBudget.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "label": constants.NginxName})
	w.PodDisruptionBudget.Client = client

	helpers.HandleCreateError(w.Configuration.Create())
	helpers.HandleCreateError(w.WebPages.Create())
	helpers.HandleCreateError(w.ServiceAccount.Create())
	helpers.HandleCreateError(w.Workload.Create())
	helpers.HandleCreateError(w.Service.Create())
	helpers.HandleCreateError(w.PodDisruptionBudget.Create())

	return w
}

// ValidateNginxWorkloadItems validates the Nginx workloads required for the core tests to run.
func (w *NginxWorkloads) ValidateNginxWorkloadItems() error {
	var err error

	err = w.Configuration.Validate()
	if err != nil {
		return err
	}
	err = w.WebPages.Validate()
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

// DeleteNginxWorkloadItems deletes the Nginx workloads required for the core tests to run.
func DeleteNginxWorkloadItems(client *kubernetes.Clientset, namespace string) {
	nginxCM := &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxConfName, Namespace: namespace}},
	}

	nginxWebPages := &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxPagesName, Namespace: namespace}},
	}

	nginxSa := &coreworkloads.ServiceAccount{
		Client:   client,
		Resource: &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxSAName, Namespace: namespace}},
	}

	nginxDeploy := &coreworkloads.Deployment{
		Client:   client,
		Resource: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	nginxSvc := &coreworkloads.Service{
		Client:   client,
		Resource: &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	nginxPdb := &coreworkloads.PodDisruptionBudget{
		Client:   client,
		Resource: &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	err := nginxCM.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxWebPages.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxSa.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxDeploy.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxSvc.Delete()
	if err != nil {
		log.Println(err)
	}
	err = nginxPdb.Delete()
	if err != nil {
		log.Println(err)
	}
}

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
	w.InitConf = sql.GeneratePostgresqlConfigMap(namespace)
	w.InitConf.Client = client

	w.Secret = sql.GeneratePostgresqlSecret(namespace)
	w.Secret.Client = client

	w.ServiceAccount = &coreworkloads.ServiceAccount{}
	w.ServiceAccount.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlSAName, "label": constants.PGSqlName})
	w.ServiceAccount.Client = client

	w.Workload = sql.GeneratePostgresStatefulSet(namespace, storageClass)
	w.Workload.Client = client

	w.Service = sql.GeneratePostgresServiceResource(namespace)
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
