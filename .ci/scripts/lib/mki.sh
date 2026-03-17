#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/install.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/teleport.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/elastic.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/k8s.sh"

#####################################
# Login in a MKI k8s cluster
# Globals:
#   None
# Arguments:
#   kibana_k8s_cluster Name of the k8s cluster
# Returns:
#   None
#####################################
function mki::login-k8s(){
  local kibana_k8s_cluster=${1:?"Missing kibana_k8s_cluster argument"}
  log::debug "Login in the k8s cluster ${kibana_k8s_cluster}"
  tsh kube login "${kibana_k8s_cluster}"
}

#####################################
# Install the required tools to login in a MKI k8s cluster
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function mki::install-tools(){
  log::debug "Install the required tools to login in a MKI k8s cluster"
  install::tsh
  install::kubectl
  install::jq
  install::yq
}

#####################################
# Run the required steps to login in a MKI k8s cluster
# Globals:
#   None
# Arguments:
#   environment Name of the environment
# Returns:
#   None
#####################################
function mki::login(){
  local environment=${1:?"Missing environment argument"}
  log::debug "Login in the environment ${environment}"
  mki::install-tools
  teleport::login "${environment}"
}

#####################################
# Get Projects in the Org
# Globals:
#   None
# Arguments:
#   console URL to the Admin console
#   project_type Type of Project (elasticsearch, observability, or security)
#   org_id Organization ID
#   console_api_key API Key to access the Admin console
# Returns:
#   A JSON list with the projects deployed in the Org
#####################################
function mki::get-projects-org(){
  local console=$1
  local project_type=$2
  local org_id=$3
  local console_api_key=$4
  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  log::debug "Get projects in the organization ${org_id}"
  # shellcheck disable=SC2086
  curl -X "GET" ${curl_args} \
    -H "Authorization: ApiKey ${console_api_key}" \
    "${console}/api/v1/admin/serverless/projects/${project_type}?organization_id=${org_id}"
}

#####################################
# Get Project JSON by Name
# Globals:
#   None
# Arguments:
#   console URL to the Admin console
#   project_type Type of Project (elasticsearch, observability, or security)
#   org_id Organization ID
#   console_api_key API Key to access the Admin console
#   project_name Name of the project
# Returns:
#   A JSON with the project details
#####################################
function mki::get-project-by-name(){
  local console=$1
  local project_type=$2
  local org_id=$3
  local console_api_key=$4
  local project_name=$5
  local projectsListJson
  projectsListJson=$(mki::get-projects-org "${console}" "${project_type}" "${org_id}" "${console_api_key}")
  local project_id
  project_id=$(echo "${projectsListJson}"|jq -r ".items[]|select(.name==\"${project_name}\")|.id")
  log::debug "Get project ${project_name} with ID ${project_id}"
  mki:get-project "${console}" "${project_type}" "${console_api_key}" "${project_id}"
}

#####################################
# Get Project JSON by ID
# Globals:
#   None
# Arguments:
#   console URL to the Admin console
#   project_type Type of Project (elasticsearch, observability, or security)
#   console_api_key API Key to access the Admin console
#   project_id ID of the project
# Returns:
#   A JSON with the project details
#####################################
function mki::get-project(){
  local console=$1
  local project_type=$2
  local console_api_key=$3
  local project_id=$4

  log::debug "Get project with ID ${project_id}"

  local curl_args=-sLSf
  if [ -n "${DEBUG}" ]; then
    curl_args="${curl_args} -v"
  fi

  # shellcheck disable=SC2086
  curl -X "GET" ${curl_args} \
    -H "Authorization: ApiKey ${console_api_key}" \
    "${console}/api/v1/admin/serverless/projects/${project_type}/${project_id}"
}

#####################################
# Get Kibana YAML file from the k8s secret
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   A YAML string with the Kibana configuration
#####################################
function mki::get-kibana-yaml(){
  local namespace=${1:?"Missing namespace argument"}
  log::debug "Get Kibana YAML from the k8s secret in the namespace ${namespace}"
  kubectl -n "${namespace}" get secret kb-kb-config --template="{{index .data \"kibana.yml\" | base64decode}}"
}

