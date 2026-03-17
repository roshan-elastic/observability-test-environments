# Observability Test Environments

Repo: **github.com/roshan-elastic/observability-test-environments**

Configurable demo environment for testing **Elastic Observability** with a relatable enterprise use case: **MegaMart** (retail / e-commerce). Supports local Kubernetes (Docker Desktop or Kind), scaling to 100+ services, and optional GCP/AWS deployment.

## Quick start (Phase 1 – local)

1. Ensure you have a Kubernetes cluster (Docker Desktop with Kubernetes enabled, or [Kind](https://kind.sigs.k8s.io/)).
2. From the repo root:
   ```bash
   ./demo/scripts/run-local.sh
   ```
3. Follow the script output for frontend URL and Kibana (if using Elastic Cloud).

Full details: [demo/README.md](demo/README.md).

## Repo layout

| Path | Purpose |
|------|--------|
| `demo/` | MegaMart demo: config, Helm chart, chaos scenarios, scripts |
| `demo/config/` | YAML config (default, local, gcp, aws) |
| `demo/chart/` | Helm chart: app services + OpenTelemetry Collector |
| `demo/scripts/` | `run-local.sh`, `run-gcp.sh`, etc. |
| `demo/chaos/` | Chaos Mesh scenario manifests (Phase 2) |

## Requirements

- `kubectl`, Helm 3.x
- Kubernetes cluster (Docker Desktop K8s or Kind)
- (Optional) Elastic Cloud deployment or ECK for Elasticsearch/Kibana

## License

Apache 2.0
