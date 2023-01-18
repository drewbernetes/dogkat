# Kubernetes End-2-End Workload Tester

An End-2-End tester that will test a variety of elements of a Kubernetes cluster.

*The resources are baked into the binary (instead of using an external Helm Chart as previous versions did).*

The tests are separated out into logical workloads so that core workloads can be tested with additional tests able to be run on top.

It is capable of testing the following:
* Workload deployments
  * With anti-affinity ensuring workloads can be split across nodes.
  * Nginx with Configmap mounts
  * Postgres database with Volume mounts using PVC
  * Cluster DNS connectivity tested by connecting the Database to the Nginx workload
* Ingress deploys, resolves and responds
    * LoadBalancer creation via Ingress resources
    * Certificate generation using CertManager (Optional)
* Scaling workloads to test cluster-autoscaler (if deployed)
* Deploying Service Monitors and Grafana Dashboards (Requires Prometheus)

See the [Helm Chart README](https://github.com/drew-viles/helm-charts/blob/main/charts/e2e-basic/README.md) for more info

# TODO
* Implement Go Testing to validate the code.
* Test cloud connectivity.
* Istio detection and testing.

See below for a comprehensive list of tests and what can be confirmed using this tool.

# Usage
```
A End-2-End tester that can be used to spin up a sandbox cluster in EKS, 
                        test all elements of a cluster rollout,
                        and then spin it down again.
                        Documentation is available here: https://github.com/drew-viles/k8s-e2e-tester/blob/main/README.md

Usage:
  k8s-e2e-tester [flags]
  k8s-e2e-tester [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of EKS E2E-Tester

Flags:
  -h, --help                     help for k8s-e2e-tester
  -k, --kubeconfig string        kubeconfig to use defaults to: /home/drew/.kube/config (default "/home/drew/.kube/config")
  -n, --namespace string         The Namespace to deploy the tests to (default "default")
  -a, --test-all                 Simply tests everything it can - invokes all test commands - won't test Istio
  -w, --test-standard-workload   Test that a workload can be deployed - this also tests Ingress, Cluster DNS, Storage and Scaling
  -v, --values string            The Helm values file to use - required

Use "k8s-e2e-tester [command] --help" for more information about a command.

```