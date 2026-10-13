#!/usr/bin/env bash
set -euo pipefail
export DISPLAY=:99 LANG=C.UTF-8 XDG_SESSION_TYPE=x11 NO_AT_BRIDGE=0
export HOME=/out/home XDG_CONFIG_HOME=/out/config XDG_DATA_HOME=/out/data XDG_STATE_HOME=/out/state XDG_RUNTIME_DIR=/out/runtime
mkdir -p "$HOME" "$XDG_RUNTIME_DIR" /out/project "$XDG_CONFIG_HOME/meu-ubuntu-sim"
mkdir -p "$XDG_CONFIG_HOME/meu-ubuntu-sim/integrations"
printf '%s\n' '"/out/project"' > "$XDG_CONFIG_HOME/meu-ubuntu-sim/integrations/officegest.json"
chmod 700 "$XDG_RUNTIME_DIR"
printf '%s\n' '{"onboarded":true,"theme":"light","showAllAreas":true,"selections":{}}' > "$XDG_CONFIG_HOME/meu-ubuntu-sim/prefs.json"
printf '%s\n' 'unchanged' > /out/project/sentinel
Xvfb :99 -screen 0 1280x900x24 >/out/xvfb.log 2>&1 &
trap 'kill "$!" 2>/dev/null || true' EXIT
sleep 1
dbus-run-session -- python3 /src/build/test/rc3-native-smoke.py
