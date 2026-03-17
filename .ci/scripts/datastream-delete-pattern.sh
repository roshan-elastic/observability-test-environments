#!/usr/bin/env bash
#
# Script to delete all red index and try to restore them from a snapshot.
# it saves a file with the name os the datastream (datastream.txt)
#
# .ci/scripts/datastreams-delete-pattern.sh dev-next "*2022*"
#
CURDIR=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${CURDIR}/lib/epr.sh"
# shellcheck source=/dev/null
. "${CURDIR}/lib/elasticsearch.sh"

log::info "usage: $0 CLUSTER PATTERN"

CLUSTER=${1:?-"missing argument"}
PATTERN=${2:?-"missing argument"}

DATASTREAMS=$(elasticsearch::api "${CLUSTER}" GET "_data_stream/${PATTERN}" | jq -c -r '.data_streams[].name')

echo "${DATASTREAMS}" > datastreams.txt
cat datastreams.txt
read -r -p "Press enter to continue"

for DATASTREAM in ${DATASTREAMS}
do
  echo "Delete ${DATASTREAM} :"
  elasticsearch::api "${CLUSTER}" DELETE "_data_stream/${DATASTREAM}"
  echo ""
done
