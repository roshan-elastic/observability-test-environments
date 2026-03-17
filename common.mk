
VENV ?= .venv
PYTHON ?= python3
PIP ?= pip3
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VALUES_FILE ?= $(CURDIR)/values.yaml
README_FILE ?= $(CURDIR)/README.md
SCHEMA_FILE ?= $(CURDIR)/values.schema.json
CONFIG_FILE ?= $(ROOT_DIR)/.ci/schema-config.json

OS = $(shell uname)

export PATH := $(CURDIR)/$(VENV)/bin:$(PATH)

ifeq ($(OS),Darwin)
	SED ?= sed -i ""
else
	SED ?= sed -i
endif

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

## @help:virtualenv:Create a Python virtual environment.
.PHONY: virtualenv
virtualenv:
	$(PYTHON) --version
	test -d $(VENV) || virtualenv -q --python=$(PYTHON) $(VENV);\
	source $(VENV)/bin/activate;\
	$(PIP) install -r requirements.txt;
	go get -u github.com/norwoodj/helm-docs/cmd/helm-docs@v1.6.0;

## @help:lint-python:Lint Python scripts.
.PHONY: lint-python
lint-python: virtualenv
	source $(VENV)/bin/activate;\
	black --diff --check --exclude='$(VENV)' $(CURDIR)

## @help:pytest:Run python tests.
.PHONY: pytest
pytest: virtualenv
	source $(VENV)/bin/activate;\
	pytest -sv --color=yes

## @help:clean:Clean the temporal files.
.PHONY: clean
clean:
	rm -fr $(VENV)
	rm -fr $(CURDIR)/.pytest_cache $(CURDIR)/node_modules $(CURDIR)/readme-generator-for-helm $(CURDIR)/package-lock.json $(CURDIR)/requirements.lock

## @help:readme-generator-for-helm:Install the readme-generator-for-helm.
readme-generator-for-helm:
	git clone https://github.com/bitnami-labs/readme-generator-for-helm; \
	cd readme-generator-for-helm; \
	git checkout 2.6.1 -b 2.6.1;
	npm install ./readme-generator-for-helm; \

## @help:readme:Generate the readme and OpenAPI compliant JSON schema using the values.yml.
.PHONY: readme
readme: readme-generator-for-helm
	@$(call check_defined, VALUES_FILE, VALUES_FILE must be defined. )
	$(CURDIR)/node_modules/.bin/readme-generator --values $(VALUES_FILE) --readme $(README_FILE) --schema $(SCHEMA_FILE) --config $(CONFIG_FILE)
