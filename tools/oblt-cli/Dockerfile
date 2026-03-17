FROM docker.io/library/alpine:3.14.2

WORKDIR /home/nonroot

RUN apk add --no-cache git openssh

COPY oblt-cli /oblt-cli

ENV HOME=/home/nonroot
ENV USER=nonroot
ENV UID=1000
ENV GID=1000

RUN addgroup -g "${GID}" "${USER}"
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "${HOME}" \
    --ingroup "$USER" \
    --uid "$UID" \
    "$USER"


RUN mkdir -p "${HOME}/.oblt-cli" "${HOME}/.ssh" \
  && chown -R "${USER}:${USER}" "${HOME}" \
  && chmod 0700 "${HOME}/.ssh"
COPY --chown=nonroot:nonroot ./static/known_hosts "${HOME}/.ssh/known_hosts"
RUN chmod 0600 "${HOME}/.ssh/known_hosts"

USER nonroot:nonroot

ENTRYPOINT ["/oblt-cli"]
