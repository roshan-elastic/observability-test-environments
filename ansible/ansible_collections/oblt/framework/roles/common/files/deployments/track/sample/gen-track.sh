#!/usr/bin/env bash
# This script counts the number of documents in the sample collection
# and prints the result to the console.

if [ ! -f documents.json ]; then
  curl -LO http://download.geonames.org/export/dump/allCountries.zip
  unzip allCountries.zip
  python3 geo2json.py > documents.json
fi

__DOCUMENT_COUNT__=$(wc -l < documents.json|tr -d ' ')

case "$OSTYPE" in
  darwin*) __DOCUMENT_SIZE_IN_BYTES__=$(stat -f %z documents.json) ;;
  linux*)  __DOCUMENT_SIZE_IN_BYTES__=$(stat -c %s documents.json) ;;
  *)        echo "unknown: $OSTYPE"; exit 1;;
esac

cp track.json.tmpl track.json
sed -i ".bak" "s/__DOCUMENT_COUNT__/${__DOCUMENT_COUNT__}/g" track.json
sed -i ".bak" "s/__DOCUMENT_SIZE_IN_BYTES__/${__DOCUMENT_SIZE_IN_BYTES__}/g" track.json
