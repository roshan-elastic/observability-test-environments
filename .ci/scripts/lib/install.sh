#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/prompt.sh"

#####################################
# Check if the brew is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::brew(){
  log::info "Checking if brew is installed"
  if [ -z "$(command -v brew)" ]; then
    log::error "You need to Install brew"
    prompt::askYN "Do you want to install brew?"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
}

#####################################
# Check if the VPB client is installed, it exits if the VPN client is not installed
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
# https://docs.elastic.dev/serverless/es-troubleshooting
function install::vpn(){
  log::info "Checking if Aviatrix VPN Client is installed"
  if [ ! -d "/Applications/Aviatrix VPN Client.app" ] && [ -z "$(command -v openvpn)" ]; then
    install::brew
    log::error "You need to Install Aviatrix VPN Client or and OpenVPN client"
    log::error http://ela.st/vpn
    prompt::askYN "Do you want to install OpenVPN Client?"
    brew install openvpn
  fi
  if [ -d "/Applications/Aviatrix VPN Client.app" ]; then
    log::info "Adding Aviatrix VPN OpenVPN Client to the PATH"
    export PATH="/Applications/Aviatrix VPN Client.app/Contents/Resources/bin/mac":${PATH}
  fi
}

#####################################
# Check if the vault client is installed, if not install it
# it uses brew to install vault
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::vault(){
  log::info "Checking if Vault is installed"
  if [ -z "$(command -v vault)" ]; then
    install::brew
    prompt::askYN "Do you want to install Vault?"
    brew tap hashicorp/tap
    brew install hashicorp/tap/vault
  fi
}

#####################################
# Check if the tsh client is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
# https://github.com/elastic/cloud/blob/master/wiki/teleport/USAGE.md
function install::tsh(){
  log::info "Checking if Teleport is installed"
  # https://cdn.teleport.dev/teleport-v14.0.3-darwin-amd64-bin.tar.gz
  # https://cdn.teleport.dev/teleport-v14.0.3-darwin-arm64-bin.tar.gz
  # https://cdn.teleport.dev/teleport-v14.0.3-linux-amd64-bin.tar.gz
  # https://cdn.teleport.dev/teleport-v14.0.3-linux-arm64-bin.tar.gz
  # https://cdn.teleport.dev/Teleport%20Connect%20Setup-14.0.3.exe
  ARCH=$(uname -m)
  OS=$(uname -s)
  TMPDIR=$(mktemp -d)
  TELEPORT_VERSION="14.1.1"

  case "${OS}-${ARCH}" in
    "Linux-x86_64")
    ARCH_BIN="linux-amd64"
    log::error  "Linux is not supported"
    exit 1
    ;;
    "Windows-x86_64")
    log::error  "Windows is not supported"
    exit 1
    ;;
    "Darwin-x86_64")
    ARCH_BIN="darwin-amd64"
    ;;
    "Darwin-aarch64")
    ARCH_BIN="darwin-arm64"
    ;;
    "Linux-aarch64")
    ARCH_BIN="linux-arm64"
    log::error  "Linux is not supported"
    exit 1
    ;;
  esac;

  if [ -z "$(command -v tsh)" ]; then
    log::error "You need to Install Teleport"
    prompt::askYN "Do you want to install Teleport?"
    curl -sSfL -o "${TMPDIR}/teleport.tgz" https://cdn.teleport.dev/teleport-v${TELEPORT_VERSION}-${ARCH_BIN}-bin.tar.gz
    mkdir -p "${HOME}/bin"
    tar -xzf "${TMPDIR}/teleport.tgz" -C "${TMPDIR}"
    cp "${TMPDIR}/teleport/tsh" "${HOME}/bin"
    rm -rf "${TMPDIR}"
  fi
}

#####################################
# Check if the nvm is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::nvm(){
  log::info "Checking if nvm is installed"
  local NVM_DIR="$HOME/.nvm"
  if [ -s "${NVM_DIR}/nvm.sh" ]; then
    # shellcheck source=/dev/null
    . "${NVM_DIR}/nvm.sh"
  fi
  if [ -z "$(command -v nvm)" ]; then
    log::error "You need to Install nvm"
    prompt::askYN "Do you want to install nvm?"
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
  fi
}

#####################################
# Check if the nvm is installed, if not install it
# Then install the node version
# Globals:
#   None
# Arguments:
#   node_version
# Returns:
#   None
#####################################
function install::node(){
  local node_version=${1:?"Missing node_version argument"}
  local NVM_DIR="$HOME/.nvm"
  log::info "Checking if Node.js (${node_version}) is installed"
  if [ -z "$(command -v node)" ] || [ "$(node --version)" != "v${node_version}" ]; then
    log::warn "Node.js (${node_version}) is not installed"
    log::warn "You have Node.js $(node --version) installed/default"
    prompt::askYN "Do you want to install/activate Node.js (${node_version}) with nvm?"
    install::nvm
    if [ -s "${NVM_DIR}/nvm.sh" ]; then
      # shellcheck source=/dev/null
      . "${NVM_DIR}/nvm.sh"
      nvm install "${node_version}"
    else
      log::error "nvm is not installed"
      exit 1
    fi
  fi
}

#####################################
# Check if the kubectl is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::kubectl(){
  log::info "Checking if kubectl is installed"
  if [ -z "$(command -v kubectl)" ]; then
    install::brew
    log::error "You need to Install kubectl"
    prompt::askYN "Do you want to install kubectl?"
    brew install kubectl
  fi
}

#####################################
# Check if the jq is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::jq(){
  log::info "Checking if jq is installed"
  if [ -z "$(command -v jq)" ]; then
    install::brew
    log::error "You need to Install jq"
    prompt::askYN "Do you want to install jq?"
    brew install jq
  fi
}

#####################################
# Check if the yq is installed, if not install it
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function install::yq(){
  log::info "Checking if yq is installed"
  if [ -z "$(command -v yq)" ]; then
    install::brew
    log::error "You need to Install yq"
    prompt::askYN "Do you want to install yq?"
    brew install yq
  fi
}

#####################################
# Check which OS is running
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   The name of the OS (linux, macos, windows)
#####################################
function install::whichOS(){
  ARCH=$(uname -m)
  OS=$(uname -s)
  log::info  "Running on ${OS}-${ARCH}"
  echo "${OS}"
}

#####################################
# Check if the OS is Linux
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   true if the OS is Linux, false otherwise
#####################################
function install::isLinux(){
  if [ "$(install::whichOS)" == "Linux" ]; then
    echo "true"
  else
    echo "false"
  fi
}

#####################################
# Check if the OS is macOS
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   true if the OS is macOS, false otherwise
#####################################
function install::isMacOS(){
  if [ "$(install::whichOS)" == "Darwin" ]; then
    echo "true"
  else
    echo "false"
  fi
}

#####################################
# Check if the OS is Windows
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   true if the OS is Windows, false otherwise
#####################################
function install::isWindows(){
  if [ "$(install::whichOS)" == "Windows" ]; then
    echo "true"
  else
    echo "false"
  fi
}
