#!/bin/env bash

set -eo pipefail

SCRIPT_LIB_PATH=$(dirname -- "${BASH_SOURCE[0]}")
# shellcheck source=/dev/null
. "${SCRIPT_LIB_PATH}/log.sh"

#####################################
# Ask a question to the user, if the user answer is not y, it exits
# Globals:
#   None
# Arguments:
#   question
# Returns:
#   None
#####################################
function prompt::askYN(){
  local question=${1:?"Missing question argument"}
  log::warn "${question}"
  read -r -p "[y/N]: " response
  if [ "${response}" != "y" ]; then
    exit 1
  fi
}

#####################################
# Warn the user, and wait for enter to configure
# Globals:
#   None
# Arguments:
#   question
# Returns:
#   None
#####################################
function prompt::waitEnter(){
  local question=${1:?"Missing question argument"}
  log::warn "${question}"
  read -r -p "[Press enter to continue]" response
}

#####################################
# Inform the user, with a dialog
# Globals:
#   None
# Arguments:
#   message
# Returns:
#   None
#####################################
function prompt::dialog(){
  local message=${1:?"Missing message argument"}
  osascript -e 'display dialog "'"${message}"'" buttons {"OK"} default button 1'
}

#####################################
# Ask the user a question, with a dialog, if the answer is not yes, it exits
# Globals:
#   None
# Arguments:
#   message
#   default_button
# Returns:
#   None
#####################################
function prompt::dialogYN(){
  local message=${1:?"Missing message argument"}
  local default_button=${2:-2}
  response=$(osascript -e 'display dialog "'"${message}"'" buttons {"Yes", "No"} default button '"${default_button}"'')
  if [ "${response}" != "button returned:Yes" ]; then
    exit 1
  fi
}

#####################################
# Configure sudo to use touch id
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   None
#####################################
function prompt::configureSudo(){
  log::info "Configuring sudo to use touch id"
  prompt::dialogYN "Do you want to configure sudo to use touch id?"
  echo "auth       sufficient     pam_tid.so" | sudo tee -a /etc/pam.d/sudo
}
