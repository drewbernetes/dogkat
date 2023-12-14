{{- define "e2e-testing.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "e2e-testing.labels" -}}
helm.sh/chart: {{ include "e2e-testing.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: "Helm"
{{- end }}

{{- define "e2e-testing.annotations" -}}
meta.helm.sh/release-namespace: {{ .Release.Namespace | quote }}
{{- end }}

{{- define "e2e-testing.nginx.labels" -}}
app: nginx-e2e
app.kubernetes.io/instance: nginx-e2e
app.kubernetes.io/name: nginx-e2e
{{- end }}

{{- define "e2e-testing.nginx.annotations" -}}
meta.helm.sh/release-name: nginx-e2e
{{- end }}

{{- define "e2e-testing.psql.labels" -}}
app: postgres
app.kubernetes.io/instance: postgres
app.kubernetes.io/name: postgres
{{- end }}

{{- define "e2e-testing.psql.annotations" -}}
meta.helm.sh/release-name: postgres
{{- end }}

{{- define "e2e-testing.gpu.labels" -}}
app: gpu
app.kubernetes.io/instance: gpu
app.kubernetes.io/name: gpu
{{- end }}

{{- define "e2e-testing.gpu.annotations" -}}
meta.helm.sh/release-name: gpu
{{- end }}
