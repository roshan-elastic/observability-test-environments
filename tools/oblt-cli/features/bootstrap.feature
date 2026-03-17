Feature: Bootstrap the Elastic Stack with known recipes.
  I want to bootstrap my Elastic Stack with some configuration option all developers share.

  Background:
    Given a file named "/tmp/foo.yaml" with:
    """
    slack-channel: '@foo'
    username: bar
    """

  Scenario: show bootstrap help
    When I successfully run `oblt-cli bootstrap --help`
    Then the output should contain:
    """
    oblt-cli bootstrap [command]
    """

  Scenario: show bootstrap cluster help
    When I successfully run `oblt-cli bootstrap cluster --help`
    Then the output should contain:
    """
    oblt-cli bootstrap cluster [flags]
    """

  Scenario: show bootstrap cluster apply all recipes
    When I successfully run `oblt-cli bootstrap cluster --cluster-name edge-oblt --dry-run --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Applying recipe : elasticsearch
    """
    And the output should contain:
    """
    Applying recipe : kibana
    """
    And the output should contain:
    """
    Enable username authentication
    """

  Scenario: show bootstrap cluster apply one Elasticsearch recipe
    When I successfully run `oblt-cli bootstrap cluster --cluster-name edge-oblt --recipes-elasticsearch '["test"]' --dry-run --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Applying recipe : elasticsearch
    """
    And the output should contain:
    """
    Applying recipe : kibana
    """

  Scenario: show bootstrap cluster apply one Kibana recipe
    When I successfully run `oblt-cli bootstrap cluster --cluster-name edge-oblt --recipes-kibana '["test"]' --dry-run --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Applying recipe : elasticsearch
    """
    And the output should contain:
    """
    Applying recipe : kibana
    """
    And the output should contain:
    """
    Enable username authentication
    """

  Scenario: show bootstrap Elasticsearch help
    When I successfully run `oblt-cli bootstrap elasticsearch --help`
    Then the output should contain:
    """
    oblt-cli bootstrap elasticsearch [flags]
    """

  Scenario: show bootstrap Elasticsearch applying one recipe using username authentication
    When I successfully run `oblt-cli bootstrap elasticsearch --url https://es.example.com --username elastic --password changeme  --recipes '["test"]' --dry-run --config /tmp/foo.yaml  --verbose`
    Then the output should contain:
    """
    Applying recipe : elasticsearch
    """
    And the output should contain:
    """
    Enable username authentication
    """

  Scenario: show bootstrap Elasticsearch applying one recipe using apiKey authentication
    When I successfully run `oblt-cli bootstrap elasticsearch --url https://es.example.com --api-key foo --recipes '["test"]' --dry-run --config /tmp/foo.yaml  --verbose`
    Then the output should contain:
    """
    Applying recipe : elasticsearch
    """
    And the output should contain:
    """
    Enable ApiKey authentication
    """

  Scenario: show bootstrap Kibana help
    When I successfully run `oblt-cli bootstrap kibana --help`
    Then the output should contain:
    """
    oblt-cli bootstrap kibana [flags]
    """

  Scenario: show bootstrap Kibana applying one recipe using username authentication
    When I successfully run `oblt-cli bootstrap kibana --url https://kibana.example.com --username elastic --password changeme --recipes '["test"]' --dry-run --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Applying recipe : kibana
    """
    And the output should contain:
    """
    Enable username authentication
    """

  Scenario: show bootstrap Kibana applying one recipe using username authentication
    When I successfully run `oblt-cli bootstrap kibana --url https://kibana.example.com --api-key foo --recipes '["test"]' --dry-run --config /tmp/foo.yaml --verbose`
    Then the output should contain:
    """
    Applying recipe : kibana
    """
    And the output should contain:
    """
    Enable ApiKey authentication
    """
