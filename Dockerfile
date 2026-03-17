FROM python:3.11-slim

RUN apt update -y \
  && apt install --yes --no-install-recommends --no-install-suggests rsync make curl gnupg git nano bsdmainutils
RUN echo "deb [signed-by=/etc/apt/trusted.gpg.d/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list >/dev/null \
  && curl -fsSLo - https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /etc/apt/trusted.gpg.d/cloud.google.gpg
RUN apt update --yes \
  && apt install --yes --no-install-recommends --no-install-suggests google-cloud-sdk-gke-gcloud-auth-plugin
# add a user to run ansible named oblt and set the workdir to /ansible
RUN apt clean -y \
  && rm -rf /var/lib/apt/lists/*
RUN useradd -ms /bin/bash oblt
COPY ansible /ansible
COPY .ci/bin /ansible/bin
COPY .ci/scripts /ansible/scripts
RUN mkdir -p /home/oblt /data \
&& chown -R oblt:oblt /ansible /home/oblt /data
USER oblt
ENV PATH=/ansible/bin:${HOME}/bin::${PATH}
ENV HOME=/home/oblt
ENV BINDIR=/ansible/bin
ENV BUILD_DIR=/data/build
ENV SCRIPTS_DIR=/ansible/scripts
ENV VENV=/ansible/.venv
WORKDIR /ansible
RUN make install-ansible
RUN make install-hermit-tools
ENV CONTAINER=true
ENTRYPOINT [ "make" ]
