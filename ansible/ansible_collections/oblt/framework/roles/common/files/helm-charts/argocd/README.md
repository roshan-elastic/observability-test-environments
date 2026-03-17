# Integrate ArgoCD in the oblt clusters

## Description

The current implementation of the bolt clusters has some limitations when deploying k8s resources. The Ansible implementation always deploys the k8s and does not allow rollbacks. We rely on Helm packages that would not update if there are no changes, but still, we face failures related to the deployment, and we cannot perform a rollback. Also, when we delete a k8s deployment from the cluster configuration, the deployments must be manually deleted after modifying the configuration.
ArgoCD is a stable solution that has well-defined k8s workflows. It allows rollback and manages the deletion of resources based on a Git repository, along with the rollbacks in case of any failure.

## Goals

* Deploy ArgoCD in an oblt-cluster
* Manage k8s deployments from a Git repo inside the k8s cluster
* Evaluate the cost of migrating k8s deployments to ArgoCD
* Evaluate manage other resources as k8s resources (Stack,VMs,...)

## ArgoCD

ArgoCD is a GitOps tool that allows us to manage k8s resources from a Git repository. It has a UI and a CLI to manage. The first step is to deploy ArgoCD in the k8s cluster. We can use the Helm chart or the ArgoCD CLI. We will use the Helm chart to deploy ArgoCD in the k8s cluster.

