Feature: list available snapshots in the Unified Release
  I want to list the snapshots

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: list snapshot
    When I successfully run `oblt-cli unified-release snapshots`
    Then the output should contain:
    """
    | VERSION |
    """

  Scenario: show help
    When I successfully run `oblt-cli unified-release snapshots --help`
    Then the output should contain:
    """
    oblt-cli unified-release snapshots [flags]
    """
