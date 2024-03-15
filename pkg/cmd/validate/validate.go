/*
Copyright 2024 EscherCloud.

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

package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/drewbernetes/dogkat/pkg/constants"
	"github.com/drewbernetes/dogkat/pkg/helm"
	"github.com/drewbernetes/dogkat/pkg/testsuite"
	"github.com/drewbernetes/dogkat/pkg/tracing"
	"github.com/drewbernetes/dogkat/pkg/util"
	"github.com/drewbernetes/dogkat/pkg/util/options"
	"github.com/drewbernetes/dogkat/pkg/workloads"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
)

//TODO: Enable these if/when they are implemented

// addPrometheusToScheme adds the Prometheus scheme to the scheme so that the clientset can use it.
//func addPrometheusToScheme() {
//	err := promscheme.AddToScheme(scheme.Scheme)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

// addIstioScheme adds the Istio scheme to the scheme so that the clientset can use it.
//func addIstioToScheme() {
//	err := istioscheme.AddToScheme(scheme.Scheme)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

func NewValidateCommand(cf *genericclioptions.ConfigFlags) *cobra.Command {
	configFlags := cf

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Create and validates test resources",
		Long: `Creates a selection of resources based on the input provided. Once all resources are deployed and confirmed as ready,
the required test suite will run against the resources to ensure everything is working as expected within a cluster.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse config items
			o := options.NewOptions(configFlags)

			// Determine test types
			testType := util.TestTypes{}
			if o.CoreOptions.Enabled {
				testType.Core = true
			}
			if o.IngressOptions.Enabled {
				testType.Ingress = true
			}
			if o.GPUOptions.Enabled {
				testType.GPU = true
			}

			m := tracing.NewGatherer()
			var chartTracer, fullTracer *tracing.Duration

			// Configure and start metrics for Chart deploy
			if m.Enabled {
				n := fmt.Sprintf("%s_chart_deployment", testType.GetType())
				d := "Times the deployment of the chart before running tests"
				chartTracer = tracing.NewCollector(m.PushGateway, n, d)
				chartTracer.Start()
			}

			//Create a new client
			client, err := helm.NewClient(*configFlags.Namespace)
			if err != nil {
				return err
			}

			// Add Kubernetes clientset to the Client
			client.KubeClient = o.Client

			//Check if chart is deployed - if so, grab the Chart
			rel, err := client.ChartDeployed()
			if err != nil {
				if !errors.Is(err, driver.ErrReleaseNotFound) {
					return err
				}
			}

			// If not, pull and deploy!
			if rel == nil {
				// Download and Load Chart
				chart, err := helm.NewChart(client, testType, o)
				if err != nil {
					return err
				}

				// Deploy Chart
				rel, err = client.Install(chart)
				if err != nil {
					return err
				}

			}

			// Let's check it actually deployed before continuing (this doesn't mean the resources are deployed of course)
			if rel.Info.Status != release.StatusDeployed {
				return fmt.Errorf("The Chart is not deployed")
			}

			log.Println("waiting for resources to be ready")
			//Check STS and Deployment are deployed - we can presume everything else is as these two make use of mounting secrets, configmaps, volumes etc.
			var coreDeployment *workloads.Deployment
			if testType.Core {
				coreDeployment, err = checkCoreReady(client)
				if err != nil {
					return err
				}
			}

			var ingressResource *workloads.Ingress
			if testType.Ingress {
				_, err = checkCoreReady(client)
				if err != nil {
					return err
				}
				ingressResource, err = checkIngressReady(client)
				if err != nil {
					return err
				}
			}

			var gpuPod *workloads.Pod
			if testType.GPU {
				gpuPod, err = checkGPUReady(client)
				if err != nil {
					return err
				}
			}

			// End metrics for Chart deploy
			if m.Enabled {
				if err = chartTracer.CompleteGathering(); err != nil {
					return err
				}
			}

			log.Println("all resources ready")

			// Configure and start metrics for tests
			if m.Enabled {
				n := fmt.Sprintf("%s_duration_seconds", testType.GetType())
				d := fmt.Sprintf("Times the %s e2e test takes to complete", testType.GetType())
				fullTracer = tracing.NewCollector(m.PushGateway, n, d)
				fullTracer.Start()
			}

			// Run tests
			tracker := []testsuite.TestTracker{}

			if testType.Core {
				t := testsuite.NewScalingTest(coreDeployment, client)
				t.Init(o.CoreOptions.ScaleTo)

				if err = t.Run(); err != nil {
					return err
				}
				if err = t.Validate(); err != nil {
					log.Println(err)
				}

				tracker = append(tracker, *t.Tracker)
			}

			if testType.Ingress {
				t := testsuite.NewEndpointTest(ingressResource, client)
				t.Init(o.IngressOptions.Host, o.IngressOptions.EnableTLS)

				if err = t.Run(); err != nil {
					return err
				}
				if err = t.Validate(); err != nil {
					log.Println(err)
				}

				tracker = append(tracker, *t.Tracker)
			}

			if testType.GPU {
				t := testsuite.NewVectorTest(gpuPod, client)

				if err = t.Run(); err != nil {
					return err
				}
				if err = t.Validate(); err != nil {
					log.Println(err)
				}

				tracker = append(tracker, *t.Tracker)

				// TODO: Implement a GPU Scale test?
			}

			if m.Enabled {
				if err = fullTracer.CompleteGathering(); err != nil {
					return err
				}
			}

			// Check to see if any tests failed validation (failed runs will result in the program exiting)
			err = checkForFailedTests(tracker)
			if err != nil {
				return err
			}

			log.Println("tests complete")
			return nil
		},
	}

	return cmd
}

func checkForFailedTests(tracker []testsuite.TestTracker) error {
	failedTests := false
	e := "the following tests failed: "

	for _, v := range tracker {
		if !v.Completed {
			e = fmt.Sprintf("%s %s", e, v.Name)
			failedTests = true
		}
	}

	data, err := json.Marshal(tracker)
	if err != nil {
		log.Println("couldn't marshall the results data")
	}

	if data != nil {
		if err = os.WriteFile("/tmp/results.json", data, 0644); err != nil {
			log.Println("couldn't create results file")
		} else {
			log.Println("results written to /tmp/results.json")
		}
	}

	if failedTests {
		return fmt.Errorf("%s\n", e)
	}

	return nil
}

// checkCoreReady validates that the Deployment, StatefulSet and the PDB for both are in a ready state.
// By confirming just these 4 things we can be confident that everything else from configMaps and Secrets
// all the way through to the PVC and PV are deployed as these 4 resources wouldn't hit a ready state without any of them.
func checkCoreReady(client *helm.Client) (*workloads.Deployment, error) {
	d, err := workloads.NewDeployment(client)
	if err != nil {
		return nil, err
	}

	pdbn, err := workloads.NewPodDisruptionBudget(client, constants.NginxName)
	if err != nil {
		return nil, err
	}

	s, err := workloads.NewStatefulSet(client)
	if err != nil {
		return nil, err
	}

	//TODO Change this before release as the chart needs an updated tag
	pdbd, err := workloads.NewPodDisruptionBudget(client, constants.PGSqlName)
	if err != nil {
		return nil, err
	}

	// Create a channel which means we don't have to wait for each resource to be ready to check the next one.
	// They'll just return as they're ready.
	resources := []workloads.Resource{d, pdbn, s, pdbd}

	checksCompleted := make(chan error, 4)
	defer close(checksCompleted)
	readyCheck := func(r workloads.Resource) {
		err = workloads.CheckReadyState(r)
		if err != nil {
			checksCompleted <- err
			return
		}

		checksCompleted <- nil
	}

	for _, r := range resources {
		go readyCheck(r)
	}

	for range resources {
		<-checksCompleted
	}

	return d, nil
}

// checkIngressReady confirms the ingress is reporting as ready with an IP and host name
func checkIngressReady(client *helm.Client) (*workloads.Ingress, error) {
	i, err := workloads.NewIngress(client)
	if err != nil {
		return nil, err
	}

	return i, workloads.CheckReadyState(i)
}

// checkGPUReady just confirms the pod is deployed.
// To be honest this will probably complete before the test is even run.
func checkGPUReady(client *helm.Client) (*workloads.Pod, error) {
	p, err := workloads.NewPod(client)
	if err != nil {
		return nil, err
	}
	err = workloads.CheckReadyState(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
