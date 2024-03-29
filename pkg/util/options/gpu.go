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

package options

import (
	"github.com/spf13/viper"
)

type GPUOptions struct {
	Enabled            bool
	NumberOfGPUs       string
	NodeLabelSelectors map[string]string
}

func newGPUOptions() GPUOptions {
	return GPUOptions{
		Enabled:            viper.GetBool("gpu.enabled"),
		NumberOfGPUs:       viper.GetString("gpu.numberOfGPUs"),
		NodeLabelSelectors: viper.GetStringMapString("gpu.nodeLabelSelectors"),
	}
}
