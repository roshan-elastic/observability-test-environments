#!/usr/bin/env bash
#
#
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

for INDEX_PREFIX in apm filebeat metricbeat heartbeat auditbeat .ds
do
  elasticsearch::api "${CLUSTER}" POST "${INDEX_PREFIX}-*/_settings" '{
      "index.routing.allocation.include._tier_preference" : "data_warm"
  }'
done
