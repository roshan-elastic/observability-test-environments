#!/usr/bin/env python3
"""
Script to process a ndjson file and create a set of files to use with bootstrap.

This script will process a ndjson file and create a set of files to use with
bootstrap. The ndjson file is expected to contain one line per request to
bootstrap. The script will create one file per line in the ndjson file.
"""

import json

import click
import yaml


def convert_headers(headers: list[str]) -> dict[str, str]:
    """Convert headers from list to dict."""
    result = {}
    for header in headers:
        key, value = header.split(':')
        result[key] = value.strip()
    return result


def save_bootstrap_file(api: str, method: str, body: str, filename: str,
                        headers: list[str], return_code: list[int],
                        ignore_errors: bool):
    """Save a bootstrap file."""
    json_body = json.dumps(json.loads(body), indent=2)
    print(json_body)

    def str_presenter(dumper, data):
        if len(data.splitlines()) > 1:  # check for multiline string
            return dumper.represent_scalar('tag:yaml.org,2002:str',
                                           data, style='|')
        return dumper.represent_scalar('tag:yaml.org,2002:str', data)

    yaml.add_representer(str, str_presenter)

    # to use with safe_dump:
    yaml.representer.SafeRepresenter.add_representer(str, str_presenter)

    bootstrap_data = {
        "api": api,
        "method": method,
        "body_format": "form-multipart",
        "body": {
            "file": {
                "content": f'{json_body}',
                "filename": "data.ndjson",
                "mime_type": "application/x-ndjson",
            },
        },
        "headers": convert_headers(headers),
        "return_code": return_code,
        "ignore_errors": ignore_errors
    }
    with open(filename, 'w', encoding='utf-8') as bootstrap_file:
        yaml.safe_dump(bootstrap_data, bootstrap_file)


def process_ndjson(filename_pattern: str, api: str, method: str,
                   headers: list[str], codes: list[int],
                   ignore_errors: bool, ndjson_file: str):
    """Process a ndjson file."""
    with open(ndjson_file, encoding='utf-8') as ndjson:
        data = ndjson.readlines()
    for i, line in enumerate(data):
        save_bootstrap_file(
            api=api,
            method=method,
            body=line,
            filename=f'{filename_pattern}-{i:03}.yml',
            headers=headers,
            return_code=codes,
            ignore_errors=ignore_errors
        )


@click.command()
@click.option('--file', help='ndjson file to process', required=True)
@click.option('--api', help='api URL to use', required=True)
@click.option('--headers',
              type=str,
              help='headers to use',
              multiple=True,
              default=["Content-Type: application/json"])
@click.option('--codes',
              type=int,
              help='error codes to use',
              multiple=True,
              default=[200])
@click.option('--ignore-errors', help='ignore errors', default=False)
@click.option('--filename-pattern', help='filename pattern', required=True)
@click.option('--method',
              help='HTTP method used to call the API', required=True)
def main(file: str, api: str, method: str, headers: list[str],
         codes: list[int], ignore_errors: bool, filename_pattern: str):
    """Run process."""
    process_ndjson(filename_pattern,
                   api, method, headers, codes, ignore_errors, file)


if __name__ == "__main__":
    main()
