# include ../../common.mk

PROJECT_NAME ?= dummy
VERSION ?= latest
DOCKER_IMAGE ?= docker.elastic.co/observability-ci/$(PROJECT_NAME):$(VERSION)
DOCKER_OPTS ?=
GITHUB_STEP_SUMMARY ?= /dev/null

## @help:docker-build:build Build docker image
.PHONY: docker-build
docker-build:
	@echo "Building docker image $(DOCKER_IMAGE)"
	@docker build -t $(DOCKER_IMAGE) $(DOCKER_OPTS) .

## @help:docker-push:push Push docker image
.PHONY: docker-push
docker-push:
	@echo "Pushing docker image $(DOCKER_IMAGE)" | tee -a $(GITHUB_STEP_SUMMARY)
	@docker push $(DOCKER_IMAGE)

## @help:docker-clean:clean Clean docker image
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning docker image $(DOCKER_IMAGE)"
	@docker rmi $(DOCKER_IMAGE)

## @help:docker-run:run Run docker image
.PHONY: docker-run
docker-run:
	@echo "Running docker image $(DOCKER_IMAGE)"
	@docker run --rm -it $(DOCKER_IMAGE)

docker-build-all:
	@echo "Building all docker images"
	for dir in integrations/package-registry integrations/dind workspace; \
	do \
		$(MAKE) -C $$dir docker-build; \
	done
