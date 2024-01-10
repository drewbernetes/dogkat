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

package workloads

import (
	"fmt"
	"log"
	"time"
)

type Resource interface {
	Name() string
	Kind() string
	IsReady() bool
}

// CheckReadyState will check every 1 second for 5 minutes to see if the resource is ready.
func CheckReadyState(r Resource) error {
	ready := r.IsReady()
	readyCount := 0
	for !ready {
		if readyCount == 300 {
			return fmt.Errorf("%s %s isn't ready\n", r.Name(), r.Kind())
		}

		time.Sleep(1 * time.Second)
		ready = r.IsReady()
		readyCount++
	}

	log.Printf("%s %s is ready\n", r.Name(), r.Kind())
	return nil
}
