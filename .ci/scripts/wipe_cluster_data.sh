#!/usr/bin/env bash
#
# Delete all data in a cluster
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

elasticsearch::api "${CLUSTER}" DELETE "/_data_stream/*?expand_wildcards=all&pretty"
elasticsearch::api "${CLUSTER}" DELETE "/*?expand_wildcards=all&pretty"
