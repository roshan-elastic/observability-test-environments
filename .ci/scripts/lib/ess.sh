#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/gcp.sh"

#####################################
# Get the Google secret for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the path to the secret in the Vault
#####################################
function ess::get-vault-secret(){
  local environment=${1:?"Missing secret argument"}
  echo "elastic-cloud-observability-team-${environment}"
}

#####################################
# Get the ESS endpoint for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the URL to the ESS endpoint
#####################################
function ess::get-ess-endpoint(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
  gcp::read-secret-field ess_endpoint "${vault_secret}"
}

#####################################
# Get the serverless endpoint for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the URL to the serverless endpoint
#####################################
function ess::get-serverless-endpoint(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
   gcp::read-secret-field serverless_endpoint "${vault_secret}"
}

#####################################
# Get the console endpoint for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the URL to the console endpoint
#####################################
function ess::get-console-endpoint(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
   gcp::read-secret-field console "${vault_secret}"
}

#####################################
# Get the API key for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the API key
#####################################
function ess::get-api-key(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
   gcp::read-secret-field apiKey "${vault_secret}"
}

#####################################
# Get the Console API key for the given environment
# Globals:
#   ESS_APIKEY: to use a different API key, by default it is get from the GCSM. (Optional)
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the API key
#####################################
function ess::get-console-api-key(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  if [ -z "${ESS_APIKEY}" ];  then
    vault_secret=$(ess::get-vault-secret "${environment}")
    gcp::read-secret-field consoleApiKey "${vault_secret}"
  else
    echo "${ESS_APIKEY}"
  fi
}

#####################################
# Get the organization ID for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the organization ID
#####################################
function ess::get-org-id(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
   gcp::read-secret-field organizationID "${vault_secret}"
}

#####################################
# Get the auth header for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the auth header
#####################################
function ess::get-auth-header(){
  local environment=${1:?"Missing environment argument"}
  local api_key
  api_key=$(ess::get-api-key "${environment}")
  echo "Authorization: ApiKey ${api_key}"
}

#####################################
# Get the console auth header for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the auth header
#####################################
function ess::get-console-auth-header(){
  local environment=${1:?"Missing environment argument"}
  local api_key
  api_key=$(ess::get-console-api-key "${environment}")
  echo "Authorization: ApiKey ${api_key}"
}

#####################################
# Send invitation to the given email in the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   email
# Returns:
#   None
#####################################
function ess::send-invitation(){
  local org_id
  local environment=${1:?"Missing environment argument"}
  local email=${2:?"Missing email argument"}
  log::info "Sending invitation to ${email} in ${environment}"
  org_id=$(ess::get-org-id "${environment}")
  ess::api "${environment}" POST "api/v1/organizations/${org_id}/invitations" "
  {
    \"emails\": [
        \"${email}\"
    ],
    \"role_assignments\": {
        \"organization\": [],
        \"deployment\": [
            {
                \"organization_id\": \"${org_id}\",
                \"role_id\": \"deployment-admin\",
                \"all\": true
            }
        ],
        \"project\": {
            \"elasticsearch\": [
                {
                    \"organization_id\": \"${org_id}\",
                    \"role_id\": \"elasticsearch-admin\",
                    \"all\": true
                }
            ],
            \"observability\": [
                {
                    \"organization_id\": \"${org_id}\",
                    \"role_id\": \"observability-admin\",
                    \"all\": true
                }
            ],
            \"security\": [
                {
                    \"organization_id\": \"${org_id}\",
                    \"role_id\": \"security-admin\",
                    \"all\": true
                }
            ]
        }
    }
  }"
}

#####################################
# Send invitation to the given email in all environments
# Globals:
#   None
# Arguments:
#   email
# Returns:
#   None
#####################################
function ess::send-invitation-all-env(){
  local email=${1:?"Missing email argument"}
  log::info "Sending invitation to ${email} in all environments"
  ess::send-invitation "staging" "${email}"
  ess::send-invitation "pro" "${email}"
  ess::send-invitation "qa" "${email}"
}

