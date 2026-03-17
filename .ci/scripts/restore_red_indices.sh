#!/usr/bin/env bash
#
# Script to delete all red index and try to restore them from a snapshot.
# it saves a file with the name os the indices (indices.txt)
#
# .ci/scripts/restore_snapshot.sh dev-next snapshot-id"
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER SNAPSHOT [REPO|external_gcs_repository]"

CLUSTER=${1:?-"missing argument"}
SNAPSHOT=${2:?-"missing argument"}
REPO=${2:-"external_gcs_repository"}

INDICES=$(elasticsearch::api "${CLUSTER}" GET "_cat/indices?health=red&h=i")

echo "${INDICES}" > indices.txt

for INDICE in ${INDICES}
do
  echo "Delete ${INDICE} :"
  elasticsearch::api "${CLUSTER}" DELETE "${INDICE}"
  echo ""
done

for INDICE in ${INDICES}
do
  echo "Restore ${INDICE} from ${SNAPSHOT} :"
  "${CURDIR}/restore_snapshot.sh" "${CLUSTER_FULLNAME}" "${SNAPSHOT}" "${INDICE}" "${REPO}"
  echo ""
done
