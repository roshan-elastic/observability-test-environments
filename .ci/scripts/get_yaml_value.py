#!/usr/bin/env python3
"""Python script to get a value from a YAML file.

Usage: python get_yaml_value.py <yaml_file> <key>
Example: python get_yaml_value.py .ci/.jenkins_job.groovy JOB_NAME
"""

import sys

import yaml


def get_yaml_value(yaml_file, key):
    """Get a value from a YAML file."""
    with open(yaml_file, 'r') as f:
        try:
            data = yaml.safe_load(f)
            for item in key.split('.'):
                data = data[item]
            return data
        except yaml.YAMLError as exc:
            print(exc)


if __name__ == '__main__':
    if len(sys.argv) < 3:
        print("Usage: python get_yaml_value.py <yaml_file> <key>")
        exit(1)
    print(get_yaml_value(sys.argv[1], sys.argv[2]))
