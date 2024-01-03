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
	"errors"
	"fmt"
	"github.com/eschercloudai/dogkat/pkg/constants"
	"github.com/eschercloudai/dogkat/pkg/util"
	"github.com/eschercloudai/dogkat/pkg/util/options"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
	"reflect"
	"strings"
	"unicode"
)

type Chart struct {
	Chart  *chart.Chart
	Values ChartValues
}

// NewChart unpacks the tgz making the Chart accessible for installation and supplied a Chart back.
func NewChart(cl *Client, testType util.TestTypes, o options.Options) (*Chart, error) {
	// Download Chart
	downloadPath := fmt.Sprintf("%s/%s-%s.tgz", "/tmp", constants.ChartName, o.Version)
	_, err := os.Stat(downloadPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = cl.PullChart(o.Version, downloadPath)
			if err != nil {
				return nil, err
			}
		}
	}

	// Load the chart
	c := &Chart{}
	unpack, err := loader.Load(downloadPath)
	if err != nil {
		return nil, err
	}
	c.Chart = unpack

	// Set the chart values into the struct
	c.loadOptionsToValues(testType, o)
	return c, nil
}

// loadOptionsToValues takes the inputted options via the config file and turns them into valid chart values.
func (c *Chart) loadOptionsToValues(testType util.TestTypes, o options.Options) {
	c.Values = ChartValues{}
	if testType.Core {
		c.Values.Core = setCoreValues(o.CoreOptions)
	}
	if testType.Ingress {
		c.Values.Ingress = setIngressValues(o.IngressOptions)
	}
	if testType.GPU {
		c.Values.Gpu = setGPUValues(o.GPUOptions)
	}
}

// lowerFirst just lowers the first letter in the supplied string.
func lowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(unicode.ToLower(rune(s[0])))) + s[1:]
}

// structToMap takes the chart.Values and turns them into a yaml file which can then be supplied during the helm install.
func structToMap(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	result := make(map[string]interface{})

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Type().Field(i)
		fieldName := lowerFirst(field.Tag.Get("yaml"))

		if fieldName == "" {
			fieldName = lowerFirst(field.Name)
		}

		fieldValue := objValue.Field(i)

		switch fieldValue.Kind() {
		case reflect.Struct:
			result[fieldName] = structToMap(fieldValue.Interface())
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.Struct {
				var sliceResult []map[string]interface{}
				for j := 0; j < fieldValue.Len(); j++ {
					sliceResult = append(sliceResult, structToMap(fieldValue.Index(j).Interface()))
				}
				result[fieldName] = sliceResult
			} else {
				// Handle the case where the field is a slice of basic types (e.g., []string)
				var sliceResult []interface{}
				for j := 0; j < fieldValue.Len(); j++ {
					sliceResult = append(sliceResult, fieldValue.Index(j).Interface())
				}
				result[fieldName] = sliceResult
			}
		default:
			result[fieldName] = fieldValue.Interface()
		}
	}

	return result
}
