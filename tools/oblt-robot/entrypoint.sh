#!/bin/sh
set -e
export GITHUB_TOKEN="${GITHUB_PASSWORD}"
gh auth setup-git
/go/bin/oblt-robot "$@"
