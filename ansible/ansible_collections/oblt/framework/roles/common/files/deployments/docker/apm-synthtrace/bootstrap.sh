#!/usr/bin/env bash
set -e

USE_BAZEL_VERSION=$(cat "${KIBANA_DIR}/.bazelversion")
USE_BASELISK_VERSION=$(cat "${KIBANA_DIR}/.bazeliskversion")
export USE_BAZEL_VERSION
echo "Using bazel version: ${USE_BAZEL_VERSION}"
echo "Using bazelisk version: ${USE_BASELISK_VERSION}"
npm install -u "@bazel/bazelisk@${USE_BASELISK_VERSION}"
export PATH=${WORKSPACE}/node_modules/.bin:$PATH

bazel --version
bazelisk --version

cd "${KIBANA_DIR}"
yarn kbn bootstrap
