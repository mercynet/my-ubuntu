#!/usr/bin/env bash
# Exercise APT and preserve synthetic XDG data in a disposable Ubuntu container.
set -euo pipefail
test "${MEU_UBUNTU_DISPOSABLE:-}" = 1
test "$(id -u)" = 0
RELEASE_VERSION="${VERSION:?defina VERSION com a candidata a testar}"
. /etc/os-release
test "$ID:$VERSION_ID" = ubuntu:26.04

OLD_VERSION="${OLD_VERSION:-0.3.2}"
OUT=/evidence
mkdir -p "$OUT"
exec > >(tee "$OUT/package-upgrade.log") 2>&1

# Minimal Ubuntu images can exclude documentation before APT runs. Keep it for
# this test so package notices are checked after installation, as on a desktop.
python3 - <<'PY'
from pathlib import Path
for path in Path('/etc/dpkg/dpkg.cfg.d').glob('*'):
    if path.is_file():
        text = path.read_text()
        path.write_text('\n'.join(line for line in text.splitlines()
                                  if not line.startswith('path-exclude=/usr/share/doc'))+'\n')
PY
apt-get update -qq
DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "/dist/meu-ubuntu_${OLD_VERSION}_amd64.deb"
test "$(meu-ubuntu --version)" = "meu-ubuntu ${OLD_VERSION}"

python3 - <<'PY'
import hashlib
import json
from pathlib import Path

home = Path('/home/ubuntu')
fixtures = {
    '.config/meu-ubuntu/prefs.json': {'theme':'dark', 'onboarded':True, 'selections':{'htop':{}}, 'manualProgress':{'example':['reviewed']}},
    '.config/meu-ubuntu/profiles/my-test.json': {'id':'my-test', 'name':'My test', 'selections':{'htop':{}}},
    '.local/state/meu-ubuntu/runs/fixture.json': {'id':'fixture', 'status':'done', 'steps':[]},
    '.local/state/meu-ubuntu/logs/fixture.log': 'synthetic history\n',
    '.local/share/meu-ubuntu/fixture.txt': 'synthetic local data\n',
}
hashes = {}
for name, value in fixtures.items():
    path = home/name
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(value if isinstance(value,str) else json.dumps(value))
    path.chmod(0o600)
    hashes[name] = hashlib.sha256(path.read_bytes()).hexdigest()
Path('/evidence/upgrade-fixtures.json').write_text(json.dumps(hashes,indent=2)+'\n')
PY

DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "/dist/meu-ubuntu_${RELEASE_VERSION}_amd64.deb"
PACKAGE_VERSION="${RELEASE_VERSION/-rc./\~rc}"
test "$(dpkg-query -W -f='${Version}' meu-ubuntu)" = "$PACKAGE_VERSION"
test "$(meu-ubuntu --version)" = "meu-ubuntu ${RELEASE_VERSION}"
meu-ubuntu --manual | grep -F 'Manual do Meu Ubuntu' >/dev/null
meu-ubuntu --licenses | grep -F 'GNU GENERAL PUBLIC LICENSE' >/dev/null
for doc in LICENSE NOTICE.md third-party/modules.json manual.md; do
  test -f "/usr/share/doc/meu-ubuntu/$doc"
done
if command -v go || command -v node || command -v npm; then
  printf 'Unexpected development tools in the package-only environment.\n' >&2
  exit 1
fi
if ldd /usr/bin/meu-ubuntu | grep -F 'not found'; then
  exit 1
fi

python3 - <<'PY'
import hashlib
import json
from pathlib import Path

hashes = json.loads(Path('/evidence/upgrade-fixtures.json').read_text())
for name, expected in hashes.items():
    actual = hashlib.sha256((Path('/home/ubuntu')/name).read_bytes()).hexdigest()
    assert actual == expected, name
print('PASS: package installed/updated; synthetic preferences, profiles, history and data unchanged.')
print('Scope: APT/CLI/files in a container; GUI migrations, GNOME/Wayland and polkit remain NOT RUN.')
PY
