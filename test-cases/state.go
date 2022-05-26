package test_cases

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"log"
	"strings"
)

//checkIfABIngressIsEnabled is an edge case detection feature.
//It recognises that an AB setup for ingress has been deployed, which is not recommended, and allows this program to compensate for a bad design.
func checkIfABIngressIsEnabled() bool {
	client := clientset.CoreV1().Namespaces()
	_, err := client.Get(context.TODO(), "ingress-private-a", metav1.GetOptions{})

	if err != nil && strings.Contains(err.Error(), strings.ToLower(field.ErrorTypeNotFound.String())) {
		return false
	}

	log.Println("ABIngress Detected - don't do this - get good instead")
	return true
}
