# dogkat

![Version: 0.1.5](https://img.shields.io/badge/Version-0.1.5-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.5](https://img.shields.io/badge/AppVersion-0.1.5-informational?style=flat-square)

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
| core.nginx.exporterImage.repo | string | `"nginx/nginx-prometheus-exporter@sha256"` | The repo to be used |
| core.nginx.exporterImage.tag | string | `"d710e0ff2505a7037dd21e47eae07025010c0de08a6247d1a704824823becfd0"` | The tag to be used |
| core.nginx.image.repo | string | `"nginx@sha256"` | The repo to be used |
| core.nginx.image.tag | string | `"02d8d94023878cedf3e3acc55372932a9ba1478b6e2f3357786d916c2af743ba"` | The tag to be used |
| core.nginx.resources | object | `{}` |  |
| core.nginx.serviceAccountName | string | `"nginx"` |  |
| core.php.image.repo | string | `"drewviles/php-pdo@sha256"` | The repo to be used |
| core.php.image.tag | string | `"253465d95c3fa68871c2ccc6c67d4ed5ee500563fbbfee3b54a9544f8025d1d6"` | The tag to be used |
| core.postgres.image.repo | string | `"postgres@sha256"` | The repo to be used |
| core.postgres.image.tag | string | `"sha256:49fd8c13fbd0eb92572df9884ca41882a036beac0f12e520274be85e7e7806e9"` | The tag to be used |
| core.postgres.statefulSet.persistentData.enabled | bool | `true` |  |
| core.postgres.statefulSet.persistentData.storageClassName | string | `"cinder"` |  |
| gpu.enabled | bool | `false` |  |
| gpu.image.repo | string | `"nvcr.io/nvidia/k8s/cuda-sample@sha256"` | The repo to be used |
| gpu.image.tag | string | `"04a20bfaf69363ec3f15fc1cdb0abc0efabeb6fb6b3a1b9cf4a575ae7b1d81d1"` | The tag to be used |
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
