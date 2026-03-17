#!/usr/bin/env bash
# increases the number of fields allowed to insert in a component, this is used when we face the error `too many fields for...` when the index is a datastream.
set -euo pipefail
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER [COMPONENT|logs-elastic_agent.apm_server@custom]"
CLUSTER=${1:?-"missing argument"}
COMPONENT=${2:"logs-elastic_agent.apm_server@custom"}

elasticsearch::api "${CLUSTER}" GET "_component_template/${COMPONENT}" '{
  "template": {
    "settings": {
      "index.mapping.total_fields.limit": 2000
    }
  },
  "_meta": {
    "package": {
      "name": "elastic_agent"
    }
  }
}'
