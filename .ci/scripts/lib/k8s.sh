#!/bin/env bash

set -eo pipefail

#####################################
# Get the current k8s context
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   The current k8s context
function k8s::current-context(){
  kubectl config current-context
}

#####################################
# Set the namespace of the current k8s context
# Globals:
#   None
# Arguments:
#   namespace Name of the k8s project namespace
# Returns:
#   None
#####################################
function k8s::set-current-context-namespace(){
  local namespace=${1:?"Missing namespace argument"}
  kubectl config set-context --current --namespace="${namespace}"
}

#####################################
# List all pods in a namespace, if a namespace is not pass, it will use the default namespace
# Globals:
#   None
# Arguments:
#   namespace Name of a namespace
# Returns:
#   None
#####################################
function k8s::get-pods(){
  local namespace=${1:-"default"}
  kubectl get po -n "${namespace}"
}

#####################################
# Delete all pods that are not running in a namespace, if a namespace is not pass, it will use the default namespace.
# Globals:
#   None
# Arguments:
#   namespace Name of a namespace
# Returns:
#   None
#####################################
function k8s::delete-pods-no-running(){
  local namespace=${1:-"default"}
  kubectl delete po -n "${namespace}" --field-selector=status.phase!=Running
}

#####################################
# Show the go-template to get the name of a k8s object
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   The go-template to get the name of a k8s object
#####################################
function k8s::get-name-go-template(){
  # shellcheck disable=SC2028
  echo '{{ range .items }}{{ .metadata.name }}{{ "\n" }}{{ end }}'
}

#####################################
# Get the names of all objects in a namespace, if a namespace is not pass, it will use the default namespace.
# Globals:
#   None
# Arguments:
#   namespace Name of a namespace
#   object_type Type of k8s object
# Returns:
#   The names of all objects in a namespace
#####################################
function k8s::get-objects-names(){
  local namespace=${1:-"default"}
  local object_type=${2:-"po"}
  kubectl get "${object_type}" -n "${namespace}" -o go-template="$(k8s::get-name-go-template)"
}
