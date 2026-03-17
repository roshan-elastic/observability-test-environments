#!/usr/bin/env bash
# delete documents from a index that has a field defined, this is for resolve migration issues with Kibana objects.
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER INDEX FIELD"
CLUSTER=${1:?-"missing argument"}
INDEX=${2:?-"missing argument"}
FIELD=${3:?-"missing argument"}


elasticsearch::api "${CLUSTER}" POST "${INDEX}/_delete_by_query" "{
\"query\": {
  \"exists\": {
    \"field\": \"${FIELD}\"
  }
 }
}"
