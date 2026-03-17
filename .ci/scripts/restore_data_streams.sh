#!/usr/bin/env bash
#
# Script to delete all RED datastreams.#
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

# get the data streams in red state
DATA_STREAMS=$(elasticsearch::api "${CLUSTER}" GET "_data_stream" | jq -r '.data_streams[] | select(.status=="RED") | .name')
echo "Data streams in RED state: ${DATA_STREAMS}"
echo "Do you want to delete them?"
read -r

# delete the data streams
for DATA_STREAM in ${DATA_STREAMS}; do
    echo "Deleting data stream ${DATA_STREAM}"
    elasticsearch::api "${CLUSTER}" DELETE "_data_stream/${DATA_STREAM}"
done
