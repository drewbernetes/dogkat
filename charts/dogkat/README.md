# dogkat

![Version: 0.1.11](https://img.shields.io/badge/Version-0.1.11-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.11](https://img.shields.io/badge/AppVersion-0.1.11-informational?style=flat-square)

End-2-End testing for GPUs and some core resources

## Installation

```shell
helm install https://drewbernetes.github.io/dogkat/dogkat --values values.yaml
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
| core.nginx.exporterImage.repo | string | `"nginx/nginx-prometheus-exporter@sha256"` |  |
| core.nginx.exporterImage.tag | string | `"6477cf3bddc4e042d3496856fb2e8e382301bac47fb18cb83924389717261cb1"` |  |
| core.nginx.image.repo | string | `"cgr.dev/chainguard/nginx"` |  |
| core.nginx.image.tag | string | `"latest"` |  |
| core.nginx.resources | object | `{}` |  |
| core.nginx.serviceAccountName | string | `"nginx"` |  |
| core.php.image.repo | string | `"drewviles/php-pdo@sha256"` |  |
| core.php.image.tag | string | `"4485f4a33423d3ca5cceb2600e72e32550ce98ce628c05dc175c7a5763faa616"` |  |
| core.postgres.image.repo | string | `"postgres@sha256"` |  |
| core.postgres.image.tag | string | `"d898b0b78a2627cb4ee63464a14efc9d296884f1b28c841b0ab7d7c42f1fffdf"` |  |
| core.postgres.statefulSet.persistentData.enabled | bool | `true` |  |
| core.postgres.statefulSet.persistentData.storageClassName | string | `"longhorn"` |  |
| gpu.enabled | bool | `false` |  |
| gpu.image.repo | string | `"nvcr.io/nvidia/k8s/cuda-sample@sha256"` | The repo to be used |
| gpu.image.tag | string | `"ac53daee629763d712e1361b77e4c4f4ad146148f9dffc6288a75732270c6e85"` | The tag to be used |
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
