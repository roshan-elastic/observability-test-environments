#!/usr/bin/env bash
# https://backstage.elastic.dev/catalog/default/api/project-api/definition#/default/deleteElasticsearchProject

CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/serverless.sh"

ENVIRONMENT=${1:?"Missing environment argument (e.g. staging)"}

for id in $(serverless::admin-get-projects "$ENVIRONMENT" elasticsearch | jq -c -r '.items[].id')
do
  serverless::delete-project "$ENVIRONMENT" "elasticsearch" "${id}"
done

for id in $(serverless::admin-get-projects "$ENVIRONMENT" observability | jq -c -r '.items[].id')
do
  serverless::delete-project "$ENVIRONMENT" "observability" "${id}"
done

for id in $(serverless::admin-get-projects "$ENVIRONMENT" security | jq -c -r '.items[].id')
do
  serverless::delete-project "$ENVIRONMENT" "security" "${id}"
done
