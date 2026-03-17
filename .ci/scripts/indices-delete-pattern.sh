#!/usr/bin/env bash
#
# Script to delete all red index and try to restore them from a snapshot.
# it saves a file with the name os the indices (indices.txt)
#
# .ci/scripts/indices-delete-pattern.sh dev-next ".ds*"
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER PATTERN"

CLUSTER=${1:?-"missing argument"}
PATTERN=${2:?-"missing argument"}

INDICES=$(elasticsearch::api "${CLUSTER}" GET "_cat/indices/${PATTERN}?expand_wildcards=all&h=i")

echo "${INDICES}" > indices.txt
cat indices.txt
read -r -p "Press enter to continue"

for INDICE in ${INDICES}
do
  echo "Delete ${INDICE} :"
  elasticsearch::api "${CLUSTER}" DELETE "${INDICE}"
  echo ""
done
