package main

// Deprecated: no longer in use
func runCoreTests(valuesFile string) {
	////Parse the values file for Helm
	//values := parseValues(valuesFile)
	//if values == nil {
	//	log.Fatalln("no values provided or couldn't parse them")
	//}
	//
	////Deploy the Helm Chart
	//actionCfg, rel, err := deployChart(cmd.namespaceName, values)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//
	////Read the manifests of the deployed Helm Chart into respective objects
	//allManifests := strings.Split(rel.Manifest, "---")
	//var res []resources.ApiResource
	//
	//for _, v := range allManifests {
	//	if v == "" {
	//		continue
	//	}
	//	res = append(res, cmd.parseResource(v))
	//}
	//
	////Thread the tests to run in parallel
	//checksCompleted := make(chan resources.ResourceReady)
	//
	//defer close(checksCompleted)
	//for _, r := range res {
	//	if r == nil {
	//		continue
	//	}
	//	go test_cases.CoreWorkloadChecks(r, checksCompleted)
	//}
	//
	//for _, _ = range res {
	//	<-checksCompleted
	//}
	//
	//log.Println("** ALL RESOURCES ARE DEPLOYED AND READY **")
	//
	//scalingTested := false
	//for _, r := range res {
	//	test_cases.ScalingValidation(r)
	//	//TODO: The RunScaling test need supdating to work with StatefulSets too.
	//	//if r.GetResourceKind() == "Deployment" || r.GetResourceKind() == "StatefulSet" {
	//	if r.GetResourceKind() == "Deployment" {
	//		if !scalingTested {
	//			scalingTested = test_cases.RunScalingTest(r, clientsets)
	//		} else {
	//			break
	//		}
	//	}
	//}
	//
	////Uninstall the Helm Chart
	//resp, err := uninstallChart(actionCfg)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//if resp != nil {
	//	log.Println(resp.Info)
	//}
	//
	//// Remove the namespace
	//err = clientsets.K8S.CoreV1().Namespaces().Delete(context.TODO(), cmd.namespaceName, metav1.DeleteOptions{})
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Printf("The namespace could not be removed. Please remove the namespace %s manually\n", cmd.namespaceName)
	//	return
	//} else {
	//	log.Printf("The namespace %s has been removed\n", cmd.namespaceName)
	//}
	//
	//log.Println("** ALL AVAILABLE TESTS COMPLETED **")
	//log.Println("See logs above for results")
}
