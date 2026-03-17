# Oblt Tools

## Release oblt repo and tools

We have a release process for the oblt repo and tools.
The release process will release a new version of the oblt repo and tools every time the `.ci/version` is modified on the `main` branch.
The release process has the following cases:

## Manual release

There is a workflow name `Release`, this workflow allows you to select the type of release to make, then will bump the version and make a PR.

* Go to GitHub Actions, to the workflow `Release`
* Run the workflow selecting the type of release [major, minor, patch]
* Wait for the workflow to finish
* Check the PRs for a PR named `Bump version X.Y.Z`
* Check that everything is fine and merge the PR
* The workflow `release-tools` is triggered
  * It creates a new GitHub Release with release notes
  * It builds and uploads the binaries to the release.

You can use the same steps the CI use, the file `.ci/Makefile` has all the target needed.
To check the targets available you can use `make -C .ci help`.
The command `make -C .ci release-patch` perform a patch release.
The command `make -C .ci release-minor` perform a minor release.
The command `make -C .ci release-major` perform a major release.

## Release after a PR is merged

The workflow `release-tools` is triggered also when a PR contains changes on `.ci/version`, so to perform a release after merge a PR, it is enough to bump the version of the file `.ci/version`, then push the changes to your PR, and finally when the PR is merged a new release would be performed.

* Make your changes on the repo
* Update `.ci/version` to a new version
* In the project root dir run: `make -C tools/oblt-cli build` and then `make -C tools/oblt-robot build`
* Create a PR
* Check all your changes are correct
* Merge the PR
* The workflow `release-tools` is triggered
  * It creates a new GitHub Release with release notes
  * It builds and uploads the binaries to the release.
