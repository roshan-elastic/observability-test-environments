PROJECT_NAME := apm-synthtrace
BRANCH ?= main
NODE_VERSION := $(shell curl -fsSL https://raw.githubusercontent.com/elastic/kibana/$(BRANCH)/.nvmrc)
BRANCH_VERSION := $(shell curl -fsSL https://raw.githubusercontent.com/elastic/kibana/$(BRANCH)/package.json | jq -r '.version')
DOCKER_OPTS := --build-arg NODE_VERSION=$(NODE_VERSION) --build-arg BRANCH=$(BRANCH)

ifeq ($(BRANCH),main)
	VERSION=${BRANCH_VERSION}-SNAPSHOT
else
	VERSION=${BRANCH_VERSION}
endif

include ../Makefile
