#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/elasticsearch.sh"


#####################################
# Get the APM endpoint
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the endpoint
#####################################
function apm::get-endpoint(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  log::debug "Getting APM endpoint for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
   gcp::read-secret-field apm_url "${vault_secret}"
}

#####################################
# Get the APM auth header
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the auth header
#####################################
function apm::get-auth-header(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  local token
  local apikey
  log::debug "Getting the auth header for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  token=$(gcp::read-secret-field apm_token "${vault_secret}")
  apikey=$(gcp::read-secret-field apm_apikey "${vault_secret}")
  if [ -n "${token}" ]; then
    echo "Authorization: Bearer ${token}"
  elif [ -n "${apikey}" ]; then
    echo "Authorization: Apikey ${apikey}"
  fi
}

#####################################
# Make a request to the APM API
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
#   method (GET, POST, PUT, DELETE)
#   api_path (e.g. _cluster/health)
#   body (JSON)
# Returns:
#   the response body
#####################################
function apm::api(){
  local cluster=${1:?"Missing cluster name argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  log::info "Making a ${method} request to ${api_path} on ${cluster}"
  local vault_secret
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  local endpoint
  endpoint=$(apm::get-endpoint "${cluster}")
  auth_header=$(apm::get-auth-header "${cluster}")
  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl -X "${method}" ${curl_args} \
  -H 'Content-Type: application/json' \
  -H "${auth_header}" \
  "${endpoint}/${api_path}" \
  -d "${body}"
}

#####################################
# Check the APM endpoint
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the response body
#####################################
function apm::check(){
  local cluster=${1:?"Missing cluster name argument"}
  local apm_endpoint
  apm_endpoint=$(apm::get-endpoint "${cluster}")
  log::info "Checking APM endpoint ${apm_endpoint}"
  apm::api "${cluster}" GET "/"
}
