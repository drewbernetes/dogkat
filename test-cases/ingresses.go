package test_cases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	networkingv1 "k8s.io/api/networking/v1"
	"log"
	"net/http"
	"strings"
	"time"
)

// testIngress is called to validate all hosts within an Ingress resource.
// Deprecated: no longer in use
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

// testHostEndpoint will check the endpoint of an individual host in an Ingress for a valid 200 response.
// Deprecated: no longer in use
func testHostEndpoints(host string, counter int) error {
	delay := time.Second * 20
	time.Sleep(delay)
	if counter >= 20 {
		return errors.New("reached the limit for checks")
	}

	resp, err := http.Get(strings.Join([]string{"https", host}, "://"))
	if err != nil {
		if strings.Contains(err.Error(), strings.ToLower("no such host")) {
			log.Printf("Dns is not resolving for %s - retrying in %s seconds\n", host, delay)
			return testHostEndpoints(host, counter+1)
		} else if strings.Contains(err.Error(), "x509: certificate") {
			log.Printf("There is a certificate error for %s - Have you got a problem with cert-manager or external DNS? Retrying in %s seconds\n", host, delay)
			return testHostEndpoints(host, counter+1)
		} else if strings.Contains(err.Error(), "No address associated with hostname") {
			log.Printf("Address error for %s - retrying in %s seconds\n", host, delay)
			return testHostEndpoints(host, counter+1)
		} else {
			return err
		}
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

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
