#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/prompt.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/install.sh"
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/gcp.sh"

OBLT_CLI_HOME=${OBLT_CLI_HOME:-"${HOME}/.oblt-cli"}

#####################################
# Check if the user is logged in Vault, if not login in Vault
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function elastic::vault-login(){
  log::debug "Checking if the user is logged in Vault"
  if ! vault token lookup > /dev/null ; then
    log::error "You need to login in Vault"
    log::error "Further information: https://docs.elastic.dev/ci/using-secrets#generating-a-github-token"
    prompt::askYN "Do you want to login in Vault?"
    vault login -method=github > /dev/null
  fi
}

#####################################
# Login in the VPN, it opens the VPN client and wait for the user to continue.
# Globals:
#   OBLT_CLI_HOME
# Arguments:
#   environment (staging or production)
# Returns:
#   None
#####################################
function elastic::vpn(){
  local environment=$1
  mkdir -p "${OBLT_CLI_HOME}/vpn"
  case "${environment}" in
    "qa")
      OVPN_FILE=${OBLT_CLI_HOME}/vpn/elastic-cloud-staging-us.ovpn
      OVPN_SECRET=secret/observability-team/ci/elastic-cloud/vpn-staging
      ;;
    "staging")
      OVPN_FILE=${OBLT_CLI_HOME}/vpn/elastic-cloud-staging-us.ovpn
      OVPN_SECRET=elastic-vpn-staging
      ;;
     "production")
      OVPN_FILE=${OBLT_CLI_HOME}/vpn/elastic-cloud-us.ovpn
      OVPN_SECRET=elastic-vpn-pro
      ;;
  esac
  gcp::read-secret "${OVPN_SECRET}" > "${OVPN_FILE}"
  log::info "Login in ${environment} VPN"
  # osascript -e 'POSIX path of (path to application "Aviatrix VPN Client")'
  log::info "Opening Aviatrix VPN Client..."
  log::warn "You have to manually add the file ${OVPN_FILE} to Aviatrix VPN Client"
  open -a "/Applications/Aviatrix VPN Client.app"
  prompt::waitEnter "When the VPN is connected (green)"
}

#####################################
# Import the VPN client CA certificate
# Globals:
#   None
# Arguments:
#   ca_cert_path: path to the CA certificate
#   configure_mark: path to the file that indicates that the CA certificate is installed
#   keychain_file: path to the keychain file
# Returns:
#   None
#####################################
function elastic::vpn-import-certs(){
  if [ "$(install::isMacOS)" == "true" ]; then
    elastic::vpn-import-certs-macOS "$@"
  fi

  if [ "$(install::isLinux)" == "true" ]; then
    elastic::vpn-import-certs-linux "$@"
  fi
}

#####################################
# Install the CA certificate in macOS
# Globals:
#   None
# Arguments:
#   ca_cert_path: path to the CA certificate
#   configure_mark: path to the file that indicates that the CA certificate is installed
#   keychain_file: path to the keychain file
# Returns:
#   None
#####################################
function elastic::vpn-import-certs-macOS(){
  local ca_cert_path=${1:-"${HOME}/.oblt-cli/vpn/rootCA.pem"}
  local configure_mark=${2:-"${HOME}/.oblt-cli/vpn/configure.mark"}
  local keychain_file=${3:-"/Library/Keychains/System.keychain"}

  if [ -f "${configure_mark}" ]; then
    prompt::dialog "CA certificate already installed

You can remove the file ${configure_mark} if you want to reimport the certificate"
    exit 0
  fi

  prompt::dialogYN "Do you want to import the VPN client CA certificate?
we will ask for your password twice" 1

  mkdir -p "$(dirname "${ca_cert_path}")"
  mkdir -p "$(dirname "${configure_mark}")"

  # this command has to run in a interactive terminal on macOS 11+
  # https://developer.apple.com/forums/thread/671582
  SCRIPT="#!/usr/bin/env bash
  sudo security add-trusted-cert -d -r trustRoot -k ${keychain_file} ${ca_cert_path} && touch ${configure_mark}
  exit \$?"

  elastic::run-on-terminal "${SCRIPT}"
}

