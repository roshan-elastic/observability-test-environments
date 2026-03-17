FROM docker/buildx-bin:0.8.1 AS buildx
FROM docker:20.10.14-dind AS docker

RUN apk --no-cache add curl git jq
# Install Bbuildx and Docker compose
COPY --from=buildx /buildx /usr/libexec/docker/cli-plugins/docker-buildx
RUN curl -sSL https://github.com/docker/compose/releases/download/v2.2.3/docker-compose-linux-x86_64 -o /usr/libexec/docker/cli-plugins/docker-compose \
  && chmod +x /usr/libexec/docker/cli-plugins/docker-compose
RUN (docker version || true) && docker buildx version && docker compose version && docker buildx ls

RUN mkdir -p /usr/shared/data
COPY run-integration.sh /run-integration.sh

ENTRYPOINT ["/run-integration.sh"]
