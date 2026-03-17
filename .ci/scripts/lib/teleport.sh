#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"

#####################################
# Login in Teleport
# Globals:
#   None
# Arguments:
#   environment (staging or production)
# Returns:
#   None
#####################################
function teleport::login(){
  local environment=$1
  local PROXY_HOST="teleport-proxy.staging.getin.cloud"
  case "${environment}" in
    "qa")
      PROXY_HOST="teleport-proxy.staging.getin.cloud"
      ;;
    "staging")
      PROXY_HOST="teleport-proxy.staging.getin.cloud"
      ;;
     "production")
      PROXY_HOST="teleport-proxy.secops.elstc.co"
      ;;
  esac
  log::info "Login in ${environment} with teleport"
  # FIXME add a few retries, it fails sometimes
  if ! tsh login --proxy="${PROXY_HOST}" --auth=okta --overwrite; then
    log::error "Login failed"
    log::error "request access to login in ${environment} MKI using Teleport in the #cloud-security-public Slack channel"
    log::error "https://cloud.elastic.dev/teleport/USAGE/#more-information"
    exit 1
  fi
}

#####################################
# Logout from Teleport
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function teleport::logout(){
  log::debug "Logout from teleport"
  tsh logout
}

#####################################
# Get Teleport session status
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function teleport::status(){
  log::debug "Get Teleport session status"
  tsh status
}

#####################################
# List MKI clusters available
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function mki::list-clusters(){
  log::debug "List MKI clusters available"
  tsh kube ls
}
