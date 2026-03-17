# Serverless Kubernetes Resources

Additionally to the [Managed Kubernetes Infrastructure (MKI)][mki] cluster,
the oblt clusters deploy some application that run on a different Kubernetes cluster.
This Kuberenes cluster is a GKE cluster managed by the [oblt framework][oblt-framework]
The resources on those Kubernetes clusters are available to the developers.

{% include '/user-guide/use-case-connect-to-k8s-cluster.md' ignore missing %}

[MKI]: https://docs.elastic.dev/mki
[oblt-framework]: /tools/oblt-framework/overview.md
