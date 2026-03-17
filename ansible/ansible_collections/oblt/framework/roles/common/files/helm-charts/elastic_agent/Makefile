HELM_TEST_NAME=elastic-agent

include ../common.mk
include ../../common.mk

default: test

## @help:lint:Lint helm templates.
.PHONY: lint
lint:
	helm lint $(HELM_OPTS) --values values-daemonset.yaml ./
	helm lint $(HELM_OPTS) --values values-deployment.yaml ./
	helm lint $(HELM_OPTS) --values values-deployment-token.yaml ./
	helm lint $(HELM_OPTS) --values values-fleet-server.yaml ./
	helm lint $(HELM_OPTS) --values values-standalone.yaml ./

## @help:test-values:Render chart templates.
.PHONY: template
test-values:
	helm template $(HELM_OPTS) --values values-daemonset.yaml ./
	helm template $(HELM_OPTS) --values values-deployment.yaml ./
	helm template $(HELM_OPTS) --values values-deployment-token.yaml ./
	helm template $(HELM_OPTS) --values values-fleet-server.yaml ./
	helm template $(HELM_OPTS) --values values-standalone.yaml ./

## @help:test-all:Run all tests.
.PHONY: test-all
test-all: lint deps template pytest
