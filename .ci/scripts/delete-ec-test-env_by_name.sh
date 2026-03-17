#!/usr/bin/env bash
#
# Delete the EC environment given the cluster name.
#
# usage: $0 CLUSTER [FORCE|false] [ENVIRONMENT|pro|qa|staging] [IS_HIDEN|false|true]
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/prompt.sh"

log::info "usage: $0 CLUSTER [FORCE|false] [ENVIRONMENT|pro|qa|staging] [IS_HIDEN|false|true] [STATUS|started|stopped|initializing]"
CLUSTER=${1:?-"missing argument"}
FORCE=${2:-"false"}
ENVIRONMENT=${3:-"pro"}
IS_HIDEN=${4:-"false"}
STATUS=${5:-"started"}

ess::search-deployments "${ENVIRONMENT}" "${CLUSTER}" "${IS_HIDEN}" "${STATUS}"| jq '.deployments[]|{"id":.id, "name":.name, "metadata": .metadata}' > deployments.json

jq -r '([.id, .name, .metadata.organization_id]) | @tsv' deployments.json
if [ "${FORCE}" == "false" ]; then
  prompt::askYN "We will delete these clusters, Is this OK?"
fi

for CLUSTER_ID in $(jq -r '.id' deployments.json);
do
  ess::delete-deployment "${ENVIRONMENT}" "${CLUSTER_ID}"
done
