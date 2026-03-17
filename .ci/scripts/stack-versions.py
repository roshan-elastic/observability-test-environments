#!/usr/bin/env python3
"""
Script to get Elastic Stack versions from artifactory.

It has several modes to work with:
    - all: get all versions
    - release: get only release versions
    - snapshot: get only snapshot versions
    - all-unique: get all versions except SNAPSHOTs that has a release
    - snapshot-unique: get all snapshot versions except SNAPSHOTs that has a release
    - is-greater-or-equal: compare two given versions and if there is a downgrade then return 1

It has several output formats:
    - json: output a json list
    - text: output a text list

It is possible to limit the number of results and get a specific position.

For example to get the last 2 versions:

    python3 stack-versions.py --mode all --limit 2 --output json

Or to get the previous version to the last one:

    python3 stack-versions.py --mode all --position -2 --output json

Or the last:

    python3 stack-versions.py --mode all --position -1 --output json

    python3 stack-versions.py --mode all --limit 1 --output json
"""

import json
import sys

import click
import requests
import semantic_version

artifactory_url = 'https://artifacts-api.elastic.co/v1/versions/'


def get_versions():
    """Get the versions from the artifactory."""
    rest = requests.get(artifactory_url, timeout=30)
    if rest.status_code != 200:
        print("Error getting versions from artifactory", file=sys.stderr)
        sys.exit(1)
    return rest.json()['versions']


def get_snapshot_versions(versions):
    """Get snapshot versions."""
    snapshot_versions = []
    for version in versions:
        if version.endswith('-SNAPSHOT'):
            snapshot_versions.append(version)
    return snapshot_versions


def get_release_versions(versions):
    """Get releases versions."""
    release_versions = []
    for version in versions:
        if not version.endswith('-SNAPSHOT'):
            release_versions.append(version)
    return release_versions


def key_func(item):
    """Get key function to order semvers."""
    clean_item = item.replace('-SNAPSHOT', '').split('+')[0]
    return tuple(int(v) for v in clean_item.split('.'))


def sort_versions(versions):
    """Sort versions."""
    return sorted(versions, key=key_func)


def output_versions(versions, limit, position, output):
    """Output versions."""
    out_versions = sort_versions(versions)
    if position:
        out_position = out_versions[position]
        out_versions = [out_position]
    else:
        out_versions = out_versions[-1 * limit:]

    if len(out_versions) == 0:
        print("No versions found", file=sys.stderr)
        sys.exit(1)

    if output == 'json':
        print(json.dumps(out_versions))
    else:
        for version in out_versions:
            print(version)


def filter_releases(versions, release_versions):
    """Filter releases versions."""
    ret = []
    for version in versions:
        if not version.replace('-SNAPSHOT', '') in release_versions:
            ret.append(version)
    return ret


def select_versions(mode):
    """Select versions to show."""
    versions = get_versions()
    snapshot_versions = get_snapshot_versions(versions)
    release_versions = get_release_versions(versions)
    ret = []
    if mode == 'snapshot':
        ret = snapshot_versions
    elif mode == 'release':
        ret = release_versions
    elif mode == 'snapshot-unique':
        ret = filter_releases(snapshot_versions, release_versions)
    elif mode == 'all-unique':
        ret = filter_releases(versions, release_versions)
        ret += release_versions
        ret = sort_versions(ret)
    else:
        ret = versions
    return ret


def is_greater_or_equal(from_, to):
    """Compare semver two versions and fail if a downgrade."""
    if (semantic_version.Version(to).truncate('patch')
            >= semantic_version.Version(from_).truncate('patch')):
        return 0
    print(f"Downgrade detected from {from_} to {to}", file=sys.stderr)
    return 1


@click.command()
@click.option('--mode',
              default='all',
              type=click.Choice([
                  'all',
                  'release',
                  'snapshot',
                  'all-unique',
                  'snapshot-unique',
                  'is-greater-or-equal']),
              help='Select mode to get versions. By default get all versions.')
@click.option('--limit',
              default=2,
              help='Limit the number of versions')
@click.option('--position',
              type=int,
              help='Return a position of the version in the list')
@click.option('--output',
              default='json',
              type=click.Choice(['json', 'text']),
              help='Select the typo of output, by default is json.'
              )
@click.option('--from', '-f', 'from_')
@click.option('--to', '-t')
def main(mode, limit, position, output, from_, to):
    """Get the versions from the artifactory."""
    if mode == 'is-greater-or-equal':
        sys.exit(is_greater_or_equal(from_, to))
    else:
        versions = select_versions(mode)
        output_versions(versions, limit, position, output)


if __name__ == "__main__":
    main()
