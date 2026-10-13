#!/usr/bin/env python3
"""Set the source version and both frontend manifests together."""
import argparse
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('version')
args = parser.parse_args()
if not re.fullmatch(r'\d+\.\d+\.\d+(?:-(?:rc|beta)\.\d+)?', args.version):
    parser.error('Use X.Y.Z, X.Y.Z-rc.N ou X.Y.Z-beta.N')
for name in ('package.json', 'package-lock.json'):
    path = ROOT/'frontend'/name
    data = json.loads(path.read_text())
    data['version'] = args.version
    if name == 'package-lock.json': data['packages']['']['version'] = args.version
    path.write_text(json.dumps(data, indent=2)+'\n')
(ROOT/'VERSION').write_text(args.version+'\n')
print('Version set:', args.version)
