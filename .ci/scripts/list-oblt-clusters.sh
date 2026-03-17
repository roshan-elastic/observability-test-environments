#!/usr/bin/env bash
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"

log::info  "usage: $0 [ENVIRONMENT|pro]"
ENVIRONMENT=${1:-"pro"}
IS_HIDEN=${2:-"false"}

ess::search-deployments-all "${ENVIRONMENT}" "${IS_HIDEN}"| jq '.deployments[]|{"id":.id, "name":.name, "metadata": .metadata}' > deployments.json

set +e
jq -r '.name' deployments.json | while IFS= read -r CLUSTER_NAME
do
  CLUSTER_NAME=$(echo "${CLUSTER_NAME}" | tr -d '"')
  CONFIG_FILE=$(grep -rl -E -e "cluster_name: \"?${CLUSTER_NAME}\"?\$" environments)
  if [ -z "${CONFIG_FILE}" ]; then
    CONFIG_FILE="none"
  fi
  echo "${CLUSTER_NAME};${CONFIG_FILE}"
done
