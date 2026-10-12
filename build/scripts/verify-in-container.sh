#!/usr/bin/env bash
set -euo pipefail
cd /src
VERSION="${VERSION:-1.0.0-rc.1}"
PACKAGE_VERSION="${VERSION/-beta./\~beta}"
PACKAGE_VERSION="${PACKAGE_VERSION/-rc./\~rc}"
deb="dist/meu-ubuntu_${VERSION}_amd64.deb"

(cd dist && sha256sum -c SHA256SUMS)
test "$(dist/meu-ubuntu --version)" = "meu-ubuntu ${VERSION}"
test "$(dist/meu-ubuntu-${VERSION}-linux-amd64 --version)" = "meu-ubuntu ${VERSION}"
test "$(dpkg-deb -f "$deb" Package)" = "meu-ubuntu"
test "$(dpkg-deb -f "$deb" Version)" = "$PACKAGE_VERSION"
test "$(dpkg-deb -f "$deb" Architecture)" = "amd64"
dist/meu-ubuntu --help | grep -F -- '--manual' >/dev/null
dist/meu-ubuntu --manual | grep -F 'Manual do Meu Ubuntu' >/dev/null
dist/meu-ubuntu --licenses | grep -F 'GNU GENERAL PUBLIC LICENSE' >/dev/null
for doc in manual.md catalogo.md desenvolvimento.md beta-usabilidade.md README.md LICENSE NOTICE.md CONTRIBUTING.md SECURITY.md CHANGELOG.md third-party/modules.json diagnostics/FPDF-LICENSE.txt diagnostics/DejaVu-LICENSE.txt; do
  dpkg-deb --fsys-tarfile "$deb" | tar -t | grep -Fx "./usr/share/doc/meu-ubuntu/$doc" >/dev/null
done
tar -tzf "dist/meu-ubuntu-${VERSION}-source.tar.gz" | grep -Fx "meu-ubuntu-${VERSION}-source/LICENSE" >/dev/null
python3 -m zipfile -t "dist/meu-ubuntu-${VERSION}.zip"
if [[ "$VERSION" == *-rc.* ]]; then
  dpkg --compare-versions "$PACKAGE_VERSION" lt "${VERSION%-rc.*}"
fi
printf 'Artefactos e versão %s conferidos.\n' "$VERSION"
