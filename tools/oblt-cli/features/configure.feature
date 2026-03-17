Feature: configure command
  As a newcomer to oblt-cli
  I want to configure oblt-cli

  Scenario: configure oblt-cli with correct parameters
    When I successfully run `oblt-cli configure --username <username> --slack-channel <slack_channel> --config /tmp/foo.yaml  --verbose`
    Then the file "/tmp/foo.yaml" should contain:
    """
    username: <username>
    """
    And the file "/tmp/foo.yaml" should contain:
    """
    slack-channel: '<slack_channel>'
    """

  Examples:
  | username | slack_channel |
  | foo      | #foo          |
  | bar      | @bar          |

  Scenario: configure oblt-cli with wrong username
    When I run `oblt-cli configure --username <username> --slack-channel <slack_channel> --config /tmp/foo.yaml  --verbose`
    Then the output should contain:
    """
    should contains only the following characters [a-z0-9]
    """

  Examples:
  | username | slack_channel |
  | foo.bar      | #foo      |

  Scenario: configure oblt-cli with wrong slack channel
    When I run `oblt-cli configure --username <username> --slack-channel <slack_channel> --config /tmp/foo.yaml  --verbose`
    Then the output should contain:
    """
    Error: the Slack channel/member ID should start with # or @
    """

  Examples:
  | username | slack_channel |
  | foo      | bar           |

  Scenario: configure oblt-cli with correct parameters
    When I successfully run `oblt-cli configure --help`
    Then the output should contain:
    """
    oblt-cli configure [flags]
    """
