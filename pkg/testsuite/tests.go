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
	"fmt"
	"github.com/eschercloudai/dogkat/pkg/helm"
	"github.com/eschercloudai/dogkat/pkg/tracing"
)


type TestCase interface {
	Init()
	Run() error
	Validate() error
}

type Test struct {
	Client      *helm.Client
	Tracing     *tracing.Duration
	Name        string
	Description string
}

func NewTest(c *helm.Client, name, description string) Test {
	t := Test{
		Client:      c,
		Name:        name,
		Description: description,
	}

	m := tracing.Gatherer()
	if m != nil {
		n := fmt.Sprintf("%s_duration_seconds", name)
		t.Tracing = tracing.NewCollector(m.PushGateway, n, description)
	}

	return t
}
