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

package helm

import (
	"github.com/drewbernetes/dogkat/pkg/util/options"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strconv"
)

type ChartValues struct {
	Core    CoreValues    `yaml:"core"`
	Gpu     GPUValues     `yaml:"gpu"`
	Ingress IngressValues `yaml:"ingress"`
}

type CoreValues struct {
	Enabled  bool        `yaml:"enabled"`
	Nginx    NginxValues `yaml:"nginx"`
	Postgres PGValues    `yaml:"postgres"`
}

type NginxValues struct {
	Resources ResourceValues `yaml:"resources"`
}

type ResourceValues struct {
	Requests RequestValues `yaml:"requests"`
}

type RequestValues struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type PGValues struct {
	StatefulSet StatefulSetValues `yaml:"statefulSet"`
}

type StatefulSetValues struct {
	PersistentData PDValues `yaml:"persistentData"`
}

type PDValues struct {
	Enabled          bool   `yaml:"enabled"`
	StorageClassName string `yaml:"storageClassName"`
}

type GPUValues struct {
	Enabled      bool `yaml:"enabled"`
	NumberOfGPUs int  `yaml:"numberOfGPUs"`
}

type IngressValues struct {
	Enabled     bool              `yaml:"enabled"`
	Annotations map[string]string `yaml:"annotations"`
	ClassName   string            `yaml:"className"`
	Host        string            `yaml:"host"`
	Tls         []TLSValues       `yaml:"tls"`
}

type TLSValues struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

func setCoreValues(o options.CoreOptions) CoreValues {
	v := CoreValues{
		Enabled: true,
		Nginx: NginxValues{
			Resources: ResourceValues{
				Requests: RequestValues{
					Cpu:    o.CPU,
					Memory: o.Memory,
				},
			},
		},
	}

	if o.StorageClass != "" {
		v.Postgres = PGValues{
			StatefulSet: StatefulSetValues{
				PersistentData: PDValues{
					Enabled:          true,
					StorageClassName: o.StorageClass,
				},
			},
		}
	}

	return v
}

func setGPUValues(o options.GPUOptions) GPUValues {
	gpus, err := strconv.Atoi(o.NumberOfGPUs)
	if err != nil {
		return GPUValues{
			Enabled:      true,
			NumberOfGPUs: 1,
		}
	}
	return GPUValues{
		Enabled:      true,
		NumberOfGPUs: gpus,
	}
}

func setIngressValues(o options.IngressOptions) IngressValues {
	v := IngressValues{
		Enabled:     true,
		Annotations: o.Annotations,
		ClassName:   o.IngressClass,
		Host:        o.Host,
		Tls: []TLSValues{
			{
				Hosts:      []string{o.TLSHost},
				SecretName: o.TLSSecretName,
			},
		},
	}

	return v
}
