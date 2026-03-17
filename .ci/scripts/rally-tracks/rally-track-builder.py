#!/usr/bin/env python3
"""This script search for known issues in a rally track files and fix them.

The files with size 0 will be deleted and the track.json will be updated.
"""

import json
import os
import shutil
import subprocess  # nosec B404

import click
import requests
from jinja2 import (BaseLoader, ChoiceLoader, DictLoader, Environment,
                    FileSystemLoader)

# https://github.com/elastic/rally/blob/master/esrally/track/loader.py#L752-L774
macros = [
    """
    {% macro collect(parts) -%}
        {% set comma = joiner() %}
        {% for part in glob(parts) %}
            {{ comma() }}
            {% include part %}
        {% endfor %}
    {%- endmacro %}
    """,
    """
    {% macro exists_set_param(setting_name, value, default_value=None, comma=True) -%}
        {% if value is defined or default_value is not none %}
            {% if comma %} , {% endif %}
            {% if default_value is not none %}
                "{{ setting_name }}": {{ value | default(default_value) | tojson }}
            {% else %}
                "{{ setting_name }}": {{ value | tojson }}
            {% endif %}
            {% endif %}
    {%- endmacro %}
    """,
]


def walk_all_files_in_dir(path: str) -> tuple[list[str], list[str], set[str]]:
    """Walk through all files in a directory."""
    zero_files = []
    data_files = []
    pipelines = set()
    for root, dirs, files in os.walk(path):
        for filename in files:
            abspath = os.path.join(root, filename)
            if (is_zero_file(abspath)
                    or filename.endswith('-documents-1k.json')
                    or filename.endswith('-documents.json.offset')):
                zero_files.append(abspath)
            else:
                if (filename.endswith('.json')
                        and filename != 'track.json'
                        and filename.find('-documents') == -1):
                    data_files.append(abspath)
                    pipelines.update(load_pipelines(abspath))
    return tuple((zero_files, data_files, pipelines))


def is_zero_file(path: str) -> bool:
    """Check if a file is empty."""
    exists = os.path.exists(path)
    size = exists and os.path.getsize(path)
    return not exists or size == 0


def load_track_file(track_file: str, track_settings: dict) -> str:
    """Parse the jinja2 template track file."""
    environment = Environment(loader=ChoiceLoader(   # nosec B701
        [
            DictLoader({"rally.helpers": "".join(macros)}),
            BaseLoader(),
            FileSystemLoader("/")
        ]),
        autoescape=False)
    template = environment.get_template(track_file)
    rendered = template.render(track_settings)
    return str(rendered)


def delete_files(files: list[str]):
    """Delete files."""
    for item in files:
        print(f'Removing {item}')
        os.remove(item)


def get_index_names_from_track_indices(track: dict) -> list[str]:
    """Get all index names from track indices."""
    indices_names = []
    for index in track['indices']:
        indices_names.append(str(index.get('name')))
    return indices_names


def get_index_names_from_files(files: list[str]) -> list[str]:
    """Get all index names from files."""
    names = []
    for index in files:
        names.append(os.path.basename(index).replace('.json', '').replace('.bz2', ''))
    return names


def save_track_file(track_file: str, track: str) -> bool:
    """Save a track file."""
    with open(track_file, 'w', encoding='utf-8') as file:
        file.write(track)
    print(f'New track file saved to {track_file}')
    return True


def load_settings_file(settings_file: str) -> dict:
    """Load a settings file."""
    if settings_file is not None:
        with open(settings_file, 'r', encoding='utf-8') as file:
            return json.load(file)
    return {}


def load_indices_files(indices_files: list[str]) -> list[dict]:
    """Load all indices files."""
    indices = []
    for index in indices_files:
        print(f'Loading {index} info')
        stats = os.stat(index)
        name = index.replace('.json', '').replace('.bz2', '').replace('-documents-1k', '').replace('-documents-1k', '')
        documents = f'{name}-documents.json'
        if not is_zero_file(documents):
            obj = {
                "name": os.path.basename(name),
                "file": documents,
                "file_1k": f'{name}-documents-1k.json',
                "name_bz2": os.path.basename(f'{name}-documents.json.bz2'),
                "size": stats.st_size,
                "documents": count_lines(documents),

            }
            indices.append(obj)
    return indices


def count_lines(index) -> int:
    """Count the number of lines in a file."""
    ret = 0
    with open(index, 'rb') as file:
        for i, line in enumerate(file):
            ret = i + 1
    return ret


def load_pipelines(pipeline_file: str) -> list[str]:
    """Load all pipelines from a file."""
    pipelines = []
    with open(pipeline_file, 'r', encoding='utf-8') as file:
        print(f'Loading pipelines from {pipeline_file}')
        obj = json.load(file)
        settings = obj.get('settings')
        if settings is not None:
            index = settings.get('index')
            if index is not None:
                default_pipeline = index.get('default_pipeline')
                final_pipeline = index.get('final_pipeline')
                if default_pipeline is not None:
                    pipelines.append(default_pipeline)
                if final_pipeline is not None:
                    pipelines.append(final_pipeline)
    return pipelines


