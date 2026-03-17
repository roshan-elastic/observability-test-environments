#!/usr/bin/env bash
# https://backstage.elastic.dev/catalog/default/api/internal-project-api/definition#/default/createProjectState

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/serverless.sh"

ENVIRONMENT=${1:?"Missing environment argument (e.g. staging)"}
PROJECT_TYPE=${2:?"Missing project type argumen (e.g. elasticsearch)"}
PROJECT_ID=${3:?"Please provide a project_id as the first argument (e.g. 1234567890abcdef1234567890abcdef)"}

serverless::admin-api "${ENVIRONMENT}" "GET" "api/internal/v1/projects/${PROJECT_TYPE}/${PROJECT_ID}/state" | jq '.'
