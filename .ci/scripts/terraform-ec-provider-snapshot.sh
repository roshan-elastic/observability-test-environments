#!/usr/bin/env bash
# this script build and copy to a GCS bucket the terraform elatic provider provider
# for darwing and linux on amd64 and arm64
#

set -eu

# GO_VERSION=1.20
# git clone git@github.com:elastic/terraform-provider-ec.git
cd terraform-provider-ec
# docker run --rm -u "$(id -u):$(id -g)" -v "$(pwd)":/app -w /app -e HOME=/app \
#             golang:${GO_VERSION} \
#             make snapshot
make release-no-publish

for ARCH in amd64 arm64; do
  for OS in darwin linux
  do
      BINARY_FOLDER="terraform-provider-ec_${OS}_${ARCH}*"
      BINARY_LOCATION="dist/${BINARY_FOLDER}"
      PLUGIN_LOCATION=~/.terraform.d/plugins
      VERSION=$(grep const "./ec/version.go"|cut -d "=" -f 2|tr -d ' "'|cut -d '-' -f 1)
      BINARY="terraform-provider-ec_v${VERSION}"
      PLUGIN_0_13="registry.terraform.io/elastic/ec/${VERSION}/${OS}_${ARCH}/${BINARY}"
      FULL_PATH=${PLUGIN_LOCATION}/${PLUGIN_0_13}


      mkdir -p "${PLUGIN_LOCATION}"
      # shellcheck disable=SC2086
      cp ${BINARY_LOCATION}/* "${PLUGIN_LOCATION}/terraform-provider-ec"
      DIR_NAME=$(dirname "${FULL_PATH}")
      mkdir -p "${DIR_NAME}"
      # shellcheck disable=SC2086
      cp ${BINARY_LOCATION}/* "${FULL_PATH}"
      echo "-> Copied terraform provider to ${FULL_PATH}"
  done
done

cd ~
find .terraform.d -name ".DS_Store" -exec rm -rf {} \;
tar -czf terraform_plugins.tar.gz .terraform.d

gsutil cp terraform_plugins.tar.gz gs://oblt-clusters/tools/terraform_plugins.tar.gz
