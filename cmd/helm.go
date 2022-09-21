package cmd

import (
	errors "errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	chartVersion = "0.0.6"
	chartPath    = "/tmp/chart.tgz"
	releaseName  = "e2e-test"
)

// fetchHelmChart will fetch the helm chart from the GitHub repo.
func fetchHelmChart(filepath string) error {
	//If the tar.gz exists, skip this step.
	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	chartURI := fmt.Sprintf("https://github.com/drew-viles/helm-charts/releases/download/End-2-End-Testing-%s/End-2-End-Testing-%s.tgz", chartVersion, chartVersion)
	resp, err := http.Get(chartURI)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("could not pull data from GitHub, status code is %d", resp.StatusCode)
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// parseValues will parse the provided values file.
func parseValues(valuesFile string) map[string]interface{} {
	values := map[string]interface{}{}
	reader, err := os.ReadFile(valuesFile)
	if err != nil {
		log.Printf("Reading values file error %s\n", err)
		return nil
	}

	err = yaml.Unmarshal(reader, values)
	if err != nil {
		log.Printf("Unmarshal values file error %s\n", err)
		return nil
	}

	return values
}

// unpackChart unpacks the tgz making the Chart accessible for installation.
func unpackChart() (*chart.Chart, error) {
	err := fetchHelmChart(chartPath)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

// initHelm prepares the helm go client for interaction with the cluster.
func initHelm(namespace string) (*chart.Chart, *action.Configuration, error) {
	settings := cli.New()
	settings.SetNamespace(namespace)

	chart, err := unpackChart()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("chart unpack error: %s\n", err))
	}

	//Set the driver: https://helm.sh/docs/topics/advanced/
	os.Setenv("HELM_DRIVER", "secret")

	actionConfig := new(action.Configuration)
	if err = actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, nil, errors.New(fmt.Sprintf("chart init error: %s\n", err))
	}
	return chart, actionConfig, nil
}

// installChart deploys the chart to the cluster.
func installChart(releaseName, namespace string, chart *chart.Chart, values map[string]interface{}, actionConfig *action.Configuration) (*release.Release, error) {
	client := action.NewInstall(actionConfig)
	client.CreateNamespace = true
	client.ReleaseName = releaseName
	client.Namespace = namespace

	release, err := client.Run(chart, values)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("chart run error: %s\n", err))
	}

	for i := 0; release.Info.Status.IsPending(); i++ {
		log.Println("Chart is deploying")
		time.Sleep(time.Second * 10)

		if i > 12 {
			break
		}
	}

	log.Printf("Successfully installed release: %s", release.Name)
	return release, nil
}

// deployChart checks if the chart is already deployed and if not deploys it.
// If the chart exists already it'll return that release.
func deployChart(namespace string, values map[string]interface{}) (*action.Configuration, *release.Release, error) {
	chart, actionConfig, err := initHelm(namespace)
	if err != nil {
		return nil, nil, err
	}

	release := isChartDeployed(releaseName, actionConfig)

	if release != nil {
		return actionConfig, release, nil
	}

	rel, err := installChart(releaseName, namespace, chart, values, actionConfig)
	return actionConfig, rel, err
}

// uninstallChart removes the chart to ensure no resources are left hanging around
func uninstallChart(actionCfg *action.Configuration) (*release.UninstallReleaseResponse, error) {
	release := isChartDeployed(releaseName, actionCfg)

	if release == nil {
		return nil, nil
	}

	client := action.NewUninstall(actionCfg)
	resp, err := client.Run(releaseName)
	if err != nil {
		return nil, err
	}

	err = os.Remove(chartPath)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// isChartDeployed checks if the chart is deployed and if so, returns it.
func isChartDeployed(releaseName string, actionConfig *action.Configuration) *release.Release {
	client := action.NewList(actionConfig)
	// Only list deployed
	client.Deployed = true
	client.AllNamespaces = false

	results, err := client.Run()
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(1)
	}

	for _, rel := range results {
		if rel.Name == releaseName {
			return rel
		}
	}
	return nil
}
