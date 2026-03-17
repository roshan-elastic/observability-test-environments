#!/usr/bin/env bash
set -ueio pipefail
USAGE="usage : $0 <pvc_name> <new_size>"
PVC_NAME=${1:?$USAGE}
NEW_SIZE=${2:?$USAGE}

kubectl patch pvc "${PVC_NAME}" -p '{"spec":{"resources":{"requests":{"storage":"'"${NEW_SIZE}"'"}}}}'
kubectl get pvc "${PVC_NAME}" -o yaml
