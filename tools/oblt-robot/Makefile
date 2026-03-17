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

REPO = elastic/observability-test-environments
REPO_VERSION ?= $$(cat $(CURDIR)/../../.ci/.version)
COMMIT ?= $$(git rev-parse --short=8 HEAD)
MAJOR = $$(echo $(REPO_VERSION) | cut -d. -f1)
MINOR = $$(echo $(REPO_VERSION) | cut -d. -f2)
PATCH = $$(echo $(REPO_VERSION) | cut -d. -f3)
REPO_RELEASE ?= $(MAJOR).$(MINOR).$(PATCH)
BUILD_DATE = $$(date -u +"%Y-%m-%dT%H:%M:%SZ")

DOCKER_REPO ?= docker.elastic.co
DOCKER_NAMESPACE ?= observability-ci
DOCKER_IMAGE ?= $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/oblt-robot
DOCKER_CMD ?= docker
GO_VERSION ?= $(shell cat ../../.go-version)

ROBOT_OPTS ?= --verbose --dry-run

GCP_PROJECT ?= elastic-observability
GCSM_CLUSTER_STATE_SECRET ?= oblt-clusters_observability-ci_cluster-state

VAULT_SECRET ?= secret/k8s/elastic-apps/apm/oblt-robot
VAULT_ELASTICSEARCH_SECRET ?= secret/k8s/elastic-apps/apm/oblt-robot-elasticsearch
VAULT_APM_SECRET ?= secret/k8s/elastic-apps/apm/oblt-robot-apm-server
GCSM_CLUSTER_ENV_SECRET ?= oblt-clusters_observability-ci_env
GCSM_GCP_CREDENTIALS ?= secret/k8s/elastic-apps/apm/oblt-robot-gcp-credentials-secret
GCP_CREDENTIALS_FILE ?= $(CURDIR)/oblt-robot-credentials.json

export GOPRIVATE = github.com/elastic

ifeq ($(shell uname), Darwin)
## See https://github.com/golang/go/issues/49138
export MallocNanoZone=0
endif

SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

.PHONY: help
help:
	@echo "Main targets:"
	@echo ""
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

build-dir:
	mkdir -p build packages

## @help:create-env:Generate the .env file using the Vault credentials, it uses production.
create-env:
	echo "OBLT_ROBOT_VERSION=$(REPO_VERSION)" > .env
	echo "OBLT_ROBOT_BUILD=dev" >> .env
	echo "GITHUB_PASSWORD=$$(vault read -field="GITHUB_PASSWORD" $(VAULT_SECRET))" >> .env
	echo "GITHUB_USER=$$(vault read -field="GITHUB_USER" $(VAULT_SECRET))" >> .env
	echo "SLACK_APP_TOKEN=$$(vault read -field="SLACK_APP_TOKEN" $(VAULT_SECRET))" >> .env
	echo "SLACK_AUTH_TOKEN=$$(vault read -field="SLACK_AUTH_TOKEN" $(VAULT_SECRET))" >> .env
	echo "SLACK_CHANNEL_ID=$$(vault read -field="SLACK_CHANNEL_ID" $(VAULT_SECRET))" >> .env
	echo "OBLT_ROBOT_VERSION=$(REPO_VERSION)" >> .env
	echo "OBLT_ROBOT_BUILD=$(BUILD_DATE)" >> .env
	echo "GOOGLE_APPLICATION_CREDENTIALS=$(GCP_CREDENTIALS_FILE)" >> .env
	echo "ELASTIC_APM_SERVICE_NAME=oblt-robot" >> .env
	echo "GIT_AUTHOR_NAME=obltmachine" >> .env
	echo "GIT_AUTHOR_EMAIL=infra-root+obltmachine@elastic.co" >> .env
	gcloud secrets versions access latest --secret="$(GCSM_CLUSTER_ENV_SECRET)" >> .env
	vault read -field=oblt-robot-credentials.json "$(GCSM_GCP_CREDENTIALS)" > "$(GCP_CREDENTIALS_FILE)"

## @help:create-env-test:Generate the .env file using the Vault credentials, the environment is the test one.
create-env-test:
	$(MAKE) create-env VAULT_SECRET=secret/observability-team/ci/oblt-robot/env-test

## @help:update-deps:Update the dependencies.
update-deps:
	go get -u
	go mod tidy

## @help:tidy:Update the go.mod
tidy:
	go mod tidy

## @help:format:Format the code.
format:
	go fmt ./...

## @help:generate:Generate the code to include static files into the binary
generate:
	@go generate ./...
	@echo "[OK] Files added to embed box!"

## @help:lint:Check the project for linting issues.
.PHONY: lint
lint:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@2023.1.3 ./...

