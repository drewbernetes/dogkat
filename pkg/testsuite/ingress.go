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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/drewbernetes/dogkat/pkg/helm"
	"github.com/drewbernetes/dogkat/pkg/tracing"
	"github.com/drewbernetes/dogkat/pkg/workloads"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type HostNamePropagationTest struct {
	Test
	Ingress        *workloads.Ingress
	Hostname       string
	ResponseBody   []byte
	ResponseStatus int
	Tracker        *TestTracker
}

func NewEndpointTest(i *workloads.Ingress, c *helm.Client) *HostNamePropagationTest {
	name := "ingress_hostname_test"
	description := "The time it takes for an ingress resource to become accessible externally"
	t := NewTest(c, name, description)

	return &HostNamePropagationTest{
		Test:    t,
		Ingress: i,
		Tracker: &TestTracker{
			Name: name,
			Description: `The ingress test makes use of the core scaling test in that it deploys the web/database service however it also deploys an Ingress resource.
This test requires External DNS and Cert-Manager to be deployed to function correctly as it will use them to create the DNS record and generate the cert.`,
			Completed: false,
		},
	}
}

// Init prepares the test with initial conditions so that a starting point can be used to validate later
func (i *HostNamePropagationTest) Init(host string, useTLS bool) {
	prefix := "http"
	if useTLS {
		prefix = "https"
	}
	i.Hostname = fmt.Sprintf("%s://%s", prefix, host)
}

// Run the actual test
func (i *HostNamePropagationTest) Run() error {
	// Start metrics gathering
	if tracing.Gatherer().Enabled {
		i.Tracing.Start()
	}
	log.Printf("Running Test: %s\n", i.Test.Name)

	resolved := false
	retries := 300

	log.Println("Waiting up to 5 minutes for endpoint to respond")
	for !resolved {
		if retries <= 0 {
			return fmt.Errorf("couldn't resolve the endpoint - maybe it's not propagated yet or hasn't been created\n")
		}

		next := func() {
			retries--
			time.Sleep(time.Second * 1)
		}

		resp, err := http.Get(i.Hostname)
		if err != nil {
			// Try and catch some common errors that may resolve themselves
			if strings.Contains(err.Error(), strings.ToLower("no such host")) || strings.Contains(err.Error(), "x509: certificate") || strings.Contains(err.Error(), "No address associated with hostname") {
				next()
				continue
			} else {
				// There's something wrong and it ain't no common error!
				return err
			}
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		i.ResponseBody = bodyBytes
		i.ResponseStatus = resp.StatusCode
		resolved = true
	}

	if tracing.Gatherer().Enabled {
		if err := i.Tracing.CompleteGathering(); err != nil {
			return err
		}
	}
	return nil
}

// Validate compares the end result of the test with the values set in the Init stage allowing a comparison to see if the test passed.
func (i *HostNamePropagationTest) Validate() error {
	if i.ResponseStatus != http.StatusOK {
		return errors.New("status was not 200")
	}
	var result struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}

	err := json.Unmarshal(i.ResponseBody, &result)
	if err != nil {
		return err
	}

	log.Printf("Response from the page was: %s\n", result.Data)
	log.Printf("Completed Test: %s\n", i.Test.Name)

	i.Tracker.Completed = true
	return nil
}
