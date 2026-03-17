#!/usr/bin/env python
"""Python script to convert Rison to JSON."""
import json
import sys

import prison

buf = sys.stdin.readlines()
print(json.dumps(prison.loads(buf)), end='')
