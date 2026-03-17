# Autopilot integration

Autopilot is a new mode of operation for GKE that removes the complexity of managing and operating a cluster. It is designed to provide a simple, secure, and cost-effective way to run containers on GKE.

The Elastic Agent can work on Autopilot clusters to collect logs and metrics from the cluster and the workloads running on it.
We provide a cluster template to deploy an Elastic Agent on Autopilot clusters and test the data on the Elastic Stack deployed on ESS.

The Elastic agent has two ways of working on Autopilot clusters:

* Managed mode: The Elastic Agent is deployed as a DaemonSet on the cluster and it is managed by Fleet in the Kibana UI.
* Standalone mode: The Elastic Agent is deployed as a Deployment on the cluster and it is managed by the user on k8s manifest resources.

The logs and metrics from the Elastic Agent are reported to the Elastic Cloud and can be visualized in the Kibana UI.

## Deploying an Elastic Agent on managed mode

It is the default mode of operation for the Elastic Agent on Autopilot clusters. The Elastic Agent is deployed as a DaemonSet on the cluster and it is managed by Fleet in the Kibana UI.

```bash
oblt-cli cluster create custom \
    --template-file autopilot \
    --parameters '{"StackVersion":"8.7.0", "AutopilotVersion":"1.24.9-gke.3200"}'
```

```bash
oblt-cli cluster create custom \
    --template-file autopilot \
    --parameters '{"StackVersion":"8.7.0", "AutopilotVersion":"1.24.9-gke.3200", "DeployMode", "managed"}'
```

## Deploying an Elastic Agent on standalone mode

The Elastic Agent is deployed as a Deployment on the cluster and it is managed by the user on k8s manifest resources.
This means the configuration of the inputs and outputs is done by the user on a configmap that contains the Elastic Agent configuration file.
In this mode the Elastic Agent does not appear in the Fleet UI.

```bash
oblt-cli cluster create custom \
    --template-file "${PWD}/environments/users/autopilot-test.yml.tmpl" \
    --parameters '{"StackVersion":"8.7.0", "AutopilotVersion":"1.24.9-gke.3200", "DeployMode", "standalone"}'
```
