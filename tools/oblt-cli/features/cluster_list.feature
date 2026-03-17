Feature: list clusters
  I want to list clusters deployed

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: list clusters
    When I successfully run `oblt-cli cluster list --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Cluster Configurations Available
    """

  Scenario: list all clusters
    When I successfully run `oblt-cli cluster list --all --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    edge-oblt
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster list --help`
    Then the output should contain:
    """
    oblt-cli cluster list [flags]
    """
