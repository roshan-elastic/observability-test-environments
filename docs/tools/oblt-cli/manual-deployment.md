# Manual deployments

It is possible to generate the configuration files and tun the Ansible playbook locally.
To do that, you can use the `--dry-run` flag, this flag will generate the cluster configuration files
but oblt-cli will not push the changes to the repository so the CI automation is not trigger.

```shell
> oblt-cli cluster create ccs --remote-cluster edge-oblt --dry-run
2021/10/04 10:39:05 Root Args: [/var/folders/rd/4lpsq9tn2h35ty0nkl9g4hlh0000gn/T/go-build4180287750/b001/exe/main cluster create ccs --remote-cluster edge-oblt --dry-run]
2021/10/04 10:39:05 Using config file: /home/my-user/.oblt-cli/config.yaml
2021/10/04 10:39:05 Reading cluster config edge-oblt-nwxrh-ccs
2021/10/04 10:39:05 Writing cluster configuration file /home/my-user/.oblt-cli/observability-test-environments/environments/home/my-user/edge-oblt-nwxrh-ccs.yml
2021/10/04 10:39:05 Template parameters: map[ClusterName:edge-oblt-nwxrh-ccs Date:2021-10-04T10:39:05+0200 ElasticsearchDockerImage:docker.elastic.co/cloud-ci/elasticsearch:8.0.0-alpha2-1b2647e3 RemoteClusterName:edge-oblt Seed:0x1369780 SlackChannel:@UCKPL50JY StackVersion:8.0.0-alpha1 User:inifc]
2021/10/04 10:39:05 git.add
2021/10/04 10:39:05
2021/10/04 10:39:05 git.commit
2021/10/04 10:39:05 [master 1b87c11] oblt-cli: Create a CCS cluster for inifc
 1 file changed, 46 insertions(+)
 create mode 100644 environments/home/my-user/edge-oblt-nwxrh-ccs.yml
```

Once the file is generated `/home/my-user/.oblt-cli/observability-test-environments/environments/home/my-user/edge-oblt-nwxrh-ccs.yml` you can use it to make the deploy.

```bash
CLUSTER_CONFIG_FILE=/home/my-user/.oblt-cli/observability-test-environments/environments/home/my-user/edge-oblt-nwxrh-ccs.yml \
make -C ansible create-cluster
```
