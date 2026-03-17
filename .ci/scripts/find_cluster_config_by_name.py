#!/usr/bin/env python3
"""Script to find the cluster config file by cluster name.

Script to crawl through the cluster config directory
and find the cluster config file by checking the cluster name in the yaml file
Usage: python find_cluster_config_by_name.py <cluster_name>
Example: python find_cluster_config_by_name.py mycluster
Returns: The path to the cluster config file
"""


import os
import sys
from collections.abc import Mapping

import click
import yaml


def read_yaml_file(file_path):
    """Read a yaml file and returns the content as a dictionary."""
    with open(file_path, 'r', encoding='UTF-8') as stream:
        try:
            return yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            print(exc)
    return None


def find_cluster_config_by_name(name, path):
    """Crawls through the cluster config directory and returns the path to the cluster config file."""
    for root, dirs, files in os.walk(os.path.abspath(path)):
        for file in files:
            if file.endswith(".yaml") or file.endswith(".yml"):
                file_path = os.path.join(root, file)
                cluster_config = read_yaml_file(file_path)
                if (cluster_config is not None
                        and isinstance(cluster_config, Mapping)
                        and cluster_config.get('cluster_name', '') == name):
                    return file_path

    return None


@click.command()
@click.argument('cluster_name', required=True)
@click.argument('environments_path', default="./environments")
@click.argument('relative', default=False)
def main(cluster_name, environments_path, relative):
    """Run the script."""
    cluster_config_path = find_cluster_config_by_name(cluster_name, environments_path)
    if cluster_config_path is not None:
        if relative:
            print(os.path.relpath(cluster_config_path, start=environments_path))
        else:
            print(cluster_config_path)


if __name__ == '__main__':
    main()
