#!/usr/bin/env python3
"""
Script to check if there are new Docker images for the Elastic Stack.

It returns the list of Elastic Stack versions that have new Docker images.
the environment variable HOURS_AGO=2 define the maximum age
of the Docker images.
"""

import json
import os
from datetime import datetime, timedelta

import requests
from gh_actions.actions import choose_output
from tenacity import retry, stop_after_attempt, wait_exponential

ARTIFACTORY_URL = 'https://artifacts-api.elastic.co/v1/versions'


def get_elastic_stack_versions() -> list:
    """Get the list of Elastic Stack versions."""
    response = get_json(ARTIFACTORY_URL)
    versions = []
    if response is not None:
        versions = response['versions']
        versions.sort()
    return versions


@retry(stop=stop_after_attempt(10),
       wait=wait_exponential(multiplier=1, min=4, max=60))
def get_json(url: str) -> dict:
    """Get the JSON from the given URL."""
    response = requests.get(url, headers={'Accept': 'application/json'},
                            timeout=30)
    ret = {}
    if response.status_code == 200:
        ret = response.json()
    return ret


def check_for_new_docker_images(elastic_version: str) -> bool:
    """Check if there are new Docker images for the given version."""
    hours_ago = int(os.environ.get('HOURS_AGO', 2))
    build_url = f'{ARTIFACTORY_URL}/{elastic_version}/builds/latest'
    build_json = get_json(build_url)
    is_new = False
    if build_json is not None:
        end_time = build_json["build"]["end_time"]
        end_time_obj = datetime.strptime(end_time, '%a, %d %b %Y %H:%M:%S %Z')
        now = datetime.now()
        delta = now - end_time_obj
        is_new = delta < timedelta(hours=hours_ago)
    return is_new


def run() -> None:
    """Run the script."""
    new_versions = []
    for version in get_elastic_stack_versions():
        if check_for_new_docker_images(version):
            new_versions.append(version)

    with choose_output() as output_file:
        output_file.write(f"versions={json.dumps(new_versions)}\n")


if __name__ == '__main__':
    run()
