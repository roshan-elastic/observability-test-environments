SHELL = /bin/bash
MAKEFLAGS += --silent --no-print-directory
.SHELLFLAGS = -ec
.SILENT:

.ONESHELL:

.PHONY: help
help:
	@grep '^## @help' Makefile|cut -d ":" -f 2-3|( (sort|column -s ":" -t) || (sort|tr ":" "\t") || (tr ":" "\t"))

## @help:update-opbeans:Build deps, README, and schema files for all Opbeans.
.PHONY: update-opbeans
update-opbeans:
	$(MAKE) -C $(CURDIR)/opbeans_main readme deps template clean
	set -e; \
	for opbean in dotnet frontend go java node php python ruby go_standalone; do \
		$(MAKE) -C $(CURDIR)/opbeans_$$opbean readme deps template clean; \
	done
