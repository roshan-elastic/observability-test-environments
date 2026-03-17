#!/usr/bin/env bash
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

CLUSTER_CONFIG_FILE=$1
STACK_VERSION=$(yq -r .stack_version "${CLUSTER_CONFIG_FILE}")
ES_VERSION=$(yq -r .elasticsearch.version "${CLUSTER_CONFIG_FILE}")
CLUSTER_NAME=$(yq -r .cluster_name "${CLUSTER_CONFIG_FILE}")
elasticsearch::api "${CLUSTER_NAME}" GET "/_license"
LICENSE_STATUS=$(elasticsearch::api "${CLUSTER_NAME}" GET "/_license" | jq -r .license.status)
echo "${LICENSE_STATUS} ${CLUSTER_NAME} ${STACK_VERSION} ${ES_VERSION}"
