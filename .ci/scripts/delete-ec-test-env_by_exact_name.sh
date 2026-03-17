#!/usr/bin/env bash
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/prompt.sh"

log::info "usage: $0 CLUSTER [FORCE|false] [ENVIRONMENT|pro|qa|staging] [IS_HIDEN|false|true]"
CLUSTER=${1:?-"missing argument"}
FORCE=${2:-"false"}
ENVIRONMENT=${3:-"pro"}
IS_HIDEN=${4:-"false"}

ess::search-deployments-by-alias "${ENVIRONMENT}" "${CLUSTER}" "${IS_HIDEN}"| jq '.deployments[]|{"id":.id, "name":.name, "alias": .alias, "metadata": .metadata}' > deployments.json

jq -r '([.id, .name, .metadata.organization_id]) | @tsv' deployments.json
if [ "${FORCE}" == "false" ]; then
  prompt::askYN "We will delete these clusters, Is this OK?"
fi

for CLUSTER_ID in $(jq -r '.id' deployments.json);
do
  ess::delete-deployment "${ENVIRONMENT}" "${CLUSTER_ID}"
done
