#!/usr/bin/env python3
"""
Python script to given a folder with molecule tests, list the tests and the scenarios for each test.

Usage: python3 list_molecule_tests.py <folder>
Example: python3 list_molecule_tests.py molecule
Output:
 [
   {
       role: "test1", scenario: "scenario1"
    },
   {
       role: "test1", scenario: "scenario2"
    },
  ]
"""

import json
import os

import click


def list_molecule_tests(folder):
    """List molecule tests in a folder."""
    tests = []
    for role in sorted(os.listdir(folder)):
        if os.path.isdir(os.path.join(folder, role, "molecule")):
            tests.append(os.path.join(folder, role, "molecule"))
    return tests


def get_test_name(test):
    """Get the name of the test from the path."""
    return test.split("/")[-2]


def get_test_scenarios(test_folder):
    """Get the scenarios for a test."""
    scenarios = []
    for scenario in sorted(os.listdir(test_folder)):
        if (os.path.isdir(os.path.join(test_folder, scenario))
                and scenario not in ["__pycache__"]):
            scenarios.append(scenario)
    return scenarios


def get_matrix(tests):
    """Get the matrix of roles and scenarios."""
    test_list = []
    for test in tests:
        for scenario in get_test_scenarios(test):
            test_list.append({
                "role": get_test_name(test),
                "scenario": scenario
            })

    return test_list


@click.command()
@click.argument('folder', required=True)
def main(folder):
    """Execute main function."""
    tests = list_molecule_tests(folder)
    test_list = get_matrix(tests)
    print(json.dumps(test_list))


if __name__ == "__main__":
    main()
