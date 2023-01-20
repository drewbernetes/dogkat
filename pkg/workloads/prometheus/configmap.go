/*
Copyright 2022 EscherCloud.
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

package prometheus

import (
	"github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
)

// GenerateGrafanaDashboardConfigMap returns a configmap containing a dashboard for Grafana.
func GenerateGrafanaDashboardConfigMap(namespace string) *coreworkloads.ConfigMap {
	cm := &coreworkloads.ConfigMap{}
	cm.Generate(map[string]string{"namespace": namespace, "name": "grafana-dashboard", "label": "nginx-e2e"})
	cm.Resource.ObjectMeta.Labels["grafana_dashboard"] = "1"
	cm.Resource.Data = map[string]string{
		"run.sh": `|-
			{
			  "annotations": {
				"list": [
				  {
					"builtIn": 1,
					"datasource": {
					  "type": "grafana",
					  "uid": "-- Grafana --"
					},
					"enable": true,
					"hide": true,
					"iconColor": "rgba(0, 211, 255, 1)",
					"name": "Annotations & Alerts",
					"target": {
					  "limit": 100,
					  "matchAny": false,
					  "tags": [],
					  "type": "dashboard"
					},
					"type": "dashboard"
				  }
				]
			  },
			  "editable": true,
			  "fiscalYearStartMonth": 0,
			  "graphTooltip": 0,
			  "id": 61,
			  "links": [],
			  "liveNow": false,
			  "panels": [
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 0,
					"y": 0
				  },
				  "id": 6,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "pluginVersion": "9.0.1",
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "code",
					  "expr": "count(nginx_up)",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					}
				  ],
				  "title": "Nginx Online",
				  "transformations": [],
				  "type": "timeseries"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 12,
					"y": 0
				  },
				  "id": 12,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_handled",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					}
				  ],
				  "title": "Connections Handled",
				  "type": "timeseries"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "fieldConfig": {
					"defaults": {
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  },
					  "unit": "short"
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 0,
					"y": 8
				  },
				  "id": 10,
				  "options": {
					"colorMode": "value",
					"graphMode": "area",
					"justifyMode": "auto",
					"orientation": "vertical",
					"reduceOptions": {
					  "calcs": [
						"lastNotNull"
					  ],
					  "fields": "",
					  "values": false
					},
					"text": {
					  "titleSize": 1
					},
					"textMode": "value"
				  },
				  "pluginVersion": "9.0.1",
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_reading",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					},
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_writing",
					  "hide": false,
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "B"
					}
				  ],
				  "title": "Reading/Writing",
				  "type": "stat"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "description": "",
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 9,
					"w": 12,
					"x": 12,
					"y": 8
				  },
				  "id": 2,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_accepted",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					}
				  ],
				  "title": "E2E Connections Accepted",
				  "type": "timeseries"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 0,
					"y": 16
				  },
				  "id": 14,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "rate(nginx_http_requests_total[$__rate_interval])",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					}
				  ],
				  "title": "Total requests",
				  "type": "timeseries"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 12,
					"y": 17
				  },
				  "id": 8,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_waiting",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					},
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "hide": false,
					  "refId": "B"
					}
				  ],
				  "title": "Connections Waiting",
				  "type": "timeseries"
				},
				{
				  "datasource": {
					"type": "prometheus",
					"uid": "prometheus"
				  },
				  "description": "",
				  "fieldConfig": {
					"defaults": {
					  "color": {
						"mode": "palette-classic"
					  },
					  "custom": {
						"axisLabel": "",
						"axisPlacement": "auto",
						"barAlignment": 0,
						"drawStyle": "line",
						"fillOpacity": 0,
						"gradientMode": "none",
						"hideFrom": {
						  "legend": false,
						  "tooltip": false,
						  "viz": false
						},
						"lineInterpolation": "linear",
						"lineWidth": 1,
						"pointSize": 5,
						"scaleDistribution": {
						  "type": "linear"
						},
						"showPoints": "auto",
						"spanNulls": false,
						"stacking": {
						  "group": "A",
						  "mode": "none"
						},
						"thresholdsStyle": {
						  "mode": "off"
						}
					  },
					  "mappings": [],
					  "thresholds": {
						"mode": "absolute",
						"steps": [
						  {
							"color": "green",
							"value": null
						  },
						  {
							"color": "red",
							"value": 80
						  }
						]
					  }
					},
					"overrides": []
				  },
				  "gridPos": {
					"h": 8,
					"w": 12,
					"x": 12,
					"y": 25
				  },
				  "id": 4,
				  "options": {
					"legend": {
					  "calcs": [],
					  "displayMode": "list",
					  "placement": "bottom"
					},
					"tooltip": {
					  "mode": "single",
					  "sort": "none"
					}
				  },
				  "targets": [
					{
					  "datasource": {
						"type": "prometheus",
						"uid": "prometheus"
					  },
					  "editorMode": "builder",
					  "expr": "nginx_connections_active",
					  "legendFormat": "__auto",
					  "range": true,
					  "refId": "A"
					}
				  ],
				  "title": "E2E Connections Active",
				  "type": "timeseries"
				}
			  ],
			  "refresh": false,
			  "schemaVersion": 36,
			  "style": "dark",
			  "tags": [],
			  "templating": {
				"list": []
			  },
			  "time": {
				"from": "2022-08-05T14:18:22.671Z",
				"to": "2022-08-05T15:03:21.727Z"
			  },
			  "timepicker": {},
			  "timezone": "",
			  "title": "End2End Testing",
			  "uid": "LOWihQk4k",
			  "version": 5,
			  "weekStart": ""
			}`,
	}

	return cm
}
