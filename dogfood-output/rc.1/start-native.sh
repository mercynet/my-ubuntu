#!/usr/bin/env bash
set -euo pipefail
cd /src
task_home=$(mktemp -d /tmp/meu-ubuntu-rc-native-XXXXXX)
mkdir -p "$task_home"/{config/meu-ubuntu-sim,state,data,runtime,exports}
chmod 700 "$task_home/runtime"
printf '%s\n' '{"theme":"dark","onboarded":true,"showAllAreas":true,"selections":{},"origins":{},"reverts":[],"manualProgress":{}}' > "$task_home/config/meu-ubuntu-sim/prefs.json"
chown -R 1000:1000 "$task_home"
Xvfb :99 -screen 0 1656x900x24 > /tmp/rc-xvfb.log 2>&1 &
echo "$!" > /tmp/rc-xvfb.pid
runuser -u ubuntu -- env HOME="$task_home" XDG_CONFIG_HOME="$task_home/config" XDG_STATE_HOME="$task_home/state" XDG_DATA_HOME="$task_home/data" XDG_RUNTIME_DIR="$task_home/runtime" DISPLAY=:99 WEBKIT_DISABLE_DMABUF_RENDERER=1 LIBGL_ALWAYS_SOFTWARE=1 WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS=1 dbus-run-session /src/dist/meu-ubuntu-1.0.0-rc.1-linux-amd64 --simular > /tmp/rc-native.log 2>&1 &
echo "$!" > /tmp/rc-native.pid
echo "$task_home" > /tmp/rc-native-home
for attempt in $(seq 1 30); do
  if DISPLAY=:99 xdotool search --onlyvisible --class . > /tmp/rc-native-windows 2>/dev/null; then break; fi
  sleep 1
done
DISPLAY=:99 xdotool search --onlyvisible --class .
