package testsuite

import (
	"context"
	"errors"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	v1 "k8s.io/api/core/v1"
	"log"
	"strings"
)

func TestGPU(pod *coreworkloads.Pod) error {
	log.Printf("checking pod logs in %s for PASSED status\n", pod.Resource.Name)
	logs := pod.Client.CoreV1().Pods(pod.Resource.Namespace).GetLogs(pod.Resource.Name, &v1.PodLogOptions{})

	logRaw, err := logs.DoRaw(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	if !strings.Contains(string(logRaw), "Test PASSED") {
		return errors.New("the test failed to complete - check the logs for more information")
	}
	log.Println("Test passed")
	return nil
}
