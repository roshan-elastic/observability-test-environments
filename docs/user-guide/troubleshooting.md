# Troubleshooting

{% include '/user-guide/how-to-report-issues.md' ignore missing %}

## oblt-cli

If it is the first time you use oblt-cli, check that you have [configured oblt-cli before use it](https://elastic.github.io/observability-test-environments/tools/oblt-cli/#configure).
The default log level can hide some errors, try to add the flag `--verbose` to have more information about the error.

### Error 128

There is some errors that are not shown in the default log,
in order to have more information about the error you must add the `--verbose` flag.
The most common issue is that your user does not have access to the [observability-test-environments repo](https://github.com/elastic/observability-test-environments), try to navigate to the previous link, and the following commands:

```bash
git clone git@github.com:elastic/observability-test-environments.git
```

or

```bash
git clone https://github.com/elastic/observability-test-environments.git
```

### Unable to update oblt-cli

The most common issue self-updating oblt-cli is a token that does not have read permission on the Elastic org.
You will need to add repo:read permission to the token. See [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

```bash
>oblt-cli update --github-token=phr_1234567890
2023/02/02 19:51:24 Binary update failed: GET https://api.github.com/repos/elastic/observability-test-environments/releases: 401 Bad credentials []
exit status 1
```

Your token must have SSO enabled and `repo:*` permissions to access to the oblt test clusters repository.
For more detail see [GitHub Token for oblt-cli](./github-token.md).

```bash
go run . update --github-token=ghp_emTy.......
2023/02/02 19:53:46 Binary update failed: GET https://api.github.com/repos/elastic/observability-test-environments/releases: 403 Resource protected by organization SAML enforcement. You must grant your Personal Access token access to this organization. []
exit status 1
```

You can use `gh` to generate the token on the fly:

```bash
oblt-cli update --github-token=$(gh auth token)
```

### Unable to retrieve secrets 403 error

This error happens when the token does not have access to the secrets in the vault, or the token is expired.
If you are in the `observability` team in GitHub, you must have permissions to access the secrets.
To check the permissions, you can use the following command:

```bash
> vault token lookup
Key                 Value
---                 -----
accessor            FFFFFFFFFF
creation_time       1681285635
creation_ttl        768h
display_name        github-my-github-username
entity_id           FFFFFFFFFF-3e73-270f-5248-FFFFFFFFFF
expire_time         2023-05-14T07:47:15.663631481Z
explicit_max_ttl    0s
id                  s.FFFFFFFFFF
issue_time          2023-04-12T07:47:15.66364637Z
meta                map[org:elastic username:my-github-username]
num_uses            0
orphan              true
path                auth/github/login
policies            [default employees observability]
renewable           true
ttl                 620h8m15s
type                service
```

The policies must contain the `observability` policy.
If this command return a 403 error, you must create a new token, because it is expired.
The following command will create a new token:

```bash
vault login -method=github token=$(gh auth token)
```

For more detail see [GitHub Token for oblt-cli](./github-token.md).

### Unable to create a new cluster it fails to push to main

In Nov of 2022 we have changed the workflow to use PRs instead of pushes to the main branch.
If you see the following error, you must update to the latest version of oblt-cli.

```bash
2022/11/24 15:28:01 remote: error: GH006: Protected branch update failed for refs/heads/main.
remote: error: Changes must be made through a pull request.
To github.com:elastic/observability-test-environments.git
 ! [remote rejected]   main -> main (protected branch hook declined)
error: failed to push some refs to 'github.com:elastic/observability-test-environments.git'
Error: exit status 1
```

### Cannot checkout using SSH

You are trying to use `oblt-cli` with SSH, but the checkout fails because is using HTTP.
The first thing to check is that you have configured the `oblt-cli` with the correct user and git mode.

```bash
$ oblt-cli configure
2023/07/12 18:07:43 Writing configuration file /home/j/.oblt-cli/config.yaml
2023/07/12 18:07:43 SlackChannel: '@FOO'
2023/07/12 18:07:43 User: 'bar'
2023/07/12 18:07:43 Git mode: 'ssh'
```

Git mode: 'ssh' is configured, but the checkout fails

```bash
$ oblt-cli cluster list
2023/07/12 18:08:04 Using config file: /home/foo/.oblt-cli/config.yaml
Username for 'https://github.com': foo
Password for 'https://foo@github.com':
2023/07/12 18:08:16 Cloning into 'observability-test-environments'...
remote: Invalid username or password.
```

The issue can be that you have a url configured in your git config file.
Check that you do not have something like this in your `~/.gitconfig` file:

```ini
[url "https://"]
    insteadOf = git://
[url "https://github.com/"]
    insteadOf = git@github.com:
```

The following command will configure git to use ssh for any `https://github.com/` GitHub URL

```bash
git config --global url.ssh://git@github.com/.insteadOf https://github.com/
```

## DNS resolution issues

Sometimes the DNS server fails to resolve the snapshot URLs. To workaround
this, add an entry to your `/etc/hosts` or
`C:\Windows\System32\Drivers\etc\hosts` file.

The IP of the host is contained in the host name.

```txt
XXX.XXX.XXX.XXX  elasticsearch.XXX.XXX.XXX.XXX.xip.io
```

## Unsupported Elasticsearch version

The following error happens when we try to connect Kibana to an unsupported Elasticsearch version.
The solution is to use the same version for Kibana and Elasticsearch.

```json
{"type":"log","@timestamp":"2019-09-02T15:51:30Z","tags":["status","plugin:grokdebugger@7.3.1","error"],"pid":8,"state":"red","message":"Status changed from green to red - This version of Kibana requires Elasticsearch v7.3.1 on all nodes. I found the following incompatible nodes in your cluster: v7.2.1 @ 10.16.0.131:9200 (10.16.0.131), v7.2.1 @ 10.16.1.184:9200 (10.16.1.184), v7.2.1 @ 10.16.0.132:9200 (10.16.0.132), v7.2.1 @ 10.16.2.118:9200 (10.16.2.118), v7.2.1 @ 10.16.2.119:9200 (10.16.2.119)","prevState":"green","prevMsg":"Ready"}
```

## Kibana index migration issue

In some cases when we change between minor or major version the .kibana* indices fail to migrate data
and it gets stuck with the following message,
the solution is to delete the index that appears in the message.

```json
{"type":"log","@timestamp":"2019-09-02T15:06:59Z","tags":["warning","migrations"],"pid":9,"message":"Another Kibana instance appears to be migrating the index. Waiting for that migration to complete. If no other Kibana instance is attempting migrations, you can get past this message by deleting index .kibana_k8s_3 and restarting Kibana."}
```

In any case, the data that we would keep in Elasticsearch is the .kibana* indices,
everything else is generated by load jobs, so it is fine to delete the whole indices but .kibana*,
[elasticdump](https://www.npmjs.com/package/elasticdump) is a nice tool to make a backup of
those indices and restore them later.

## Fix beats indices mappings

Sometimes the index templates, index mapping, ILM policies, or other related index settings are wrong or corrupt, in those cases, we can restore the original settings by using the `setup` command from the beat binary, the following examples show you how:

This example will restore indices settings for Metricbeat using a Metricbeat pod to run the command and
connecting to the internal DNS of Elasticsearch.

```bash
kubectl exec -n kube-system metricbeat-deploy-bf9fdfd7f-7qnl9  -- metricbeat setup --index-management -E output.logstash.enabled=false -E 'output.elasticsearch.hosts=["elasticsearch-master.default.svc.cluster.local:9200"]'
```

This example will restore indices settings for Filebeat using a Filebeat pod to run the command and
connecting to the internal DNS of Elasticsearch.

```bash
kubectl exec -n kube-system filebeat-p4t7c -- filebeat setup --index-management -E output.logstash.enabled=false -E 'output.elasticsearch.hosts=["elasticsearch-master.default.svc.cluster.local:9200"]'
```

Finally, this example would use a Docker container on a fixed version to load the indices settings
using the https URL of an Elasticsearch protected with username and password.

```bash
docker run docker.elastic.co/beats/metricbeat:7.3.0 setup --index-management -E output.logstash.enabled=false -E 'output.elasticsearch.hosts=["https://elasticsearch.35.189.242.157.ip.es.io:443"]' -E 'output.elasticsearch.username=admin' -E 'output.elasticsearch.password=INGRESS_PASSWORD' -E 'output.elasticsearch.ssl.verification_mode=none'
```

## Unable to enable Profiling

I have deployed a cluster and when I try to enable profiling I get the following error:

```log
[kibana.log][ERROR] Saved object [ingest-agent-policies/policy-elastic-agent-on-cloud] not found
Error: Saved object [ingest-agent-policies/policy-elastic-agent-on-cloud] not found
```

The error is because the policy does not exist, the policy does not exist because the Integration server is not configured in that cluster.
To fix this issue, you must deploy the Integration server in the cluster.

You can edit the cluster configuration file and add the following configuration and make a PR to the observability-test-environments repository:

```yaml
stack:
  ess:
    integrations:
      enabled: true
```

Or update the cluster using the following command:

```shell
CLUSTER_NAME=my-cluster
oblt-cli cluster update --cluster-name=${CLUSTER_NAME} --parameters '{"stack":{"ess":{"integrations":{"enabled": true}}}}'
```
