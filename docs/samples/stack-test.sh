#!/usr/bin/env bash
#{% raw %}
# Sample script to provision the stack an make a call to the API of Elasticsearch, APM, and Kibana
set -euo pipefail

# Configure environment variables
PATH=${HOME}/bin:${PATH}
GITHUB_TOKEN=$(gh auth token)
GITHUB_USERNAME=your-username
SLACK_CHANNEL=#slack-ckannel
SERVICE_ACCOUNT_EMAIL=my-service-account-email@elastic-observability.iam.gserviceaccount.com
GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json

CLUSTER_INFO_FILE=${PWD}/cluster-info.json
ENV_FILE=${PWD}/env.sh
BIN_DIR=${HOME}/bin
OBLT_CLI_INSTALL_SCRIPT=${PWD}/install.sh

export GITHUB_TOKEN

echo  "Configure cleanup"
function cleanup() {
  echo "Cleanup"
  if [ -f "${CLUSTER_INFO_FILE}" ]; then
    CLUSTER_NAME=$(jq -r .ClusterName "${CLUSTER_INFO_FILE}")
    oblt-cli cluster destroy --cluster-name "${CLUSTER_NAME}" --force
  fi
  rm -f "${CLUSTER_INFO_FILE}"
  rm -f "${ENV_FILE}"
  rm -f "${OBLT_CLI_INSTALL_SCRIPT}"
}

trap cleanup EXIT

echo "Install oblt-cli"
SCRIPT_URL=https://raw.githubusercontent.com/elastic/observability-test-environments/main/tools/oblt-cli/scripts/install.sh
curl -fsSL -o "${OBLT_CLI_INSTALL_SCRIPT}" -H "Authorization: Bearer ${GITHUB_TOKEN}" "${SCRIPT_URL}"
chmod ugo+x "${OBLT_CLI_INSTALL_SCRIPT}"
"${OBLT_CLI_INSTALL_SCRIPT}" "${BIN_DIR}"
# Go way
# GOPRIVATE=github.com/elastic go install github.com/elastic/observability-test-environments/tools/oblt-cli@6.2.1

echo "Provision the stack"
ELASTIC_STACK_VERSION=8.10.0-SNAPSHOT
oblt-cli cluster create ess \
  --stack-version "${ELASTIC_STACK_VERSION}" \
  --cluster-name-prefix my-job \
  --output-file "${CLUSTER_INFO_FILE}" \
  --wait 15 \
  --username "${GITHUB_USERNAME}" \
  --slack-channel "${SLACK_CHANNEL}" \
  --git-http-mode \
  --save-config

echo "Activate service account"
gcloud auth activate-service-account "${SERVICE_ACCOUNT_EMAIL}" --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"

echo "Get credntials"
CLUSTER_NAME=$(jq -r .ClusterName "${CLUSTER_INFO_FILE}")
oblt-cli cluster secrets env --cluster-name "${CLUSTER_NAME}" --output-file "${ENV_FILE}"

echo "Source the credentials"
# SC1090 SC1091
# shellcheck source=/dev/null
source "${ENV_FILE}"

echo "Make some calls to the API"
curl -X GET -sSfL \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  "${ELASTICSEARCH_HOST}/_cluster/health" | jq .

curl -X GET -sSfL \
  -H "kbn-xsrf: true" \
  -H 'Content-Type: application/json' \
  -u "${KIBANA_USERNAME}:${KIBANA_PASSWORD}" \
  "${KIBANA_HOST}/api/status" | jq .

curl -X GET -sSfL \
  -H "Authorization: Bearer ${ELASTIC_APM_SECRET_TOKEN}" \
  "${ELASTIC_APM_SERVER_URL}/" | jq .
#{% endraw %}
