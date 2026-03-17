#!/usr/bin/env bash
# return the info from the APM service
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/apm.sh"

log::info "usage: $0 CLUSTER"
CLUSTER=${1:?-"missing argument"}

apm::check "${CLUSTER}"
