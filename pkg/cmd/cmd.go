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

package cmd

import (
	"github.com/drewbernetes/dogkat/pkg/cmd/delete"
	"github.com/drewbernetes/dogkat/pkg/cmd/util/config"
	"github.com/drewbernetes/dogkat/pkg/cmd/validate"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// init is auto run by cobra - all commands should be added here.
func newRootCommand() *cobra.Command {
	cobra.OnInitialize(config.InitConfig)

	cmd := &cobra.Command{
		Use: "dogkat",
		Long: `Deploys resources to allow End-2-End testing to be conducted.
It can be used to test most elements of a cluster to ensure consistent stability and functionality.
Documentation is available here: https://github.com/drewbernetes/dogkat/blob/main/README.md`,
	}

	configFlags := genericclioptions.NewConfigFlags(true)

	commands := []*cobra.Command{
		validate.NewValidateCommand(configFlags),
		delete.NewDeleteCommand(configFlags),
		NewVersionCmd(),
	}

	// Add default k8s flags
	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(commands...)
	return cmd
}

func Generate() *cobra.Command {
	return newRootCommand()
}
