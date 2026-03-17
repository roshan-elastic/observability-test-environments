#!/usr/bin/env bash
#
# Script to create a new Index an asociate it to the alias.
# this operation force a rollover.
#
# .ci/scripts/rollover_indices_force.sh dev-next"
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

INDICES=$(elasticsearch::api "${CLUSTER}" GET "_aliases" | jq -r '.[].aliases|paths[0]' | grep -v "^\." | sort | uniq)
echo "${INDICES}" > indices.txt

for INDICE in ${INDICES}
do
  ROLLOVER="${INDICE}-$(date +%Y.%m.%d)-000001"
  echo "Rollover ${INDICE} :"
  elasticsearch::api "${CLUSTER}" PUT "${ROLLOVER}" "{\"aliases\":{\"${INDICE}\":{\"is_write_index\": true}}}"
  echo ""
done
