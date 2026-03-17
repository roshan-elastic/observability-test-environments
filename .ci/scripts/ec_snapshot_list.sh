#!/usr/bin/env bash
#
#
# .ci/scripts/ec_snapshot_list.sh dev-next
#
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"
CLUSTER=${1:?-"missing argument"}

# snapshot info
elasticsearch::api "${CLUSTER}" GET "/_snapshot/found-snapshots/_all"|jq '.snapshots[]| {state: .state, name: .snapshot, start_time: .start_time, end_time: .end_time, failures: .failures, shards: .shards}'
