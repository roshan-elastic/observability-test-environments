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
UPDATECLI_OPTS ?= --clean=true --debug=true
JQ_QUERY ?= .
JSON_FILE ?= test.json

CLUSTER_FINDER=$(shell realpath $(CURDIR)/../../.ci/scripts/find_cluster_config_by_name.py)
READ_YAML_KEY=$(shell realpath $(CURDIR)/../../.ci/scripts/get_yaml_value.py)
ENVIRONMENTS_DIR=$(shell realpath $(CURDIR)/../../environments)
SCRIPTS_DIR=$(shell realpath $(CURDIR)/../../.ci/scripts)

ELASTICSEARCH_DOKER_IMAGE=docker.elastic.co/observability-ci/elasticsearch-cloud-ess
KIBANA_DOCKER_IMAGE=docker.elastic.co/observability-ci/kibana-cloud
ELASTIC_AGENT_DOCKER_IMAGE=docker.elastic.co/observability-ci/elastic-agent-cloud
PROFILING_DOCKER_IMAGE=docker.elastic.co/observability-ci/profiling-agent

export GIT_USER ?= $(shell git config user.name)
export GIT_EMAIL ?= $(shell git config user.email)
export PATH := ${PATH}:${HOME}/bin

PYTHON ?= python3
PIP ?= pip3
VENV ?= .venv
VENV_ACTIVATE = $(VENV)/bin/activate

SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

.PHONY: help
help:
	@echo "Environment variables:"
	@echo ""
	@echo "UPDATECLI_OPTS additional flags to pass to updatecli (e.g. --debug)"
	@echo "GIT_USER User used for the git operations."
	@echo "GIT_EMAIL Email used for the git operations."
	@echo "GITHUB_TOKEN Token to perform GitHub authentication."
	@echo "ELASTIC_DOCKER_USERNAME The user to authenticate on the Elastic Docker registry."
	@echo "ELASTIC_DOCKER_PASSWORD The password to authenticate on the Elastic Docker registry."
	@echo "Main targets:"
	@echo ""
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

## @help:venv:Prepare a virtualenv for Python.
.PHONY: venv
venv:
	$(PYTHON) --version 1>&2;
	$(PYTHON) -m venv $(VENV) 1>&2;\
	source $(VENV_ACTIVATE);\
	$(PYTHON) -m pip install -U pip 1>&2;\
	$(PIP) install --prefer-binary -r $(CURDIR)/../requirements.txt 1>&2;

## @help:update-tools:Check if there is new Tools versions and make a PR with the changes.
.PHONY: update-tools
update-tools: venv
	source $(VENV_ACTIVATE); \
	updatecli apply $(UPDATECLI_OPTS) --config  $(CURDIR)/tools.d

## @help:update-elastic-stack-ess:Update Elastic Stack version on a ESS cluster.(requires CLUSTER_NAME)
.PHONY: update-elastic-stack-ess
update-elastic-stack-ess: venv
	source $(VENV_ACTIVATE); \
	export CLUSTER_CONFIG_FILE=$$($(CLUSTER_FINDER) ${CLUSTER_NAME} $(ENVIRONMENTS_DIR) "true"); \
	export UPDATE_CHANNEL=$$($(READ_YAML_KEY) "$(ENVIRONMENTS_DIR)/$${CLUSTER_CONFIG_FILE}" 'stack.update_channel'); \
	export STACK_MODE=$$($(READ_YAML_KEY) "$(ENVIRONMENTS_DIR)/$${CLUSTER_CONFIG_FILE}" 'stack.mode'); \
	echo "CLUSTER_CONFIG_FILE=$${CLUSTER_CONFIG_FILE}"; \
	echo "UPDATE_CHANNEL=$${UPDATE_CHANNEL}"; \
	echo "STACK_MODE=$${STACK_MODE}"; \
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-elastic-stack.yml

## @help:stack-version-unstable:Returns the version we use on main branch.
.PHONY: stack-version-unstable
stack-version-unstable: venv
	source $(VENV_ACTIVATE); \
	$(CURDIR)/../../.ci/scripts/stack-versions.py --output text --mode snapshot-unique --limit 1

## @help:stack-version-development:Returns the version we in active development.
.PHONY: stack-version-development
stack-version-development: venv
	source $(VENV_ACTIVATE); \
	$(CURDIR)/../../.ci/scripts/stack-versions.py --output text --mode snapshot --position -2

## @help:stack-version-release:Returns the latest release or the current build candidate, it depends of the date is close to a FF/release.
.PHONY: stack-version-release
stack-version-release:
	source $(VENV_ACTIVATE); \
	$(CURDIR)/../../.ci/scripts/stack-versions.py  --output text --mode release --limit 1

## @help:stack-build-%:Returns the version we should use on '%' environment.
.PHONY: stack-build-%
stack-build-%:
	ELASTIC_STACK_VERSION=$$(make stack-version-$*) \
	ELASTIC_STACK_BUILD=latest \
	$(MAKE) artifacts-versions |jq -r '.build.build_id'

