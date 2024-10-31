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

package constants

import (
	"fmt"
	"os"
	"path"
)

var (

	// Version is the release version
	Version string

	// Revision is the git revision set via the Makefile
	Revision string

	// Application is the name of the application set via the Makefile
	Application = path.Base(os.Args[0])
)

const (
	ChartName   = "dogkat"
	ReleaseName = "dogkat-testing"
	RepoURL     = "https://drewbernetes.github.io/dogkat"

	TestCore    = "core"
	TestGPU     = "gpu"
	TestIngress = "ingress"
	TestAll     = "all"

	GPUName   = "gpu-test"
	NginxName = "nginx-e2e"
	PGSqlName = "web-database"

	MetricsPrefix  = "dogkat_test"
	MetricsJobName = "dogkat_workloads"
)

func VersionPrint() {
	fmt.Printf("%s/%s (revision/%s)\n", Application, Version, Revision)
}
