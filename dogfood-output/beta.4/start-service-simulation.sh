#!/bin/sh
set -eu
mkdir -p /tmp/beta4-sim/home /tmp/beta4-sim/config/meu-ubuntu-sim /tmp/beta4-sim/state /tmp/beta4-sim/data /tmp/beta4-sim/run
chmod 700 /tmp/beta4-sim/run
printf '%s\n' '{"onboarded":true,"theme":"light","selections":{"apache2":{}},"origins":{"apache2":["manual"]}}' > /tmp/beta4-sim/config/meu-ubuntu-sim/prefs.json
chown -R 1000:1000 /tmp/beta4-sim
runuser -u ubuntu -- env HOME=/tmp/beta4-sim/home XDG_CONFIG_HOME=/tmp/beta4-sim/config XDG_STATE_HOME=/tmp/beta4-sim/state XDG_DATA_HOME=/tmp/beta4-sim/data XDG_RUNTIME_DIR=/tmp/beta4-sim/run DISPLAY=:99 WEBKIT_DISABLE_DMABUF_RENDERER=1 LIBGL_ALWAYS_SOFTWARE=1 WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS=1 dbus-run-session /src/dist/meu-ubuntu-0.6.0-beta.4-linux-amd64 --simular >/tmp/native-sim.log 2>&1 &