#####################################
# Get the ESS console endpoint for the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the URL to the Serverless endpoint
#####################################
function ess::get-console(){
  local environment=${1:?"Missing environment argument"}
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
   gcp::read-secret-field console "${vault_secret}"
}

#####################################
# Make a request to the ESS API
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   method (GET, POST, PUT, DELETE)
#   api_path (e.g. api/v1/organizations)
#   body (JSON)
# Returns:
#   the response body
#####################################
function ess::api(){
  local environment=${1:?"Missing environment argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  log::info "Making and API call to ESS ${environment}"
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
  local endpoint
  endpoint=$(ess::get-ess-endpoint "${environment}")
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
# Make a request to the ESS Console API
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   method (GET, POST, PUT, DELETE)
#   api_path (e.g. api/v1/organizations)
#   body (JSON)
# Returns:
#   the response body
#####################################
function ess::console-api(){
  local environment=${1:?"Missing environment argument"}
  local method=${2:?"Missing method argument"}
  local api_path=${3:?"Missing API path argument"}
  local body=${4:-""}
  log::info "Making and API call to ESS ${environment}"
  local vault_secret
  vault_secret=$(ess::get-vault-secret "${environment}")
  local endpoint
  endpoint=$(ess::get-console-endpoint "${environment}")
  auth_header=$(ess::get-console-auth-header "${environment}")

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
# Search for deployments in the given environment and pattern
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   pattern (e.g. "observability-team-pro")
#   is_hidden (true or false)
#   status (started, stopped, initializing)
# Returns:
#   the response body
#####################################
function ess::search-deployments(){
  local environment=${1:?"Missing environment argument"}
  local pattern=${2:?"Missing pattern argument"}
  local is_hidden=${3:-"false"}
  local status=${4:-"started"}
  local org_id
  org_id=$(ess::get-org-id "${environment}")
  local query="{
    \"size\": 150,
    \"query\": {
      \"bool\": {
        \"must\": [
          {\"match\": { \"name\": { \"query\": \"${pattern}\", \"operator\": \"and\" } } },
          {\"match\": { \"metadata.organization_id\": { \"query\": \"${org_id}\" } } },
          {\"match\": { \"metadata.hidden\": { \"query\": \"${is_hidden}\" } } },
          {\"match\": { \"resources.elasticsearch.info.status\": { \"query\": \"${status}\" } } }
        ]
      }
    }
  }"
  ess::api "${environment}" POST "api/v1/deployments/_search" "${query}"
}

#####################################
# Search for All deployments in the given environment
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   pattern (e.g. "observability-team-pro")
#   is_hidden (true or false)
#   status (started, stopped, initializing)
# Returns:
#   the response body
#####################################
function ess::search-deployments-all(){
  local environment=${1:?"Missing environment argument"}
  local is_hidden=${2:-"false"}
  local status=${3:-"started"}
  local org_id
  org_id=$(ess::get-org-id "${environment}")
  local query="{
    \"size\": 150,
    \"query\": {
      \"bool\": {
        \"must\": [
          {\"match\": { \"metadata.organization_id\": { \"query\": \"${org_id}\" } } },
          {\"match\": { \"metadata.hidden\": { \"query\": \"${is_hidden}\" } } },
          {\"match\": { \"resources.elasticsearch.info.status\": { \"query\": \"${status}\" } } }
        ]
      }
    }
  }"
  ess::api "${environment}" POST "api/v1/deployments/_search" "${query}"
}

#####################################
# Search for deployments in the given environment and alias
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   alias (e.g. "observability-team-pro")
#   is_hidden (true or false)
# Returns:
#   the response body
#####################################
function ess::search-deployments-by-alias(){
  local environment=${1:?"Missing environment argument"}
  local alias=${2:?"Missing alias argument"}
  local is_hidden=${3:-"false"}
  local org_id
  org_id=$(ess::get-org-id "${environment}")
  local query="{
    \"size\": 150,
    \"query\": {
      \"bool\": {
        \"filter\": [
          {
            \"match_all\": {
              \"metadata.organization_id\": \"${org_id}\",
              \"metadata.hidden\": \"${is_hidden}\"
            }
          },
          {
            \"term\" : {
              \"alias\" : {
                \"value\" : \"${alias}\"
              }
            }
          }
        ]
      }
    }
  }"
  ess::api "${environment}" POST "api/v1/deployments/_search" "${query}"
}

