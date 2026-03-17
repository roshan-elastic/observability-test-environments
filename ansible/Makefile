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

ROLE ?= common
SCENARIO ?=
VENV ?= .venv
PYTHON ?= python3
ANSIBLE_OPTS ?=
MOLECULE_OPTS ?=
MOLECULE_TEST_OPTS ?=
STAGE ?= test
NO_CLEAN ?=
DISABLE_HERMIT ?=
BINDIR ?= $(CURDIR)/../.ci/bin
CONTAINER ?= false

export CLOUDSDK_PYTHON = $(PYTHON)
export PATH := $(CURDIR)/$(VENV)/bin:$(PATH)
export ANSIBLE_CONFIG := $(CURDIR)/ansible.cfg
export USE_GKE_GCLOUD_AUTH_PLUGIN := True
export PATH := ${PATH}:${HOME}/bin

ifeq ($(CONTAINER),false)
export OBJC_DISABLE_INITIALIZE_FORK_SAFETY := YES

ifeq ($(origin GCP_AUTH_KIND),undefined)
	export GCP_PROJECT := $(shell gcloud config get-value project)
	export GCP_AUTH_KIND := accesstoken
	export GCP_ACCESS_TOKEN := $(shell gcloud auth print-access-token)
	export CLOUDSDK_AUTH_ACCESS_TOKEN := $(GCP_ACCESS_TOKEN)
  export GOOGLE_OAUTH_ACCESS_TOKEN := $(GCP_ACCESS_TOKEN)
endif

export CLOUDSDK_CONFIG := ${HOME}/.config/gcloud
export GCP_ACCOUNT ?= $(shell gcloud config get-value account)

SERVICE_ACCOUNT_EMAIL ?= oblt-clusters@elastic-observability.iam.gserviceaccount.com
SERVICE_ACOUNT_PRIVATE_KEY ?= $(CURDIR)/private_key.pem
SERVICE_ACOUNT_PUBLIC_KEY ?= $(CURDIR)/public_key.pem

ifeq ($(origin DOCKER_USERNAME),undefined)
	DOCKER_GCLOUD_SECRET ?= elastic-oci-registry
	export DOCKER_USERNAME ?= $(shell gcloud secrets versions access latest --secret="$(DOCKER_GCLOUD_SECRET)" --project="$(GCP_PROJECT)"|yq -r .username)
	export DOCKER_PASSWORD ?= $(shell gcloud secrets versions access latest --secret="$(DOCKER_GCLOUD_SECRET)" --project="$(GCP_PROJECT)"|yq -r .password)
	export DOCKER_NAMESPACE ?= observability-ci
	export DOCKER_REGISTRY ?= docker.elastic.co
endif

endif

TEST_CLUSTER_NAME = test-creation-$(shell date +%s)
SCRIPTS_DIR = $(CURDIR)/../.ci/scripts

SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

.PHONY: help
help:
	@echo "Environment variables:"
	@echo ""
	@echo "ANSIBLE_OPTS additional flags for Ansible"
	@echo "MOLECULE_OPTS additional flags for Molecule"
	@echo "MOLECULE_TEST_OPTS additional flags for Molecule test"
	@echo "ANSIBLE_VERBOSITY=3 set the verbosity level for Ansible in molecule tests"
	@echo "STAGE=test|converge|destroy|clean|verify the stage to use in molecule (default: test)"
	@echo "NO_CLEAN=true do not clean the build folder"
	@echo "DESTROY_CLUSTER=false do not destroy deployed clusters"
	@echo "CLUSTER_CONFIG_FILE=/path/to/file.yml change the cluster config file"
	@echo "ELASTIC_STACK_VERSION=7.8.0 this force to use the stack version you set"
	@echo ""
	@echo "Main targets:"
	@echo ""
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

## @help:create-role:Creates a role using molecule. It's mandatory to add ROLE=name to make invocation
.PHONY: create-role
create-role: install-ansible
	@$(call check_defined, ROLE, ROLE must be defined. )
	source $(VENV)/bin/activate;\
	cd ansible/ansible_collections/oblt/framework/roles;\
	ansible-galaxy init $(ROLE);\
	cd $(ROLE);\
	molecule init scenario;

