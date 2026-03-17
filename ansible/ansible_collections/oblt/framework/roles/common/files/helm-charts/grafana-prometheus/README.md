# Grafana+Prometheus Helm chart

This Role deploy [Grafana][] and [Prometheus][] into the k8s cluster. [Grafana][] is the open source analytics & monitoring solution for every database [Prometheus][] is a metrics and alerting open-source monitoring solution. It has also tasks to uninstall the services installed on a Kubernetes cluster.

[Grafana]: https://grafana.com/
[Prometheus]: https://prometheus.io/

## Parameters

### prometheus Prometheus configuration

see [Prometheus Helm Chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus)

| Name                 | Description       | Value |
| -------------------- | ----------------- | ----- |
| `prometheus.enabled` | Enable Prometheus |       |


### grafana Grafana configuration

see [Grafana Helm Chart](https://github.com/grafana/helm-charts/blob/main/charts/grafana/README.md)

| Name              | Description    | Value |
| ----------------- | -------------- | ----- |
| `grafana.enabled` | Enable Grafana |       |
