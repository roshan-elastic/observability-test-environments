#!/usr/bin/env bash

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/kibana.sh"

log::info "Usage: $0 <cluster-name> [package-name|apm] [environment|snapshot]"
CLUSTER_NAME=${1:?"Missing cluster name argument (e.g. edge-oblt)"}
PACKAGE_NAME=${2:-"apm"}
EPR_ENVIRONMENT=${3:-"snapshot"}
PACKAGE_VERSION=$(epr::search "${EPR_ENVIRONMENT}" "package=${PACKAGE_NAME}&prerelease=true" | jq -r .[].version)

epr::download snapshot "${PACKAGE_NAME}" "${PACKAGE_VERSION}"
kibana::fleet::upload-package "${CLUSTER_NAME}" "${PWD}/${PACKAGE_NAME}-${PACKAGE_VERSION}.zip"
