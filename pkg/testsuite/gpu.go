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
	"context"
	"errors"
	"github.com/drewbernetes/dogkat/pkg/helm"
	"github.com/drewbernetes/dogkat/pkg/workloads"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"log"
	"strings"
)

type VectorTest struct {
	Test
	Pod     *workloads.Pod
	Logs    *rest.Request
	Tracker *TestTracker
}

// NewVectorTest allows the logs for the vector_add program to be checked to make sure it passed.
func NewVectorTest(p *workloads.Pod, c *helm.Client) *VectorTest {
	name := "gpu_vector_test"
	description := "Times the testing of the GPU resource to bring a GPU node in and run a vector add program"
	t := NewTest(c, name, description)

	return &VectorTest{
		Test: t,
		Pod:  p,
		Tracker: &TestTracker{
			Name: name,
			Description: `The GPU tester will deploy an NVIDIA CUDA Pod which runs a simple vector add. 
The success of this Pod will determine if the GPU is available and functioning. 
This will not confirm the validity of licenses, only functionality of the GPU at the time of the test running. If a valid
license is not deployed the performance of the GPU will degrade over time.`,
			Completed: false,
		},
	}
}

// Init prepares the test with initial conditions so that a starting point can be used to validate later
func (v *VectorTest) Init() {}

// Run the actual test
func (v *VectorTest) Run() error {
	log.Printf("Running Test: %s\n", v.Test.Name)
	v.Logs = v.Pod.PodInterface.GetLogs(v.Pod.Pod.Name, &v1.PodLogOptions{})
	return nil
}

// Validate compares the end result of the test with the values set in the Init stage allowing a comparison to see if the test passed.
// There is no metrics gathering here because by the time the logs are checked, the chart will have deployed and the
// test run already so those metrics are more relevant than any that could be gathered here.
func (v *VectorTest) Validate() error {
	raw, err := v.Logs.DoRaw(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	if !strings.Contains(string(raw), "Test PASSED") {
		return errors.New("the test failed to complete - check the logs for more information")
	}
	log.Printf("Completed Test: %s\n", v.Test.Name)

	v.Tracker.Completed = true
	return nil
}
