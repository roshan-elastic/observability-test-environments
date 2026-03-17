#!/usr/bin/env bash
# Licensed to Elasticsearch B.V. under one or more contributor
# license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright
# ownership. Elasticsearch B.V. licenses this file to you under
# the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
set -euo pipefail
set +x

if [ $# -lt 1 ]; then
  echo "${0} VERSION PREFIX [SHA|'latest'] [PUSH_VERSION|\${VERSION}] [REGISTRY|'docker.elastic.co']"
  echo FORCE=true force to make the pull and push
  exit 1
fi

VERSION=${1:?-"missing argument"}
PREFIX=${2:-"observability-ci"}
SHA=${3:-"latest"}
PUSH_VERSION=${4:-"${VERSION}"}
REGISTRY=${5:-"docker.elastic.co"}
DAY_OF_YEAR=$(date +%j)
MAJOR=$(echo "${VERSION}"|cut -d "." -f 1)
NO_KPI_URL_PARAM="x-elastic-no-kpi=true"
LIST_OF_IMAGE=()
MANIFEST="./metadata-${VERSION}.txt"
BUILDS="./builds-${VERSION}.txt"
DOCKER_IMAGES_LIST="./docker-images-${VERSION}.txt"
FORCE=${FORCE:-""}
MINOR=$(echo "${VERSION}"|cut -d "." -f 2)
PATCH=$(echo "${VERSION}"|cut -d "." -f 3|cut -d "-" -f 1)


rm -f "${MANIFEST}" "${BUILDS}" "${DOCKER_IMAGES_LIST}"

function saveManifest(){
  local version=$1
  local build_id=$2
  curl -sSf "https://artifacts-api.elastic.co/v1/versions/${version}/builds/${build_id}/?${NO_KPI_URL_PARAM}" | jq 'del(.manifests)' > "${MANIFEST}"
  curl -sSf "https://artifacts-api.elastic.co/v1/versions/${version}?${NO_KPI_URL_PARAM}" | jq 'del(.manifests)' > "${BUILDS}"
}

function tagAndPush(){
  local image=${1}
  local new_image=${2}
  docker tag "${image}" "${new_image}"
  docker push "${new_image}"
  LIST_OF_IMAGE+=("${new_image}")
}

function publish() {
  local image=${1}
  IMAGE_NO_PREFIX=$(basename "${image}")
  IMAGE_NO_TAG=${REGISTRY}/${PREFIX}/${IMAGE_NO_PREFIX%:*}
  tagAndPush "${image}" "${IMAGE_NO_TAG}:${PUSH_VERSION}"
  tagAndPush "${image}" "${IMAGE_NO_TAG}:${PUSH_VERSION}-${DAY_OF_YEAR}"

  if [ "${MAJOR}" -gt 7 ]; then
    SHA=$(docker inspect --format='{{index .Config.Labels "org.label-schema.vcs-ref"}}' "${image}")
    if [ -n "${SHA}" ]; then
      SHA_SHORT=${SHA:0:7}
      tagAndPush "${image}" "${IMAGE_NO_TAG}:${PUSH_VERSION}-${SHA_SHORT}"
      tagAndPush "${image}" "${IMAGE_NO_TAG}:${SHA}"
    fi
  fi
  tagAndPush "${image}" "${IMAGE_NO_TAG}:${BUILD_ID}"
}

function processDockerImage(){
  local image="${1}"
  local pull_tag="${2}"
  DOCKER_IMAGE_PULL="${REGISTRY}/${image}:${pull_tag}"
  DOCKER_IMAGE="${REGISTRY}/${image}:${VERSION}"
  DOCKER_IMAGE_ID="${REGISTRY}/${image}:${BUILD_ID}"
  docker pull "${DOCKER_IMAGE_PULL}"
  docker tag "${DOCKER_IMAGE_PULL}" "${DOCKER_IMAGE_ID}"
  docker tag "${DOCKER_IMAGE_PULL}" "${DOCKER_IMAGE}"
  publish "${DOCKER_IMAGE}"
}

BUILD_ID=$(curl -sSL "https://artifacts-api.elastic.co/v1/versions/${VERSION}/builds/${SHA}?${NO_KPI_URL_PARAM}"|jq -r .build.build_id)

if [ -z "${FORCE}" ]; then
  (docker pull "${REGISTRY}/${PREFIX}/stack-done:${BUILD_ID}" > /dev/null 2>&1 && (touch .skip_push_and_pull)) || true
fi

if [ -f ".skip_push_and_pull" ]; then
  echo "🐳 Stack already processed"
  exit 0
fi

saveManifest "${VERSION}" "${BUILD_ID}"
if [ "${MAJOR}" -lt 8 ]; then
  PROFILE_PATH=""
else
  if [ "${MINOR}" -ge 13 ] || { [ "${MINOR}" -ge 12 ] && [ "${PATCH}" -ge 1 ]; } then
    PROFILE_PATH="staging/profiling-agent"
  else
    PROFILE_PATH="observability/profiling-agent"
  fi
fi
if grep -q "SNAPSHOT" <<< "${VERSION}"; then
  for image in elasticsearch/elasticsearch \
    kibana/kibana \
    apm/apm-server \
    beats/auditbeat \
    beats/filebeat \
    beats/heartbeat \
    beats/metricbeat \
    beats/packetbeat \
    elastic-agent/elastic-agent \
    cloud-release/elasticsearch-cloud \
    cloud-release/elasticsearch-cloud-ess \
    cloud-release/kibana-cloud \
    cloud-release/elastic-agent-cloud \
    observability/profiling-agent
  do
    processDockerImage "${image}" "${BUILD_ID}-SNAPSHOT"
  done
else
  for image in staging/elasticsearch \
    staging/kibana \
    staging/apm-server \
    staging/auditbeat \
    staging/filebeat \
    staging/heartbeat \
    staging/metricbeat \
    staging/packetbeat \
    staging/elastic-agent \
    cloud-release/elasticsearch-cloud \
    cloud-release/elasticsearch-cloud-ess \
    cloud-release/kibana-cloud \
    cloud-release/elastic-agent-cloud \
    ${PROFILE_PATH}
  do
    processDockerImage "${image}" "${BUILD_ID}"
  done
fi

docker pull alpine:latest
docker tag alpine:latest "${REGISTRY}/${PREFIX}/stack-done:${BUILD_ID}"
docker push "${REGISTRY}/${PREFIX}/stack-done:${BUILD_ID}"

echo "🐳 List of images to pushed" > "${DOCKER_IMAGES_LIST}"
echo "------------------------" >> "${DOCKER_IMAGES_LIST}"
for value in "${LIST_OF_IMAGE[@]}"
do
    echo "* ${value}" >> "${DOCKER_IMAGES_LIST}"
done

if [ "${CI}" == "true" ]; then
  cat "${DOCKER_IMAGES_LIST}" >> "${GITHUB_STEP_SUMMARY}"
fi