We created the k8s cluster using the [oblt clusters framework](https://github.com/elastic/observability-test-environments/pull/16703). This saved time configuring and deploying ArgoCD in this PoC. For the rest, we follow the instructions at [Getting started with ArgoCD](https://argo-cd.readthedocs.io/en/stable/getting_started).
Then we install the [Helm chart in the k8s cluster]( https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd) using the oblt clusters framework too.
Locally we installed th [ArgoCD CLI](https://argo-cd.readthedocs.io/en/stable/cli_installation/), but we did not use it. We use regular kubectl commands.

```shell
# 2.9.1 and 2.9.2 are broken
VERSION=v2.8.7
curl -sSL -o argocd-darwin-amd64 https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-darwin-amd64
```

To use the AdgoCD UI we can expose the port in an ingress or use port-forward. We use port-forward to access the UI.

```shell
kubectl config set-context --current --namespace argocd
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

or

```shell
argocd admin dashboard --namespace argocd
```

The ArgoCD UI requires authentication. We can use the CLI to get the admin password.

```shell
argocd admin initial-password -n argocd
```

In the oblt clusters, we usually deploy nginx ingress controller and cert-manager. So to test if it is possible to do the same in ArgoCD, we create a Helm Chart to deploy the nginx ingress controller and cert-manager.
We added and improvements over the deployments we usually make, we added an external secrets operator to manage secrets in GCP Secret Manager.
The external secrets operator requires a service account with permission to access the secrets in [GCP Secret Manager](https://external-secrets.io/latest/provider/google-secrets-manager/). We created a service account and added the permissions to access the secrets in GCP Secret Manager.
We configured it to use [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) to access the secrets in GCP Secret Manager.

```shell
PROJECT_ID=elastic-observability
CLUSTER_NAME=argocd
COMPUTE_REGION=us-central1-c
GSA_NAME=k8s-esternal-secrets
GSA_PROJECT=${PROJECT_ID}
ROLES_NAME="roles/secretmanager.secretAccessor roles/iam.serviceAccountTokenCreator roles/iam.workloadIdentityUser"
# service account name allow to impersonate
KSA_NAME=external-secrets-default
NAMESPACE=external-secrets

gcloud container clusters update "${CLUSTER_NAME}" \
  --region="${COMPUTE_REGION}" \
  --workload-pool="${PROJECT_ID}.svc.id.goog"

gcloud container node-pools update" ${CLUSTER_NAME}-pool" \
  --cluster="${CLUSTER_NAME}" \
  --region="${COMPUTE_REGION}" \
  --workload-metadata=GKE_METADATA

gcloud iam service-accounts create "${GSA_NAME}" \
  --project "${GSA_PROJECT}" \
  --display-name "${GSA_NAME}" \
  --description="Service account to access to Google Secrets Manager from k8s clusters"

for ROLE_NAME in ${ROLES_NAME}; do
  gcloud projects add-iam-policy-binding "${GSA_PROJECT}" \
    --member "serviceAccount:${GSA_NAME}@${GSA_PROJECT}.iam.gserviceaccount.com" \
    --role "${ROLE_NAME}" \
    --condition=None
done

# Allow the k8s service account to impersonate the GCP service account
gcloud iam service-accounts add-iam-policy-binding "${GSA_NAME}@${GSA_PROJECT}.iam.gserviceaccount.com" \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[${NAMESPACE}/${KSA_NAME}]"
```

In the same PR, we added the bootstrap-k8s Helm Chart and the deploy using helm into the ArgoCD namespace.
When the helm chart is deployed in the argo namespace, we see the deployments and their status in the ArgoCD UI.

It is important to use the latest versions of the CDRs to avoid issues like the one we faced with [Helm Charts deployments](https://github.com/argoproj/argo-cd/issues/16466).

In the test, external secrets were really easy to use and really useful to manage secrets in GCP Secret Manager. We can use the same approach to manage secrets in Hashicorp Vault, AWS Secrets Manager, or Azure Key Vault.
The secrets are stored in an external Vault, and the k8s secrets are refreshed every hour(configurable).
GitHub actions have support for [GCP Secret Manager](https://github.com/google-github-actions/get-secretmanager-secrets), so we can use it to manage secrets in GCP Secret Manager as we want to make in this [issue](https://github.com/elastic/observability-test-environments/issues/11075).
The secrests are stored as plain text or as JSON sees [All keys, One Secret](https://external-secrets.io/latest/guides/all-keys-one-secret/). We used JSON secrets to store secrets.

```shell
SECRET_NAME=oblt-clusters-github-token-ro
SECRET_VALUE='
{
  "name": "oblt-clusters-github-token-ro",
  "password": "github_pat_REDACTED",
  "username": "kuisathaverat",
  "url": "https://github.com/elastic/observability-test-environments.git",
  "type": "git"
}
'
echo "${SECRET_VALUE}"|gcloud secrets create "${SECRET_NAME}" --data-file=-
gcloud secrets versions access latest --secret="${SECRET_NAME}" 2>/dev/null
#gcloud secrets delete "${SECRET_NAME}"
```

Some of the features we used in ArgoCD are:

* GitOps: it gives us the possibility to synchronize the k8s resources with a Git repository.
* Multi-cluster: we can manage multiple clusters from the same ArgoCD instance.
* Control plane: we can manage the k8s resources from the ArgoCD UI or the CLI.
* Templating JSONNET: it is possible to use JSONNET to template the k8s resources.
* Orchestration: Argo CD support to manage deployments with [hooks](https://argoproj.github.io/argo-cd/user-guide/resource_hooks/) and [waves](https://argoproj.github.io/argo-cd/user-guide/sync-waves/). Hooks can be implemented with scripts, containers, argo workflows, or k8s jobs. Waves allow us to deploy resources in a specific order.

There is a repository with [some examples](https://github.com/argoproj/argocd-example-apps/tree/master) of the features of ArgoCD.

## Crossplane

Crossplane creates Kubernetes Custom Resource Definitions (CRDs) to represent the external resources as native Kubernetes objects. As native Kubernetes objects, you can use standard commands like kubectl create and kubectl describe. The full Kubernetes API is available for every Crossplane resource.

We can use Crossplane to manage k8s resources and other resources as k8s resources. We can use it to manage VMs, databases, storage, etc. We can use it to manage the resources in the cloud or on-premises.
Crossplane has a [marketplace](https://marketplace.upbound.io/) with providers to manage resources in the cloud and on-premises.

Crossplane has a [CLI tool](https://docs.crossplane.io/v1.14/cli/) we can use to manage the k8s resources associated with the providers.

```shell
curl -sL "https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh" | sh
```

To install Crossplane, we can use the [Helm chart](https://github.com/crossplane/crossplane/tree/master/cluster/charts/crossplane), and we follow the instructions at the [install Crossplane](https://docs.crossplane.io/v1.14/software/install/) using the oblt clusters framework, and using the ArgoCD namespace to manage the resources with ArgoCD. You can find important information about how Crossplain works at [Crossplane Introduction](https://docs.crossplane.io/v1.14/getting-started/introduction/)

References:

* Group resources https://docs.crossplane.io/v1.14/concepts/compositions/
* Configuration Packages https://docs.crossplane.io/v1.14/concepts/packages/
* ArgoCD+Crossplane https://docs.crossplane.io/knowledge-base/integrations/argo-cd-crossplane/
* Any terraform provider can be converted into a Crossplane provider https://github.com/crossplane/upjet

### GCP provider

Crossplane can manage GCP resources as k8s resources by using the [GCP provider](https://docs.upbound.io/providers/provider-gcp). Clossplane providers are deployed as a k8s resource. In the case of GCP provider,
we have to install a provider for each type of resource we want to manage (Buckets, VMs, GKE,...).
For the authentication, we used the same approach we used in ArgoCD, we created a service account and we added the permissions to access the resources in GCP.
Then we configured it to use [Workload Identity](https://docs.upbound.io/providers/provider-gcp/authentication/).

```shell
PROJECT_ID=elastic-observability
CLUSTER_NAME=argocd
COMPUTE_REGION=us-central1-c
GSA_NAME=k8s-crossplane
GSA_PROJECT=${PROJECT_ID}
# in the PoC we use one service account with all the permissions, it is better to use one service account per provider with the required permissions only.
ROLES_NAME="roles/container.clusterAdmin roles/compute.admin roles/storage.admin roles/secretmanager.admin roles/iam.serviceAccountTokenCreator roles/iam.workloadIdentityUser roles/iam.serviceAccountUser roles/iam.serviceAccountCreator roles/iam.securityAdmin"
# service account name allow to impersonate
KSA_NAME=crossplane-default
NAMESPACE=crossplane-system

gcloud iam service-accounts create "${GSA_NAME}" \
  --project "${GSA_PROJECT}" \
  --display-name "${GSA_NAME}" \
  --description="Service account to access to manage k8s clusters"

for ROLE_NAME in ${ROLES_NAME}; do
  gcloud projects add-iam-policy-binding "${GSA_PROJECT}" \
    --member "serviceAccount:${GSA_NAME}@${GSA_PROJECT}.iam.gserviceaccount.com" \
    --role "${ROLE_NAME}" \
    --condition=None
done

# Allow the k8s service account to impersonate the GCP service account
gcloud iam service-accounts add-iam-policy-binding "${GSA_NAME}@${GSA_PROJECT}.iam.gserviceaccount.com" \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[${NAMESPACE}/${KSA_NAME}]"

# list roles
gcloud projects get-iam-policy "${PROJECT_ID}"  \
--flatten="bindings[].members" \
--format='table(bindings.role)' \
--filter="bindings.members:${GSA_NAME}@${GSA_PROJECT}.iam.gserviceaccount.com"
```

The provider must be installed before the resources. This means that they can not be in the same Helm Chart, or they should be installed as CDRs.

We put all the providers in the Crossplane Helm Chart and the resources in different Helm Charts.

#### Buckets

We can create a bucket using the [Bucket](https://marketplace.upbound.io/providers/upbound/provider-gcp-storage/) resource.
It is easy to manage. The bucket is created when the resource is created. The status of the Bucket is stored in the k8s resource in the status section. When the k8s resource is deleted, the Bucket is deleted too.

#### Secrets

We can create a secret using the [Secret](https://marketplace.upbound.io/providers/upbound/provider-gcp-secretmanager/) resource. Like the bucket, the secret is created when the resource is created, and the status of the secret is stored in the k8s resource in the status section. When the k8s resource is deleted, the secret is deleted too.

#### VMs

We can create a VM using the [Instance](https://marketplace.upbound.io/providers/upbound/provider-gcp-compute/) resource. The VM is created when the resource is created, and the status of the VM is stored in the k8s resource in the status section. When the k8s resource is deleted, the VM is deleted too.

#### GKE

We can create a GKE cluster using the [Cluster](https://marketplace.upbound.io/providers/upbound/provider-gcp-container) resource. The cluster is created when the resource is created, and the status of the cluster is stored in the k8s resource in the status section. When the k8s resource is deleted, the cluster is deleted too.
We tested to create a regular GKE cluster with a machine pool and an Autopilot cluster.

#### IAM service account

We can create a service account using the [ServiceAccount](https://marketplace.upbound.io/providers/upbound/provider-gcp-cloudplatform) resource. When the resource is created, the service account is created, and the status of the service account is stored in the k8s resource in the status section. When the k8s resource is deleted, the service account is deleted too. After creating the user account, you need to add permissions to the user account.
We can get the member from the user account to use to assign the permissions to the service account.

```shell
kubectl get ServiceAccount.cloudplatform.gcp.upbound.io --output=jsonpath='{.items[0].status.atProvider.member}'
```

### Terraform provider

Crossplane can manage Terraform plans as k8s resources by using the [Terraform provider](https://marketplace.upbound.io/providers/upbound/provider-terraform). The Terraform provider can use an inline Terraform plan or a Git repository. The provider documentation is not well structured, and we found [issues](https://github.com/upbound/provider-terraform/issues/218) when we tried to use it.
It is important to set a backend other than the default. Terraform plans without state don't work.
Some parameters are a mess because the terraform provider is imported and is not well documented. Worth reviewing the [issues](https://github.com/upbound/provider-gcp/issues/62) to find help.
The variables passed to the Terraform plan must be strings. No other types are supported.
Once we solve the issues, it is easy to run Terraform plans as k8s resources. It is possible to inject secrets as variables in the Terraform plan. This allows not to have the secrets in Git because we use the external secrets operator for that.

References:

* [Quickstart](https://marketplace.upbound.io/providers/upbound/provider-terraform/v0.12.0/docs/quickstart)
* [Configuration](https://marketplace.upbound.io/providers/upbound/provider-terraform/v0.12.0/docs/configuration)

### Ansible provider

Crossplane can manage Ansible roles as k8s resources by using the [Ansible provider](https://github.com/crossplane-contrib/provider-ansible). The Ansible provider can use an inline Ansible role or a Git repository that contains a role or an Ansible collection. It is not well documented, and the limitation of only importing collections means extra work is needed to build a collection with the roles we want to use. Overall it is easy to use, and it is a good option to manage k8s resources.

References:

* [Building an Ansible Collection](https://docs.ansible.com/ansible/latest/dev_guide/developing_collections_distributing.html#building-your-collection-tarball)
