#!/usr/bin/env python3
"""Update a JSON file with a new value for a key."""
import json

import click


def load_json(filename):
    """Load a JSON file and return a dictionary."""
    with open(filename, 'r', encoding="UTF-8") as f:
        return json.load(f)


def write_json(filename, data):
    """Write a dictionary to a JSON file."""
    with open(filename, 'w', encoding="UTF-8") as f:
        json.dump(data, f, indent=4)


def replace_key_value(data, key, value):
    """Replace the value of a key in a dictionary.

    The key can be a nested key, e.g. 'a.b.c' or 'a.b[0].c'
    """
    keys = key.split('.')
    dcopy = data
    for key in keys[:-1]:
        # if the key is an array get the key and the index
        if '[' in key:
            key, i = key.split('[')
            i = int(i[:-1])
            dcopy = dcopy[key][i]
        else:
            dcopy = dcopy[key]
    dcopy[keys[-1]] = value


@click.command()
@click.option('--filename', required=True, help='JSON file to update')
@click.option('--key', required=True, help='Key to update')
@click.option('--value', required=True, help='New value')
@click.option('--backup', is_flag=True, help='Create a backup of the file')
def main(filename, key, value, backup):
    """Update a JSON file with a new value for a key."""
    data = load_json(filename)
    if backup:
        write_json(filename + '.bak', data)
    replace_key_value(data, key, value)
    write_json(filename, data)


# main
if __name__ == '__main__':
    main()
