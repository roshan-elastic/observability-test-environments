# Oblt Framework

## Overview

The Oblt Framework is a collection of tools and libraries that are used to build and deploy the [Oblt test clusters](/). The framework is designed to be modular and extensible, allowing developers to easily add new features and functionality to their applications.

The framework is built on top of several open-source technologies, including [Ansible](https://www.ansible.com/), [Terraform](https://www.terraform.io/), and [Kubernetes](https://kubernetes.io/). These tools are used to automate the deployment and management of the test clusters.

The main component is a set of [Ansible playbooks and roles](./ansible-collection.md) that define the configuration of the test clusters. These playbooks are used to deploy the various components of the test clusters, such as Elasticsearch, Kibana, and other applications. Every role is designed to be isolated, reusable, and easy to extend. The playbooks define a set of tasks that are executed to perform more complex operations, such as create a new cluster, destroy an existing cluster, or update an existing cluster. All these components are published as an [Ansible collection](./ansible-collection.md) with each release of the project in GitHub.

To facilitate the use of the framework, the [Ansible Collection](ansible-collection.md) is packaged as a Docker image that can be run on any machine with Docker installed. This Docker container contains all the necessary dependencies to run the playbooks, including Ansible, Terraform, and the Kubernetes CLI.

Finally, the framework includes GitHub Actions and GitHub Actions Workflows that are used to automate the deployment of the test clusters. These workflows are triggered by events in the GitHub repository, such as a new pullrequest or a new release. The workflows use the Docker image to run the playbooks and deploy the test clusters.
