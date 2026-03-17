#!/usr/bin/env bash
# script to list the relevant to migrate indices on a Elasticsearch
kubectl exec -it elasticsearch-master-0 -- curl -X GET localhost:9200/_cat/indices| \
  cut -d " " -f 3|\
  grep -v \
    -e '^apm' \
    -e 'task_manager' \
    -e '.monitoring' \
    -e '.watcher-' \
    -e 'heartbeat' \
    -e 'auditbeat' \
    -e 'task-manager' \
    -e '^ilm' \
    -e '.watches' \
    -e '.apm-agent-configuration' \
    -e 'triggered_watches' \
    -e '.tasks'


kubectl get secret elastic-basic-auth  -o jsonpath='{.data.plaintext}'|base64 --decode
kubectl get ingress elasticsearch-master -o jsonpath='{.spec.rules[0].host}'
