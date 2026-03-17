#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/elastic.sh"

#####################################
# Make a request to the Elasticsearch API
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
function elasticsearch::api(){
  local cluster=${1:?"Missing cluster name argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  log::info "Making a ${method} request to ${api_path} on ${cluster}"
  local vault_secret
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  local endpoint
  endpoint=$(elasticsearch::get-endpoint "${cluster}")
  auth_header=$(elasticsearch::get-auth-header "${cluster}")
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
# Get the Elasticsearch auth header
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the auth header
#####################################
function elasticsearch::get-auth-header(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  local username
  local password
  log::debug "Getting the auth header for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  username=$(gcp::read-secret-field elasticsearch_username "${vault_secret}")
  password=$(gcp::read-secret-field elasticsearch_password "${vault_secret}")
  echo "Authorization: Basic $(echo -n "${username}:${password}" | base64)"
}

#####################################
# Get the cluster state secret Path in Vault
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the secret path
#####################################
function elasticsearch::get-vault-secret(){
  local cluster=${1:?"Missing secret argument"}
  log::debug "Getting the secret path for ${cluster}"
  echo "oblt-clusters_${cluster}_cluster-state"
}

#####################################
# Get the Elasticsearch endpoint
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the endpoint
#####################################
function elasticsearch::get-endpoint(){
  local cluster=${1:?"Missing cluster argument"}
  local vault_secret
  log::debug "Getting Elasticsearch endpoint for ${cluster}"
  vault_secret=$(elasticsearch::get-vault-secret "${cluster}")
  gcp::read-secret-field elasticsearch_url "${vault_secret}"
}

#####################################
# Increase the number of shards per node
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
#   value (number of shards per node)
# Returns:
#   None
# https://www.elastic.co/guide/en/elasticsearch/reference/current/increase-cluster-shard-limit.html
#####################################
function elasticsearch::set-max-shards-per-node(){
  local cluster=${1:?"Missing cluster argument"}
  local value=${2:?"Missing value argument"}
  log::debug "Setting max shards per node to ${value} for ${cluster}"
  elasticsearch::api "${cluster}" PUT "_cluster/settings" "{
    \"persistent\": {
      \"cluster.max_shards_per_node\": ${value}
      }
    }"
}

#####################################
# Get the cluster health
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
# Returns:
#   the cluster health
#####################################
function elasticsearch::get-cluster-health(){
  local cluster=${1:?"Missing cluster argument"}
  log::debug "Getting cluster health for ${cluster}"
  elasticsearch::api "${cluster}" GET "_cluster/health"
}

#####################################
# Get the explanation for a shard allocation
# Globals:
#   None
# Arguments:
#   cluster (edge-oblt)
#   index (e.g. .kibana)(Optional)
# Returns:
#   the shard explanation
#####################################
function elasticsearch::explain(){
  local cluster=${1:?"Missing cluster argument"}
  local index=${2:-""}
  log::debug "Getting shard explanation for ${cluster}/${index}"
  if [ -n "${index}" ]; then
    index="{\"index\": \"${index}\"}"
  fi
  echo elasticsearch::api "${cluster}" GET "_cluster/allocation/explain?pretty" "${index}"
}
