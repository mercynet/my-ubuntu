#!/usr/bin/env bash
# Runs inside meu-ubuntu-build:26.04. Produces dist/meu-ubuntu, .deb and SHA256SUMS.
set -euo pipefail
cd /src
VERSION="${VERSION:-$(cat VERSION)}"
PACKAGE_VERSION="${VERSION/-beta./\~beta}"
PACKAGE_VERSION="${PACKAGE_VERSION/-rc./\~rc}"
export SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct 2>/dev/null || echo 1790000000)}"

python3 build/scripts/check-version.py
python3 build/scripts/generate-docs.py --check
(cd frontend && npm ci --no-audit --no-fund && npm run build)

mkdir -p dist
python3 build/scripts/collect-licenses.py
go build -trimpath -buildvcs=false -tags desktop,production,webkit2_41 \
  -ldflags "-s -w -X main.version=${VERSION}" -o dist/meu-ubuntu .

# .deb declaring the runtime libraries the binary links against.
PKG=dist/deb/meu-ubuntu_${VERSION}_amd64
rm -rf dist/deb && mkdir -p "$PKG/DEBIAN" "$PKG/usr/bin" "$PKG/usr/share/applications" "$PKG/usr/share/icons/hicolor/512x512/apps"
install -m755 dist/meu-ubuntu "$PKG/usr/bin/meu-ubuntu"
install -m644 build/appicon.png "$PKG/usr/share/icons/hicolor/512x512/apps/meu-ubuntu.png"
install -m644 build/linux/meu-ubuntu.desktop "$PKG/usr/share/applications/meu-ubuntu.desktop"
mkdir -p "$PKG/usr/share/doc/meu-ubuntu/logos"
install -m644 README.md docs/manual.md docs/catalogo.md docs/desenvolvimento.md docs/beta-usabilidade.md "$PKG/usr/share/doc/meu-ubuntu/"
install -m644 LICENSE NOTICE.md CONTRIBUTING.md SECURITY.md CHANGELOG.md RELEASING.md VERSION "$PKG/usr/share/doc/meu-ubuntu/"
cp -R dist/third-party "$PKG/usr/share/doc/meu-ubuntu/"
mkdir -p "$PKG/usr/share/doc/meu-ubuntu/diagnostics"
install -m644 internal/diag/assets/*LICENSE.txt "$PKG/usr/share/doc/meu-ubuntu/diagnostics/"
install -m644 frontend/src/assets/apps/*LICENSE.txt frontend/src/assets/apps/sources.json frontend/src/assets/apps/README.md "$PKG/usr/share/doc/meu-ubuntu/logos/"
cat > "$PKG/DEBIAN/control" <<CTRL
Package: meu-ubuntu
Version: ${PACKAGE_VERSION}
Architecture: amd64
Maintainer: Paulo Silva <paulo.silva@guisoft.net>
Depends: libgtk-3-0t64, libwebkit2gtk-4.1-0, libjavascriptcoregtk-4.1-0, libglib2.0-0t64, libc6, pkexec, gnupg, ca-certificates
Recommends: flatpak, gnome-shell-extension-prefs
Section: admin
Priority: optional
Description: Central pessoal de pós-instalação do Ubuntu (Wails/Go)
 Escolher aplicações, ajustes GNOME, extensões e perfis reutilizáveis,
 com plano revisto antes de executar e operações administrativas via pkexec.
CTRL
find "$PKG" -exec touch -h -d "@${SOURCE_DATE_EPOCH}" {} +
dpkg-deb --root-owner-group -Zxz --build "$PKG" "dist/meu-ubuntu_${VERSION}_amd64.deb" >/dev/null
cp dist/meu-ubuntu "dist/meu-ubuntu-${VERSION}-linux-amd64"
python3 build/scripts/source-archive.py
(cd dist && sha256sum "meu-ubuntu-${VERSION}-linux-amd64" "meu-ubuntu_${VERSION}_amd64.deb" "meu-ubuntu-${VERSION}-source.tar.gz" > SHA256SUMS)
python3 build/scripts/release-bundle.py
rm -rf dist/deb
cat dist/SHA256SUMS
[ -n "${HOST_UID:-}" ] && chown -R "${HOST_UID}:${HOST_GID:-$HOST_UID}" dist frontend/dist frontend/node_modules || true
