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

package testhelpers

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	// Mux is a multiplexer that can be used to register handlers.
	Mux *http.ServeMux

	// Server is an in-memory HTTP server for testing.
	Server *httptest.Server
)

// GenerateCloudsFile creates a cloud file for testing
func GenerateCloudsFile() error {
	var err error
	f, err := os.Create(CloudPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(cloud))
	if err != nil {
		return err
	}

	err = os.Setenv("OS_CLIENT_CONFIG_FILE", CloudPath)
	if err != nil {
		return err
	}

	return err
}

// SetupPersistentPortHTTP prepares the Mux and Server listening specific port.
func SetupPersistentPortHTTP(t *testing.T, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Errorf("Failed to listen to 127.0.0.1:%d: %s", port, err)
	}
	Mux = http.NewServeMux()
	Server = httptest.NewUnstartedServer(Mux)
	Server.Listener = l
	Server.Start()
}

// TeardownHTTP releases HTTP-related resources.
func TeardownHTTP() {
	Server.Close()
}

// Endpoint returns a fake endpoint that will actually target the Mux.
func Endpoint() string {
	return Server.URL + "/"
}

// ServiceClient returns a generic service client for use in tests.
func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{TokenID: TokenID},
		Endpoint:       Endpoint(),
	}
}

func CommonServiceClient() *gophercloud.ServiceClient {
	sc := ServiceClient()
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc
}