#####################################
# Install the CA certificate in Linux
# Globals:
#   None
# Arguments:
#   ca_cert_path: path to the CA certificate
#   configure_mark: path to the file that indicates that the CA certificate is installed
#   keychain_file: path to the keychain file
# Returns:
#   None
#####################################
function elastic::vpn-import-certs-linux(){
  local ca_cert_path=${1:-"${HOME}/.oblt-cli/vpn/rootCA.pem"}
  local configure_mark=${2:-"${HOME}/.oblt-cli/vpn/configure.mark"}
  local keychain_file=${3:-"/Library/Keychains/System.keychain"}
  prompt::askYN "Do you want to import the VPN client CA certificate?
we will ask for your password"
  if [ -d "/etc/pki/ca-trust/source/anchors" ]; then
    log::info "Installing CA certificate in /etc/pki/ca-trust/source/anchors"
    sudo cp "${ca_cert_path}" /etc/pki/ca-trust/source/anchors/
    sudo update-ca-trust
  fi
  if [ -d "/usr/local/share/ca-certificates" ]; then
    log::info "Installing CA certificate in /usr/local/share/ca-certificates"
    sudo cp "${ca_cert_path}" /usr/local/share/ca-certificates/
    sudo update-ca-certificates
  fi
}

#####################################
# Run the VPN client in macOS
# Globals:
#   None
# Arguments:
#   openvpn_config_file: path to the openvpn config file
#   openvpn_auth_file: path to the openvpn auth file
#   openvpn_log_file: path to the openvpn log file
#   openvpn_pid_file: path to the openvpn pid file
# Returns:
#   None
#####################################
function elastic::vpn-run(){
  local openvpn_config_file=${1:-"${HOME}/.oblt-cli/vpn/production.ovpn"}
  local openvpn_auth_file=${2:-"${HOME}/.oblt-cli/vpn/auth.txt"}
  local openvpn_log_file=${3:-"${HOME}/.oblt-cli/vpn/openvpn.log"}
  local openvpn_pid_file=${4:-"${HOME}/.oblt-cli/vpn/openvpn.pid"}

  install::vpn

  log::warn "Any other VPN client must be closed before running this command"
  log::warn "Starting VPN client it requires your password"

  # shellcheck disable=SC2086
  elastic::vpn-stop ${openvpn_pid_file}

  mkdir -p "$(dirname "${openvpn_config_file}")"
  mkdir -p "$(dirname "${openvpn_auth_file}")"
  mkdir -p "$(dirname "${openvpn_log_file}")"
  mkdir -p "$(dirname "${openvpn_pid_file}")"

  if [ "$(install::isMacOS)" == "true" ]; then
    elastic::vpn-run-macos "$@"
    return
  fi
  if [ "$(install::isLinux)" == "true" ]; then
    elastic::vpn-run-linux "$@"
    return
  fi
}

#####################################
# Run the VPN client in macOS
# Globals:
#   None
# Arguments:
#   openvpn_config_file: path to the openvpn config file
#   openvpn_auth_file: path to the openvpn auth file
#   openvpn_log_file: path to the openvpn log file
#   openvpn_pid_file: path to the openvpn pid file
# Returns:
#   None
#####################################
function elastic::vpn-run-macos(){
  local openvpn_config_file=${1:-"${HOME}/.oblt-cli/vpn/production.ovpn"}
  local openvpn_auth_file=${2:-"${HOME}/.oblt-cli/vpn/auth.txt"}
  local openvpn_log_file=${3:-"${HOME}/.oblt-cli/vpn/openvpn.log"}
  local openvpn_pid_file=${4:-"${HOME}/.oblt-cli/vpn/openvpn.pid"}
  SCRIPT="#!/usr/bin/env bash
  function cleanup()
  {
      rm -f ${openvpn_auth_file} || true
  }
  function showlogs()
  {
      tail -n 20 ${openvpn_log_file} || true
  }
  trap cleanup EXIT
  trap showlogs ERR
  sudo openvpn --daemon --config \"${openvpn_config_file}\" --auth-user-pass \"${openvpn_auth_file}\" --log \"${openvpn_log_file}\" --writepid \"${openvpn_pid_file}\"
  sudo openvpn --daemon --config \"${openvpn_config_file}\" --auth-user-pass \"${openvpn_auth_file}\" --log \"${openvpn_log_file}\" --writepid \"${openvpn_pid_file}\"
  "

  eval "${SCRIPT}"
  sleep 5
  log::info "VPN client started in the background"
  log::info "To close the VPN client run: sudo kill \$(cat ${openvpn_pid_file})"
}

