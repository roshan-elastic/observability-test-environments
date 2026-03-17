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
MAJOR = $$(echo $(REPO_VERSION) | cut -d. -f1)
MINOR = $$(echo $(REPO_VERSION) | cut -d. -f2)
PATCH = $$(echo $(REPO_VERSION) | cut -d. -f3)
REPO_RELEASE ?= $(MAJOR).$(MINOR).$(PATCH)
GO_VERSION ?= $(shell cat ../../.go-version)
GO_WORKDIR = "/go/src/github.com/elastic/observability-test-environments/tools/oblt-cli"
UID = $$(id -u)
GUID = $$(id -g)
DOCKER_REPO ?= docker.elastic.co
DOCKER_NAMESPACE ?= observability-ci
DOCKER_IMAGE ?= $(DOCKER_REPO)/$(DOCKER_NAMESPACE)/oblt-cli
DOCKER_CMD ?= docker
SSH_KEY ?= $(HOME)/.ssh/id_rsa
TEST_PAKCAGE ?= ./...
OAPI_CODEGEN=github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0

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
	@echo "Environment variables:"
	@echo ""
	@echo ""
	@echo "Main targets:"
	@echo ""
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

build-dir:
	mkdir -p build packages

## @help:lint:Check the project for linting issues.
.PHONY: lint
lint:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@2023.1.7 ./...

## @help:sync-version:Copy version file to the tool
sync-version:
	mkdir -p $(CURDIR)/static
	cp $(CURDIR)/../../.ci/.version $(CURDIR)/static/.version
	cp $(CURDIR)/../../docs/cluster-config-schema.json $(CURDIR)/static/cluster-config-schema.json

## @help:tidy:Update the go.mod
tidy:
	go mod tidy

## @help:format:Format the code.
format:
	go fmt ./...

## @help:generate:Generate the code to include static files into the binary
generate: sync-version
	@go generate ./...
	@echo "[OK] Files added to embed box!"

## @help:update-deps:Update the dependencies.
update-deps:
	go get -u
	go mod tidy

## @help:build:Build the CLI binary.
build: tidy format gen-cluster-config-schema generate
	go build -o build/oblt-cli

## @help:documentation:Generate the documentation for the CLI.
documentation: build
	$(CURDIR)/build/oblt-cli documentation --output-folder $(CURDIR)/../../docs/tools/oblt-cli/generated

## @help:test:Run all the tests.
.PHONY: test
test: build-dir
	go run gotest.tools/gotestsum --format testname --junitfile build/junit-report.xml -- $(TEST_PAKCAGE) -race -coverprofile=build/coverage.txt
	go run github.com/boumenot/gocover-cobertura < build/coverage.txt > build/coverage.xml

## @help:test-e2e:Run e2e tests. you can add FEATURE env var to run a subset of tests. SSH_KEY contains the SSH key to use for git.
test-e2e: clean-all
	GOOS=linux GOARCH=amd64 go build -o build/oblt-cli
	echo "git-user:x:$$(id -u):$$(id -g):Git User:/tmp:/bin/bash" > $(CURDIR)/fake_passwd
	docker run -i --rm \
		-v $(CURDIR):/app \
		-v ${SSH_KEY}:/tmp/.ssh/id_rsa \
		-v ${SSH_KEY}.pub:/tmp/.ssh/id_rsa.pub \
		-v $(CURDIR)/fake_passwd:/etc/passwd  \
		-e GIT_AUTHOR_NAME="$$(git config user.name)" \
	  -e GIT_AUTHOR_EMAIL="$$(git config user.email)" \
	  -e GIT_COMMITTER_NAME="$$(git config user.name)" \
	  -e GIT_COMMITTER_EMAIL="$$(git config user.email)" \
		-e GIT_SSH_COMMAND='ssh -i /tmp/.ssh/id_rsa -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \
		-w /app \
		-u $$(id -u):$$(id -g) \
		-e HOME=/tmp \
		-e BUNDLE_PATH=/tmp \
		-e GEM_HOME=/tmp \
		-e PATH="/app/build:/app/build:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
		ruby:3.0.2 /bin/bash -c "git --version && bundle config set --local path '.vendor' && bundle install && bundle exec cucumber --publish-quiet ${FEATURE}"

## @help:clean:Clean the build folder.
.PHONY: clean
clean:
	rm -fr build .cache dist

## @help:clean-all:Clean the build folder and packages folder.
.PHONY: clean-all
clean-all: clean
	rm -fr build packages Gemfile.lock .vendor .bundle fake_password

## @help:release:Upload the packages generated to the latest release on the repository.
.PHONY: release
release: generate
	go run github.com/goreleaser/goreleaser release --clean

## @help:snapshot:Generate local binaries based on the current code.
.PHONY: snapshot
snapshot: generate
	go run github.com/goreleaser/goreleaser release --skip=publish --snapshot --clean

## @help:gen-cluster-config-schema:Generate the cluster config schema object, marshalls, and validators
gen-cluster-config-schema: sync-version
	go run $(OAPI_CODEGEN) --config=$(CURDIR)/codegen.config $(CURDIR)/static/clusters-api.yaml