#####################################
# Delete a deployment in the given environment identified by the given id
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   id (e.g. "1f9bc79e780cb636ca3907397e1dfeaa")
# Returns:
#   the response body
#####################################
function ess::delete-deployment(){
  local environment=${1:?"Missing environment argument"}
  local id=${2:?"Missing deployment id argument"}

  ess::api "${environment}" POST "api/v1/deployments/${id}/_shutdown?skip_snapshots=true&hide=true"
}

#####################################
# Get the deployment definition in the given environment identified by the given id
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   id (e.g. "1f9bc79e780cb636ca3907397e1dfeaa")
# Returns:
#   the response body
#####################################
function ess::get-deployment(){
  local environment=${1:?"Missing environment argument"}
  local id=${2:?"Missing deployment id argument"}

  ess::api "${environment}" GET "api/v1/deployments/${id}?convert_legacy_plans=true&=enrich_with_template=false&show_metadata=false&show_plan_defaults=false&show_plan_history=false&show_plan_logs=false&show_plans=false&show_security=false&show_settings=false&show_system_alerts=0"
}

#####################################
# Change a deployment maintenance mode
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   id (e.g. "1f9bc79e780cb636ca3907397e1dfeaa")
#   mode (start or stop)
# Returns:
#   the response body
#####################################
function ess::manteinance-mode(){
  local environment=${1:?"Missing environment argument"}
  local id=${2:?"Missing deployment id argument"}
  local mode=${3:?"Missing mode argument"}

  if [ "${mode}" != "start" ] && [ "${mode}" != "stop" ]; then
    log::error "mode must be start or stop"
    exit 1
  fi

  ess::api "${environment}" POST "api/v1/deployments/${id}/elasticsearch/main-elasticsearch/instances/maintenance-mode/_${mode}"
}

#####################################
# Restore a deployment in the given environment identified by the given id
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   id (e.g. "1f9bc79e780cb636ca3907397e1dfeaa")
# Returns:
#   the response body
#####################################
function ess::restore-cluster(){
  local environment=${1:?"Missing environment argument"}
  local id=${2:?"Missing deployment id argument"}

  ess::api "${environment}" POST "api/v1/deployments/${id}/_restore"
}

#####################################
# Create a user on ESS
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
#   email electronic mail address
#   password password
#   first_name first name
#   last_name last name
# Returns:
#   the response body
#####################################
function ess::create-user(){
  local environment=${1:?"Missing environment argument"}
  local email=${2:?"Missing email argument"}
  local password=${3:?"Missing password argument"}
  local first_name=${4:?"Missing first name argument"}
  local last_name=${5:?"Missing last name argument"}
  local payload
  payload=$(cat <<EOF
{
    "email": "${email}",
    "password": "${password}",
    "first_name": "${first_name}",
    "last_name": "${last_name}",
    "is_paying": false,
    "allowed_not_to_pay": true,
    "domain": "found",
    "company_info": {},
    "marketplace": {}
}
EOF
)
  ess::console-api "${environment}" POST "api/v1/saas/users" "${payload}"
}

#####################################
# List users in ESS
# Globals:
#   None
# Arguments:
#   environment (staging, pro, or qa)
# Returns:
#   the response body
#####################################
function ess:list-users(){
  local environment=${1:?"Missing environment argument"}
  local org_id
  org_id=$(ess::get-org-id "${environment}")
  log::debug "Listing users in ${environment}"
  ess::console-api "${environment}" GET "api/v1/organizations/${org_id}/members"
}
