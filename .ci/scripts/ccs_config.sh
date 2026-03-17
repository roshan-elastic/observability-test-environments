#!/usr/bin/env bash
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER REMOTE_CLUSTER_NAME REMOTE_CLUSTER_HOST_PORT [REMOTE_CLUSTER_PORT|9400] "
CLUSTER=${1:?-"missing argument"}
REMOTE_CLUSTER_NAME=${2:?-"missing argument"}
REMOTE_CLUSTER_HOST=${3:?-"missing argument"}
REMOTE_CLUSTER_PORT=${4:-"9400"}

elasticsearch::api "${CLUSTER}" PUT "/_cluster/settings?pretty" "
{
  \"persistent\": {
    \"cluster\": {
      \"remote\": {
        \"${REMOTE_CLUSTER_NAME}\": {
          \"skip_unavailable\": false,
          \"mode\": \"proxy\",
          \"proxy_address\": \"${REMOTE_CLUSTER_HOST}:${REMOTE_CLUSTER_PORT}\",
          \"proxy_socket_connections\": 3,
          \"server_name\": \"${REMOTE_CLUSTER_HOST}\",
          \"seeds\": null,
          \"node_connections\": null
        }
      }
    }
  }
}"
