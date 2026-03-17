# Rotate Ingress credentials

generate a random password

```bash
dd if=/dev/urandom count=15 bs=1 2> /dev/null |base64
TK6zDnZ4dweACVar57HY
```

update the cluster manually, we enable only deploys that use ingress.

```bash
NEW_PASSWOD=TK6zDnZ4dweACVar57HY
CLUSTER_CONFIG_FILE=/path/to/the/cluster/config.yml
ANSIBLE_OPTS="--extra-vars 'password_plaintext=${NEW_PASSWORD} update_k8s=true update_stack=false update_elastic_agent=false update_beats=false update_curator=false update_security=false'" \
make -C ansible update-cluster
```

Generate credentials and utils files

```bash
CLUSTER_CONFIG_FILE=/path/to/the/cluster/config.yml
make -C ansible tasks-pr-changes
```

create a new entry in https://p.elstc.co/ with the credentials with the generated file `ansible/build/credentials`, the expiration should be `never`

update the short link at https://links.elastic.dev/ (e.g https://links.elastic.dev/alias/edge-oblt-credentials)
