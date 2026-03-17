# Accessing to Kubernetes clusters

## Using oblt-cli

The `oblt-cli cluster k8s` command will give you access to the Kubernetes cluster.
This command will open a shell authenticated and in the context of the k8s cluster configured.

```bash
CLUSTER_NAME=edge-oblt
oblt-cli cluster k8s --cluster-name ${CLUSTER_NAME}
```

```bash
2023/09/14 19:12:02 SlackChannel: '@MYSLACKID'
2023/09/14 19:12:02 User: 'myUser'
2023/09/14 19:12:02 Git mode: 'ssh'
2023/09/14 19:12:15 Writing output file /Users/myUser/.oblt-cli/edge-lite-oblt-activate.sh
Hermit environment /Users/myUser/.oblt-cli/observability-test-environments/.ci activated
/Users/myUser/bin/google-cloud-sdk/bin/gke-gcloud-auth-plugin
Fetching cluster endpoint and auth data.
kubeconfig entry generated for edge-lite-oblt.
Context "gke_elastic-observability_us-central1-c_edge-lite-oblt" modified.
Switched to context "gke_elastic-observability_us-central1-c_edge-lite-oblt".
Welcome to the k8s cluster shell
```

Then we can use `kubectl` and `helm` to access the Kubernetes cluster.

## Using gcloud to access

The oblt clusters use the `elastic-observability` GCP project to store the GKE clusters,
any user in the observability group can access the clusters.
You nee [`gcloud`](https://cloud.google.com/sdk/docs/install) and [`kubectl`](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl) installed to access the cluster and work with k8s resources.
If you plant to manage Helm deployments you need also [`helm`](https://helm.sh/docs/intro/install/).
These are the commands to authenticate on GCP, get the GKE cluster credentials and check for some k8s resources.

```bash
CLUSTER_NAME=edge-oblt
gcloud auth login
gcloud container clusters get-credentials --project elastic-observability --zone us-central1-c CLUSTER_NAME
kubectl get pods -A
helm list
```

## Access Kubernetes resources

At this point you can use `kubectl` and `helm` to access the Kubernetes cluster.
The terminal where you executed `source activate` has all the tools and secrets needed
to access the Kubernetes cluster.

```bash
kubeclt get pods -A
kubectl get services -A
kubectl get deployment -A
kubectl get statefulset -A
kubectl get daemonset -A
kubectl get cronjobs -A
kubectl get ingress -A

helm list
```

All the kubernetes nodes, and pods logs are ingested in Elasticsearch,
so it is also possible to check the logs in Kibana, see [How to check the logs](./use-case-check-logs.md).

## Diagnose why a pod does not start

There are many reasons why a pod does not start,
it is possible to diagnose the cause by using a couple of `kubeclt` commands.

### Describe your k8s resource

The `kubectl describe` command will give you some info a bout the sequence of events that happen when the pod starts.
Checking the events output of the `pod` and the `deployment` or `statefulset` you will have the root cause why the pod does not start.

```bash
kubectl describe pod opbeans-go-86b56b7c9f-lrvxp
kubectl describe deployment opbeans-go
```

### Access Kubernetes logs

If the pods restarts or the services are failing you should take a look to the logs.
You can use the kubeclt command to grab the logs from your pods.
First list the pod you want to check the logs,
in the example we use the label `app` to filter.

```bash
❯ kubectl get pod --selector app=opbeans-go
NAME                          READY   STATUS    RESTARTS   AGE
opbeans-go-86b56b7c9f-lrvxp   1/1     Running   0          13d
```

Then you can request the logs

```bash
❯ kubectl logs -f opbeans-go-86b56b7c9f-lrvxp
```

## Forward ports

once you are authenticated, you can forward local ports to the k8s cluster and
after that you can connect to localhost ports to connect to the k8s cluster services.
Kibana and Elasticsearch services will be HTTP and unauthenticated on that port.

`kubectl port-forward service/opbeans 3000`
