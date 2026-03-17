FROM docker.elastic.co/package-registry/package-registry:v1.20.0

# Disable package validation (already done).
ENV EPR_DISABLE_PACKAGE_VALIDATION=true

ENV EPR_FEATURE_PROXY_MODE=true
ENV EPR_PROXY_TO=https://epr-snapshot.elastic.co/

COPY config.yml /package-registry/config.yml
RUN mkdir -p /usr/shared/data/packages
COPY development/ /usr/shared/data/packages/development
