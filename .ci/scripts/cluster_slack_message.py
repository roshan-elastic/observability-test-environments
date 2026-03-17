#!/usr/bin/env python3
"""
Python script to compose a Slack message to send with GitHub Actions.

it uses the oblt cluster info files.
The message is URL encode to be possible to put it on a GitHub Action output.
https://github.com/orgs/community/discussions/26288
https://api.slack.com/reference/messaging/payload
"""

import os
from urllib.parse import quote

from gh_actions.actions import choose_output


def load_file(file_path: str) -> str:
    """Load a file."""
    ret = ""
    if os.path.isfile(file_path):
        with open(file_path, encoding="UTF-8") as file:
            ret = file.read()
    return ret


def run() -> None:
    """Run the script."""
    build_dir = os.environ.get('BUILD_DIR', './build')
    kibana_config = load_file(f'{build_dir}/kibana.yml')
    users = load_file(f'{build_dir}/users.yml')
    deploy_info = load_file(f'{build_dir}/deploy-info.md')

    kibana_message = quote(f"""
    ```
    {kibana_config}
    ```
    """)
    users_message = quote(f"{users}")
    deploy_info_message = quote(f"{deploy_info}")

    with choose_output() as output:
        output.write(f"usersMessage={users_message}\n")
        output.write(f"deployInfoMessage={deploy_info_message}\n")
        output.write(f"kibanaMessage={kibana_message}\n")
        # Add the messages to the GitHub Actions mask
        # should be in this with because the stdout is not available outside
        print(f"::add-mask::{users_message}")
        print(f"::add-mask::{deploy_info_message}")
        print(f"::add-mask::{kibana_message}")


if __name__ == '__main__':
    run()
