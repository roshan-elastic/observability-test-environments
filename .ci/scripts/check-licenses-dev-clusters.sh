#!/usr/bin/env bash
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

echo > licenses.json

set +e

# shellcheck disable=SC2013,SC2044
for cluster in $(find "${PWD}/environments/users" -name "*.yml"|grep -v deploykibana|sort)
do
  log::info "Checking licenses for ${cluster}"
  CLUSTER_NAME=$(yq -r .cluster_name "${cluster}")
  if [ -z "${CLUSTER_NAME}" ]; then
    # it is not a cluster config file
    continue
  fi
  OBLT_MANAGED=$(yq -r .oblt_managed "${cluster}")
  if [ "${OBLT_MANAGED}" == "false" ]; then
    # it is not an oblt managed cluster
    continue
  fi
  STACK_MODE=$(yq -r .stack.mode "${cluster}")
  if [ "${STACK_MODE}" != "ess" ] && [ -n "${STACK_MODE}" ]; then
    # it is not a ESS cluster
    continue
  fi
  STACK_VERSION=$(yq -r .stack.version "${cluster}")
  SLACK_CHANNEL=$(yq -r .slack_channel "${cluster}")
  REMOTE_CLUSTER=$(yq -r .stack.ess.elasticsearch.ccs_remote_cluster "${cluster}")
  LICENSE_STATUS=$(elasticsearch::api "${CLUSTER_NAME}" GET "_license" | jq -r .license.status)

  if [ "${LICENSE_STATUS}" == "" ] || [ "${LICENSE_STATUS}" == "null" ]; then
    LICENSE_STATUS="unknown"
  fi
  cat << EOF
    {
      "status": "${LICENSE_STATUS}",
      "name": "${CLUSTER_NAME}",
      "stack_version": "${STACK_VERSION}",
      "slack_channel": "${SLACK_CHANNEL}",
      "config": "${cluster}",
      "remote_cluster": "${REMOTE_CLUSTER}"
    }
EOF
done >> licenses.json

jq . -s -c licenses.json > licenses_all.json
jq '.[]|select(.status=="expired")|{"name": .name, "version": .stack_version, "slack_channel": .slack_channel, "config": .config}' licenses_all.json | jq -s -c . > expired.json
jq '.[]|select(.status=="unknown")|{"name": .name, "version": .stack_version, "slack_channel": .slack_channel, "config": .config}' licenses_all.json | jq -s -c . > unknown.json
