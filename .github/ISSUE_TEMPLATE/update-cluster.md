---
name: Update oblt cluster
about: Issue to track cluster updates.
title: "[update][CLUSTER_NAME][STACK_VERSION] Updated edge cluster to STACK_VERSION"
labels: area:clusters, team:Automation, triage
assignees: ''
projects: 469

---
## Update process

This issue is to track the update of CLUSTER_NAME to STACK_VERSION

* [ ] Notify the process starts at https://status.obs.elastic.dev/
* [ ] Run the update job : https://apm-ci.elastic.co/job/apm-shared/job/oblt-test-env/job/update-test-oblt-cluster/BUILD_NUMBER
* [ ] Update DNS and secrets
* [ ] Update developer clusters
* [ ] Notify the process ends at https://status.obs.elastic.dev/
