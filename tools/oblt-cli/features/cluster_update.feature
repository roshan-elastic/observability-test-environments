Feature: update clusters
  I want to update a developer clusters.

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster update --help`
    Then the output should contain:
    """
    oblt-cli cluster update [flags]
    """
