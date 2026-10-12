#!/usr/bin/env bash
set -euo pipefail
cd /src
task_home=$(mktemp -d /tmp/meu-ubuntu-about-XXXXXX)
mkdir -p "$task_home"/{config/meu-ubuntu-sim,state,data,runtime}
chmod 700 "$task_home/runtime"
cat > "$task_home/config/meu-ubuntu-sim/prefs.json" <<'JSON'
{"theme":"dark","onboarded":true,"showAllAreas":true,"selections":{},"origins":{},"reverts":[],"manualProgress":{}}
JSON
chown -R 1000:1000 "$task_home"
Xvfb :99 -screen 0 1656x900x24 >/tmp/about-xvfb.log 2>&1 &
echo "$!" >/tmp/about-xvfb.pid
runuser -u ubuntu -- env HOME="$task_home" XDG_CONFIG_HOME="$task_home/config" XDG_STATE_HOME="$task_home/state" XDG_DATA_HOME="$task_home/data" XDG_RUNTIME_DIR="$task_home/runtime" DISPLAY=:99 WEBKIT_DISABLE_DMABUF_RENDERER=1 LIBGL_ALWAYS_SOFTWARE=1 WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS=1 dbus-run-session /src/dist/meu-ubuntu-0.6.0-beta.5-linux-amd64 --simular >/tmp/about-native.log 2>&1 &
echo "$!" >/tmp/about-native.pid
echo "$task_home" >/tmp/about-native-home
for attempt in $(seq 1 30); do
  if DISPLAY=:99 xdotool search --onlyvisible --class . >/tmp/about-windows 2>/dev/null; then break; fi
  sleep 1
done
DISPLAY=:99 xdotool search --onlyvisible --class .
