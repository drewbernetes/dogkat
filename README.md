# Kubernetes End-2-End Workload Tester
This project was started as a way of testing as many components of a cluster as possible.

It will use my [e2e Helm Chart](https://github.com/drew-viles/helm-charts/tree/main/charts/e2e-basic) and as a result will require a values.yaml to be provided to it via the --values -v flag

It will test the following:
* Certificate generation using CertManager
* Workload deployments
    * With anti-affinity ensuring workloads can be split across nodes.
    * Nginx
      * With Configmap mounts
    * Postgres database 
      * With Volume mounts using PVC
* Ingress deploys, resolves and responds
* Scaling workloads to test cluster-autoscaler (if deployed)
* Cluster DNS connectivity tested by connecting the Database to the Nginx workload

See the [Helm Chart README](https://github.com/drew-viles/helm-charts/blob/main/charts/e2e-basic/README.md) for more info

# TODO
* AWS Connectivity - Confirms OIDC is working as it should
* Istio detection and testing
* Code the tests workloads into the application rather than relying on a 3rd party set of yaml files.
  * This allows for in-app testing
  * Detecting ingress controllers and deploying workloads to detect them rather than programmatically guessing based on **-a** or **-b** being on the end of the name.


See below for a comprehensive list of tests and what can be confirmed using this tool.

# Usage
```
A End-2-End tester that can be used to spin up a sandbox cluster in EKS, 
                test all elements of a cluster rollout,
				and then spin it down again.
				Documentation is available here: SOME-WEB-LINK

Usage:
  k8s-e2e-test [flags]
  k8s-e2e-test [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of EKS E2E-Tester

Flags:
  -n, --cluster-name string      Cluster Name - required
  -c, --create                   Supply this flag if a new cluster is required
  -h, --help                     help for k8s-e2e-test
  -k, --kubeconfig string        Cluster Name - required (default "filepath.Join(home, \".kube\", \"config\")")
  -t, --run-terraform            Run Terraform against the cluster
  -a, --test-all                 Simply tests everything it can - invokes all test commands - won't test Istio
  -m, --test-istio               Test that the istio service mesh is working at a basic level
  -o, --test-oidc                Test that the AWS connectivity works via OIDC
  -w, --test-standard-workload   Test that a workload can be deployed - this also tests Ingress, Storage and Scaling

Use "k8s-e2e-test [command] --help" for more information about a command.
```