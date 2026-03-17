ARG GO_VERSION=1.21
############################
# STEP 1 build executable binary
############################
FROM golang:${GO_VERSION}-alpine as builder
ARG BUILD_DATE
ARG VERSION
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata curl && update-ca-certificates

# Create oblt-robot user
ENV HOME=/home/oblt-robot
ENV USER=oblt-robot
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/github.com/elastic/observability-test-environments/tools/oblt-robot

# copy oblt-cli for local replacement
COPY ./oblt-cli $GOPATH/src/github.com/elastic/observability-test-environments/tools/oblt-cli

# use modules
COPY ./oblt-robot .

# go mod tidy is necessary because otherwise the build will fail with 'go: updates to go.mod needed;'
RUN go mod tidy && \
    go mod download && \
    go mod verify

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags \"-static\" -X main.Version=$VERSION -X main.Build=$BUILD_DATE" -a \
    -o /go/bin/oblt-robot .

WORKDIR /tmp
RUN curl -sSLo ghlinux.tar.gz \
  https://github.com/cli/cli/releases/download/v2.10.1/gh_2.10.1_linux_amd64.tar.gz \
  && tar --strip-components=1 -xf ghlinux.tar.gz \
  && rm ghlinux.tar.gz

############################
# STEP 2 build a small image
############################
FROM alpine:3.15
ARG BUILD_DATE
ARG VERSION
ARG COMMIT
LABEL maintainer="observability-robots@elastic.co"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date="$BUILD_DATE"
LABEL org.label-schema.name="docker.elastic.co/observability-ci/oblt-robot"
LABEL org.label-schema.description="Slack bot for Observability Robots"
LABEL org.label-schema.version="$VERSION"
LABEL org.label-schema.url="https://elastic.co/"
LABEL org.label-schema.vcs-ref="$COMMIT"
LABEL org.label-schema.vcs-url="https://github.com/elastic/observability-test-environments"

RUN apk update && apk add --no-cache git openssh

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy default certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy user data
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /go/bin/oblt-robot /go/bin/oblt-robot
# Copy gh cli binary.
COPY --from=builder /tmp/bin/gh /usr/local/bin/gh

COPY oblt-robot/entrypoint.sh /entrypoint.sh

# Create oblt-robot user dir
ENV HOME=/home/oblt-robot
RUN mkdir -p /home/oblt-robot && chown oblt-robot:oblt-robot /home/oblt-robot
WORKDIR $HOME

# Use an unprivileged user.
USER oblt-robot:oblt-robot
# Setup GitHub authentication
# RUN gh auth setup-git
# Configure APM Agent
ENV ELASTIC_APM_LOG_LEVEL=debug
ENV ELASTIC_APM_SERVICE_NAME=oblt-robot
ENV GIT_AUTHOR_NAME=obltmachine
ENV GIT_AUTHOR_EMAIL=infra-root+obltmachine@elastic.co
ENV GIT_COMMITTER_NAME=obltmachine
ENV GIT_COMMITTER_EMAIL=infra-root+obltmachine@elastic.co
# Run the oblt-robot binary.
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "--verbose", "--dry-run" ]
