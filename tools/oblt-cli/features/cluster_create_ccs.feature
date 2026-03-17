Feature: create CCS clusters
  I want to create a developer cluster configured with CCS to an existing oblt cluster.

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster create ccs --help`
    Then the output should contain:
    """
    oblt-cli cluster create ccs [flags]
    """
