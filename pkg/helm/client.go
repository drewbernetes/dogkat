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

package helm

import (
	"fmt"
	"github.com/eschercloudai/dogkat/pkg/constants"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"log"
	"os"
	"strings"
	"time"
)

type Client struct {
	KubeClient    *kubernetes.Clientset
	Configuration *action.Configuration
	Username      string
	Password      string
	Settings      *cli.EnvSettings
	Url           string
}

// NewClient registers and configures a new kubernetes and helm client
func NewClient(namespace string) (*Client, error) {
	// Configure settings for Helm
	settings := cli.New()
	settings.SetNamespace(namespace)
	helmDriver := os.Getenv("HELM_DRIVER")

	conf := new(action.Configuration)
	if err := conf.Init(settings.RESTClientGetter(), settings.Namespace(), helmDriver, log.Printf); err != nil {
		return nil, fmt.Errorf("chart init error: %s\n", err)
	}

	// Generate a new Client
	client := &Client{
		Configuration: conf,
		Settings:      settings,
		Url:           constants.RepoURL,
	}

	return client, nil
}

func (c *Client) ChartDeployed() (*release.Release, error) {
	client := action.NewGet(c.Configuration)

	return client.Run(constants.ReleaseName)
}

func (c *Client) PullChart(chartVersion, path string) error {
	conf := action.WithConfig(c.Configuration)
	a := action.NewPullWithOpts(conf)
	a.Settings = c.Settings
	a.RepoURL = c.Url
	a.Version = chartVersion
	a.DestDir = fmt.Sprintf("/%s", strings.Split(path, "/")[1])

	_, err := a.Run(constants.ChartName)
	if err != nil {
		return err
	}

	return nil
}

// Install the chart as a release onto the cluster.
func (c *Client) Install(chart *Chart) (*release.Release, error) {

	client := action.NewInstall(c.Configuration)
	client.CreateNamespace = true
	client.ReleaseName = constants.ReleaseName
	client.Namespace = c.Settings.Namespace()

	release, err := client.Run(chart.Chart, structToMap(chart.Values))
	if err != nil {
		return nil, err
	}

	for i := 0; release.Info.Status.IsPending(); i++ {
		log.Println("Chart is deploying")
		time.Sleep(time.Second * 10)

		if i > 12 {
			break
		}
	}

	log.Printf("Successfully installed release: %s\n", release.Name)
	return release, nil
}

func (c *Client) Uninstall() error {
	client := action.NewUninstall(c.Configuration)
	client.Wait = true
	client.IgnoreNotFound = true

	_, err := client.Run(constants.ReleaseName)
	if err != nil {
		return err
	}

	log.Printf("successfully removed release: %s\n", constants.ReleaseName)
	return nil
}
