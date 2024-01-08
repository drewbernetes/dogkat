# dogkat

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.1](https://img.shields.io/badge/AppVersion-0.1.1-informational?style=flat-square)

End-2-End testing for GPUs and some core resources

## Installation

```shell
helm install https://eschercloudai.github.io/dogkat/dogkat --values values.yaml
```

## Updating the Chart

Make sure you've run and resolved any issues using the following as failures of these will cause the pipeline/actions to fail.
```
yamllint charts/dogkat/Chart.yaml  --config-file .github/ct.yaml
yamllint charts/dogkat/values.yaml  --config-file .github/ct.yaml

helm-docs
```

## Configuration

The following table lists the configurable parameters of the chart and the default values.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| core.enabled | bool | `false` |  |
| core.nginx.image.repo | string | `"nginx"` | The repo to be used |
| core.nginx.image.tag | string | `"1.25-alpine"` | The tag to be used |
| core.nginx.resources | object | `{}` |  |
| core.nginx.serviceAccountName | string | `"nginx"` |  |
| core.php.image.repo | string | `"drewviles/php-pdo"` | The repo to be used |
| core.php.image.tag | string | `"v1.1.0"` | The tag to be used |
| core.postgres.image.repo | string | `"postgres"` | The repo to be used |
| core.postgres.image.tag | string | `"16-alpine"` | The tag to be used |
| core.postgres.statefulSet.persistentData.enabled | bool | `true` |  |
| core.postgres.statefulSet.persistentData.storageClassName | string | `"cinder"` |  |
| gpu.enabled | bool | `false` |  |
| gpu.image.repo | string | `"nvidia/samples"` | The repo to be used |
| gpu.image.tag | string | `"vectoradd-cuda11.2.1"` | The tag to be used |
| gpu.nodeLabelSelectors."nvidia.com/gpu.present" | string | `"true"` |  |
| gpu.numberOfGPUs | int | `1` |  |
| gpu.resources | object | `{}` |  |
| ingress.annotations | object | `{}` |  |
| ingress.className | string | `"nginx"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.host | string | `"test.example.uk"` |  |
| ingress.tls[0].hosts[0] | string | `"test.example.uk"` |  |
| ingress.tls[0].secretName | string | `"test-secret"` |  |
| monitoring.grafana.dashboards | bool | `false` |  |
| monitoring.serviceMonitor.enabled | bool | `false` |  |
