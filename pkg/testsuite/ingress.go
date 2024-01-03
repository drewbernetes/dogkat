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

package testsuite

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eschercloudai/dogkat/pkg/helm"
	"github.com/eschercloudai/dogkat/pkg/tracing"
	"github.com/eschercloudai/dogkat/pkg/workloads"
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
}

func NewEndpointTest(i *workloads.Ingress, c *helm.Client) *HostNamePropagationTest {
	name := "ingress_hostname_test"
	description := "The time it takes for an ingress resource to become accessible externally"
	t := NewTest(c, name, description)

	return &HostNamePropagationTest{
		Test:    t,
		Ingress: i,
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
	retries := 20

	for !resolved {
		if retries <= 0 {
			return fmt.Errorf("couldn't resolve the endpoint - maybe it's not propagated yet or hasn't been created\n")
		}

		next := func() {
			retries--
			time.Sleep(time.Second * 5)
		}

		resp, err := http.Get(i.Hostname)
		if err != nil {
			// Try and catch some common errors that may resolve
			if strings.Contains(err.Error(), strings.ToLower("no such host")) {
				log.Printf("dns not propagated for %s\n", i.Hostname)
				next()
				continue
			} else if strings.Contains(err.Error(), "x509: certificate") {
				log.Printf("There is a certificate error for %s - Have you got a problem with cert-manager or external DNS?\n", i.Hostname)
				next()
				continue
			} else if strings.Contains(err.Error(), "No address associated with hostname") {
				log.Printf("Address error for %s\n", i.Hostname)
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
	return nil
}
