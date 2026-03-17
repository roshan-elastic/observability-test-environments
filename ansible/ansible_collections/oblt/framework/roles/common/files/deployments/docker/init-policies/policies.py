#!/usr/bin/env python
"""Policies for the Elastic Agent Helm chart."""

import json
from typing import Any

import click
import requests
import yaml


def create_or_update_policy(url: str, username: str, password: str,
                            policy_type: str, data: Any) -> None:
    """Create a policy.

    It calls the Fleet API to create a policy.
    """
    headers = {
        'Content-Type': 'application/json',
        'kbn-xsrf': 'true'}
    response = None
    policy_id = data.get('id')
    base_url = f'{url}/api/fleet/{policy_type}'
    exists_response = requests.get(f'{base_url}/{policy_id}',
                                   auth=(username, password),
                                   timeout=30,
                                   headers=headers)
    if exists_response.status_code == 200:
        response = update_policy(username, password, data, headers, policy_id, base_url)
    elif exists_response.status_code == 404:
        response = create_policy(username, password, data, headers, policy_id, base_url)
    if response.status_code == 200:
        print('Policy created successfully')
    else:
        print(f'Failed to create policy: {response.status_code} {response.text}')


def create_policy(username, password, data, headers, policy_id, base_url):
    """Create a policy."""
    print(f'Policy {policy_id} does not exist, creating')
    return requests.post(f'{base_url}',
                         json=json.dumps(data),
                         auth=(username, password),
                         timeout=30,
                         headers=headers)


def update_policy(username, password, data, headers, policy_id, base_url):
    """Update a policy."""
    print(f'Policy {policy_id} already exists, updating')
    data.pop('id')
    return requests.put(f'{base_url}/{policy_id}',
                        json=json.dumps(data),
                        auth=(username, password),
                        timeout=30,
                        headers=headers)


def get_policies(filename, policy_type):
    """Get a field from the file."""
    field = f'kibana.{policy_type}'
    field_value = None
    with open(filename, encoding="UTF-8") as yaml_file:
        field_yaml = yaml.safe_load(yaml_file)
        field_value = field_yaml
        for item in field.split('.'):
            field_value = field_value.get(item, {})
    return field_value


def process_policies(kibana_url, kibana_username, kibana_password,
                     policy_type_item, field_value):
    """Iterate over the policies and create or update them."""
    for policy in field_value:
        print(policy.get('name'))
        create_or_update_policy(kibana_url,
                                kibana_username,
                                kibana_password,
                                policy_type_item,
                                policy)


def process_policy_types(kibana_url, kibana_username, kibana_password,
                         filename, policy_type):
    """Iterate over the policy types."""
    for policy_type_item in policy_type.split(','):
        print(f'Extracting field {policy_type_item} from {filename}')
        field_value = get_policies(filename, policy_type_item)
        process_policies(kibana_url, kibana_username, kibana_password,
                         policy_type_item, field_value)


@click.command()
@click.option('--kibana-url',
              envvar='KIBANA_HOST',
              default='https://kibana.example.com',
              show_default=True,
              help='The Elasticsearch URL')
@click.option('--kibana-username',
              envvar='KIBANA_USERNAME',
              default='elastic',
              show_default=True,
              help='The Elasticsearch user')
@click.option('--kibana-password',
              envvar='KIBANA_PASSWORD',
              default='changeme',
              show_default=True,
              help='The Elasticsearch password')
@click.option('--filename',
              default='policies.yaml',
              show_default=True,
              help='The file to read from')
@click.option('--policy-type',
              default='agent_policies,package_policies',
              show_default=True,
              help='Extract agent policies')
def main(kibana_url, kibana_username, kibana_password,
         filename, policy_type) -> None:
    """Extract a field from the policies.yaml file."""
    process_policy_types(kibana_url, kibana_username, kibana_password, filename, policy_type)


if __name__ == '__main__':
    main()
