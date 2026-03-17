# ESS

ESS template allow to deploy full Elastic Stack in ESS.
The following command will deploy a new cluster using the 8.10.2 version of the stack:

```shell
oblt-cli cluster create ess --stack-version 8.10.2 --release
```

To deploy a snapshot version of the stack, use the following command:

```shell
oblt-cli cluster create ess --stack-version 8.11.0-SNAPSHOT
```

!!! Note

    The `--release` flag is important when you deploy releases. If you do not set it the Stack versions is considered as a snapshot version.

!!! Note

    This template only deploy the Elastic Stack, check other templates if you need data.
