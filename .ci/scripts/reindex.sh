#!/usr/bin/env bash
set -euo pipefail
set +x

echo "usage: $0 CLUSTER_SRC INDEX_SRC CLUSTER_DST INDEX_DST"
CLUSTER_SRC=${1:?-"missing argument"}
INDEX_SRC=${2:?-"missing argument"}
CLUSTER_DST=${3:?-"missing argument"}
INDEX_DST=${4:?-"missing argument"}


CURDIR=$(dirname "$0")
# shellcheck disable=SC1091
source "${CURDIR}/cluster-state.sh" "${CLUSTER_SRC}"

USERNAME_SRC=${USERNAME}
PASSWORD_SRC=${PASSWORD}
ES_URL_SRC=${ES_URL}

# shellcheck disable=SC1091
source "${CURDIR}/cluster-state.sh" "${CLUSTER_DST}"

USERNAME_SRC=${USERNAME}
PASSWORD_DST=${PASSWORD}
ES_URL_DST=${ES_URL}

echo "You have to add the setting 'reindex.remote.whitelist: [\"*.staging.foundit.no:*\",\"*.example.com:*\"]' to the cluster from you copy the data see https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html#reindex-from-remote"

curl -X POST -u "${USERNAME_DST}:${PASSWORD_DST}" "${ES_URL_DST}/_reindex?pretty" -H 'Content-Type: application/json' -d"
{
  \"source\": {
    \"index\": \"${INDEX_SRC}\",
    \"remote\": {
      \"host\": \"${ES_URL_SRC}\",
      \"username\": \"${USERNAME_SRC}\",
      \"password\": \"${PASSWORD_SRC}\"
    }
  },
  \"dest\": {
    \"index\": \"${INDEX_DST}\"
  }
}
"
