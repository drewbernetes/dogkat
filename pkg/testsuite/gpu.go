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
	"context"
	"errors"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v1 "k8s.io/api/core/v1"
	"log"
	"strings"
)

func TestGPU(pod *coreworkloads.Pod) error {
	log.Printf("checking pod logs in %s for PASSED status\n", pod.Resource.Name)
	logs := pod.Client.CoreV1().Pods(pod.Resource.Namespace).GetLogs(pod.Resource.Name, &v1.PodLogOptions{})

	logRaw, err := logs.DoRaw(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	if !strings.Contains(string(logRaw), "Test PASSED") {
		return errors.New("the test failed to complete - check the logs for more information")
	}
	log.Println("Test passed")
	return nil
}
