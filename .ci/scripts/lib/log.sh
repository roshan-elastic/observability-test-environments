#!/bin/env bash

set -eo pipefail

SHELL_ESCAPE="\033[";
SHELL_RESET="${SHELL_ESCAPE}0m";
SHELL_COLOR_RED="${SHELL_ESCAPE}31m";
SHELL_COLOR_YELLOW="${SHELL_ESCAPE}33m";
SHELL_COLOR_GREEN="${SHELL_ESCAPE}32m";
DEBUG="${DEBUG:-}"

#####################################
# Print an error message in the stderr
# Globals:
#   None
# Arguments:
#   message
# Returns:
#   None
#####################################
function log::error() {
  echo -e "${SHELL_COLOR_RED}[$(date +'%Y-%m-%dT%H:%M:%S%z')]:${SHELL_RESET} ${*}" >&2
}

#####################################
# Print a warning message in the stderr
# Globals:
#   None
# Arguments:
#   message
# Returns:
#   None
#####################################
function log::warn() {
  echo -e "${SHELL_COLOR_YELLOW}[$(date +'%Y-%m-%dT%H:%M:%S%z')]:${SHELL_RESET} ${*}" >&2
}

#####################################
# Print an info message in the stderr
# Globals:
#   None
# Arguments:
#   message
# Returns:
#   None
#####################################
function log::info() {
  echo -e "${SHELL_COLOR_GREEN}[$(date +'%Y-%m-%dT%H:%M:%S%z')]:${SHELL_RESET} ${*}" >&2
}

#####################################
# Print a debug message in the stderr
# Globals:
#   DEBUG
# Arguments:
#   message
# Returns:
#   None
#####################################
function log::debug() {
  if [ -n "${DEBUG}" ]; then
    echo -e "${SHELL_COLOR_GREEN}[$(date +'%Y-%m-%dT%H:%M:%S%z')]:${SHELL_RESET} ${*}" >&2
  fi
}
