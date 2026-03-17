#!/usr/bin/env bash
# MegaMart demo – run locally (Docker Desktop K8s or Kind)
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEMO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG="${DEMO_ROOT}/config/local.yaml"
CHART="${DEMO_ROOT}/chart"
NAMESPACE="${NAMESPACE:-megamart}"

echo "=== MegaMart demo (local) ==="
echo "Config: $CONFIG"
echo "Chart: $CHART"
echo "Namespace: $NAMESPACE"
echo ""

if ! kubectl cluster-info &>/dev/null; then
  echo "No Kubernetes cluster reachable. Start Docker Desktop K8s or run: kind create cluster"
  exit 1
fi

# Helm dependency update (opentelemetry-demo)
if ! command -v helm &>/dev/null; then
  echo "Helm is required. Install from https://helm.sh"
  exit 1
fi
if [[ ! -d "$CHART/charts" ]] || [[ -z "$(ls -A "$CHART/charts" 2>/dev/null)" ]]; then
  echo "Updating Helm dependencies..."
  helm dependency update "$CHART"
fi

echo "Installing MegaMart demo..."
helm upgrade --install megamart "$CHART" \
  -f "$CONFIG" \
  -n "$NAMESPACE" \
  --create-namespace \
  --wait --timeout 5m

echo ""
echo "Done. To access the storefront:"
echo "  kubectl port-forward svc/frontend-proxy 8080:8080 -n $NAMESPACE"
echo "  Then open http://localhost:8080"
echo ""
echo "If you set ELASTICSEARCH_URL and ELASTICSEARCH_API_KEY, configure them in $CONFIG under otelCollector.elasticsearch and re-run this script, or upgrade:"
echo "  helm upgrade megamart $CHART -f $CONFIG -n $NAMESPACE --set otelCollector.elasticsearch.endpoint=<URL> --set otelCollector.elasticsearch.apiKey=<KEY>"
echo ""
