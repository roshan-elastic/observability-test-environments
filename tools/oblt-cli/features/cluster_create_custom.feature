Feature: create custom clusters
  I want to create a developer cluster configured using a cluster template.

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster create custom --help`
    Then the output should contain:
    """
    oblt-cli cluster create custom [flags]
    """
