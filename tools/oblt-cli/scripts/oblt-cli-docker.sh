#!/usr/bin/env bash
set -eu

GIT_AUTHOR_NAME="$(git config user.name)"
GIT_AUTHOR_EMAIL="$(git config user.email)"
OBLT_CLI_CONFIG="${HOME}/.oblt-cli/config.yaml"
SSH_KEY="${HOME}/.ssh/id_rsa"
USER_HOME="/home/nonroot"
VERSION="latest"

docker --version
touch "${OBLT_CLI_CONFIG}"

# shellcheck disable=SC2048,SC2086
docker run -it --rm \
  -e GIT_AUTHOR_NAME="${GIT_AUTHOR_NAME}"\
  -e GIT_AUTHOR_EMAIL="${GIT_AUTHOR_EMAIL}" \
  -e GIT_COMMITTER_NAME="${GIT_AUTHOR_NAME}" \
  -e GIT_COMMITTER_EMAIL="${GIT_AUTHOR_EMAIL}" \
  -v "${OBLT_CLI_CONFIG}:${USER_HOME}/.oblt-cli/config.yaml" \
  -v "${SSH_KEY}:${USER_HOME}/.ssh/id_rsa" \
  -v "${SSH_KEY}.pub:${USER_HOME}/.ssh/id_rsa.pub" \
  "docker.elastic.co/observability-ci/oblt-cli:${VERSION}" $*
