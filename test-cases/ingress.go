package test_cases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1Typed "k8s.io/client-go/kubernetes/typed/networking/v1"
	"log"
	"net/http"
	"strings"
	"time"
)

type IngressResource struct {
	Client   v1Typed.IngressInterface
	Resource *v1.Ingress
	Error    error
}

func (r *IngressResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *IngressResource) GetError() error {
	return r.Error
}

func (r *IngressResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *IngressResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *IngressResource) IsReady() bool {
	for _, v := range r.Resource.Status.LoadBalancer.Ingress {
		if v.Hostname == "" && v.IP == "" {
			return false
		}
	}
	return true
}

func (r *IngressResource) GetClient(namespace string) {
	r.Client = clientset.NetworkingV1().Ingresses(namespace)
}

func (r *IngressResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *IngressResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *IngressResource) Update() {
}
func (r *IngressResource) Delete() {
}

func testIngress(hosts []networkingv1.IngressTLS) error {
	for _, v := range hosts {
		for _, host := range v.Hosts {
			err := testHostEndpoints(host, 0)
			if err != nil {
				return errors.New(fmt.Sprintf("cannot continue, %s", err.Error()))
			}
		}
	}
	return nil
}

func testHostEndpoints(host string, counter int) error {
	delay := time.Second * 20
	time.Sleep(delay)
	if counter >= 20 {
		return errors.New("reached the limit for checks")
	}

	resp, err := http.Get(strings.Join([]string{"https", host}, "://"))
	if err != nil {
		if strings.Contains(err.Error(), strings.ToLower("no such host")) {
			log.Printf("dns is not resolving for %s - retrying in %s seconds\n", host, delay)
			testHostEndpoints(host, counter+1)
		} else if strings.Contains(err.Error(), "x509: certificate") {
			log.Printf("There is a certificate error for %s - retrying in %s seconds\n", host, delay)
			testHostEndpoints(host, counter+1)
		} else if strings.Contains(err.Error(), "No address associated with hostname") {
			log.Printf("Address error for %s - retrying in %s seconds\n", host, delay)
			testHostEndpoints(host, counter+1)
		} else {
			return err
		}
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return errors.New("status was not 200")
	}
	var result struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Response from the page was: %s\n", result.Data)
	return nil
}
