#!/bin/bash
set -euxo pipefail

# renovate: datasource=github-releases depName=yannh/kubeconform
KUBECONFORM_VERSION=0.6.4

#CHART_DIRS="$(git diff --find-renames --name-only "$(git rev-parse --abbrev-ref HEAD)" remotes/origin/main -- charts | cut -d '/' -f 2 | uniq)"
CHART_DIRS=$(ls charts)

# install kubeconform
curl --silent --show-error --fail --location --output /tmp/kubeconform.tar.gz "https://github.com/yannh/kubeconform/releases/download/v${KUBECONFORM_VERSION}/kubeconform-linux-amd64.tar.gz"
tar -xf /tmp/kubeconform.tar.gz kubeconform

# validate charts
for CHART_DIR in ${CHART_DIRS}; do
  helm template --values charts/"${CHART_DIR}"/ci/test-values.yaml charts/"${CHART_DIR}" | ./kubeconform --strict --ignore-missing-schemas --kubernetes-version "${KUBERNETES_VERSION#v}"
done
