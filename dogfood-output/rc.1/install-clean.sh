#!/usr/bin/env bash
set -euo pipefail
mkdir -p /evidence
apt-get update -qq > /evidence/install.log 2>&1
DEBIAN_FRONTEND=noninteractive apt-get install -y /dist/meu-ubuntu_1.0.0-rc.1_amd64.deb >> /evidence/install.log 2>&1
test "$(dpkg-query -W -f='${Version}' meu-ubuntu)" = "1.0.0~rc1"
test "$(meu-ubuntu --version)" = "meu-ubuntu 1.0.0-rc.1"
meu-ubuntu --manual | grep -F 'Manual do Meu Ubuntu' >/dev/null
meu-ubuntu --licenses | grep -F 'GNU GENERAL PUBLIC LICENSE' >/dev/null
test -f /usr/share/doc/meu-ubuntu/third-party/modules.json
printf 'RC installed in a clean Ubuntu 26.04 container; package, CLI, manual and licenses verified.\n' | tee /evidence/install-result.txt
chown -R 1000:1000 /evidence
