#!/usr/bin/env bash
CURDIR=$(dirname "$0")
# shellcheck source=/dev/null
. "${CURDIR}/activate-hermit"
"${CURDIR}/hermit" upgrade --quiet
