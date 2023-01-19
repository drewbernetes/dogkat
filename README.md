# Kubernetes End-2-End Workload Tester

An End-2-End tester that will test a variety of elements of a Kubernetes cluster.

*The resources are baked into the binary (instead of using an external Helm Chart as previous versions did).*

The tests are separated out into logical workloads so that core workloads can be tested with additional tests able to be run on top.

| Run              | Description                                                                                                                                                                                                                                                                                                                                                                                         |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Core             | Deploys an Nginx deployment and SQL StatefulSet. Nginx Pod is a collection of Nginx, PHP and Nginx Prometheus Exporter. It has ConfigMap mounts and SQL password mounts and has Affinity with SQL and anti-affinity with itself. SQL contains Postgres with CM and Secret Mounting along with PersistentVolumeClaims for testing the CSI. It has affinity with Nginx and anti-affinity with itself. |
| Ingress          | Contains the Core workload and also deploys an Ingress resource. The Ingress can be configured to support TLS too allowing testing of things like Cert-Manager.                                                                                                                                                                                                                                     |
| GPU              | Deploys an Nvidia sample application for adding Vectors. This test will be targeted at a GPU node and will confirm the function of a GPU in a cluster.                                                                                                                                                                                                                                              |
| **Coming Soon ** |                                                                                                                                                                                                                                                                                                                                                                                                     |
| Monitoring       | Will deploy the core workload and also deploy a ServiceMonitor and a Grafana Dashboard for the Nginx service to display Grafana Functionality. The Service Monitor will need to be manually checked until a decision on vaildation that this is working can be achieved.                                                                                                                            |
| Istio            | This will deploy the core workload and then add a Virtual Service and Gateway on. It will then be validated in the same was as the Ingress test. This will confirm basic functionality of Istio.                                                                                                                                                                                                    |


# TODO
* Detect availability of things like an Ingress Controller and StorageClass before running tests.
* Implement Go Testing.
* Monitoring detection and testing.
* Istio detection and testing.

See below for a comprehensive list of tests and what can be confirmed using this tool.

# Usage
For details on how to use the tool, run `e2e-test --help`

# Examples

Test core workloads with a defined storage class:
```shell
e2e-test validate core --storage-class longhorn
```

Test Ingress with tls:
```shell
e2e-test validate ingress --storage-class longhorn --ingress-class nginx --enable-tls --annotations cert-manager.io/cluster-issuer=letsencrypt
```

Test GPU
```shell
e2e-test validate gpu --number-of-gpus 1
```