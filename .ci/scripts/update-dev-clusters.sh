#!/usr/bin/env bash
set -euo pipefail
MSG="parameter missing."
REMOTE_CLUSTER_NAME=${1:?$MSG}
CURDIR=$(dirname "$0")
REMOTE_CLUSTER_CONFIG_FILE=$("${CURDIR}/find_cluster_config_by_name.py" "${REMOTE_CLUSTER_NAME}" "${CURDIR}/../../environments" false)

# shellcheck disable=SC2013
for cluster in $(grep -rl "ccs_remote_cluster\: \"${REMOTE_CLUSTER_NAME}\""  environments)
do
  CLUSTER_NAME=$(yq .cluster_name "${cluster}") \
  REMOTE_CLUSTER_CONFIG_FILE=${REMOTE_CLUSTER_CONFIG_FILE} \
  CLUSTER_CONFIG_FILE=$("${CURDIR}/find_cluster_config_by_name.py" "${CLUSTER_NAME}" "${CURDIR}/../../environments" true) \
  STACK_MODE=ess \
  make -C "${CURDIR}/../update-versions" update-dev-cluster || echo "::warning title=Failed to update::The cluster ${CLUSTER_NAME} failed to update."
done
