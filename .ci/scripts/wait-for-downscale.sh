#!/usr/bin/env bash
# Licensed to Elasticsearch B.V. under one or more contributor
# license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright
# ownership. Elasticsearch B.V. licenses this file to you under
# the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
set -euo pipefail

NAMESPACE=${1:?"NAMESPACE is missing"}
OBJ_NAME=${2:?"OBJ_NAME is missing"}

TYPE=${OBJ_NAME%/*}

if [ "${TYPE}" = "deployments" ]; then
  FIELD="{.spec.replicas}"
fi

if [ "${TYPE}" = "statefulset" ]; then
  FIELD="{.status.replicas}"
fi

# shellcheck disable=SC2086
until [ "$(kubectl get -n ${NAMESPACE} ${OBJ_NAME} -o jsonpath=${FIELD})" = "0" ]
do
  echo "Waiting to stop ${OBJ_NAME}"
  sleep 5
done
