#!/usr/bin/env bash
#
# https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-allocation-explain.html
# .ci/scripts/allocation_explain.sh dev-next [index]
#
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER INDEX [SHARD_NUMBER|0] [PRIMARY|true]"

CLUSTER=${1:?-"missing argument"}
INDEX=${2:?-"missing argument"}
SHARD_NUMBER=${3:-"0"}
PRIMARY=${4:-"true"}

elasticsearch::api "${CLUSTER}" GET "_cluster/allocation/explain?pretty" "{
  \"index\": \"${INDEX}\",
  \"shard\": ${SHARD_NUMBER},
  \"primary\": ${PRIMARY}
}"|jq .
