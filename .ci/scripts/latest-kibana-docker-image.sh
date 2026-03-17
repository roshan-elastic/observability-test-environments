#!/usr/bin/env bash

TEMP_DIR=$(mktemp -d)
DONE_FILE="${TEMP_DIR}/done"

docker login docker.elastic.co -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}" 1>&2
for i in {0..30}
do
    SHA=$(gh api repos/elastic/kibana/commits --jq ".[${i}].sha"|cut -b 1-12)
    export DOCKER_IMAGE="docker.elastic.co/kibana-ci/kibana-serverless:git-${SHA}"
    echo "Checking image ${DOCKER_IMAGE}" 1>&2
    (docker manifest inspect "${DOCKER_IMAGE}" 1>&2 && touch "${DONE_FILE}") || true
    if [ -f "${DONE_FILE}" ]; then
        echo "Image ${DOCKER_IMAGE} exists" 1>&2
        echo "${DOCKER_IMAGE}"
        break
    fi
done

# TODO use recreate update strategy and an update on an existing cluster
# export CI=true
# export GITHUB_TOKEN=$(gh auth token)
# PARAMATER_FILE="${TEMP_DIR}/params.json"
#
# cat <<EOF > "${PARAMATER_FILE}"
# {
#     "ProjectType": "observability",
#     "KibanaDockerImage": "${DOCKER_IMAGE}"
# }
# EOF
#
# oblt-cli cluster destroy --cluster-name lite-serverless-oblt --wait 15 --force || true
# oblt-cli cluster create custom --template serverless --cluster-name lite-serverless-oblt --parameters-file "${PARAMATER_FILE}" --wait 15
