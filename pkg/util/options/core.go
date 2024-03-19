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

type ImageOptions struct {
	Repo string
	Tag  string
}

type CoreOptions struct {
	Enabled       bool
	CPU           string
	Memory        string
	StorageClass  string
	ScaleTo       int32
	Nginx         ImageOptions
	NginxExporter ImageOptions
	PHP           ImageOptions
	Postgres      ImageOptions
}

func newCoreOptions() CoreOptions {
	return CoreOptions{
		Enabled:      viper.GetBool("core.enabled"),
		CPU:          viper.GetString("core.cpu"),
		Memory:       viper.GetString("core.memory"),
		StorageClass: viper.GetString("core.storageClassName"),
		ScaleTo:      viper.GetInt32("core.scaleTo"),
		Nginx: ImageOptions{
			Repo: viper.GetString("core.nginx.image.repo"),
			Tag:  viper.GetString("core.nginx.image.tag"),
		},
		NginxExporter: ImageOptions{
			Repo: viper.GetString("core.nginxExporter.image.repo"),
			Tag:  viper.GetString("core.nginxExporter.image.tag"),
		},
		PHP: ImageOptions{
			Repo: viper.GetString("core.php.image.repo"),
			Tag:  viper.GetString("core.php.image.tag"),
		},
		Postgres: ImageOptions{
			Repo: viper.GetString("core.postgres.image.repo"),
			Tag:  viper.GetString("core.postgres.image.tag"),
		},
	}
}
