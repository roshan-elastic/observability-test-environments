# Cert Manager Helm chart

This Helm Chart deploy a Cert manager service.

## Parameters

### cert-manager Configuration parameters for cert-manager




### cert-manager.global global configuration for cert-manager

| Name                                           | Description                              | Value                 |
| ---------------------------------------------- | ---------------------------------------- | --------------------- |
| `cert-manager.global.leaderElection.namespace` | namespace to use for leader election     | `default`             |
| `cert-manager.startupapicheck.startupapicheck` | time to wait for the webhook to be ready | `5m`                  |
| `cert-manager.installCRDs`                     | install CRDs                             | `true`                |
| `cert-manager.ingressShim.defaultIssuerName`   | name of the default issuer               | `letsencrypt-staging` |
| `cert-manager.ingressShim.defaultIssuerKind`   | kind of the default issuer               | `ClusterIssuer`       |


### cert-manager.resources requests and limits

| Name                                     | Description                         | Value   |
| ---------------------------------------- | ----------------------------------- | ------- |
| `cert-manager.resources.limits.memory`   | The memory limit for the Opbean     | `128Mi` |
| `cert-manager.resources.requests.memory` | The requested memory for the Opbean | `64Mi`  |
| `cert-manager.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`  |


### cert-manager.webhook Configuration parameters for cert-manager webhook




### cert-manager.webhook.resources requests and limits

| Name                                             | Description                         | Value   |
| ------------------------------------------------ | ----------------------------------- | ------- |
| `cert-manager.webhook.resources.limits.memory`   | The memory limit for the Opbean     | `128Mi` |
| `cert-manager.webhook.resources.requests.memory` | The requested memory for the Opbean | `64Mi`  |
| `cert-manager.webhook.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`  |


### cert-manager.cainjector Configuration parameters for cert-manager cainjector




### cert-manager.cainjector.resources requests and limits

| Name                                                | Description                         | Value   |
| --------------------------------------------------- | ----------------------------------- | ------- |
| `cert-manager.cainjector.resources.limits.memory`   | The memory limit for the Opbean     | `256Mi` |
| `cert-manager.cainjector.resources.requests.memory` | The requested memory for the Opbean | `64Mi`  |
| `cert-manager.cainjector.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`  |