## @help:install-ansible:Install ansible and the dependencies needed in a virtual environment.
.PHONY: install-ansible
install-ansible:
ifeq ($(CONTAINER),false)
	test -d $(VENV) || $(PYTHON) -m venv $(VENV);\
	source $(VENV)/bin/activate;\
	$(PYTHON) --version; \
	$(PYTHON) -m pip install -U pip;\
	$(PYTHON) -m pip install --quiet -r requirements.txt --no-compile; \
	ansible-galaxy collection install --upgrade -r $(CURDIR)/requirements.yaml;
endif
	source $(VENV)/bin/activate; \
	ansible-galaxy collection list

.PHONY: install-hermit-tools
## @help:install-hermit-tools:Install the hermit tools.
install-hermit-tools:
	curl -fsSL https://github.com/cashapp/hermit/releases/download/stable/install.sh | /bin/bash
	source $(BINDIR)/activate-hermit \
  && gcloud --version \
  && gh --version \
  && helm version --short \
  && kubectl version --client \
  && terraform version \
  && jq --version \
  && yq --version \
  && k9s version \
  && kind --version

## @help:check-cluster-config:Check if the CLUSTER_CONFIG_FILE is a valid cluster config file.
check-cluster-config:
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	CLUSTER_NAME=$$( cat $(CLUSTER_CONFIG_FILE)| python -c "import sys, yaml; print(yaml.safe_load(sys.stdin).get('cluster_name', ''))" || ""); \
	if [ -z "$${CLUSTER_NAME}" ]; then \
		echo "CLUSTER_NAME is not defined in $(CLUSTER_CONFIG_FILE), it is not a Cluster config file"; \
		exit 1; \
	fi;

## @help:create-cluster:Read the deploy-config.xml file and deploy the test environment described in it.
.PHONY: create-cluster
create-cluster: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) create-cluster.yml

## @help:update-cluster:Read the deploy-config.xml file and deploy/update the Elastic stack and auxiliary apps basic described in it. it requires a k8s cluster deployed before.
.PHONY: update-cluster
update-cluster: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) update-cluster.yml

fake-cluster: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) fake-cluster.yml

## @help:update-SUFFIX:Execute the tasks-SUFFIX.yml in a playbook (deprecated use tasks-%).
.PHONY: update-%
update-%: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	ansible --version;\
	TASKS="$*" ansible-playbook $(ANSIBLE_OPTS) tasks-runner.yml

## @help:tasks-SUFFIX:Execute the tasks-SUFFIX.yml in a playbook.
.PHONY: tasks-%
tasks-%: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	ansible --version;\
	TASKS="$*" ansible-playbook $(ANSIBLE_OPTS) tasks-runner.yml

## @help:tasks-list:List the tasks available.
.PHONY: tasks-list
tasks-list:
	for i in $$(ls $(CURDIR)/tasks-*); \
	do \
		TASK=$$(basename $$i ".yml");\
		echo "$${TASK}";\
	done

## @help:destroy-cluster:Read the deploy-config.xml file and destroy the enviroent described in it, it also delete the secrets related.
.PHONY: destroy-cluster
destroy-cluster: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) destroy-cluster.yml

## @help:rotate-ess-api-key:Rotate the API key for an ESS service account. https://www.elastic.co/guide/en/cloud/current/Authentication.html#get-api-keys
.PHONY: rotate-ess-api-key
rotate-ess-api-key: clean install-ansible
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate;\
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) rotate-ess-api-key.yml

## @help:rotate-ess-api-key-pro:Rotate the API key for an ESS service account.
.PHONY: rotate-ess-api-key-pro
rotate-ess-api-key-pro:
	$(MAKE) rotate-ess-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-pro'"

## @help:rotate-ess-api-key-qa:Rotate the API key for an ESS service account.
.PHONY: rotate-ess-api-key-qa
rotate-ess-api-key-qa:
	$(MAKE) rotate-ess-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-qa'"

