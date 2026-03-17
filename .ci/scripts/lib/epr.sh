#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"

#######################################
# Get the URL for the EPR API
# Globals:
#   None
# Arguments:
#   Environment
# Returns:
#   URL
#######################################
function epr::environment-api-url(){
  local environment=${1:?"Environment is required"}
  local url="https://epr.elastic.co"
  case "${environment}" in
    "prod")
      ;;
    "production")
      ;;
    "staging")
      url="https://epr-staging.elastic.co"
      ;;
    "snapshot")
      url="https://epr-snapshot.elastic.co"
      ;;
    "experimental")
      url=" https://epr-experimental.elastic.co"
      ;;
    *)
      log::error "Unknown environment: ${environment}"
      exit 1
      ;;
  esac
  log::info "Using EPR API URL: ${url}"
  echo "${url}"
}

#######################################
# Make a call to the EPR API
# Globals:
#   None
# Arguments:
#   Environment
#   Method
#   API
# Returns:
#   API call response
#######################################
function epr::api(){
  local environment=${1:?"Environment is required"}
  local method=${2:?"Method is required"}
  local api=${3:?"API is required"}
  local endpoint
  log::info "Calling ${method} ${api} in ${environment}"
  endpoint=$(epr::environment-api-url "${environment}")
  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl -X "${method}" ${curl_args} "${endpoint}/${api}"
}

#######################################
# Get the EPR API info
# Globals:
#   None
# Arguments:
#   Environment
# Returns:
#   API call response
#######################################
function epr::info() {
  local environment=${1:?"Environment is required"}
  log::info "Getting info from ${environment}"
  epr::api "${environment}" "GET" "/"
}

#######################################
# Search packages
# https://github.com/elastic/package-registry?tab=readme-ov-file#search
# Globals:
#   None
# Arguments:
#   Environment (e.g. production, staging, snapshot, experimental)
#   Query (e.g. all=true, category=security, category=security&category=siem, prerelease=true, package=apm)
# Returns:
#   API call response
#######################################
function epr::search(){
  local environment=${1:?"Environment is required"}
  local query=${2:?"Query is required"}
  log::info "Searching packages in ${environment} with query: ${query}"
  epr::api "${environment}" "GET" "/search?${query}"
}

#######################################
# Get the list of all packages
# Globals:
#   None
# Arguments:
#   Environment (e.g. production, staging, snapshot, experimental)
# Returns:
#   API call response
#######################################
function epr::all-packages(){
  local environment=${1:?"Environment is required"}
  log::info "Getting all packages in ${environment}"
  epr::search "${environment}" "all=true"
}

#######################################
# Download a package
# Globals:
#   None
# Arguments:
#   Environment
#   Package
#   Version
# Returns:
#   Downloads the package
#######################################
function epr::download(){
  local environment=${1:?"Environment is required"}
  local package=${2:?"Package is required"}
  local version=${3:?"Version is required"}
  local endpoint
  log::info "Downloading package ${package} version ${version} from ${environment}"
  endpoint=$(epr::environment-api-url "${environment}")
  local curl_args=-sLfO
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl ${curl_args} -X "GET" "${endpoint}/epr/${package}/${package}-${version}.zip"
}
