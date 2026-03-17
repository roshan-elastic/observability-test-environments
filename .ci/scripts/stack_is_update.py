#!/usr/bin/env python3
"""Script to check if the version is a downgrade."""

import os
import sys

import click
import semantic_version
import yaml


def read_cluster_config(cluster_config: str) -> dict:
    """Read the cluster configuration."""
    if not os.path.exists(cluster_config):
        print(f"The file does not exists {cluster_config}.")
        sys.exit(1)
    with open(cluster_config, mode="r", encoding="UTF-8") as f:
        return yaml.safe_load(f)


def is_greater_or_equal(from_, to) -> bool:
    """Compare semver two versions and fail if a downgrade."""
    ret = True
    to_trucate = semantic_version.Version(to).truncate("patch")
    from_trucate = semantic_version.Version(from_).truncate("patch")
    if (to_trucate > from_trucate):
        ret = True
        print(f"Upgrade detected from {from_} to {to}", file=sys.stderr)
    if (to_trucate == from_trucate):
        ret = True
        print(f"The version are equal {from_} to {to}", file=sys.stderr)
    if (to_trucate < from_trucate):
        ret = False
        print(f"Downgrade detected from {from_} to {to}", file=sys.stderr)
    return ret


def is_same_build(from_, to) -> bool:
    """Compare two versions and return True is they are the same build."""
    ret = True
    to_build = semantic_version.Version(to).prerelease
    from_build = semantic_version.Version(from_).prerelease
    if (to_build == from_build):
        ret = True
        print(f"The build are equal {from_} to {to}", file=sys.stderr)
    if (to_build != from_build):
        ret = False
        print(f"The build are different {from_} to {to}", file=sys.stderr)
    return ret


def get_current_build(cluster_config: str) -> str:
    """Get the current build from the cluster configuration."""
    config = read_cluster_config(cluster_config)
    stack = config.get("stack", {})
    current_version = stack.get("version", None)
    current_build = stack.get("build", '')
    stack_mode = stack.get("mode", "ess")
    stack_mode_config = stack.get(stack_mode, {})
    if current_version is not None and len(current_build) == 0:
        es_config = stack_mode_config.get("elasticsearch", {})
        current_build = es_config.get("image", "").split(":")[-1]
        if len(current_build) == 0:
            kbn_config = stack_mode_config.get("kibana", {})
            current_build = kbn_config.get("image", "").split(":")[-1]
        if len(current_build) == 0:
            current_build = current_version
    return current_build


def is_update(cluster_config: str, version: str) -> bool:
    """Check if the version is a update."""
    current_build = get_current_build(cluster_config)
    ret = (is_greater_or_equal(current_build, version)
           and not is_same_build(current_build, version))
    return ret


@click.command()
@click.option("--cluster-config",
              help="Cluster configuration file",
              required=True)
@click.option("--version",
              help="Version to check",
              required=True)
def main(cluster_config: str, version: str):
    """Run script."""
    if version == "" or version is None:
        print("Version is empty", file=sys.stderr)
        sys.exit(1)
    if cluster_config == "" or cluster_config is None:
        print("Cluster config is empty", file=sys.stderr)
        sys.exit(1)
    if is_update(cluster_config, version):
        sys.exit(0)
    sys.exit(1)


if __name__ == "__main__":
    main()