#####################################
# Run the VPN client in Linux
# Globals:
#   None
# Arguments:
#   openvpn_config_file: path to the openvpn config file
#   openvpn_auth_file: path to the openvpn auth file
#   openvpn_log_file: path to the openvpn log file
#   openvpn_pid_file: path to the openvpn pid file
# Returns:
#   None
#####################################
function elastic::vpn-run-linux(){
  local openvpn_config_file=${1:-"${HOME}/.oblt-cli/vpn/production.ovpn"}
  local openvpn_auth_file=${2:-"${HOME}/.oblt-cli/vpn/auth.txt"}
  local openvpn_log_file=${3:-"${HOME}/.oblt-cli/vpn/openvpn.log"}
  local openvpn_pid_file=${4:-"${HOME}/.oblt-cli/vpn/openvpn.pid"}
  SCRIPT="#!/usr/bin/env bash
  function cleanup()
  {
      rm -f ${openvpn_auth_file} || true
  }
  function showlogs()
  {
      tail -n 20 ${openvpn_log_file} || true
  }
  trap cleanup EXIT
  trap showlogs ERR
  sudo openvpn --daemon --config \"${openvpn_config_file}\" --auth-user-pass \"${openvpn_auth_file}\" --log \"${openvpn_log_file}\" --writepid \"${openvpn_pid_file}\"
  "

  eval "${SCRIPT}"
  sleep 5
  log::info "VPN client started in the background"
  log::info "To close the VPN client run: sudo kill \$(cat ${openvpn_pid_file})"
}

#####################################
# Run a script in a new macOS terminal
# Globals:
#   None
# Arguments:
#   script: script to run
# Returns:
#   Script output
#####################################
function elastic::run-on-macOS-terminal(){
  local script=$1
  local script="tell application \"Terminal\"
    activate
    do script \"${script}\"
  end tell"
  osascript -e "${script}"
}

#####################################
# Run a script in a terminal
# Globals:
#   None
# Arguments:
#   script: script to run
# Returns:
#   Script output
#####################################
function elastic::run-on-terminal(){
  local SCRIPT=$1
  if [ "$(install::isMacOS)" == "true" ]; then
    elastic::run-on-macOS-terminal "${SCRIPT}"
    return
  fi
  if [ "$(install::isLinux)" == "true" ]; then
    eval "${SCRIPT}"
    return
  fi
}

#####################################
# Stop the VPN client
# Globals:
#   None
# Arguments:
#   openvpn_pid_file: path to the openvpn pid file
# Returns:
#   None
#####################################
function elastic::vpn-stop(){
  local openvpn_pid_file=${1:-"${HOME}/.oblt-cli/vpn/openvpn.pid"}
  log::warn "Killing previous VPN client"
  local openvpn_pid
  if [ -f "${openvpn_pid_file}" ]; then
    openvpn_pid=$(head -1 "${openvpn_pid_file}")
    if [ -n "${openvpn_pid}" ]; then
      # shellcheck disable=SC2086
      sudo kill "${openvpn_pid}" || true
    fi
  fi
}
