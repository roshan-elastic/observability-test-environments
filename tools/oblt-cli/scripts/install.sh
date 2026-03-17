#!/usr/bin/env bash
# Script to install the latest version of oblt-cli
# usage: install.sh [INSTALL_DIR|.]

DEST_DIR=${1:-"."}
OWNER="elastic"
REPO="observability-test-environments"
LATEST=$(gh release list --repo "${OWNER}/${REPO}" --limit 1|cut -d$'\t' -f 1)
#LATEST=2.0.26
ARCH=$(uname -m)
OS=$(uname -s)
TMPDIR=$(mktemp -d)

case "${OS}-${ARCH}" in
  "Linux-x86_64")
  PATTERN="*linux_amd64*"
  ;;
  "Windows-x86_64")
  PATTERN="*windows_amd64*"
  ;;
  "Darwin-x86_64")
  PATTERN="*darwin_amd64*"
  ;;
  "Darwin-aarch64")
  PATTERN="*darwin_arm64*"
  ;;
  "Linux-aarch64")
  PATTERN="*linux_arm64*"
  ;;
esac;

gh release download "${LATEST}" --repo "${OWNER}/${REPO}" --dir "${TMPDIR}" --pattern "${PATTERN}"
# shellcheck disable=SC2086
tar -xzf "${TMPDIR}/"oblt-cli_${PATTERN}.tar.gz -C "${DEST_DIR}"
