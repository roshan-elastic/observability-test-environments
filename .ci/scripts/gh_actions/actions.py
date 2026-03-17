"""Python module to common functions to use in GitHub Actions."""

import os
import sys
from io import TextIOWrapper


def choose_output() -> TextIOWrapper:
    """Choose the output to use."""
    gitub_actions_output = os.getenv("GITHUB_OUTPUT", None)
    if gitub_actions_output is not None:
        return open(gitub_actions_output, "a", encoding="UTF-8")
    return sys.stdout
