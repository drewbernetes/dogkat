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
func CreateNginxWorkloadItems(client *kubernetes.Clientset, namespace string) {
	nginxCM := web.GenerateNginxConfigMap(namespace)
	nginxCM.Client = client

	nginxWebPages := web.GenerateWebpageConfigMap(namespace)
	nginxWebPages.Client = client

	nginxSa := &coreworkloads.ServiceAccount{}
	nginxSa.Generate(map[string]string{"namespace": namespace, "name": constants.NginxSAName, "label": constants.NginxName})
	nginxSa.Client = client

	nginxDeploy := web.GenerateNginxDeploy(namespace)
	nginxDeploy.Client = client

	nginxSvc := web.GenerateNginxServiceResource(namespace)
	nginxSvc.Client = client

	nginxPdb := &coreworkloads.PodDisruptionBudget{}
	nginxPdb.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "label": constants.NginxName})
	nginxPdb.Client = client

	helpers.HandleCreateError(nginxCM.Create())
	helpers.HandleCreateError(nginxWebPages.Create())
	helpers.HandleCreateError(nginxSa.Create())
	helpers.HandleCreateError(nginxDeploy.Create())
	helpers.HandleCreateError(nginxSvc.Create())
	helpers.HandleCreateError(nginxPdb.Create())
}

// ValidateNginxWorkloadItems validates the Nginx workloads required for the core tests to run.
func ValidateNginxWorkloadItems(client *kubernetes.Clientset, namespace string) (*NginxWorkloads, error) {
	var err error
	workload := &NginxWorkloads{}

	workload.Configuration = &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxConfName, Namespace: namespace}},
	}

	workload.WebPages = &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxPagesName, Namespace: namespace}},
	}

	workload.ServiceAccount = &coreworkloads.ServiceAccount{
		Client:   client,
		Resource: &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxSAName, Namespace: namespace}},
	}

	workload.Workload = &coreworkloads.Deployment{
		Client:   client,
		Resource: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	workload.Service = &coreworkloads.Service{
		Client:   client,
		Resource: &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	workload.PodDisruptionBudget = &coreworkloads.PodDisruptionBudget{
		Client:   client,
		Resource: &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: constants.NginxName, Namespace: namespace}},
	}

	err = workload.Configuration.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.WebPages.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.ServiceAccount.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.Workload.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.Service.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.PodDisruptionBudget.Validate()
	if err != nil {
		return nil, err
	}

	return workload, nil
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
func CreateSQLWorkloadItems(client *kubernetes.Clientset, namespace, storageClass string) {
	sqlCM := sql.GeneratePostgresqlConfigMap(namespace)
	sqlCM.Client = client

	sqlSecret := sql.GeneratePostgresqlSecret(namespace)
	sqlSecret.Client = client

	sqlSa := &coreworkloads.ServiceAccount{}
	sqlSa.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlSAName, "label": constants.PGSqlName})
	sqlSa.Client = client

	sqlSts := sql.GeneratePostgresStatefulSet(namespace, storageClass)
	sqlSts.Client = client

	sqlSvc := sql.GeneratePostgresServiceResource(namespace)
	sqlSvc.Client = client

	sqlPdb := &coreworkloads.PodDisruptionBudget{}
	sqlPdb.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlName, "label": constants.PGSqlName})
	sqlPdb.Client = client

	helpers.HandleCreateError(sqlCM.Create())
	helpers.HandleCreateError(sqlSecret.Create())
	helpers.HandleCreateError(sqlSa.Create())
	helpers.HandleCreateError(sqlSts.Create())
	helpers.HandleCreateError(sqlSvc.Create())
	helpers.HandleCreateError(sqlPdb.Create())
}

// ValidateSQLWorkloadItems validates the Postgres workloads required for the core tests to run.
func ValidateSQLWorkloadItems(client *kubernetes.Clientset, namespace string) (*PostgresWorkloads, error) {
	var err error
	workload := &PostgresWorkloads{}

	workload.InitConf = &coreworkloads.ConfigMap{
		Client:   client,
		Resource: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlConfName, Namespace: namespace}},
	}

	workload.Secret = &coreworkloads.Secret{
		Client:   client,
		Resource: &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlPasswdName, Namespace: namespace}},
	}

	workload.ServiceAccount = &coreworkloads.ServiceAccount{
		Client:   client,
		Resource: &v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlSAName, Namespace: namespace}},
	}

	workload.Workload = &coreworkloads.StatefulSet{
		Client:   client,
		Resource: &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	workload.Service = &coreworkloads.Service{
		Client:   client,
		Resource: &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	workload.PodDisruptionBudget = &coreworkloads.PodDisruptionBudget{
		Client:   client,
		Resource: &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: constants.PGSqlName, Namespace: namespace}},
	}

	err = workload.InitConf.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.Secret.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.ServiceAccount.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.Workload.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.Service.Validate()
	if err != nil {
		return nil, err
	}
	err = workload.PodDisruptionBudget.Validate()
	if err != nil {
		return nil, err
	}

	return workload, nil
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
