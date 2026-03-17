#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"

#####################################
# Login in GCP
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function gcp::gcp-login(){
  log::info "Login in GCP"
  gcloud auth login
}

#####################################
# Read a secret from GCP
# Globals:
#   None
# Arguments:
#   secret (string): The secret name
# Returns:
#   The secret value
#####################################
function gcp::read-secret(){
  local secret=$1
  log::debug "Reading secret ${secret}"
  gcloud secrets versions access latest --secret="${secret}" --project="elastic-observability"
}

#####################################
# Read a field from a secret in GCP
# Globals:
#   None
# Arguments:
#   field (string): The field to read
#   secret (string): The secret name
# Returns:
#   the secret field value
#####################################
function gcp::read-secret-field(){
  local field=$1
  local secret=$2
  log::debug "Reading field ${field} from secret ${secret}"
  gcp::read-secret "${secret}"|yq ".${field}"
}

#####################################
# list all the networks using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   List of networks
#####################################
function gcp::listNetworks(){
    local filter=$1
    log::debug "Listing networks with filter ${filter}"
    gcloud compute networks list --format json --filter="name~${filter}"|jq -c -r '.[].name'
}

#####################################
# Delete all the networks using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   None
#####################################
function gcp::deleteNetworks(){
    local filter=$1
    log::debug "Deleting networks with filter ${filter}"
    for i in $(gcp::listNetworks "${filter}"); do
        log::info "Deleting network ${i}"
        gcloud compute networks delete --quiet "${i}"
    done
}

#####################################
# list all VMs using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   List of VMs [{name,zone},{name,zone}]
#####################################
function gcp::listVMs(){
    local filter=$1
    gcloud compute instances list --format json --filter="name~${filter}"|jq -c -r '.[]|[.name,.zone]'
}

#####################################
# Delete all the VMs using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   None
#####################################
function gcp::deleteVMs(){
    local filter=$1
    log::debug "Deleting VMs with filter ${filter}"
    for i in $(gcp::listVMs "${filter}"); do
        log::info "Deleting VM ${i}"
        name=$(echo "${i}"|jq -r '.[0]')
        zone=$(echo "${i}"|jq -r '.[1]'|cut -d '/' -f 9)
        gcloud compute instances delete --quiet "${name}" --zone "${zone}"
    done
}


#####################################
# list all firewall rules using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   List of firewall rules
#####################################
function gcp::listFirewallRules(){
    local filter=$1
    log::debug "Listing firewall rules with filter ${filter}"
    gcloud compute firewall-rules list --format json --filter="name~${filter}"|jq -c -r '.[].name'
}

#####################################
# Delete all the firewall rules using a filter
# Globals:
#   None
# Arguments:
#   filter (regexp): The filter to use
# Returns:
#   None
#####################################
function gcp::deleteFirewallRules(){
    local filter=$1
    log::debug "Deleting firewall rules with filter ${filter}"
    for i in $(gcp::listFirewallRules "${filter}"); do
        log::info "Deleting firewall rule ${i}"
        gcloud compute firewall-rules delete --quiet "${i}"
    done
}
