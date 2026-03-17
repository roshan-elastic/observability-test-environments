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

include ../common.mk

PYTHON ?= python3
PIP ?= pip3
VENV ?= .venv
VENV_ACTIVATE = $(VENV)/bin/activate
VERSION_FILE = $(CURDIR)/.version
VERSION_ANSIBLE_COLELLECTION_FILE = $(CURDIR)/../ansible/ansible_collections/oblt/framework/galaxy.yml
VERSION = $(shell cat $(VERSION_FILE))
VERSION_GH_ACTIONS_FILE = $(CURDIR)/../.github/actions/cluster-target/action.yml
NEW_BRANCH ?= release-$(VERSION)
COMMIT_MESSAGE ?= Bump version: $(VERSION)
REMOTE_URL ?= https://github.com/elastic/observability-test-environments.git
OBLT_CLI_DIR = $(CURDIR)/../tools/oblt-cli
OBLT_ROBOT_DIR = $(CURDIR)/../tools/oblt-robot
CLUSTER_FINDER = $(shell realpath $(CURDIR)/scripts/find_cluster_config_by_name.py)
ENVIRONMENTS_DIR = $(CURDIR)/../environments
DOCS_DIR = $(CURDIR)/../docs
ANSIBLE_DIR = $(CURDIR)/../ansible
ROLES_DIR = $(ANSIBLE_DIR)/ansible_collections/oblt/framework/roles

export REMOTE := tmp-$(shell date +'%Y%m%d-%H%M%S')
export PATH := $(PATH):$(CURDIR)/scripts:$(CURDIR)/bin
export GOPRIVATE = github.com/elastic

SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

.PHONY: help
help:
	@echo "Environment variables:"
	@echo ""
	@echo "Main targets:"
	@echo ""
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

## @help:venv:Prepare a virtualenv for Python.
.PHONY: venv
venv:
	$(PYTHON) --version;
	$(PYTHON) -m venv $(VENV);\
	source $(VENV_ACTIVATE);\
	$(PYTHON) -m pip install -U pip;\
	$(PIP) install --prefer-binary -r requirements.txt;

## @help:lint:Check the project for linting issues.
.PHONY: lint
lint: venv
	source $(VENV_ACTIVATE);\
	$(PYTHON) -m flake8 --config .flake8 .;

## @help:clean:Clean the project.
.PHONY: clean
clean:
	rm -fr $(VENV);

## @help:release-patch:Increase a patch version.
.PHONY: release-patch
release-patch: venv
	source $(VENV_ACTIVATE);\
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) patch $(VERSION_FILE)
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) patch "$(VERSION_ANSIBLE_COLELLECTION_FILE)"
	GIT_FILES='$(VERSION_FILE) $(OBLT_CLI_DIR) $(VERSION_ANSIBLE_COLELLECTION_FILE)' $(MAKE) commit-changes

## @help:release-minor:Increase a minor version.
.PHONY: release-minor
release-minor: venv
	source $(VENV_ACTIVATE);\
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) minor $(VERSION_FILE)
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) minor "$(VERSION_ANSIBLE_COLELLECTION_FILE)"
	GIT_FILES='$(VERSION_FILE) $(OBLT_CLI_DIR) $(VERSION_ANSIBLE_COLELLECTION_FILE)' $(MAKE) commit-changes

## @help:release-major:Increase a major version.
.PHONY: release-major
release-major: venv
	source $(VENV_ACTIVATE);\
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) major $(VERSION_FILE)
	bump2version --allow-dirty --no-commit --no-tag --current-version $(VERSION) major "$(VERSION_ANSIBLE_COLELLECTION_FILE)"
	GIT_FILES='$(VERSION_FILE) $(OBLT_CLI_DIR) $(VERSION_ANSIBLE_COLELLECTION_FILE)' $(MAKE) commit-changes

## @help:oblt-framework-bump:Bump the oblt-framework version.
oblt-framework-bump:
	yq e '.inputs.frameworkVersion.default = "$(VERSION)"' -i "$(VERSION_GH_ACTIONS_FILE)"
	GIT_FILES='$(VERSION_GH_ACTIONS_FILE)' $(MAKE) commit-changes

## @help:create-tag:Create a tag for the current version.
.PHONY: create-tag
create-tag:
	git fetch --force --tags
	git tag $(VERSION)
	git push --force --tags

## @help:release-packages:Release the packages for the latest tag.
.PHONY: release-packages
release-packages:
	$(MAKE) -C $(OBLT_CLI_DIR) release
	$(MAKE) -C $(OBLT_ROBOT_DIR) release

