package cmd

import (
	"context"
	"e2e-test/resources"
	test_cases "e2e-test/test-cases"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

func runCoreTests(valuesFile string) {
	//Parse the values file for Helm
	values := parseValues(valuesFile)
	if values == nil {
		log.Fatalln("no values provided or couldn't parse them")
	}

	//Deploy the Helm Chart
	actionCfg, rel, err := deployChart(namespaceName, values)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Read the manifests of the deployed Helm Chart into respective objects
	allManifests := strings.Split(rel.Manifest, "---")
	var res []resources.ApiResource

	for _, v := range allManifests {
		if v == "" {
			continue
		}
		res = append(res, parseResource(v))
	}

	//Thread the tests to run in parallel
	checksCompleted := make(chan struct {
		Ready    bool
		Resource resources.ApiResource
	})

	defer close(checksCompleted)
	for _, r := range res {
		if r == nil {
			continue
		}
		go test_cases.CoreWorkloadChecks(r, checksCompleted)
	}

	for _, _ = range res {
		<-checksCompleted
	}

	log.Println("**ALL RESOURCES ARE DEPLOYED**")

	scalingTested := false
	for _, r := range res {
		test_cases.ScalingValidation(r)
		if r.GetResourceKind() == "Deployment" || r.GetResourceKind() == "StatefulSet" {
			if !scalingTested {
				scalingTested = test_cases.RunScalingTest(r, clientsets)
			} else {
				break
			}
		}
	}

	//Uninstall the Helm Chart
	resp, err := uninstallChart(actionCfg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if resp != nil {
		log.Println(resp.Info)
	}

	// Remove the namespace
	err = clientsets.K8S.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		log.Printf("The namespace could not be removed. Please remove the namespace %s manually\n", namespaceName)
		return
	} else {
		log.Printf("The namespace %s has been removed\n", namespaceName)
	}

	//TODO: Remove the namespace
	log.Println("**ALL AVAILABLE TESTS COMPLETED**")
	log.Println("See logs above for results")
}
