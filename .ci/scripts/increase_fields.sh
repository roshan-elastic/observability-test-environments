#!/usr/bin/env bash
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER INDEX"
CLUSTER=${1:?-"missing argument"}
INDEX=${2:?-"missing argument"}

elasticsearch::api "${CLUSTER}" PUT "/${INDEX}/_settings" '{"index.mapping.total_fields.limit": 4000}'