## @help:build-date-%:Returns the build date of the version we should use on '%' update channel.
.PHONY: build-date-%
build-date-%:
	ELASTIC_STACK_VERSION=$$(make stack-version-$*) \
	ELASTIC_STACK_BUILD=$$(make stack-build-$*) \
	$(MAKE) artifacts-versions |jq -r '.build.end_time'

## @help:artifacts-versions:Returns the artifacts versions for the given ELASTIC_STACK_VERSION and ELASTIC_STACK_BUILD.
artifacts-versions:
	[ -n "$${ELASTIC_STACK_VERSION}" ] || (echo "ELASTIC_STACK_VERSION is not defined or empty"; exit 1)
	[ -n "$${ELASTIC_STACK_BUILD}" ] || (echo "ELASTIC_STACK_BUILD is not defined or empty"; exit 1)
	curl -sSL "https://artifacts-api.elastic.co/v1/versions/$${ELASTIC_STACK_VERSION}/builds/$${ELASTIC_STACK_BUILD}"

## @help:is-upgrade-version-%:Returns if the given version is a upgrade compared with the '%' environment.
.PHONY: is-upgrade-version-%
is-upgrade-version-%:
	[ -n "$$TO_VERSION" ] || (echo "TO_VERSION is not defined"; exit 1)
	$(SCRIPTS_DIR)/stack_is_update.py --cluster-config "$(ENVIRONMENTS_DIR)/$${CLUSTER_CONFIG_FILE}" --version "$${TO_VERSION}"

## @help:is-downgrade-version:Returns if the given version is a downgrade compared with the current cluster. Requires environment variables FROM_VERSION and TO_VERSION.
is-downgrade-version:
	test -n "$$FROM_VERSION" || (echo "FROM_VERSION is not defined"; exit 1)
	test -n "$$TO_VERSION" || (echo "TO_VERSION is not defined"; exit 1)
	echo "Update $$FROM_VERSION to $$TO_VERSION"
	$(CURDIR)/../../.ci/scripts/stack-versions.py --mode is-greater-or-equal --from "$$FROM_VERSION" --to "$$TO_VERSION"

## @help:update-channel-value-%:Returns the value of a JSON path from the update channels file (e.g release.version, release.build).
update-channel-value-%:
	jq -r -c '.$*' $(CURDIR)/update_channels.json

## @help:is-downgrade-channel-%:Returns if the given version is a downgrade compared with the '%' channel.
is-downgrade-channel-%:
	FROM_VERSION=$$(make update-channel-value-$*.version) \
	TO_VERSION=$$TO_VERSION \
	$(MAKE) is-downgrade-version

## @help:update-hermit:Update Hermit and packages.
.PHONY: update-hermit
update-hermit:
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-hermit-packages.yml

## @help:update-opbeans: Update Opbeans and dependencies.
.PHONY: update-opbeans
update-opbeans:
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-opbeans.yml --values ./opbeans.yaml

## @help:update-helm-charts:Check if there is new Helm Chart versions and make a PR with the changes.
.PHONY: update-helm-charts
update-helm-charts:
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-helm-charts.yaml --values ./helm-charts.yaml

## @help:gke-latest-version:Returns the latest version available on GKE.
.PHONY: gke-latest-version
gke-latest-version:
	gcloud container get-server-config --format "json(channels)" --zone us-central1 | jq -r '.channels[]|select(.channel=="STABLE")|.defaultVersion'

## @help:update-dev-cluster:Update the dev cluster.
.PHONY: update-dev-cluster
update-dev-cluster:
	@echo CLUSTER_NAME=${CLUSTER_NAME}
	@echo REMOTE_CLUSTER_CONFIG_FILE=${REMOTE_CLUSTER_CONFIG_FILE}
	@echo CLUSTER_CONFIG_FILE=${CLUSTER_CONFIG_FILE}
	@echo STACK_MODE=${STACK_MODE}
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-dev-cluster.yml

## @help:set-key-value:Update a key value on a YAML file, uses FILE for the relative file path, KEY and SET_VALUE to set the value in the YAML file.
.PHONI: set-key-value
set-key-value:
	@echo 	yq -i 'eval(strenv(KEY)) = strenv(SET_VALUE)' ${CURDIR}/../../${FILE}
	yq -i 'eval(strenv(KEY)) = strenv(SET_VALUE)' ${CURDIR}/../../${FILE}

## @help:update-channels:Check if there is new update stack versions and make a PR with the changes in the update channels.
.PHONY: update-channels
update-channels: venv
	source $(VENV_ACTIVATE); \
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/channels.d

## @help:update-json:Update a JSON_FILE file with the given JQ_QUERY.
update-json:
	@echo "JSON_FILE=$(JSON_FILE)"
	@echo "JQ_QUERY=$(JQ_QUERY)"
	jq '$(JQ_QUERY)' $(CURDIR)/$(JSON_FILE) > $(CURDIR)/$(JSON_FILE).tmp
	mv $(CURDIR)/$(JSON_FILE).tmp $(CURDIR)/$(JSON_FILE)

