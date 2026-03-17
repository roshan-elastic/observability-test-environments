# Licensed to Elasticsearch B.V. under one or more contributor
# license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright
# ownership. Elasticsearch B.V. licenses this file to you under
# the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

include common.mk

DOCKER_REPO ?= docker.elastic.co
DOCKER_NAMESPACE ?= observability-ci
DOCKER_IMAGE ?= $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/oblt-framework
DOCKER_IMAGE_RELEASE ?= $(shell cat .ci/.version)
DOCKER_PLATFORM ?= linux/amd64
DOCKER_OPTS ?=
CLUSTER_CONFIG_FILE ?= $(CURDIR)/ansible/test-cluster.yml
VERSION ?= $(shell cat .ci/.version)

.PHONY: docker-build
## @help:docker-build:Build the Docker image.
docker-build:
	docker build ${DOCKER_OPTS} --progress=plain -t $(DOCKER_IMAGE):$(DOCKER_IMAGE_RELEASE) --platform "${DOCKER_PLATFORM}" .

.PHONY: docker-push
## @help:docker-push:Push the Docker image.
docker-push:
	docker push ${DOCKER_OPTS} $(DOCKER_IMAGE):$(DOCKER_IMAGE_RELEASE)

## @help:docker-pull:Pull the Docker image.
docker-pull:
	docker pull ${DOCKER_OPTS} $(DOCKER_IMAGE):$(DOCKER_IMAGE_RELEASE)

.PHONY: docker-run
## @help:docker-run:Run the Docker image. TARGET is the make target to run.
docker-run:
	$(MAKE) -C .github/actions/cluster-target docker-run BUILD_DIR=$(CURDIR)/build VERSION=$(VERSION) CLUSTER_CONFIG_FILE=$(CLUSTER_CONFIG_FILE)

.PHONY: docker-copy-build
## @help:docker-copy-build:Copy the build artifacts from the container to the host.
docker-copy-build:
	$(MAKE) -C .github/actions/cluster-target docker-copy-build BUILD_DIR=$(CURDIR)/build

. PHONY: get-gcp-env
## @help:get-gcp-env:Print the GCP environment variables.
get-gcp-env:
	@echo export GCP_PROJECT=$(shell gcloud config get-value project)
	@echo export GCP_AUTH_KIND=accesstoken
	@echo export GCP_ACCESS_TOKEN=$(shell gcloud auth print-access-token)

.PHONY: docs-portable
## @help:docs-portable:Build the portable docs server in a Docker container.
docs-portable:
	docker build -f Dockerfile_docs -t docker.elastic.co/observability-ci/oblt-clusters-docs-serve:0.0.1 .
