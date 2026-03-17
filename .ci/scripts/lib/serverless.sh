#!/bin/env bash
# https://backstage.elastic.dev/catalog/default/api/admin-project-api
# https://backstage.elastic.dev/catalog/default/api/project-api

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/ess.sh"

#####################################
# Get the ESS endpoint for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the URL to the Serverless endpoint
#####################################
function serverless::get-endpoint(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  log::debug "Getting Serverless endpoint for ${environment}"
  vault_secret=$(ess::get-vault-secret "${environment}")
  gcp::read-secret-field serverless_endpoint "${vault_secret}"
}

#####################################
# Make a rest call to the admin project API
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   method (GET, POST, PUT, DELETE)
#   api_path the path of the API (e.g. "api/v1/admin/serverless/projects/elasticsearch")
#   body (optional body for POST and PUT requests)
# Returns:
#   the response body
#####################################
function serverless::admin-api(){
  local environment=${1:?"Missing environment argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  log::warn "Access to the admin API requirre an VPN connection to the Elastic network in some environments."
  log::debug "Making ${method} request to ${api_path} with body ${body}"
  ess::console-api "${environment}" "${method}" "${api_path}" "${body}"
}

#####################################
# Make a rest call to the serverless API
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   method (GET, POST, PUT, DELETE)
#   api_path the path of the API (e.g. "api/v1/serverless/projects/elasticsearch")
#   body (optional body for POST and PUT requests)
# Returns:
#   the response body
#####################################
function serverless::api(){
  local environment=${1:?"Missing environment argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  local endpoint
  local auth_header
  log::debug "Making ${method} request to ${api_path} with body ${body}"
  endpoint=$(serverless::get-endpoint "${environment}")
  auth_header=$(ess::get-auth-header "${environment}")

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
# Get the Serverless projects list of a type for the given environment
# if a project id is provided, get the project details
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   id (Optional the project id)
# Returns:
#   the response body
#####################################
function serverless::get-projects(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local id=${3:-""}
  log::debug "Getting Serverless projects for ${environment} of type ${type} with id '${id}'"
  if [ -n "$id" ]; then
    id="/${id}"
  fi
  serverless::api "${environment}" "GET" "api/v1/serverless/projects/${type}${id}"
}

#####################################
# Get the Serverless projects list of a type for the given environment from the Admin project API
# if a project id is provided, get the project details
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   id (Optional the project id)
# Returns:
#   the response body
#####################################
function serverless::admin-get-projects(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local id=${3:-""}
  log::debug "Getting Serverless projects for ${environment} of type ${type} with id '${id}'"
  if [ -n "$id" ]; then
    id="/${id}"
  fi
  serverless::admin-api "${environment}" "GET" "api/v1/admin/serverless/projects/${type}${id}"
}

#####################################
# Update a Serverless project of a type for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   id (the project id)
#   body (the project body)
# Returns:
#   the response body
#####################################
function serverless::admin-update-project(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local id=${3:?"Missing project id argument"}
  local body=${4:?"Missing body argument"}
  log::debug "Updating Serverless project for ${environment} of type ${type} with id '${id}'"
  serverless::admin-api "${environment}" "PUT" "api/v1/admin/serverless/projects/${type}/${id}" "${body}"
}

#####################################
# Delete a Serverless project of a type for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   id (the project id)
# Returns:
#   the response body
#####################################
function serverless::delete-project(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local id=${3:?"Missing project id argument"}
  log::debug "Deleting Serverless project for ${environment} of type ${type} with id '${id}'"
  serverless::api "${environment}" "DELETE" "api/v1/serverless/projects/${type}/${id}"
}

#####################################
# Create a Serverless project of a type for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   name (the project name)
#   region_id (the region id)
# Returns:
#   the response body
#####################################
function serverless::create-project(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local name=${3:?"Missing name argument"}
  local region_id=${4:?"Missing name argument"}
  serverless::api "${environment}" POST "api/v1/serverless/projects/observability" "{\"name\": \"${name}\", \"region_id\": \"${region_id}\", \"overrides\": {\"elasticsearch\": {}, \"kibana\": {}, \"fleet\": {}}}"
}

#####################################
# Delete all Serverless projects that its alias starts with the given string
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   type (observability, elasticsearch, or security)
#   pattern (the pattern to match)
# Returns:
#   the response body
#####################################
function serverless::delete-projects-startswith(){
  local environment=${1:?"Missing environment argument"}
  local type=${2:?"Missing type argument"}
  local pattern=${3:?"Missing pattern argument"}
  local projects
  projects=$(serverless::get-projects "${environment}" "${type}")
  for project in $(echo "${projects}" | jq -r ".items[]|select(.alias | startswith(\"${pattern}\"))|.id" ); do
    serverless::delete-project "${environment}" "${type}" "${project}"
  done
}
