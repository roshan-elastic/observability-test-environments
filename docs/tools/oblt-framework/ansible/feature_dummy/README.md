# feature_dummy role

## Overview

This role is used to install and configure `feature_dummy` on a oblt-cluster.
This `dummy` feature is an example of a `feature` that can be installed on a oblt-cluster.

## Requirements

none.

## Dependencies

none.

## Example Playbook

```yaml
- name: Test feature_dummy role
  hosts: localhost
  connection: local
  gather_facts: true
  roles:
    - role: feature_dummy
```

## License

Apache License 2.0

## Parameters

### General

| Name            | Description                   | Value |
| --------------- | ----------------------------- | ----- |
| `dummy_version` | Default value for a parameter | `""`  |
