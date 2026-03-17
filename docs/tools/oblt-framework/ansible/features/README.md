# Features role

## Overview

This role is used to install and configure `features` on a oblt-cluster.
A `feature` is a well defined workflow that can be executed on a oblt-cluster.
A `feature` can be a `service`, `application`, `tool`, `library`, etc.

## Feature

A `feature` is a well defined workflow that can be executed on a oblt-cluster.
The implementation of a `feature` is in a `feature_FEATURE_NAME` role.
The `feature` role is responsible for installing and configuring the `feature` on a oblt-cluster.
The `feature` must have a `main.yml`, `deploy.yml` and `undeploy.yml` file in the `tasks` directory.
The `main.yml` file is in charge to make the whole deployment and configuration of the `feature`.
The `deploy.yml` file is in charge to deploy the `feature`.
The `undeploy.yml` file is in charge to undeploy the `feature`.
The `feature_dummy` role is an example of a `feature` role.

## Use features role

To use the `features` role, you need to define a list of `features` to install.
The `features` list is a list of dictionaries with the following keys:

```yaml
features:
  - name: dummy
```

The `name` key is the name of the `feature` to install.
The feature can have additional parameters that are passed to the `feature` role.

```yaml
features:
  - name: dummy
    version: 1.0.0
```

Those parameters are passed to the `feature_dummy` role and are used as variables in the `feature_dummy` role.

## Requirements

A oblt cluster deployed.

## Dependencies

none.

## Example Playbook

```yaml
- name: Test features role
  hosts: localhost
  connection: local
  gather_facts: true
  vars:
    features:
      - name: dummy
  roles:
    - role: features
```

## License

Apache License 2.0

## Parameters

### General

| Name       | Description                 | Value |
| ---------- | --------------------------- | ----- |
| `features` | List of features to install | `[]`  |
