/*
Copyright 2022 EscherCloud.
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
	"github.com/eschercloudai/k8s-e2e-tester/pkg/tracing"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// TestIngress is called to validate all hosts within an Ingress resource.
func TestIngress(host, pushGateway string) error {
	tracer := tracing.Duration{JobName: "e2e_workloads", PushURL: pushGateway}
	tracer.SetupMetricsGatherer("dogkat_test_ingress_time_to_live_duration_seconds", "The time it takes for an ingress resource to become accessible")
	tracer.Start()

	err := testHostEndpoint(host, 30)
	if err != nil {
		return fmt.Errorf("cannot continue, %s", err.Error())
	}
	tracer.CompleteGathering()

	return nil
}

// testHostEndpoint will check the endpoint of an individual host in an Ingress for a valid 200 response.
func testHostEndpoint(host string, retries int) error {
	delay := time.Second * 20
	time.Sleep(delay)
	if retries <= 0 {
		return errors.New("reached the limit for checks")
	}

	resp, err := http.Get(strings.Join([]string{"https", host}, "://"))
	if err != nil {
		if strings.Contains(err.Error(), strings.ToLower("no such host")) {
			log.Printf("Dns is not resolving for %s - retrying in %s seconds\n", host, delay)
			return testHostEndpoint(host, retries-1)
		} else if strings.Contains(err.Error(), "x509: certificate") {
			log.Printf("There is a certificate error for %s - Have you got a problem with cert-manager or external DNS? Retrying in %s seconds\n", host, delay)
			return testHostEndpoint(host, retries-1)
		} else if strings.Contains(err.Error(), "No address associated with hostname") {
			log.Printf("Address error for %s - retrying in %s seconds\n", host, delay)
			return testHostEndpoint(host, retries-1)
		} else {
			return err
		}
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("status was not 200")
	}
	var result struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return err
	}

	log.Printf("Response from the page was: %s\n", result.Data)
	return nil
}
