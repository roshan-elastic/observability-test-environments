# MegaMart Demo

Observability demo environment using the **MegaMart** use case (retail / e-commerce): storefront, cart, checkout, payments, inventory, fulfillment, and platform services. Config-driven for scale, instrumentation, and backend (Elasticsearch via OpenTelemetry).

## Phase 1: Run locally

### Prerequisites

- **Kubernetes**: Docker Desktop with Kubernetes enabled, or [Kind](https://kind.sigs.k8s.io/) (`kind create cluster`).
- **kubectl** and **Helm 3.x**.
- (Optional) **Elastic Cloud** or ECK for Elasticsearch/Kibana; otherwise the collector will be ready and you can add the backend later.

### Deploy

From the repo root:

```bash
./demo/scripts/run-local.sh
```

Or with explicit config:

```bash
helm install megamart demo/chart -f demo/config/local.yaml -n megamart --create-namespace
```

### Config

- **Default**: [demo/config/default.yaml](config/default.yaml) – service count, instrumentation, Elasticsearch placeholders.
- **Local**: [demo/config/local.yaml](config/local.yaml) – overrides for local K8s (Docker Desktop / Kind).

Override Elasticsearch (e.g. Elastic Cloud) via env or values:

```bash
export ELASTICSEARCH_URL="https://..."
export ELASTICSEARCH_API_KEY="..."
./demo/scripts/run-local.sh
```

### Access

After deploy, the script prints URLs. Typical:

- **Storefront**: `http://localhost:8080` (after `kubectl port-forward svc/storefront 8080:8080 -n megamart`)
- **Kibana**: your Elastic Cloud or ECK URL (data appears once the collector is configured with valid credentials).

## Architecture (Phase 1)

- **Namespaces**: `megamart` (app + collector).
- **Services**: Storefront, Cart, Checkout, Payment, Inventory, Order-Orchestrator, Fulfillment (MegaMart naming; some backed by OpenTelemetry Demo images for quick traces).
- **Telemetry**: In-cluster OpenTelemetry Collector; OTLP receiver → Elasticsearch exporter (or OTLP to Elastic). Apps set `OTEL_EXPORTER_OTLP_ENDPOINT` to the collector service.

## Later phases

- **Phase 2**: Chaos Mesh + failure scenarios (see `demo/chaos/`).
- **Phase 3**: Scale to 100+ services via config.
- **Phase 4+**: Two local clusters, GCP, AWS.
