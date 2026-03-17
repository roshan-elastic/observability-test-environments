## About

These k8s resources are meant for usage in GKE Autopilot and Standard clusters, but might be fitting for other kinds of clusters as well.
It's following the guide from https://cloud.google.com/kubernetes-engine/docs/how-to/capacity-provisioning

## TL;DR

10 spare pods will be deployed with a low priority class.

When a test or oblt-cli is creating resources with a higher or default priority, the spare pods will be evicted to make space
for the created resources.

This will speed up provisioning and lower the chance of timeouts during e.g. helm installations.
