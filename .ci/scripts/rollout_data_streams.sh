#!/usr/bin/env bash
#
# Script to rollover all datastreams.
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

DATA_STREAMS=$(elasticsearch::api "${CLUSTER}" GET "_data_stream/_all" | jq -r '.data_streams[].name')

for DATA_STREAM in ${DATA_STREAMS}
do
  echo "Rollover ${DATA_STREAM} :"
  elasticsearch::api "${CLUSTER}" POST "${DATA_STREAM}/_rollover?pretty"
  echo ""
done
