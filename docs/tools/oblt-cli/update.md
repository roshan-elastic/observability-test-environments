# Update oblt-cli

Since v2.0.30, the oblt-cli comes with a `oblt-cli version` command that prints out the current version of the tool so that you can check out what version you are using.

There are several ways to update oblt-cli.
The easy way is to use the `oblt-cli update` command.
You only need to provide a GitHub token with the `--github-token` option (see [GitHub Token for oblt-cli](../../user-guide/github-token.md)).


The `oblt-cli update` command will download the latest version of oblt-cli and replace the current binary.

```bash
oblt-cli update --github-token "$(gh auth token)"
```

> **Note:**
> `macOS` users will have to manually verify the binary using this command:

```shell
xattr -r -d com.apple.quarantine oblt-cli
```

## macOS via Homebrew

> **NOTE:**
> Since `v6.3.1`

To check, later, if a newer version is available, use `brew outdated oblt-cli`. Should you then want to upgrade to the latest version, use `brew upgrade oblt-cli`.


If you cannot install it, please be sure you set the `HOMEBREW_GITHUB_API_TOKEN` environment variable:

```bash
export HOMEBREW_GITHUB_API_TOKEN=$(gh auth token)
brew outdated oblt-cli
brew upgrade oblt-cli
```

## Update oblt-cli downloading the latest release

You can also download the latest release from the [GitHub releases page](https://github.com/elastic/observability-test-environments/releases) and replace the current binary.
