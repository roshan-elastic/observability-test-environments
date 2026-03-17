Feature: list available releases in the Unified Release
  I want to list the releases

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: list releases
    When I successfully run `oblt-cli unified-release releases`
    Then the output should contain:
    """
    | VERSION |
    """

  Scenario: show help
    When I successfully run `oblt-cli unified-release releases --help`
    Then the output should contain:
    """
    oblt-cli unified-release releases [flags]
    """
