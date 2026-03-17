# Run in Docker

Oblt-cli is also available in a Docker image,
the Docker image makes git checkouts from GitHub by SSH,
so it needs you SSH key (`-v ${HOME}/.ssh/id_rsa:/home/nonroot/.ssh/id_rsa`, `-v ${HOME}/.ssh/id_rsa.pub:/home/nonroot/.ssh/id_rsa.pub`), your GitHub author name (`-e GIT_AUTHOR_NAME="$(git config user.name)"`, `-e GIT_COMMITTER_NAME="$(git config user.name)"`), and your email (`-e GIT_AUTHOR_EMAIL="$(git config user.email)"`, `-e GIT_COMMITTER_EMAIL="$(git config user.email)"`).
I we want to persist the oblt-cli configuration,
we will need to pass also a configuration line as volume (`-v ${HOME}/.oblt-cli/config.yaml:/home/nonroot/.oblt-cli/config.yaml`).

```bash
docker --version
touch "${HOME}/.oblt-cli/config.yaml"

docker run -it --rm \
  -e GIT_AUTHOR_NAME="$(git config user.name)" \
  -e GIT_AUTHOR_EMAIL="$(git config user.email)" \
  -e GIT_COMMITTER_NAME="$(git config user.name)" \
  -e GIT_COMMITTER_EMAIL="$(git config user.email)" \
  -v "${HOME}/.oblt-cli/config.yaml:/home/nonroot/.oblt-cli/config.yaml" \
  -v "${HOME}/.ssh/id_rsa:/home/nonroot/.ssh/id_rsa" \
  -v "${HOME}/.ssh/id_rsa.pub:/home/nonroot/.ssh/id_rsa.pub" \
  docker.elastic.co/observability-ci/oblt-cli:latest --help
```

**NOTE** This image is in our Docker registry in a private namespace,
so it require authentication, for more details check [Accessing the Docker Registry](https://github.com/elastic/infra/blob/master/docs/container-registry/accessing-the-docker-registry.md)
