#!/usr/bin/env bash
#
# Script to delete all UNASSIGNED shards and try to restore them from a snapshot.
# it saves a file with the name os the indices (indices.txt)
#
# .ci/scripts/restore_unassigned_shards.sh dev-next snapshot-id"
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER SNAPSHOT [REPO|external_gcs_repository] [INDICES_FILE]"

CLUSTER=${1:?-"missing argument"}
SNAPSHOT=${2:?-"missing argument"}
REPO=${3:-"external_gcs_repository"}
INDICES_FILE=${4:-""}

if [ -z "${INDICES_FILE}" ]; then
  INDICES=$(elasticsearch::api "${CLUSTER}" GET "_cat/shards?h=i,st&s=st,i"|grep "UNASSIGNED"|cut -d " " -f 1|sort|uniq)
  echo "${INDICES}"
  echo "${INDICES}" > indices.txt
else
  INDICES=$(cat "${INDICES_FILE}")
fi

N=0
for INDICE in ${INDICES}
do
  echo "Delete ${INDICE} :"
  elasticsearch::api "${CLUSTER}" DELETE "${INDICE}" &
  echo ""
  N=$((N+1))
  [ $((N%10)) -eq 0 ] && sleep 5
done

for INDICE in ${INDICES}
do
  echo "Restore ${INDICE} from ${SNAPSHOT} :"
  "${CURDIR}/restore_snapshot.sh" "${CLUSTER}" "${SNAPSHOT}" "${INDICE}" "${REPO}"
  echo ""
done