def load_pipeline_definitions(pipeline_names: list[str], es_url: str,
                              es_username: str,
                              es_password: str) -> list[dict]:
    """Load all pipeline definitions from Elasticseacr."""
    pipelines = []
    for name in pipeline_names:
        print(f'Loading pipeline {name}')
        url = f'{es_url}/_ingest/pipeline/{name}'
        headers = {
            'Content-Type': 'application/json',
        }
        basic = requests.auth.HTTPBasicAuth(es_username, es_password)
        result = requests.get(url,
                              headers=headers,
                              auth=basic,
                              timeout=30)
        if result.status_code == 200:
            body = result.json().get(name)
            body.pop('_meta')
            bodys = json.dumps(body, indent=2)
            pipeline = {
                "name": name,
                "body": bodys
            }
            pipelines.append(pipeline)
    return pipelines


def sort_records(indices: list[dict]):
    """Sort all records in all files by timestamp attribute."""
    for index in indices:
        file_name = index.get("file")
        print(f'Sorting {file_name}')
        with open(f"{file_name}.sorted", 'w', encoding='utf-8') as file:
            subprocess.run(["jlsort", "-k", "'@timestamp'", file_name],  # nosec B607, B603
                           check=True,
                           stdout=file,
                           stderr=subprocess.DEVNULL)
        os.remove(file_name)
        os.rename(f"{file_name}.sorted", file_name)


def generate_1k_files(indices: list[dict]):
    """Generate a 1k records file from each file."""
    for index in indices:
        file_name = index.get("file")
        file_1k_name = index.get("file_1k")
        print(f'Generating 1k file for {file_name}')
        with open(file_name, 'r', encoding='utf-8') as source:
            documents_to_read = min(index.get("documents"), 1000)
            head = [next(source) for _ in range(documents_to_read)]
            with open(file_1k_name, 'w', encoding='utf-8') as outfile:
                outfile.writelines(head)


def choose_bzip2_command(file_name: str) -> list[str]:
    """Choose the bzip2 command to use."""
    if shutil.which('7z') is not None:
        return ["7z", "a", "-y", f"{file_name}.bz2", file_name]
    elif shutil.which('pbzip2') is not None:
        return ['pbzip2', "-k", "-f", file_name]
    elif shutil.which('lbzip2') is not None:
        return ['lbzip2', "-k", "-f", file_name]
    else:
        return ['bzip2', "-k", "-f", file_name]


def bzip2_files(indices: list[dict]):
    """Bzip2 all files."""
    for index in indices:
        file_name = index.get("file")
        print(f'Compressing {file_name}')
        command = choose_bzip2_command(file_name)
        subprocess.run(command, check=True)  # nosec B603


@click.command()
@click.option('--track-path',
              required=True,
              help='Absolute path to the rally track to fix.')
@click.option('--settings',
              required=False,
              help='Absolute path to the rally JSON settings file.')
@click.option('--track-name',
              required=True,
              help='Name of the track.')
@click.option('--indices-pattern',
              required=False,
              default='.ds*',
              help='Pattern to match indices.')
@click.option('--es-url',
              required=False,
              default='http://localhost:9200',
              help='Elasticsearch URL.')
@click.option('--es-username',
              required=False,
              default='elastic',
              help='Elasticsearch username.')
@click.option('--es-password',
              required=False,
              default='changeme',
              help='Elasticsearch password.')
def main(track_path: str, settings: str,
         track_name: str, indices_pattern: str,
         es_url: str, es_username: str, es_password: str):
    """Run the main function."""
    zero_files, data_files, pipelines = walk_all_files_in_dir(track_path)
    if len(zero_files) == 0:
        print('No files to fix found. Checking trackfile for missing indices.')
    print(f'Found {len(zero_files)} zero files.')
    delete_files(zero_files)
    indices_files = load_indices_files(data_files)

    track_settings = load_settings_file(settings)
    track_settings['indices'] = indices_files
    track_settings['track_name'] = track_name
    track_settings['track_path'] = track_path
    track_settings['indices_pattern'] = indices_pattern
    track_settings['pipelines'] = load_pipeline_definitions(pipelines,
                                                            es_url,
                                                            es_username,
                                                            es_password)

    track_template = os.path.join(os.path.dirname(os.path.abspath(__file__)),
                                  'track.json.j2')
    track_file = os.path.join(track_path, 'track.json')

    track = load_track_file(track_template, track_settings)
    save_track_file(track_file, track)

    sort_records(indices_files)
    generate_1k_files(indices_files)
    bzip2_files(indices_files)


if __name__ == '__main__':
    main()
