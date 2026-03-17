Feature: destroy clusters
  I want to destroy clusters

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster destroy --help`
    Then the output should contain:
    """
    oblt-cli cluster destroy [flags]
    """

  Scenario Outline: Destroying Golden Clusters is forbidden
    Given I successfully run `oblt-cli configure --username therobots --slack-channel "@1234356"  --verbose`
    When I run `oblt-cli cluster destroy --cluster-name <cluster_name>`
    Then the exit status should be 1
      And the output should contain:
    """
    It's forbidden to destroy the following clusters: [dev-oblt edge-oblt monitoring-oblt release-oblt]. You wanted to destroy [<cluster_name>]
    """
  Examples:
  | cluster_name    |
  | dev-oblt        |
  | edge-oblt       |
  | monitoring-oblt |
  | release-oblt    |
