---
render_macros: false
---
# k8s-apps Role

## Overview

This Role deploy a k8s manifest files specified in a configuration.

It also has tasks to uninstall the services installed on a Kubernetes cluster.

## Requirements

It is necessary to have gcloud, kubectl, and the Helm CLI installed.

## Role Variables

* apps.k8s: List of manifest files location.
* apps.k8s.name: Name of the k8s app.
* apps.k8s.src: folder where are the k8s manifest files.

```yaml
apps:
  k8s:
    - name: App foo 01
      src: {{ playbook_dir }}/k8s/my_app_01
    - name: App foo 02
      src: {{ playbook_dir }}/k8s/my_app_02
```

The folder that contains the k8s YAML files can contain plain k8s YAML files (*.yml or *.yaml)
and jinja2 templates (*.j2) the templates will be evaluated with the variables defined in the playbook.
check the following examples for a plan yaml file and a jinja2 template.

[plain_file.yaml]

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: hello-node
  name: hello-node
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello-node
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: hello-node
    spec:
      containers:
      - image: k8s.gcr.io/echoserver:1.4
        name: echoserver
        resources: {}
status: {}
```

[template.j2]

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: hello-node-template
  name: hello-node-template
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello-node-template
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: hello-node-template
    spec:
      containers:
      - image: k8s.gcr.io/echoserver:{{ app_config.version }}
        name: echoserver-{{ template_value | default("foo", true) }}
        resources: {}
status: {}
```

in the template the expression `{{ app_config.version }}` will be replaced
with the value passed in the variable `version` in the app definition.
All values defines under `apps.k8s[ITEM].*`
will be translated to `app_config.*` in the template.

```yaml
vars:
  cluster_name: oblt-test
  apps:
    k8s:
      - name: App foo 01
        src: {{ playbook_dir }}/k8s/my_app_01
        version: 1.4
        custom_var00: foo
        custom_var01: bar
```

## Dependencies

It include [common][] and [k8s][] role.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    apps:
      k8s:
        - name: App foo 01
          src: {{ playbook_dir }}/k8s/my_app_01
          version: 1.4
          custom_var00: foo
          custom_var01: bar
        - name: App foo 02
          src: {{ playbook_dir }}/k8s/my_app_02
          version: 1.4
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.k8s
    - role: oblt.framework.k8s_apps
```

## License

Apache License 2.0

[common]:../common/README.md
[k8s]:../k8s/README.md

## Parameters

### General

| Name       | Description                 | Value |
| ---------- | --------------------------- | ----- |
| `k8s_apps` | List of K8s apps to install | `[]`  |
