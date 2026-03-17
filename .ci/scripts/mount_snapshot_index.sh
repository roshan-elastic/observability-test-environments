#!/usr/bin/env bash
#
# https://www.elastic.co/guide/en/elasticsearch/reference/current/searchable-snapshots.html
# https://www.elastic.co/guide/en/elasticsearch/reference/current/searchable-snapshots-api-mount-snapshot.html
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER SNAPSHOT [REPO|'external_gcs_repository']"

CLUSTER=${1:?-"missing argument"}
SNAPSHOT=${2:?-"missing argument"}
REPO=${3:-"external_gcs_repository"}

INDICES=$(elasticsearch::api "${CLUSTER}" GET "_snapshot/${REPO}/${SNAPSHOT}" |jq '.snapshots[]|.indices[]|select(.|test(".ds-metrics*") or test(".ds-log*") or test(".ds-heartbeat*") or test(".ds-auditbeat*") or test(".ds-packetbeat*") or test(".ds-filebeat*") or test(".ds-metricbeat*"))')

elasticsearch::api "${CLUSTER}" GET "/_cluster/settings?pretty" '{
  "persistent" : {
    "cluster.routing.allocation.node_concurrent_recoveries" : "20",
    "indices.recovery.max_concurrent_file_chunks" : "8",
    "indices.recovery.max_concurrent_operations" : "4",
    "indices.recovery.max_bytes_per_sec" : "400mb"
  }
}'

for i in ${INDICES}
do
  elasticsearch::api "${CLUSTER}" POST "/_snapshot/${REPO}/${SNAPSHOT}/_mount?pretty" "
    {
    \"index\": ${i},
    \"index_settings\": {
      \"index.number_of_replicas\": 0
    },
    \"ignore_index_settings\": [ \"index.refresh_interval\" ]
  }"
done
