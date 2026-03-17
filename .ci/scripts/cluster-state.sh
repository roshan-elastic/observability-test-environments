#!/usr/bin/env bash
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/ess.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

CLUSTER=${1:?-"missing argument"}

vault_secret=$(elasticsearch::get-vault-secret "${CLUSTER}")
ess_pro=$(ess::get-vault-secret "pro")
ELASTICSEARCH_API_KEY=$(gcp::read-secret-field elasticsearch_apikey "${vault_secret}")
USERNAME=$(gcp::read-secret-field elasticsearch_username "${vault_secret}")
PASSWORD=$(gcp::read-secret-field elasticsearch_password "${vault_secret}")
ES_URL=$(gcp::read-secret-field elasticsearch_url "${vault_secret}")
APM_URL=$(gcp::read-secret-field apm_url "${vault_secret}")
APM_TOKEN=$(gcp::read-secret-field apm_token "${vault_secret}")
TF_VAR_ess_apikey=$(gcp::read-secret-field apiKey "${ess_pro}")

export USERNAME
export PASSWORD
export ELASTICSEARCH_API_KEY
export ES_URL
export APM_URL
export APM_TOKEN
export TF_VAR_ess_apikey
