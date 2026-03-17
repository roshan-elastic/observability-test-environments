FROM gitpod/workspace-base

#RUN brew install gh

RUN id
USER root
RUN mkdir -p /opt/openvscode-server
# CONFIGURE NON-INTERACTIVE APT
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
      gh \
      rsync \
    && apt-get clean autoclean \
    && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}

RUN curl -sSL https://github.com/gitpod-io/openvscode-server/releases/download/openvscode-server-v1.77.3/openvscode-server-v1.77.3-linux-x64.tar.gz \
      -o openvscode-server.tgz \
    && tar -xzf openvscode-server.tgz --strip-components=1 -C /opt/openvscode-server \
    && rm openvscode-server.tgz
RUN curl -sSL https://github.com/elastic/elastic-package/releases/download/v0.79.0/elastic-package_0.79.0_linux_amd64.tar.gz \
      -o elastic-package.tgz \
    && tar -xzf elastic-package.tgz -C /usr/local/bin elastic-package \
    && rm elastic-package.tgz

RUN useradd -l -u 1000 -G sudo -md /home/gitpod -s /bin/bash -p elastic elastic \
    && mkdir -p /workspace \
    && chown -R elastic:elastic /workspace \
    && chown -R elastic:elastic /home/gitpod

USER elastic
WORKDIR /workspace

ENTRYPOINT ["/bin/bash"]
CMD ["-c", "cd '/opt/openvscode-server' && ./bin/openvscode-server --without-connection-token --port 3000 --host 127.0.0.1"]
EXPOSE 3000
