#!/usr/bin/env bash

function installApt() {
  REPO_URL="https://packages.cloud.google.com/apt"
  sudo touch /etc/apt/sources.list.d/google-cloud-sdk.list
  echo "deb [signed-by=/etc/apt/trusted.gpg.d/cloud.google.gpg] $REPO_URL cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list >/dev/null
  curl -fsSLo - https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/cloud.google.gpg

  sudo apt update -y
  sudo apt install -y google-cloud-sdk-gke-gcloud-auth-plugin
}

OS=$(uname -s)

echo "Checking gke-gcloud-auth-plugin is installed"
(command -v gke-gcloud-auth-plugin &> /dev/null && touch /tmp/gke-gcloud-auth-plugin-installed) || true

if [ ! -f /tmp/gke-gcloud-auth-plugin-installed ] && { [ -z "${CI}" ] || [ "${OS}" == "Darwin" ]; }; then
  echo "Installing gke-gcloud-auth-plugin component"
  (gcloud components install gke-gcloud-auth-plugin && touch /tmp/gke-gcloud-auth-plugin-installed)|| true
fi

if [ ! -f /tmp/gke-gcloud-auth-plugin-installed ] && [ -n "${CI}" ] && [ "${OS}" == "Linux" ]; then
  echo "Installing gke-gcloud-auth-plugin via apt-get"
  installApt && touch /tmp/gke-gcloud-auth-plugin-installed
fi
