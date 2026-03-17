Feature: list templates
  I want to list clusters templates

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: list templates
    When I successfully run `oblt-cli cluster templates --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    [ccs]
    """

  Scenario: show help
    When I successfully run `oblt-cli cluster templates --help`
    Then the output should contain:
    """
    oblt-cli cluster templates [flags]
    """
