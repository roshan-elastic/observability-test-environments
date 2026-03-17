#!/usr/bin/env bash
# https://backstage.elastic.dev/catalog/default/api/project-api

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/serverless.sh"

ENVIRONMENT=${1:?"Missing environment argument (e.g. staging)"}
PROJECT_TYPE=${2:?"Missing project type argumen (e.g. elasticsearch)"}


serverless::get-projects "${ENVIRONMENT}" "${PROJECT_TYPE}" | jq '.'
