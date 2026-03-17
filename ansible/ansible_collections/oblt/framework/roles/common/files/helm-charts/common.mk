HELM_OPTS ?= --values values.yaml
HELM_TEST_NAME ?= my-helm-chart

GITHUB_STEP_SUMMARY ?= /dev/null
OCI_REGISTRY_URL ?= docker.elastic.co
OCI_NAMESPACE ?= observability-ci

## @help:deps:Update helm charts dependencies
.PHONY: deps
deps:
	helm dependency update
	helm dependency build

## @help:lint:Lint helm templates.
.PHONY: lint
lint:
	helm lint $(HELM_OPTS) $(CURDIR)

## @help:template:Render chart templates.
.PHONY: template
template: deps
	helm template $(HELM_TEST_NAME) $(HELM_OPTS) $(CURDIR)

## @help:helm-unittest:Run helm-unittest tests (https://github.com/helm-unittest/helm-unittest/blob/main/DOCUMENT.md).
.PHONY: helm-unittest
helm-unittest:
	docker run -ti --rm -v $(CURDIR):/apps helmunittest/helm-unittest .

## @help:test-all:Run all tests.
.PHONY: test-all
test-all: install-helm-unittest lint deps template helm-unittest pytest

## @help:update-initdb-scripts:Update the initdb scripts.
.PHONY: update-initdb-scripts
update-initdb-scripts:
	mkdir -p $(CURDIR)/config; \
	cd /tmp; \
  curl -sS -OL https://github.com/elastic/opbeans-node/archive/main.tar.gz ; \
  tar -xzf main.tar.gz; \
	cp opbeans-node-main/db/*.sql $(CURDIR)/config; \
	rm -fr main.tar.gz opbeans-node-main
	sed -i '.bak' -e '/^--.*$$/d' $(CURDIR)/config/*.sql
	sed -i '.bak' -e 's/\r$$//g' $(CURDIR)/config/*.sql
	sed -i '.bak' -e '/^$$/d'  $(CURDIR)/config/*.sql

## @help::helm-package:Package and push the helm chart to the OCI registry.
.PHONY: helm-package
helm-package:
	chart_dir=$(CURDIR); \
	chart_name=$$(basename "$(CURDIR)"); \
	chart_version=$$(yq '.version' "$(CURDIR)/Chart.yaml"); \
	tmp_dir=$$(mktemp -d); \
	echo "::group::Releasing $${chart_name} version $${chart_version}"; \
	helm package --dependency-update "$${chart_dir}" -d "$${tmp_dir}"; \
	helm push "$${tmp_dir}/"*.tgz "oci://$(OCI_REGISTRY_URL)/$(OCI_NAMESPACE)"; \
	echo "oci://$(OCI_REGISTRY_URL)/$(OCI_NAMESPACE)/$${chart_name}:$${chart_version}" >> "$(GITHUB_STEP_SUMMARY)"; \
	rm -fr "$${tmp_dir}"; \
	echo "::endgroup::"
