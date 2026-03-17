#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/install.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/elasticsearch.sh"

#####################################
# Run Kibana en dev mode with a custom config file
# Globals:
#   PWD Current working directory
# Arguments:
#   kibana_folder full path of the kibana folder
#   kibana_file full path of the kibana.yml file
# Returns:
#   None
#####################################
function kibana::run-dev(){
  local kibana_folder=${1:?"Missing kibana_folder argument"}
  local kibana_file=${2:?"Missing kibana_file argument"}

  log::info Run Kibana
  OLD_PWD="${PWD}"
  trap 'cd "${OLD_PWD}"' EXIT

  cd "${kibana_folder}" || exit 1
  install::node "$(cat .nvmrc)"
  yarn kbn bootstrap
  yarn serverless-oblt --config "${kibana_file}"
}

#####################################
# Get the kibana auth header
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the auth header
#####################################
function kibana::get-auth-header(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  local username
  local password
  log::debug "Getting the auth header for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  username=$( gcp::read-secret-field elasticsearch_username "${vault_secret}")
  password=$( gcp::read-secret-field elasticsearch_password "${vault_secret}")
  echo "Authorization: Basic $(echo -n "${username}:${password}" | base64)"
}

#####################################
# Get the kibana endpoint
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the endpoint
#####################################
function kibana::get-endpoint(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  log::debug "Getting kibana endpoint for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
   gcp::read-secret-field kibana_url "${vault_secret}"
}

#####################################
# Make a call to the Kibana API
# Globals:
#   None
# Arguments:
#   cluster (e.g. edge-oblt)
#   method (GET, POST, PUT, DELETE)
#   api (e.g. api/status)
#   data path to a file containing the data to send
#   headers (e.g. "content-type: application/zip,kbn-xsrf: true")
# Returns:
#   API call response
#####################################
function kibana::api(){
  local cluster=${1:?"Cluster is required"}
  local method=${2:?"Method is required"}
  local api=${3:?"API is required"}
  local data=${4:-""}
  local headers=${5:-""}
  local auth_header
  local endpoint
  endpoint=$(kibana::get-endpoint "${cluster}")
  auth_header=$(kibana::get-auth-header "${cluster}")

  log::debug "Calling ${method} ${api}${endpoint}"
  if [[ -n "${headers}" ]]; then
    headers="-H ${headers//,/ -H }"
  fi

  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl -X "${method}" ${curl_args} \
    -H "${auth_header}" ${headers} \
    -H 'kbn-xsrf: true' "${endpoint}/${api}" -d "${data}"
}

#####################################
# Get the Kibana status
# Globals:
#   None
# Arguments:
#   cluster (e.g. edge-oblt)
# Returns:
#   API call response
#####################################
function kibana::status(){
  local cluster=${1:?"Cluster is required"}
  log::debug "Getting status from ${cluster}"
  kibana::api "${cluster}" GET "api/status"
}

#####################################
# Upload a package to Kibana
# Globals:
#   None
# Arguments:
#   cluster (e.g. edge-oblt)
#   package path to the package to upload
# Returns:
#   API call response
#####################################
function kibana::fleet::upload-package(){
  local cluster=${1:?"Cluster is required"}
  local package=${2:?"Package is required"}
  local auth_header
  local endpoint
  local api
  log::debug "Uploading package ${package} to ${cluster}"
  endpoint=$(kibana::get-endpoint "${cluster}")
  auth_header=$(kibana::get-auth-header "${cluster}")
  api="api/fleet/epm/packages"
  if [ ! -f "${package}" ]; then
    log::error "Package ${package} not found"
    exit 1
  fi
  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl ${curl_args} \
    -X "POST" \
    -H "${auth_header}" \
    -H "content-type: application/zip" \
    -H 'kbn-xsrf: true' \
    "${endpoint}/${api}" --data-binary "@${package}"
}
