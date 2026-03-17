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

echo "usage: $0 PRODUCT VERSION [DAY_OF_YEAR]"
PRODUCT=${1:?-"missing argument"}
VERSION=${2:?-"missing argument"}
DAY_OF_YEAR=${3:-"$(date +%j)"}
YESTERDAY=$(echo "${DAY_OF_YEAR}" | awk '{printf "%03.0f", $1 - 1 }')
ROLLBACK_DAY_OF_YEAR=${3:-"${YESTERDAY}"}
PREFIX='observability-ci'
REGISTRY="docker.elastic.co"

echo "Rollback ${PRODUCT}:${VERSION} to ${ROLLBACK_DAY_OF_YEAR} day of the year"

IMAGE="${REGISTRY}/${PREFIX}/${PRODUCT}:${VERSION}"

docker pull "${IMAGE}-${ROLLBACK_DAY_OF_YEAR}"
docker tag "${IMAGE}-${ROLLBACK_DAY_OF_YEAR}" "${IMAGE}"
docker push "${IMAGE}"
