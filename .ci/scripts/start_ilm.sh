#!/usr/bin/env bash
#
#
# .ci/scripts/start_ilm.sh dev-next
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

elasticsearch::api "${CLUSTER}" POST "_ilm/start"
