#!/usr/bin/env bash
#
# This script restores all resources in an Elastic Cloud deployment
#
# More on this: https://www.elastic.co/guide/en/cloud-enterprise/current/restore-deployment.html
#
set -eo pipefail

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"

log::info "usage: $0 CLUSTER_ID"
CLUSTER_ID=${1:-?"Missing CLUSTER_ID argument"}

ess::restore-cluster "pro" "${CLUSTER_ID}"
