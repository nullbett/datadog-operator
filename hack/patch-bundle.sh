#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)
YQ="$ROOT_DIR/bin/yq"

# Apply patch-bundle-csv-scc.yaml on CSV file
# (cannot be done using Kustomize as the final CSV is generated by operator-sdk binary)
$YQ m -i -a "$ROOT_DIR/bundle/manifests/datadog-operator.clusterserviceversion.yaml" "$ROOT_DIR/hack/patch-bundle-csv-scc.yaml"

# Remove ServiceAccount bundled in SCC (but required for Kustomize installation to work)
rm -f "$ROOT_DIR/bundle/manifests/datadog-operator-manager_v1_serviceaccount.yaml"

# Remove defaultOverride section in DatadogAgent status due to the error: "datadoghq.com_datadogagents.yaml bigger than total allowed limit"
$YQ d -i "$ROOT_DIR/bundle/manifests/datadoghq.com_datadogagents.yaml" 'spec.validation.openAPIV3Schema.properties.status.properties.defaultOverride'
