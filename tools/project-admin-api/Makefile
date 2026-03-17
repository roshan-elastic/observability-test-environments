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

GO_VERSION ?= $(shell cat ../../.go-version)

GO_SWAGGER_VERSION = v0.30.5
PROJECT_API_SPEC_URL = https://raw.githubusercontent.com/elastic/project-api/main/api/user-v1-spec.yml
PROJECT_API_SPEC_FILE = $(CURDIR)/user-v1-spec.yml
PROJECT_ADMIN_API_SPEC_URL = https://raw.githubusercontent.com/elastic/project-api/main/api/admin-v1-spec.yml
PROJECT_ADMIN_API_SPEC_FILE = $(CURDIR)/admin-v1-spec.yml
PKG_ADMIN_FOLDER = pkg/api/admin
EXT_MAPPING=./user-v1-spec.yml:github.com/elastic/observability-test-environments/tools/project-api/${PKG_FOLDER}
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

## @help:update-deps:Update the dependencies.
update-deps:
	go get -u
	go get -u ./tools
	go mod tidy

## @help:build:Build the CLI binary.
build:
	go mod tidy
	go fmt ./...
	go build -o build/oblt-project-api

## @help:clean:Clean the build folder.
.PHONY: clean
clean:
	rm -fr build .cache dist

## @help:clean-all:Clean the build folder and packages folder.
.PHONY: clean-all
clean-all: clean
	rm -fr build packages Gemfile.lock .vendor .bundle fake_password

## @help:generate-project-api:Generate the code for the project API.
.PHONY: generate-project-api
generate-project-api:
	[ -n "$${GITHUB_TOKEN}" ] || (echo "GITHUB_TOKEN is not set" && exit 1)
	curl -sSfL -o "${PROJECT_API_SPEC_FILE}" -H "Authorization: Bearer $${GITHUB_TOKEN}" "${PROJECT_API_SPEC_URL}"
	curl -sSfL -o "${PROJECT_ADMIN_API_SPEC_FILE}" -H "Authorization: Bearer $${GITHUB_TOKEN}" "${PROJECT_ADMIN_API_SPEC_URL}"
	mkdir -p "$(CURDIR)/${PKG_ADMIN_FOLDER}"
	go run $(OAPI_CODEGEN) --config=$(CURDIR)/codegen.config "${PROJECT_ADMIN_API_SPEC_FILE}"
