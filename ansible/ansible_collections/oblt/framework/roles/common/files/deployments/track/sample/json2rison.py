#!/usr/bin/env python
"""Python script to convert JSON to Rison."""
import json
import sys

import prison

buf = json.load(sys.stdin)
print(prison.dumps(buf), end='')