## @help:update-channel-json:Update the update channel JSON file with the given UPDATE_CHANNEL, JSON_FILE, UPDATE_CHANNEL_VERSION, UPDATE_CHANNEL_BUILD and UPDATE_CHANNEL_DATE.
update-channel-json:
	test -n "$${UPDATE_CHANNEL}" || (echo "UPDATE_CHANNEL is not defined"; exit 1)
	test -n "$${JSON_FILE}" || (echo "JSON_FILE is not defined"; exit 1)
	test -f "$${JSON_FILE}" || (echo "JSON_FILE does not exist"; exit 1)
	test -n "$${UPDATE_CHANNEL_VERSION}" || (echo "UPDATE_CHANNEL_VERSION is not defined"; exit 1)
	test -n "$${UPDATE_CHANNEL_BUILD}" || (echo "UPDATE_CHANNEL_BUILD is not defined"; exit 1)
	test -n "$${UPDATE_CHANNEL_DATE}" || (echo "UPDATE_CHANNEL_DATE is not defined"; exit 1)
	JQ_QUERY=".$${UPDATE_CHANNEL}.version = \"$${UPDATE_CHANNEL_VERSION}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.build = \"$${UPDATE_CHANNEL_BUILD}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.date = \"$${UPDATE_CHANNEL_DATE}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.docker.elasticsearch = \"$(ELASTICSEARCH_DOKER_IMAGE):$${UPDATE_CHANNEL_BUILD}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.docker.kibana = \"$(KIBANA_DOCKER_IMAGE):$${UPDATE_CHANNEL_BUILD}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.docker.fleet = \"$(ELASTIC_AGENT_DOCKER_IMAGE):$${UPDATE_CHANNEL_BUILD}\"" \
	$(MAKE) update-json
	JQ_QUERY=".$${UPDATE_CHANNEL}.docker.profiling = \"$(PROFILING_DOCKER_IMAGE):$${UPDATE_CHANNEL_BUILD}\"" \
	$(MAKE) update-json

## @help:update-cluster-config:Update the FILE cluster config file with the given STACK_VERSION, ELASTICSEARCH_IMAGE, KIBANA_IMAGE, INTEGRATIONS_IMAGE, PROFILING_IMAGE, and UPDATE_AT.
update-cluster-config:
	test -n "$${STACK_VERSION}" || (echo "STACK_VERSION is not defined"; exit 1)
	test -n "$${ELASTICSEARCH_IMAGE}" || (echo "ELASTICSEARCH_IMAGE is not defined"; exit 1)
	test -n "$${KIBANA_IMAGE}" || (echo "KIBANA_IMAGE is not defined"; exit 1)
	test -n "$${INTEGRATIONS_IMAGE}" || (echo "INTEGRATIONS_IMAGE is not defined"; exit 1)
	test -n "$${PROFILING_IMAGE}" || (echo "PROFILING_IMAGE is not defined"; exit 1)
	test -n "$${UPDATE_AT}" || (echo "UPDATE_AT is not defined"; exit 1)
	test -n "$${FILE}" || (echo "FILE is not defined"; exit 1)

	$(MAKE) set-key-value KEY=".stack.version" SET_VALUE="$${STACK_VERSION}"
	$(MAKE) set-key-value KEY=".stack.ess.elasticsearch.image" SET_VALUE="$${ELASTICSEARCH_IMAGE}"
	$(MAKE) set-key-value KEY=".stack.ess.kibana.image" SET_VALUE="$${KIBANA_IMAGE}"
	$(MAKE) set-key-value KEY=".stack.ess.integrations.image" SET_VALUE="$${INTEGRATIONS_IMAGE}"
	$(MAKE) set-key-value KEY=".stack.ess.profiling.image" SET_VALUE="$${PROFILING_IMAGE}"
	$(MAKE) set-key-value KEY=".updated_at" SET_VALUE="$${UPDATE_AT}"

.PHONY: update-elastic-stack-serverless
## @help:update-elastic-stack-serverless:Update Elastic Stack version on a Serverless cluster.(requires CLUSTER_NAME)
update-elastic-stack-serverless: venv
	source $(VENV_ACTIVATE); \
	export CLUSTER_CONFIG_FILE=$$($(CLUSTER_FINDER) ${CLUSTER_NAME} $(ENVIRONMENTS_DIR) "true"); \
	export UPDATE_CHANNEL=$$($(READ_YAML_KEY) "$(ENVIRONMENTS_DIR)/$${CLUSTER_CONFIG_FILE}" 'stack.update_channel'); \
	export STACK_MODE=$$($(READ_YAML_KEY) "$(ENVIRONMENTS_DIR)/$${CLUSTER_CONFIG_FILE}" 'stack.mode'); \
	echo "CLUSTER_CONFIG_FILE=$${CLUSTER_CONFIG_FILE}"; \
	echo "UPDATE_CHANNEL=$${UPDATE_CHANNEL}"; \
	echo "STACK_MODE=$${STACK_MODE}"; \
	updatecli apply $(UPDATECLI_OPTS) --config $(CURDIR)/update-elastic-stack-serverless.yml