## @help:rotate-ess-api-key-staging:Rotate the API key for an ESS service account.
.PHONY: rotate-ess-api-key-staging
rotate-ess-api-key-staging:
	$(MAKE) rotate-ess-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-staging'"

## @help:rotate-ess-console-api-key:Rotate the Console API key for an ESS service account.
rotate-ess-console-api-key: clean install-ansible
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate;\
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) rotate-ess-console-api-key.yml

## @help:rotate-ess-console-api-key-qa:Rotate the Console API key for an ESS service account.
rotate-ess-console-api-key-pro:
	$(MAKE) rotate-ess-console-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-pro'"

## @help:rotate-ess-console-api-key-qa:Rotate the Console API key for an ESS service account.
rotate-ess-console-api-key-qa:
	$(MAKE) rotate-ess-console-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-qa'"

## @help:rotate-ess-console-api-key-staging:Rotate the Console API key for an ESS service account.
rotate-ess-console-api-key-staging:
	$(MAKE) rotate-ess-console-api-key ANSIBLE_OPTS="$(ANSIBLE_OPTS) --extra-vars='api_key_secret=elastic-cloud-observability-team-staging'"

## @help:restore-snapshot:Restore snapshot to cluster, needs SOURCE_REPOSITORY env var.
.PHONY: restore-snapshot
restore-snapshot: clean install-ansible
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit; \
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version; \
	ansible-playbook $(ANSIBLE_OPTS) restore-snapshot.yml

## @help:unit-test ROLE=role_name [SCENARIO=""] [MOLECULE_OPTS=""] [MOLECULE_TEST_OPTS=""]:Executes Unit tests for Ansible roles using molecule
.PHONY: unit-test
unit-test: install-ansible
	@$(call check_defined, ROLE, ROLE must be defined. )
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate;\
	echo 'ROLE=$(ROLE)';\
	echo 'SCENARIO=$(SCENARIO)';\
	echo 'STAGE=$(STAGE)';\
	echo $(PATH);\
	cd ./ansible_collections/oblt/framework/roles/$(ROLE);\
	if [ -z "$${NO_CLEAN}" ]; then \
		rm -fr molecule/default/build molecule/default/test-results | true;\
		rm -fr molecule/$(SCENARIO)/build molecule/$(SCENARIO)/test-results;\
	fi; \
	if [ -z "$(SCENARIO)" ]; then \
		molecule $(MOLECULE_OPTS) $(STAGE) --all $(MOLECULE_TEST_OPTS);\
	else \
		molecule $(MOLECULE_OPTS) $(STAGE) --scenario-name $(SCENARIO) $(MOLECULE_TEST_OPTS);\
	fi;

## @help:create-daily:Create a daily ephemeral cluster, it destroys the previous cluster (deprecated).
create-daily:
	@$(MAKE) recreate-cluster

## @help:recreate-cluster:Destroys the clusters and creates it again.
recreate-cluster:
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	@$(MAKE) destroy-cluster || true
	$(SCRIPTS_DIR)/delete-ec-test-env_by_exact_name.sh $(TEST_CLUSTER_NAME) "true" || true;
	@$(MAKE) create-cluster

## @help:test-creation:Test the creation of a cluster, it disables kubernetes and uses a temporary name.
.PHONY: test-creation
test-creation: ANSIBLE_OPTS="--extra-vars='cluster_name=${TEST_CLUSTER_NAME} k8s_enabled=false ephemeral_cluster=true golden_cluster=false features=[]'"
test-creation:
	@$(call check_defined, CLUSTER_CONFIG_FILE, CLUSTER_CONFIG_FILE must be defined. )
	@echo "Test the creation of a cluster $(ANSIBLE_OPTS)"
	export ANSIBLE_OPTS=$(ANSIBLE_OPTS); \
		$(MAKE) create-cluster || touch failed; \
		$(MAKE) destroy-cluster || true; \
		$(SCRIPTS_DIR)/delete-ec-test-env_by_name.sh $(TEST_CLUSTER_NAME) "true" || true; \
		if [ -f failed ]; then \
			rm failed; \
			exit 1; \
		fi

## @help:effective-config:Print the effective configuration of the oblt cluster.
.PHONY:
effective-config: tasks-effective-config

