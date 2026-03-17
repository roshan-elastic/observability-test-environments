#!/usr/bin/env bash
#
# This script put a cluster in maintenance mode.
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"

log::info "usage: $0 CLUSTER_ID [ENVIRONMENT|pro|qa|staging]"
CLUSTER_ID=${1:?-"missing argument"}
ENVIRONMENT=${2:-"pro"}

ess::manteinance-mode "${ENVIRONMENT}" "${CLUSTER_ID}" "start"
