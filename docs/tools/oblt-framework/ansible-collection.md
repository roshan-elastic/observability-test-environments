# Oblt Framework Ansible Collection

## Overview

The Oblt Framework is a collection of tools and libraries that are used to build and deploy the [Oblt test clusters](/). The framework is designed to be modular and extensible, allowing developers to easily add new features and functionality to their applications.
The framework is build with [Ansible](https://www.ansible.com/), and uses other tools like [Terraform](https://www.terraform.io/), and [Kubernetes](https://kubernetes.io/).
The input of the framework is a YAML file that defines the settings of the cluster,
The cluster config file is created manually or with other tools like the [oblt-cli][] and [oblt-robot].
The syntax of the cluster config file is defined by the [oblt-schema][].

## Contributing

Anyone can contribute to the Oblt Framework, the process is simple, just create a feature branch, make your changes, and create a pull request.

The mos common contributions are:

* Bug fixes: If you find a bug in the framework, please reportit by creating an issue in the [issue tracker][], and if you can fix it, create a pull request. We will love that.
* Documentation: If you find that the documentation is not clear, or you find a typo, please create a pull request with the changes.
* New features: If you have an idea for a new feature, please create an issue in the [issue tracker][], and if you can implement it, create a pull request.
* Bootstrap scripts: If you have a script that can help to bootstrap the environment, please create a pull request with the script.

## Testing

The Oblt Framework uses [Molecule](https://molecule.readthedocs.io/en/latest/) to test the Ansible roles. Molecule is a testing framework for Ansible roles that allows you to test your roles in an isolated environment. Each role has a set of scenarios that define the tests that will be run.
Every scenario test only one use case. The stages of each scenario ussially are:

* Linting: Check the syntax of the role.
* Destroy: Destroy the test environment.
* Converge: Apply the role to the test environment.
* Idempotence: Check if the role is idempotent.
* Verify: Run the tests.

To run the test you can symply run the following command:

```shell
make -C ansible unit-test ROLE=common SCENARIO=no_config
```

ROLE is the name of the role that you want to test, and SCENARIO is the name of the scenario that you want to test.

## Bug fixes

The way to fix an issue it is to have a way to replitate the issue, and then create a test that fails because of the issue, and then fix the issue and create a test that passes.
The easy way to do this is to create a new scenario in the molecule tests, and then create a test that fails, and then fix the issue and create a test that passes.

## Documentation

The documentation is written in markdown and is located in the `docs/tools/oblt-framework/ansible/ROLE` directory. In that folder you will find a `README.md` file that is the main ROLE documentation.
You must link the `README.md` file in the `mkdocs.yml` file to include it in the documentation.
The `README.md` must have the basic structure:

```yaml
---
render_macros: false
---
# Ansible Role: Name of the role

Some description of the role.

## Requirements

List of requirements of the role.

## Example Playbook

Some examples of how to use the role.

## License

Apache License 2.0

## Parameters
```

`Parameters` is mandatory section that it is generated from the `defaults/main.yml` file of the role.
To generate that documentation we use [Readme Generator For Helm][]

The documentation is generated using [MkDocs](https://www.mkdocs.org/), a static site generator that's geared towards project documentation. You can see plenty of examples by checking [cluster config][] file.
The basic definition looks like the following:

```yaml
## @section Describe the Role
## @descriptionStart
## Edtended description of the configuration.
## @descriptionEnd

## @param param00 The description of the parameter
param00: ""
## @param param01 The description of the parameter
param01: 1
## @param param02 The description of the parameter
param02: []
## @param param03 The description of the parameter
param03: {}
## @param param04 The description of the parameter
## @param param04.subparam01 The description of the subparameter
## @param param04.subparam02 The description of the subparameter
## @section param04.complexlist The description of the subparameter
## @section param04.complexlist[0].name The description of the subparameter
## @skip param04.complexlist[1].name
param04:
  subparam01: ""
  subparam02: 1
  complexlist:
    - name: ""
    - name: ""
```

To generate the documentation you can run the following command when you finish the documentation:

```shell
make -C .ci docs-ansible
```

Finally you can test the documentation by running an MKDocs servere locally (localhost:8000) with the following command:

```shell
make -C .ci docs-serve
```

## Creating a role

Roles are the main building blocks of an Ansible playbook. A role is a collection of tasks, handlers, templates, and variables that can be used to configure a system. Roles are organized into directories, with each directory containing a specific type of file.
Our roles are located in the `ansible/ansible_collections/oblt/framework/roles` directory.
To create a new role you can use the following command:

```shell
make -C ansible create-role ROLE=MY-ROLE-NAME
```

This command will create a new role in the `ansible/ansible_collections/oblt/framework/roles` directory with the name `MY-ROLE-NAME`. It creates the basic structure of the role, with directories for tasks, handlers, templates, and variables. Remove those directories if you don't need them. Modify the files to fit your needs.

Implement you new feature in the role, and then create a new scenario in the molecule tests to test the new feature. Check other roles to see how to create a new scenario.

[Introduction to Molecule][]

## Releases

The release process is simple, just change the version number in the `ansible/ansible_collections/oblt/framework/galaxy.yml` and `.ci/.version` files. Then create a new pull request with the changes.
This will trigger a new release that it is published as a Docker image in the [Elastic Docker registry][] and in the [release assets][] as a tar file.

The Docker image is published with the following tags:

* `latest`: The latest version of the framework from the branch `main`.
* `X.Y.Z`: The version of the framework.
* `pr-NUMVER`: When you create a pull request a new image is published with the PR number.

to pull the imagen for the version `X.Y.Z` you can use the following command, it requires authentication with the Elastic Docker registry:

```shell
docker pull docker.elastic.co/observability-ci/oblt-framework:X.Y.Z
```

To download the tar file you can use the following command:

```shell
gh release download X.Y.Z --pattern 'oblt-framework-*.tar.gz'
ansible-galaxy collection install oblt.framework-X.Y.Z.tar.gz
```

## Troubleshooting

### Python fork errors

Mac OS Mojave has some security option enabled that makes Python crash when makes a fork, to fix it you have to disable a security option by exporting the environment variables `OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES` and `no_proxy="*"`

```log
objc[2066]: +[__NSPlaceholderDate initialize] may have been in progress in another thread when fork() was called.
objc[2066]: +[__NSPlaceholderDate initialize] may have been in progress in another thread when fork() was called. We cannot safely call it or ignore it in the fork() child process. Crashing instead. Set a breakpoint on objc_initializeAfterForkError to debug.
```

[oblt-schema]: ../../user-guide/cluster-config.md
[issue tracker]: https://github.com/elastic/observability-test-environments/issues
[oblt-cli]: ela.st/oblt-cli
[Readme Generator For Helm]: https://github.com/bitnami-labs/readme-generator-for-helm
[cluster config]: https://github.com/elastic/observability-test-environments/blob/main/docs/config-cluster.yml
[Introduction to Molecule]: https://www.youtube.com/watch?v=DAnMyBZ8-Qs
[Elastic Docker registry]: https://container-library.elastic.co/
[release assets]: https://github.com/elastic/observability-test-environments/releases