## @help:commit-changes:Commit changes.
.PHONY: commit-changes
commit-changes: generate-blobs
	git remote add ${REMOTE} "${REMOTE_URL}"
	git checkout -b $(NEW_BRANCH)
	git add $$(echo $${GIT_FILES})
	git diff --cached --exit-code || THERE_ARE_CHANGES="true"; \
	if [ "$${THERE_ARE_CHANGES}" == "true" ]; then \
		git commit -m "$(COMMIT_MESSAGE)"; \
		git push -u ${REMOTE} HEAD; \
		$(MAKE) create-pr; \
	else \
		echo "There are no changes to commit"; \
	fi
	git remote remove "${REMOTE}"

## @help:create-pr:DEPRECATED Create a PR.The message is the last commit message, the branch is the current branch. For extra options, use GH_EXTRA_ARGS.
.PHONY: create-pr
create-pr:
	gh auth status;
	COMMIT_MESSAGE="$(shell git log --format=%B -n 1|head -1)"; \
	BRANCH="$(shell git rev-parse --abbrev-ref HEAD)"; \
	echo "gh pr create --title "$${COMMIT_MESSAGE}" --body "$${COMMIT_MESSAGE}" --head "elastic:$${BRANCH}" --base "main" ${GH_EXTRA_ARGS} --repo "elastic/observability-test-environments""; \
	gh pr create --title "$${COMMIT_MESSAGE}" --body "$${COMMIT_MESSAGE}" --head "elastic:$${BRANCH}" --base "main" ${GH_EXTRA_ARGS} --repo "elastic/observability-test-environments"

## @help:generate-blobs:Generate blobs for the current version.
.PHONY: generate-blobs
generate-blobs:
	 $(MAKE) -C $(OBLT_CLI_DIR) generate

## @help:popeye-install:Install Popeye locally.
.PHONY: popeye-install
popeye-install:
	curl -sSL -o popeye.tgz https://github.com/derailed/popeye/releases/download/$(POPEYE_VERSION)/popeye_$(OS)_$(ARCH).tar.gz
	tar -xzf popeye.tgz
	mkdir -p $(BIN_DIR)
	mv popeye $(BIN_DIR)

## @help:popeye:Run the Popeye checks.
.PHONY: popeye
popeye: popeye-install
	popeye -f spinach.yml -o junit > popeye.junit

.PHONY: create-x509-cert
## @help:create-x509-cert:Create a x509 selfsigned certificate to use for TLS connections.
create-x509-cert:
	openssl req -x509 -newkey rsa:4096 -keyout k8s/key.pem -out k8s/cert.crt -days 365 -subj '/CN=*' -nodes

.PHONY: docs-publish
## @help:docs-publish:Publish the documentation.
docs-publish: venv
	source $(VENV_ACTIVATE); \
	pip install -r requirements-mkdocs.txt; \
	cd $(CURDIR)/..; \
  mkdocs gh-deploy --force

.PHONY: docs-serve
## @help:docs-serve:Serve the documentation locally.
docs-serve: venv
	source $(VENV_ACTIVATE); \
	pip install -r requirements-mkdocs.txt; \
	cd $(CURDIR)/..; \
	mkdocs serve

.PHONY: docs-cluster-config
## @help:docs-cluster-config:Generate the cluster configuration documentation.
docs-cluster-config:
	VALUES_FILE=$(DOCS_DIR)/config-cluster.yml \
	SCHEMA_FILE=$(DOCS_DIR)/cluster-config-schema.json \
	README_FILE=$(DOCS_DIR)/user-guide/cluster-config.md \
	$(MAKE) readme

.PHONY: docs-ansible
## @help:docs-ansible:Generate the Ansible roles documentation.
docs-ansible:
	set -e; \
	for role in $$(ls $(ROLES_DIR)); \
	do \
		export VALUES_FILE=$(ROLES_DIR)/$${role}/defaults/main.yml; \
		export SCHEMA_FILE=$(DOCS_DIR)/tools/oblt-framework/ansible/$${role}/schema.json; \
		export README_FILE=$(DOCS_DIR)/tools/oblt-framework/ansible/$${role}/README.md; \
		if [ -f "$${README_FILE}" ] && [ -f "$${VALUES_FILE}" ]; then \
			echo "****************************************************************************************************"; \
			echo "Generating documentation for $${role}"; \
			echo "****************************************************************************************************"; \
			$(MAKE) readme; \
		fi; \
	done

## @help:clean-ephemeral-clusters:Deletes the ephemeral clusters.
.PHONY: clean-ephemeral-clusters
clean-ephemeral-clusters:
	for item in "test-creation" "oblt-test" "smoke-upgrade-" "benchci-" "gatekeeper-", "apm-mis-ess", "apm-aws-lambda-smoke-testing-"; \
	do \
		$(CURDIR)/scripts/delete-ec-test-env_by_name.sh "$${item}" "true" "pro"; \
	done
