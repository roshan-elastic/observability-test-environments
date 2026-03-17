#!/usr/bin/env bash
#
# This script gets the JSON definition for an existing Elastic Cloud Cluster
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"

log::info "usage: $0 CLUSTER_ID"
CLUSTER_ID=${1:-?"Missing CLUSTER_ID argument"}
ENVIRONMENT=${2:-"pro"}

ess::get-deployment "${ENVIRONMENT}" "${CLUSTER_ID}" | jq .
