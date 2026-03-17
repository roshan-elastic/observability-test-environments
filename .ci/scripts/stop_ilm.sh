#!/usr/bin/env bash
#
# Stops ILM recommended before restore an snapshot.
# .ci/scripts/stop_ilm.sh dev-next
#
# https://www.elastic.co/guide/en/elasticsearch/reference/current/index-lifecycle-and-snapshots.html

set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

elasticsearch::api "${CLUSTER}" POST "_ilm/stop"