## @help:build:Build the CLI binary.
build: tidy format generate
	go build -ldflags="-X main.Version=$(REPO_VERSION) -X main.Build=$(BUILD_DATE) -X main.Commit=$(COMMIT)" -o build/oblt-robot

## @help:build-docker:Build the Docker image.
build-docker: generate
	docker build -t $(DOCKER_IMAGE):$(REPO_VERSION) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VERSION=$(REPO_VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-f Dockerfile ..
	docker tag $(DOCKER_IMAGE):$(REPO_VERSION) $(DOCKER_IMAGE):latest

## @help:run-docker-test:Build and Run the Docker container in the test environment.
run-docker-test: build-docker create-env-test
	docker run \
		--env-file .env \
		-it \
		--rm \
		-v $(GCP_CREDENTIALS_FILE):$(GCP_CREDENTIALS_FILE) \
		$(DOCKER_IMAGE):$(REPO_VERSION) $(ROBOT_OPTS)

## @help:clean:Clean the build folder.
.PHONY: clean
clean:
	rm -fr build

## @help:push-docker:Push the Docker image.
push-docker:
	docker push $(DOCKER_IMAGE):$(REPO_VERSION)
	docker push $(DOCKER_IMAGE):latest

## @help:test:Run all the tests.
.PHONY: test
test: build-dir
	go run gotest.tools/gotestsum --format testname --junitfile build/junit-report.xml -- ./... -race -coverprofile=build/coverage.txt
	go run github.com/boumenot/gocover-cobertura < build/coverage.txt > build/coverage.xml

see-coverage:
	go tool cover -html=build/coverage.txt

## @help:release:Upload the packages generated to the latest release on the repository.
.PHONY: release
release: generate
	go run github.com/goreleaser/goreleaser release --clean
	$(MAKE) push-docker

## @help:snapshot:Generate local binaries based on the current code.
.PHONY: snapshot
snapshot: generate
	go run github.com/goreleaser/goreleaser release --skip=publish --snapshot --clean

## @help:deploy:Deploy the application on Elastic Apps
.PHONY: delete
deploy:
	kubectl apply -f manifest.yml

## @help:undeploy:Undeploy the application on Elastic Apps
.PHONY: undeploy
undeploy:
	kubectl delete -f manifest.yml

## @help:restart:Restart the deployment
.PHONY: restart
restart:
	kubectl rollout restart deployment/oblt-robot-deployment

## @help:configure-k8s:Configure the k8s cluster where to deploy the app
.PHONY: configure-k8s
configure-k8s:
	kubectl config set-context gke_elastic-apps-163815_us-central1-a_apps --namespace=apm
	kubectl config use-context gke_elastic-apps-163815_us-central1-a_apps

## @help:k8s-status:Help with debugging what's the Elastic App doing.
.PHONY: k8s-status
k8s-status:
	@echo "See current pods for the slackbot"
	kubectl get pods -l app=oblt-robot-slackbot
	@read -n 1 -p "See the current slackbot deployment pod"
	kubectl get deployments oblt-robot-deployment
	@read -n 1 -p "See the deployment describe"
	kubectl describe deployments oblt-robot-deployment
	@read -n 1 -p "See the logs"
	kubectl logs deployment/oblt-robot-deployment

## @help:elastic-apps-env:Configure the k8s elastic apps cluster
.PHONY: elastic-apps-env
elastic-apps-env:
	gcloud container clusters get-credentials apps --zone us-central1-a --project elastic-apps-163815
	kubectl config set-context $$(kubectl config current-context) --namespace=apm

update-vault-secrets:
	gcloud secrets versions access latest --secret=$(GCSM_CLUSTER_STATE_SECRET) --project=$(GCP_PROJECT)|sed 's/: /=/g' > .env;
	. .env; \
	vault write "$(VAULT_ELASTICSEARCH_SECRET)" \
		ELASTICSEARCH_PASSWORD="$${elasticsearch_password}" \
		ELASTICSEARCH_URL="$${elasticsearch_url}" \
		ELASTICSEARCH_USERNAME="$${elasticsearch_username}" \
		description="Secret consumed by the oblt-robot slack app" \
		ticket=https://github.com/elastic/observability-test-environments/tree/main/tools/oblt-robot; \
	vault write "$(VAULT_APM_SECRET)" \
		ELASTIC_APM_LOG_LEVEL=error \
		ELASTIC_APM_SERVER_URL="$${apm_url}" \
		ELASTIC_APM_SERVICE_NAME=oblt-robot \
		description="Secret consumed by the oblt-robot slack app" \
		ticket=https://github.com/elastic/observability-test-environments/tree/main/tools/oblt-robot

## @help:gcp-login:Login to GCP
gcp-login:
	gcloud auth login --update-adc
