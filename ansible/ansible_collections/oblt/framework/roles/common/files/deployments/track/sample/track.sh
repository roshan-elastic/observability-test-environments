#!/usr/bin/env bash
# This script run a race with the track against a cluster.
# It expects the following environment variables to be set:
# - ELASTICSEARCH_HOSTS: The list of hosts to use to connect to Elasticsearch (e.g. "https://es.example.com:9200").
# - ELASTICSEARCH_USERNAME: The optional user name to use when connecting to Elasticsearch.
# - ELASTICSEARCH_PASSWORD: The optional password to use when connecting to Elasticsearch.
# - MONITORING_ELASTICSEARCH_HOSTS: The list of hosts to use to connect to the monitoring cluster (e.g. "https://monitoring.example.com:9200").
# - MONITORING_ELASTICSEARCH_USERNAME: The optional user name to use when connecting to the monitoring cluster.
# - MONITORING_ELASTICSEARCH_PASSWORD: The optional password to use when connecting to the monitoring cluster.

set -euo pipefail

MONITORING_KIBANA_HOST=${MONITORING_ELASTICSEARCH_HOST/.es./.kb.}
CLUSTER_UUID=$(curl -sL -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" "${ELASTICSEARCH_HOST}" | python -c "import sys, json; print(json.load(sys.stdin)['cluster_uuid'])")
CLUSTER_NAME=$(curl -sL -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" "${ELASTICSEARCH_HOST}" | python -c "import sys, json; print(json.load(sys.stdin)['cluster_name'])")
RALLY_TRACK=sample-track
TRANSACTION_NAME="GET /{index}/_search"
TRANSACTION_URL_PATH="/geonames/_search"

START_TIMESTAMP_ISO=$(date -u "+%Y-%m-%dT%H:%M:%S.000Z")
case "$OSTYPE" in
  darwin*) RACE_ID=$(uuidgen) ;;
  linux*)  RACE_ID=$(cat /proc/sys/kernel/random/uuid) ;;
  *)        echo "unknown: $OSTYPE"; exit 1;;
esac

esrally race \
  --track-path ./ \
  --target-hosts "${ELASTICSEARCH_HOST}" \
  --client-options "basic_auth_user:'${ELASTICSEARCH_USERNAME}',basic_auth_password:'${ELASTICSEARCH_PASSWORD}'" \
  --race-id "${RACE_ID}" \
  --cluster-name "${CLUSTER_NAME}" \
  --pipeline benchmark-only

STOP_TIMESTAMP_ISO=$(date -u "+%Y-%m-%dT%H:%M:%S.000Z")

cat <<EOF > /tmp/race.json
{
  "cluster_uuid": "${CLUSTER_UUID}",
  "cluster_name": "${CLUSTER_NAME}",
  "start": "${START_TIMESTAMP_ISO}",
  "stop": "${STOP_TIMESTAMP_ISO}",
  "dashboard": "${MONITORING_KIBANA_HOST}/app/dashboards#/view/a7d4baea-5956-4749-bffa-d289e5d660b4?_g=(time:(from:'${START_TIMESTAMP_ISO}',to:'${STOP_TIMESTAMP_ISO}'))&_a=(query:(language:kuery,query:'cluster_uuid:${CLUSTER_UUID}%20OR%20(labels.deploymentName:%22${CLUSTER_NAME}%22%20and%20data_stream.type:traces%20AND%20transaction.name:%20%22${TRANSACTION_NAME}%22%20AND%20url.path:%22${TRANSACTION_URL_PATH}%22)%20OR%20(labels.deploymentName:%22${CLUSTER_NAME}%22%20and%20data_stream.type:metrics)'))",
  "report": "${MONITORING_KIBANA_HOST}/api/reporting/generate/pngV2?jobParams=",
  "race_id": "${RACE_ID}",
  "track": "${RALLY_TRACK}",
  "report_job_params": {
    "browserTimezone": "Europe/Madrid",
    "layout": {
        "dimensions": {
            "height": 5440,
            "width": 660
        },
        "id": "preserve_layout"
    },
    "locatorParams": {
        "id": "DASHBOARD_APP_LOCATOR",
        "params": {
            "dashboardId": "a7d4baea-5956-4749-bffa-d289e5d660b4",
            "preserveSavedFilters": true,
            "query": {
                "language": "kuery",
                "query": "cluster_uuid:${CLUSTER_UUID} OR (labels.deploymentName:${CLUSTER_NAME} and data_stream.type:traces AND transaction.name: \"${TRANSACTION_NAME}\" AND url.path:\"${TRANSACTION_URL_PATH}\") OR (labels.deploymentName:${CLUSTER_NAME} and data_stream.type:metrics)"
            },
            "timeRange": {
                "from": "2024-01-18T18:14:18.000Z",
                "to": "2024-01-18T18:23:19.000Z"
            },
            "useHash": false,
            "viewMode": "view"
        }
    },
    "objectType": "dashboard",
    "title": "Query Analyser",
    "version": "8.12.0-SNAPSHOT"
  }
  "dashboard_params_g": {
    "time": {
      "from":"${START_TIMESTAMP_ISO}",
      "to":"${STOP_TIMESTAMP_ISO}"
    }
  },
  "dashboard_params_a":{
    "query":{
      "language":"kuery",
      "query":"cluster_uuid:${CLUSTER_UUID} OR (labels.deploymentName:"${CLUSTER_NAME}" and data_stream.type:traces AND transaction.name: "${TRANSACTION_NAME}" AND url.path:"${TRANSACTION_URL_PATH}") OR (labels.deploymentName:"${CLUSTER_NAME}" and data_stream.type:metrics)"
    }
  }
}
EOF

cat /tmp/race.json

JOB_PARAMS=$(python -c "import sys, json; print(json.dumps(json.load(sys.stdin)['report_job_params']), end='')" < /tmp/race.json| ./json2rison.py|python -c "import sys, urllib.parse; print(urllib.parse.quote(sys.stdin.read()), end='')")
echo "${JOB_PARAMS}"
curl -sSL -X POST \
  -u "${MONITORING_ELASTICSEARCH_USERNAME}:${MONITORING_ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  -H "kbn-xsrf: true" \
  "${MONITORING_KIBANA_HOST}/api/reporting/generate/pngV2?jobParams=${JOB_PARAMS}" > /tmp/report.json


python -c "import sys, json; print(json.dumps(json.load(sys.stdin), indent=2))" < /tmp/report.json
REPORT_URL=$(python -c "import sys, json; print(json.load(sys.stdin)['path'])" < /tmp/report.json)
sleep 30
curl -sSL -X GET \
  -u "${MONITORING_ELASTICSEARCH_USERNAME}:${MONITORING_ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  -H "kbn-xsrf: true" \
  "${MONITORING_KIBANA_HOST}${REPORT_URL}"
