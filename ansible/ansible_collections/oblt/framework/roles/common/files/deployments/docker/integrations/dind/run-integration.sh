#!/bin/sh
# This script is used to run the integration tests for the docker integration.

set -ex

PACKAGE=${1:?-"Please specify the package name"}
BRANCH=${2:-"main"}
REPO=${3:-"elastic/integrations"}

PACKAGE_DIR="/usr/shared/data/integrations/packages/${PACKAGE}/_dev/deploy/docker"
CHECKOUT_DIR="/usr/shared/data/integrations"

if [ -d "${CHECKOUT_DIR}" ]; then
  echo "Pulling the latest changes from the integrations repo"
  cd  "${CHECKOUT_DIR}"
  git pull || true
else
  echo "Cloning the integrations repo"
  git clone -b "${BRANCH}" --depth 1 --single-branch "https://github.com/${REPO}.git" "${CHECKOUT_DIR}"
fi

if [ -d "${PACKAGE_DIR}" ]; then
  echo "Changing to folder ${PACKAGE_DIR}"
  cd "${PACKAGE_DIR}"
else
  echo "No docker integration found for '${PACKAGE}'"
  exit 1
fi

echo "Configure the docker-compose file"
# edit a docker-compose file to export the same port as the one used by the integration
PORTS=$(docker compose config --format json | jq -r -c '.services[].ports[].target' || true)
for p in ${PORTS};
do
  sed -i -e "s/- ${p}\$/- ${p}:${p}/" docker-compose.yml
done

echo "Starting the docker compose"
docker compose up