## @help:lint:Check the syntax of the Ansible playbooks
.PHONY: lint
lint: install-ansible
	export CLUSTER_CONFIG_FILE=$(CURDIR)/config-cluster.yml; \
	ansible-playbook create-cluster.yml --syntax-check; \
	ansible-playbook destroy-cluster.yml --syntax-check; \
	ansible-playbook update-cluster.yml --syntax-check; \

## @help:clean:Deletes the build folder.
.PHONY: clean
clean:
	if [ -z "$${NO_CLEAN}" ] && [ -d $(CURDIR)/build ]; then \
		chmod -R u+w $(CURDIR)/build; \
		rm -fr $(CURDIR)/build; \
	fi

## @help:clean-env:Deletes the virtualenv folder.
.PHONY: clean-env
clean-env:
	if [ -z "$${NO_CLEAN}" ]; then \
		rm -fr $(VENV); \
	fi

## @help:readme-generator-for-helm:Install the readme-generator-for-helm.
readme-generator-for-helm:
	git clone https://github.com/bitnami-labs/readme-generator-for-helm; \
	cd readme-generator-for-helm; \
	git checkout 2.4.3 -b 2.4.3;
	npm install ./readme-generator-for-helm; \

## @help:readme:Generate the readme and OpenAPI compliant JSON schema using the values.yml. VALUES_FILE=values.yml README=README.md SCHEMA=values.schema.json
.PHONY: readme
readme: readme-generator-for-helm
	$(CURDIR)/node_modules/.bin/readme-generator --values $(VALUES_FILE) --readme $(README) --schema $(SCHEMA)

.PHONY: generate-x509-cert
## @help:generate-x509-cert:Generate a selfsigned x509 certificate and store it in vault.
generate-x509-cert:
	openssl req -x509 -newkey rsa:2048 \
	-keyout $(CURDIR)/key.pem \
	-out $(CURDIR)/cert.crt \
	-days 365 \
	-subj '/CN=*' -nodes
	echo "{\"crt\":\"$$(cat $(CURDIR)/cert.crt)\",\"key\":\"$$(cat $(CURDIR)/key.pem)\",\"description\":\"General purpose x509 selfsigned certificate\"}" \
	| yq -p json | gcloud secrets versions add oblt-clusters_x509-selfsigned --data-file=-
	rm $(CURDIR)/key.pem $(CURDIR)/cert.crt

.PHONY: run-on-all-clusters
## @help:run-on-all-clusters:Run a task on all clusters. it requires the env var TASK=task-name, and optionally FILTER=filter
run-on-all-clusters:
	@$(call check_defined, TASK, TASK must be defined. )
	if [ -z "$(FILTER)" ]; then \
		oblt-cli cluster list --all --disable-banner --output-file "$(CURDIR)/clusters.json"; \
	else \
		oblt-cli cluster list --all --disable-banner --output-file "$(CURDIR)/clusters.json" --filter $(FILTER); \
	fi; \
	for cluster in $$(jq -r '.items[].clusterConfigPath' "$(CURDIR)/clusters.json"); \
	do \
		echo "Running $(TASK) on $${cluster}";\
		export CLUSTER_CONFIG_FILE=$${cluster};\
		$(MAKE) "$(TASK)";\
	done

