#!/usr/bin/env bash
# https://backstage.elastic.dev/catalog/default/api/region-api/definition

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/serverless.sh"

ENVIRONMENT=${1:?"Missing environment argument (e.g. staging)"}

serverless::api "${ENVIRONMENT}" "GET" "api/v1/serverless/regions" | jq '.'
