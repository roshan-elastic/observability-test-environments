# Github Actions

module "your-project-gh-actions" {
  source            = "../../../../modules/google-github-actions-workflow-identity-provider" # path may vary
  repository        = "your-repository"
  workflow_filenames = [
    "you-workflow-filename.yml"
  ]
}

module "google-oblt-cluster-secrets-access-gh" {
  source  = "../../../../modules/google-oblt-cluster-secrets-access" # path may vary
  members = values(module.your-project-gh-actions.principals)
}


# Buildkite Pipelines

module "your-project-bk-pipelines" {
  source            = "../../../../modules/google-buildkite-workload-identity-provider" # path may vary
  repository        = "your-repository"
  pipeline_slugs = [
    "your-pipeline-slug",
  ]
}

module "google-oblt-cluster-secrets-access-bk" {
  source  = "../../../../modules/google-oblt-cluster-secrets-access" # path may vary
  members = values(module.your-project-bk-pipelines.principals)
}
