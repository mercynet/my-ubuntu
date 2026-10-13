#!/usr/bin/env python3
"""Keep the build version and frontend manifests consistent."""
import json
import os
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
version = (ROOT/'VERSION').read_text().strip()
if not re.fullmatch(r'\d+\.\d+\.\d+(?:-(?:rc|beta)\.\d+)?', version):
    raise SystemExit('VERSION inválida; use X.Y.Z, X.Y.Z-rc.N ou X.Y.Z-beta.N')
if os.environ.get('VERSION', version) != version:
    raise SystemExit('VERSION do ambiente difere das fontes; use set-version.py antes de construir')
package = json.loads((ROOT/'frontend/package.json').read_text())
lock = json.loads((ROOT/'frontend/package-lock.json').read_text())
if any(v != version for v in (package['version'], lock['version'], lock['packages']['']['version'])):
    raise SystemExit('Versões frontend inconsistentes; execute python3 build/scripts/set-version.py '+version)
print('Version:', version)
