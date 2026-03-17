#!/usr/bin/env bash
# script to automate the generation of corpora for the k8s cluster
# https://github.com/elastic/observability-dev/blob/main/docs/infraobs/cloudnative-monitoring/dev-docs/elastic-generator-tool-with-rally.md
set -e

CORPORA_GENERATED_FOLDER="${HOME}/Library/Application Support/elastic-integration-corpus-generator-tool/corpora"
OUTPUT_FOLDER="${PWD}/generated"
# git clone git@github.com:elastic/rally-tracks.git
RALLY_TRACK_SOURCE="${HOME}/src/rally-tracks/tsdb_k8s_queries"
# git clone git@github.com:elastic/elastic-integration-corpus-generator-tool.git
CORPUS_TOOL_SOURCE="${HOME}/src/elastic-integration-corpus-generator-tool"
CORPUS_TOOL="${CORPUS_TOOL_SOURCE}/elastic-integration-corpus-generator-tool"
SIZE=$((1000*1024))
NOW=$(gdate -Iseconds)
mkdir -p "${OUTPUT_FOLDER}"
cp -r "${RALLY_TRACK_SOURCE}"/* "${OUTPUT_FOLDER}"
rm -fr  "${CORPORA_GENERATED_FOLDER}"
${CORPUS_TOOL} generate-with-template \
  "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.pod/gotext.tpl" \
  "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.pod/fields.yml" \
  -c "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.pod/configs.yml" \
  -y gotext \
  -t "${SIZE}" \
  -n "${NOW}"
TPL_FILE=$(ls -1 "${CORPORA_GENERATED_FOLDER}")
mv "${CORPORA_GENERATED_FOLDER}/${TPL_FILE}" "${OUTPUT_FOLDER}/doc-ds-metrics_k8s-pod.json"
rm -fr  "${CORPORA_GENERATED_FOLDER}"
${CORPUS_TOOL} generate-with-template \
  "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.container/gotext.tpl" \
  "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.container/fields.yml" \
  -c "${CORPUS_TOOL_SOURCE}/assets/templates/kubernetes.container/configs.yml" \
  -y gotext \
  -t "${SIZE}" \
  -n "${NOW}"
TPL_FILE=$(ls -1 "${CORPORA_GENERATED_FOLDER}")
mv "${CORPORA_GENERATED_FOLDER}/${TPL_FILE}" "${OUTPUT_FOLDER}/doc-ds-metrics_k8s-container.json"

cd "${OUTPUT_FOLDER}" || exit 1

cat <<EOF > track.json
{% import "rally.helpers" as rally with context %}

{% set search_iterations = 20 %}
{% set search_warmup_iterations = 50 %}
{% set target_interval = 4 %}
{% set touch_bulk_indexing_clients = touch_bulk_indexing_clients | default(3) %}
{% set touch_bulk_size = touch_bulk_size | default(50) %}

{% set end_time = "2023-05-16T21:00:00.000Z" %}
{% set time_intervals = {"15_minutes": ["30s", "2023-05-16T20:45:00.000Z"], "2_hours": ["1m", "2023-05-16T19:00:00.000Z"], "24_hours": ["30m", "2023-05-15T21:00:00.000Z"]} %}
{
  "version": 2,
  "description": "metricbeat information for elastic-app k8s cluster",
  "composable-templates": [
    {
      "name": "pod-template",
      "index-pattern": "k8s-pod*",
      "delete-matching-indices": true,
      "template": "pod-template.json"
    },
    {
      "name": "container-template",
      "index-pattern": "k8s-container*",
      "delete-matching-indices": true,
      "template": "container-template.json"
    }
  ],
  "data-streams": [
    {
      "name": "k8s-pod"
    },
    {
      "name": "k8s-container"
    }
  ],
  "corpora": [
EOF

N=0
for file in doc-ds-*.json; do
  NAME=$(basename "${file}" .json)
  DATASTREAM=$(echo "${NAME}" | cut -d'_' -f2)
  sed -i .bak '/^$/d' "${file}"
  rm -f "${file}.bak"
  jlsort -k '@timestamp' "${file}" > "${file}.sorted"
  mv "${file}.sorted" "${file}"
  head -n 1000 "${file}" > "${file}.1k"
  mv "${file}.1k" "${NAME}-1k.json"
  7z a "${file}.bz2" "${file}"
  DOCUMENTS=$(wc -l < "${file}")
  if [ "${N}" -gt 0 ]; then
    echo "    ,">> track.json
  fi
  cat <<EOF >> track.json
    {
      "name": "${DATASTREAM}",
      "documents": [
        {
          "target-data-stream": "${DATASTREAM}",
          "source-file": "${file}.bz2",
          "document-count": ${DOCUMENTS}
        }
      ]
    }
EOF
  N=$((N+1))
done

cat <<EOF >> track.json
  ],
  "operations": [
    {{ rally.collect(parts="operations/*.json") }}
  ],
  "challenges": [
    {{ rally.collect(parts="challenges/*.json") }}
  ]
}
EOF
