#!/usr/bin/env bash
set -euo pipefail
exec > /src/private/build-setup.log 2>&1
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y --no-install-recommends nodejs npm python3 curl ca-certificates git gcc libc6-dev pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev dpkg-dev
curl --fail --silent --show-error --location https://go.dev/dl/go1.26.9.linux-amd64.tar.gz -o /tmp/go.tar.gz
printf '%s  /tmp/go.tar.gz\n' 42d158b4d8f7b61ac0a830567c940a86098fb7aac52e467a5ebec03ef5cc2f8d | sha256sum -c -
tar -C /usr/local -xzf /tmp/go.tar.gz
touch /tmp/packs-build-ready
sleep 3600
