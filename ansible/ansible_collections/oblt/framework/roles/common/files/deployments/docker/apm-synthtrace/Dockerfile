ARG NODE_VERSION=18.17.1-bullseye
FROM node:${NODE_VERSION}-bullseye

ARG BRANCH=main
USER node
WORKDIR /workspace
ENV WORKSPACE=/workspace
ENV KIBANA_DIR=${WORKSPACE}/kibana

RUN git clone --depth 1 --single-branch --branch "${BRANCH}" https://github.com/elastic/kibana.git

COPY --chown=node:0 bootstrap.sh ${WORKSPACE}
RUN chmod ugo+rx "${WORKSPACE}/bootstrap.sh" \
  && "${WORKSPACE}/bootstrap.sh"

ENV ELASTICSEARCH_USERNAME=elastic
ENV ELASTICSEARCH_PASSWORD=changeme
ENV ELASTICSEARCH_HOST=localhost:9200
WORKDIR ${KIBANA_DIR}
ENTRYPOINT [ "node" ]
CMD [ "scripts/synthtrace", "simple_trace.ts",  "--target=https://${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}@${ELASTICSEARCH_HOST}", "--live"]