#####################################
# Get Kibana Reporting Key from the kibana yaml k8s secret
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   A string with the Kibana Reporting Key (xpack.reporting.encryptionKey)
#####################################
function mki:get-kibana-reporting-key(){
  local namespace=${1:?"Missing namespace argument"}
  local kibana_yaml
  kibana_yaml=$(mki::get-kibana-yaml "${namespace}")
  log::debug "Get Kibana Reporting Key from the k8s secret in the namespace ${namespace}"
  echo "${kibana_yaml}"|yq .xpack.reporting.encryptionKey
}

#####################################
# Get Kibana Security Key from the kibana yaml k8s secret
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   A string with the Kibana Security Key (xpack.security.encryptionKey)
#####################################
function mki::get-kibana-security-key(){
  local namespace=${1:?"Missing namespace argument"}
  local kibana_yaml
  kibana_yaml=$(mki::get-kibana-yaml "${namespace}")
  log::debug "Get Kibana Security Key from the k8s secret in the namespace ${namespace}"
  echo "${kibana_yaml}"|yq .xpack.security.encryptionKey
}

#####################################
# Get Kibana Encrypted Saved Objects Key from the kibana yaml k8s secret
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   A string with the Kibana Encrypted Saved Objects Key (xpack.encryptedSavedObjects.encryptionKey)
#####################################
function mki::get-kibana-encrypted-save-objects-key(){
  local namespace=${1:?"Missing namespace argument"}
  local kibana_yaml
  kibana_yaml=$(mki::get-kibana-yaml "${namespace}")
  log::debug "Get Kibana Encrypted Saved Objects Key from the k8s secret in the namespace ${namespace}"
  echo "${kibana_yaml}"|yq .xpack.encryptedSavedObjects.encryptionKey
}

#####################################
# Get Kibana Service Account Token from the kibana yaml k8s secret
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   A string with the Kibana Service Account Token (elasticsearch.serviceAccountToken)
#####################################
function mki::get-kibana-service-account-token(){
  local namespace=${1:?"Missing namespace argument"}
  kibana_yaml=$(mki::get-kibana-yaml "${namespace}")
  log::debug "Get Kibana Service Account Token from the k8s secret in the namespace ${namespace}"
  echo "${kibana_yaml}"|yq .elasticsearch.serviceAccountToken
}

#####################################
# Generate Kibana YAML file to connect to the servreless project
# Globals:
#   ELASTICSEARCH_HOST Hostname of the Elasticsearch cluster
#   ELASTICSEARCH_PASSWORD Password of the Elasticsearch cluster
# Arguments:
#   namespace Name of the k8s project namespace
#   kibana_file full path of the generated kibana.yml file
# Returns:
#   None
#####################################
function mki::generate-kibana-yml(){
  local namespace=${1:?"Missing namespace argument"}
  local kibana_file=${2:?"Missing kibana_folder argument"}
  local reporting_key
  reporting_key=$(mki:get-kibana-reporting-key "${namespace}")
  local security_key
  security_key=$(mki::get-kibana-security-key "${namespace}")
  local encrypted_saved_objects_key
  encrypted_saved_objects_key=$(mki::get-kibana-encrypted-save-objects-key "${namespace}")
  local service_account_token
  service_account_token=$(mki::get-kibana-service-account-token "${namespace}")

  log::debug "Generate Kibana YAML file in the namespace ${namespace}"
  cat > "${kibana_file}" << EOF
elasticsearch.hosts: ${ELASTICSEARCH_HOST}
elasticsearch.serviceAccountToken: ${service_account_token}
elasticsearch.ssl.verificationMode: none
elasticsearch.ignoreVersionMismatch: true
migrations:
  skip: true
server.host: 0.0.0.0
xpack:
  reporting:
    encryptionKey: ${reporting_key}
  security:
    authc:
      providers:
        basic:
          cloud-basic:
            order: 150
    encryptionKey: ${security_key}
    loginAssistanceMessage: 'Credentials: ${ELASTICSEARCH_USERNAME} / ${ELASTICSEARCH_PASSWORD}'
  encryptedSavedObjects:
    encryptionKey: ${encrypted_saved_objects_key}
EOF
}
