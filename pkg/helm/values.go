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
	PHP      PHPValues   `yaml:"php"`
	Postgres PGValues    `yaml:"postgres"`
}

type RepoTagValues struct {
	Repo string `yaml:"repo"`
	Tag  string `yaml:"tag"`
}

type NginxValues struct {
	Image         RepoTagValues  `yaml:"image"`
	ExporterImage RepoTagValues  `yaml:"exporterImage"`
	Resources     ResourceValues `yaml:"resources"`
}

type PHPValues struct {
	Image RepoTagValues `yaml:"image"`
}

type ResourceValues struct {
	Requests RequestValues `yaml:"requests"`
}

type RequestValues struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type PGValues struct {
	Image       RepoTagValues     `yaml:"image"`
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
	TLS         []TLSValues       `yaml:"tls"`
}

type TLSValues struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

func setCoreValues(o options.CoreOptions) *CoreValues {
	if o.Enabled {
		v := &CoreValues{
			Enabled: o.Enabled,
			Nginx: NginxValues{
				Image: RepoTagValues{
					Repo: o.Nginx.Repo,
					Tag:  o.Nginx.Tag,
				},
				ExporterImage: RepoTagValues{
					Repo: o.NginxExporter.Repo,
					Tag:  o.NginxExporter.Tag,
				},
				Resources: ResourceValues{
					Requests: RequestValues{
						Cpu:    o.CPU,
						Memory: o.Memory,
					},
				},
			},
			PHP: PHPValues{
				Image: RepoTagValues{
					Repo: o.PHP.Repo,
					Tag:  o.PHP.Tag,
				},
			},
			Postgres: PGValues{
				Image: RepoTagValues{
					Repo: o.Postgres.Repo,
					Tag:  o.Postgres.Tag,
				},
			},
		}

		if o.StorageClass != "" {
			v.Postgres.StatefulSet = StatefulSetValues{
				PersistentData: PDValues{
					Enabled:          true,
					StorageClassName: o.StorageClass,
				},
			}
		}
		return v
	}

	return nil
}

func setGPUValues(o options.GPUOptions) *GPUValues {
	if o.Enabled {
		gpus, err := strconv.Atoi(o.NumberOfGPUs)
		if err != nil {
			return &GPUValues{
				Enabled:      o.Enabled,
				NumberOfGPUs: 1,
			}
		}
		return &GPUValues{
			Enabled:      o.Enabled,
			NumberOfGPUs: gpus,
		}
	}

	return nil
}

func setIngressValues(o options.IngressOptions) *IngressValues {
	if o.Enabled {
		v := &IngressValues{
			Enabled:     o.Enabled,
			Annotations: o.Annotations,
			ClassName:   o.IngressClass,
			Host:        o.Host,
			TLS: []TLSValues{
				{
					Hosts:      []string{o.TLSHost},
					SecretName: o.TLSSecretName,
				},
			},
		}

		return v
	}

	return nil
}
