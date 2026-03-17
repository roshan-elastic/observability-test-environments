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

DOCKER_CMD ?= docker compose
DOCKER_OPTS ?=
BUILD_DIR ?= $(CURDIR)/build
VERSION ?= latest
CLUSTER_CONFIG_FILE ?= no-configuration-file-passed

SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

.ONESHELL:

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

.PHONY: help
help:
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

.PHONY: docker-run
## @help:docker-run:Run the Docker image. TARGET is the make target to run.
docker-run:
	@$(call check_defined, GCP_PROJECT, GCP_PROJECT must be defined. )
	@$(call check_defined, GCP_ACCESS_TOKEN, GCP_ACCESS_TOKEN must be defined. )
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
# TODO: remove these variables
	@$(call check_defined, DOCKER_USERNAME, DOCKER_USERNAME must be defined. )
	@$(call check_defined, DOCKER_PASSWORD, DOCKER_PASSWORD must be defined. )
	@$(call check_defined, DOCKER_REGISTRY, DOCKER_REGISTRY must be defined. )
	rm -rf $(BUILD_DIR) || true
	mkdir -p $(BUILD_DIR)
	echo "BUILD_DIR=$(BUILD_DIR)"
	echo "CLUSTER_CONFIG_FILE=$(CLUSTER_CONFIG_FILE)"
	VERSION="$(VERSION)" \
	$(DOCKER_CMD) run -i --rm ${DOCKER_OPTS} \
		-e GCP_PROJECT \
		-e GCP_AUTH_KIND=accesstoken \
		-e GCP_ACCESS_TOKEN \
		-e "CLOUDSDK_AUTH_ACCESS_TOKEN=$${GCP_ACCESS_TOKEN}" \
		-e "GOOGLE_OAUTH_ACCESS_TOKEN=$${GCP_ACCESS_TOKEN}" \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e DOCKER_REGISTRY \
		-e ANSIBLE_OPTS \
		-e GITHUB_ENABLED \
		-v "$(CLUSTER_CONFIG_FILE):/ansible/cluster-config.yml" \
		oblt-framework $${TARGET}
	$(MAKE) docker-copy-build BUILD_DIR=$(BUILD_DIR)
	$(DOCKER_CMD) down -v

.PHONY: docker-copy-build
## @help:docker-copy-build:Copy the build artifacts from the container to the host.
docker-copy-build:
	NEW_UID=$(id -u) \
	NEW_GID=$(id -g) \
	BUILD_DIR=$(BUILD_DIR) \
	$(DOCKER_CMD) run -i --rm copy-files