.PHONY: gcp-authenticate
## @help:gcp-authenticate:Authenticate with GCP and create a key for the service account.
gcp-authenticate:
	gcloud auth application-default login
	gcloud iam service-accounts keys create $(GCP_SERVICE_ACCOUNT_FILE) --iam-account $(SERVICE_ACCOUNT_EMAIL)

	# Create temporal keys version it needs some review
	# openssl req -x509 -nodes -newkey rsa:4096 -days 1 -subj '/CN=*' -keyout $(SERVICE_ACOUNT_PRIVATE_KEY) -out $(SERVICE_ACOUNT_PUBLIC_KEY)
	# gcloud iam service-accounts keys upload $(SERVICE_ACOUNT_PUBLIC_KEY) --iam-account $(SERVICE_ACCOUNT_EMAIL)
	# touch $(GCP_SERVICE_ACCOUNT_FILE)
	# chmod 0600 $(GCP_SERVICE_ACCOUNT_FILE)
	# PRIVATE_KEY=$$(echo -ne $$(cat $(SERVICE_ACOUNT_PRIVATE_KEY))); \
	# PROJECT_ID=$$(gcloud config get-value project); \
	# CLIENT_EMAIL=$(SERVICE_ACCOUNT_EMAIL); \
	# CLIENT_ID=$$(gcloud iam service-accounts describe $(SERVICE_ACCOUNT_EMAIL) --format='value(uniqueId)'); \
	# jq -n \
	# 	--arg PRIVATE_KEY "$$PRIVATE_KEY" \
	# 	--arg PROJECT_ID $$PROJECT_ID \
	# 	--arg CLIENT_EMAIL $$CLIENT_EMAIL \
	# 	--arg CLIENT_ID $$CLIENT_ID \
	# 	--arg PRIVATE_KEY_ID $$PRIVATE_KEY_ID \
	# 	'{"type": "service_account","project_id": $$PROJECT_ID,"private_key_id": $$PRIVATE_KEY_ID,"private_key": $$PRIVATE_KEY,"client_email": $$CLIENT_EMAIL,"client_id": $$CLIENT_ID,"auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri": "https://oauth2.googleapis.com/token","auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url": @uri "https://www.googleapis.com/robot/v1/metadata/x509/\($$CLIENT_EMAIL)"}' > $(GCP_SERVICE_ACCOUNT_FILE)

.PHONY: gcp-delete-keys
## @help:gcp-delete-keys:Delete the keys for the service account.
gcp-delete-keys:
	for key in $$(gcloud iam service-accounts keys list  --format='value(KEY_ID)' --iam-account $(SERVICE_ACCOUNT_EMAIL) ); \
	do \
		gcloud iam service-accounts keys delete $$key --iam-account $(SERVICE_ACCOUNT_EMAIL) --quiet; \
	done

## @help:test-template:Performs a basic test on a template
.PHONY: test-template
test-template: clean install-ansible
	@$(call check_defined, TEMPLATE_FILE, TEMPLATE_FILE must be defined. )
	@$(call check_defined, OUTPUT_FILE, OUTPUT_FILE must be defined. )
	@$(call check_defined, PARAMETERS, PARAMETERS must be defined. )
	oblt-cli cluster create custom \
		--git-http-mode \
		--username "nobody" \
		--slack-channel "#no-channel" \
		--disable-banner \
		--cluster-name-prefix "test-template" \
		--template-file "$${TEMPLATE_FILE}" \
		--output-file "$${OUTPUT_FILE}" \
		--dry-run \
		--parameters "$${PARAMETERS}"
	CLUSTER_CONFIG_FILE=$$(jq -r '.ClusterConfigFile' "$${OUTPUT_FILE}"); \
	export CLUSTER_CONFIG_FILE; \
	[ -n "$(DISABLE_HERMIT)" ] || source $(BINDIR)/activate-hermit;\
	source $(VENV)/bin/activate; \
	$(MAKE) check-cluster-config; \
	ansible --version;\
	ansible-playbook $(ANSIBLE_OPTS) test-common.yml

.PHONY: collection-build
## @help:collection-build:Build the collection
collection-build: clean install-ansible
	mkdir -p $(CURDIR)/build
	source $(VENV)/bin/activate;\
		ansible --version;\
		ansible-galaxy collection build $(ANSIBLE_OPTS) $(CURDIR)/ansible_collections/oblt/framework --output-path $(CURDIR)/build

.PHONY: collection-install
## @help:collection-install:Build and Install the collection in the current virtual environment
collection-install: collection-build
	source $(VENV)/bin/activate;\
		ansible-galaxy collection install $(CURDIR)/build/oblt-framework-*.tar.gz

.PHONY: collection-publish
## @help:collection-publish:Build and Publish the collection to the Galaxy
collection-publish: collection-build
	source $(VENV)/bin/activate;\
		ansible-galaxy collection publish $(CURDIR)/build/oblt-framework-*.tar.gz
