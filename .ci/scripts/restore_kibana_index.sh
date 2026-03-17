#!/usr/bin/env bash
# Workaround to reindex Kibana index for the new mapping in Stack.
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"
CLUSTER=${1:?-"missing argument"}

echo "Stop Kibana, then press enter to continue"
read -r

elasticsearch::api "${CLUSTER}" PUT "/.kibana_reindex_temp?pretty" '{
    "mappings": {
        "dynamic": false,
        "properties": {}
    }
}'

elasticsearch::api "${CLUSTER}"  POST "_reindex?pretty" "{
  \"source\": {
    \"index\": \".kibana\"
  },
  \"dest\": {
    \"index\": \"kibana_reindex_temp\"
  }
}"

elasticsearch::api "${CLUSTER}" DELETE "/.kibana_*"

echo "Restart Kibana, then press enter to continue"
read -r

elasticsearch::api "${CLUSTER}" POST "/_reindex?pretty" "{
  \"source\": {
    \"index\": \"kibana_reindex_temp\"
  },
  \"dest\": {
    \"index\": \".kibana\"
  }
}
"

elasticsearch::api "${CLUSTER}" DELETE "/kibana_reindex_temp?pretty"
