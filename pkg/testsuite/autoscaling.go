/*
Copyright 2024 Drewbernetes.

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

package testsuite

import (
	"fmt"
	"github.com/drewbernetes/dogkat/pkg/helm"
	"github.com/drewbernetes/dogkat/pkg/helpers"
	"github.com/drewbernetes/dogkat/pkg/tracing"
	"github.com/drewbernetes/dogkat/pkg/workloads"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"log"
	"time"
)

// ScalingTest is used to validate the functionality of the autoscaler.
// TODO: This can be improved on using the following process
//   - create resources
//   - scale resources
//   - check if new node exists
//   - if not, check if any errors exist with regards to scaling
//   - if not, scale again and repeat last two steps
//   - if an error occurs - throw it back to the user
//   - if a new node appears, it's successful
type ScalingTest struct {
	Test
	StartingNodes    int
	EndNodes         int
	StartingReplicas int32
	TargetReplicas   int32
	Deployment       *workloads.Deployment
	Tracker          *TestTracker
}

// NewScalingTest returns a new Scaling test which will check the cluster autoscaler works.
func NewScalingTest(d *workloads.Deployment, c *helm.Client) *ScalingTest {
	name := "cluster_autoscaler_test"
	description := "The time it takes from scaling workloads to the autoscaler bringing in a new node and the new workloads being deployed"
	t := NewTest(c, name, description)

	return &ScalingTest{
		Test:       t,
		Deployment: d,
		Tracker: &TestTracker{
			Name: name,
			Description: `The scaling test has been designed to test node auto-scaling as well as some other core components.
It deploys a basic web service as a Deployment and database as a StatefulSet along with some Secrets, Configmaps, Pod Disruption Budgets and some other core resources.
The purpose of this test is to ensure that a basic workload can be deployed. When queried, the website returns an 'ok' value which is pulled from the database otherwise it throws and error.
This confirms the functionality of CoreDNS within the cluster. Once deployed, it scales the nginx (web) workload up to the specified value in the config supplied which, if enough, will trigger a node scale.
If the nodes do not scale up, then an error is returned.`,
			Completed: false,
		},
	}
}

// Init prepares the test with initial conditions so that a starting point can be used to validate later
func (c *ScalingTest) Init(targetReplicas int32) {
	c.StartingNodes = countNodes(c.Client.KubeClient)
	c.StartingReplicas = *c.Deployment.Deployment.Spec.Replicas
	c.TargetReplicas = targetReplicas
}

// Run the actual test
func (c *ScalingTest) Run() error {
	// Start metrics gathering
	if tracing.Gatherer().Enabled {
		c.Tracing.Start()
	}

	log.Printf("Running Test: %s\n", c.Test.Name)
	log.Printf("Node count before Scale %v\n", c.StartingNodes)

	re := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploy, err := c.Deployment.Get()
		if err != nil {
			return err
		}

		// Increase the replicas
		deploy.Spec.Replicas = helpers.IntPtr(c.TargetReplicas)

		d, err := c.Deployment.Update(context.Background(), deploy, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		c.Deployment.Deployment = d
		return nil
	})

	if re != nil {
		return re
	}

	err := checkForUpdatedDeployment(c, c.TargetReplicas)
	if err != nil {
		return err
	}

	log.Printf("Waiting for Deployment to scale\n")

	err = workloads.CheckReadyState(c.Deployment)
	if err != nil {
		// It seems the resource never became ready - let's print out why
		reason := c.Deployment.Deployment.Status.String()
		return fmt.Errorf("%s - %s", err, reason)
	}

	//The Deployment scaled... lets count the nodes, did they?
	c.EndNodes = countNodes(c.Client.KubeClient)

	// Complete metrics gathering

	if tracing.Gatherer().Enabled {
		if err = c.Tracing.CompleteGathering(); err != nil {
			return err
		}
	}
	return nil
}

// Validate compares the end result of the test with the values set in the Init stage allowing a comparison to see if the test passed.
func (c *ScalingTest) Validate() error {
	// validate the scaling to see if the nodes actually scaled
	if c.EndNodes <= c.StartingNodes {
		//the nodes didn't scale
		return fmt.Errorf("the nodes didn't scale. old value: %d, new value %d\n", c.StartingNodes, c.EndNodes)
	}

	log.Printf("Completed Test: %s\n", c.Test.Name)

	c.Tracker.Completed = true
	return nil
}

// countNodes returns the current number of nodes in the cluster.
func countNodes(client *kubernetes.Clientset) int {
	allNodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return len(allNodes.Items)
}

// checkForUpdatedDeployment gets the Deployment, checks for an update and will repeat this until it's detected.
// Once detected the scaling test value for deployment is updated
func checkForUpdatedDeployment(c *ScalingTest, target int32) error {
	deploy, err := c.Deployment.Get()
	if err != nil {
		return err
	}

	// Check if UnavailableReplicas is 0, if it is then we're waiting for the change to be reflected in the cluster before continuing
	for target != deploy.Status.Replicas {
		log.Printf("Waiting for Deployment to update\n")
		time.Sleep(500 * time.Millisecond)

		deploy, err = c.Deployment.Get()
		if err != nil {
			return err
		}
	}

	// Set the Deployment to the updated resources
	c.Deployment.Deployment, err = c.Deployment.Get()
	if err != nil {
		return err
	}
	return nil
}
