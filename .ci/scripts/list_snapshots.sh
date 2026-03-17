#!/usr/bin/env bash
#
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER [REPO|'external_gcs_repository']"

CLUSTER=${1:?-"missing argument"}
REPO=${2:-"external_gcs_repository"}

# snapshot info
elasticsearch::api "${CLUSTER}" POST "_snapshot/${REPO}/_verify"
elasticsearch::api "${CLUSTER}" GET "_snapshot/${REPO}/_all" |jq '.snapshots[]|{"snapshot":.snapshot, "state":.state, "end_time":.end_time}'
