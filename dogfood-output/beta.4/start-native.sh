#!/bin/sh
set -eu
mkdir -p /tmp/beta4/home /tmp/beta4/config /tmp/beta4/state /tmp/beta4/data /tmp/beta4/run
chmod 700 /tmp/beta4/run
chown -R 1000:1000 /tmp/beta4
Xvfb :99 -screen 0 1656x900x24 >/tmp/xvfb.log 2>&1 &
runuser -u ubuntu -- env HOME=/tmp/beta4/home XDG_CONFIG_HOME=/tmp/beta4/config XDG_STATE_HOME=/tmp/beta4/state XDG_DATA_HOME=/tmp/beta4/data XDG_RUNTIME_DIR=/tmp/beta4/run DISPLAY=:99 WEBKIT_DISABLE_DMABUF_RENDERER=1 LIBGL_ALWAYS_SOFTWARE=1 WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS=1 dbus-run-session /src/dist/meu-ubuntu-0.6.0-beta.4-linux-amd64 >/tmp/native.log 2>&1 &
