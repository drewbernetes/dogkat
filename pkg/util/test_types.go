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

package util

import (
	"fmt"
	"github.com/drewbernetes/dogkat/pkg/constants"
)

type TestTypes struct {
	Core    bool
	Ingress bool
	GPU     bool
}

func (t *TestTypes) GetType() string {
	test := ""
	if t.Core {
		if test == "" {
			test = constants.TestCore
		} else {
			test = fmt.Sprintf("%s_%s", test, constants.TestCore)
		}
	}
	if t.Ingress {
		if test == "" {
			test = constants.TestIngress
		} else {
			test = fmt.Sprintf("%s_%s", test, constants.TestIngress)
		}
	}
	if t.GPU {
		if test == "" {
			test = constants.TestGPU
		} else {
			test = fmt.Sprintf("%s_%s", test, constants.TestGPU)
		}
	}
	if t.Core && t.Ingress && t.GPU {
		test = constants.TestAll
	}

	return test
}
