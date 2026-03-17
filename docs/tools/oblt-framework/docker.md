# Oblt framework on Docker

## Overview

The Oblt framework on Docker is a containerized version of the [Oblt framework][].
It is basically a Python container with the Oblt framework installed and all the necessary dependencies.

## Prerequisites

* Access to the Elastic Docker registry is required to pull the image.
* Access to the Google Cloud Platform (GCP) is required to create/access GCP resources.
* Access to the Vault is required to create/retrieve secrets (deprecated).

## Run

The Oblt framework on Docker can be run using a simple `docker run` command.

```shell
gcloud auth login

CLUSTER_CONFIG_FILE=cluster-config.yml
GCP_PROJECT=$(gcloud config get-value project)
GCP_AUTH_KIND=accesstoken
GCP_ACCESS_TOKEN=$(gcloud auth print-access-token)
CLOUDSDK_AUTH_ACCESS_TOKEN=${GCP_ACCESS_TOKEN}
GOOGLE_OAUTH_ACCESS_TOKEN=${GCP_ACCESS_TOKEN}

docker run -i --rm \
  -e GCP_PROJECT \
  -e GCP_AUTH_KIND \
  -e GCP_ACCESS_TOKEN \
  -e CLOUDSDK_AUTH_ACCESS_TOKEN \
  -e GOOGLE_OAUTH_ACCESS_TOKEN \
  -e DOCKER_USERNAME \
  -e DOCKER_PASSWORD \
  -e DOCKER_REGISTRY \
  -e ANSIBLE_OPTS \
  -e GITHUB_ENABLED \
  -v "${CLUSTER_CONFIG_FILE}:/ansible/cluster-config.yml" \
  docker.elastic.co/observability-ci/oblt-framework:latest help
```

## Build

The observability test environments repository contains a Makefile with the following targets:

* `docker-build`: Build the Docker image.
* `docker-push`: Push the Docker image to the Elastic Docker registry.
* `docker-run`: Run the Docker image locally.
* `docker-pull`: Pull the Docker image from the Elastic Docker registry.

```shell
make docker-build
```

```shell
make docker-push
```

## Release

The release process is automated using GitHub Actions.
Every time a new tag is pushed to the repository, the Docker image is built and pushed to the Elastic Docker registry.
This make that everityme a new version of [oblt-cli][] is release a new version of the Docker image is also release.

[Release workflow][]

## Troubleshooting

[Oblt framework]: ./ansible-collection.md
[Release workflow]: https://github.com/elastic/observability-test-environments/blob/main/.github/workflows/release-oblt-framework.yml
[oblt-cli]: http://ela.st/oblt-cli
