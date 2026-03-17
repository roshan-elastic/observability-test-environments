#!/usr/bin/env bash
"""Script to process the countries and generate a ndJSON file."""

import json

cols = (("geonameid", "int", True),
        ("name", "string", True),
        ("asciiname", "string", False),
        ("alternatenames", "string", False),
        ("latitude", "double", True),
        ("longitude", "double", True),
        ("feature_class", "string", False),
        ("feature_code", "string", False),
        ("country_code", "string", True),
        ("cc2", "string", False),
        ("admin1_code", "string", False),
        ("admin2_code", "string", False),
        ("admin3_code", "string", False),
        ("admin4_code", "string", False),
        ("population", "long", True),
        ("elevation", "int", False),
        ("dem", "string", False),
        ("timezone", "string", False))


def main():
    """Run script."""
    with open("allCountries.txt", "rt", encoding="UTF-8") as f:
        for line in f:
            tup = line.strip().split("\t")
            record = {}
            for i in range(len(cols)):
                name, type, include = cols[i]
                if tup[i] != "" and include:
                    if type in ("int", "long"):
                        record[name] = int(tup[i])
                    elif type == "double":
                        record[name] = float(tup[i])
                    elif type == "string":
                        record[name] = tup[i]
            print(json.dumps(record, ensure_ascii=False))


if __name__ == "__main__":
    main()
