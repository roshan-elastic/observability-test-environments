Feature: create clusters
  I want to create a developer clusters.

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster create --help`
    Then the output should contain:
    """
    oblt-cli cluster create [command]
    """
