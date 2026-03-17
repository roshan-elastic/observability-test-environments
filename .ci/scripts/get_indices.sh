#!/usr/bin/env bash
#
# Script to get the list of indices and aliases from a cluster.
# it saves a file with the name of the indices (indices.txt) and aliases (aliases.txt)
#
# .ci/scripts/get_indices.sh dev-next"
#
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER"

CLUSTER=${1:?-"missing argument"}

INDICES=$(elasticsearch::api "${CLUSTER}" GET "_cat/indices?h=i"|sort)

echo "${INDICES}" > indices.txt

ALIASES=$(elasticsearch::api "${CLUSTER}" GET "_cat/aliases?h=alias"|sort|uniq)

echo "${ALIASES}" > aliases.txt
