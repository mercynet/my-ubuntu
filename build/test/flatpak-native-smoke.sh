#!/usr/bin/env bash
# Run only in the native QA container; the application uses its simulation executor.
set -euo pipefail
export DISPLAY=:99 LANG=C.UTF-8 XDG_SESSION_TYPE=x11 NO_AT_BRIDGE=0
export HOME=/out/home XDG_CONFIG_HOME=/out/config XDG_DATA_HOME=/out/data XDG_STATE_HOME=/out/state XDG_RUNTIME_DIR=/out/runtime
mkdir -p "$HOME" "$XDG_RUNTIME_DIR" "$XDG_CONFIG_HOME/meu-ubuntu-sim"
chmod 700 "$XDG_RUNTIME_DIR"
printf '%s\n' '{"onboarded":true,"theme":"light","showAllAreas":true,"selections":{}}' > "$XDG_CONFIG_HOME/meu-ubuntu-sim/prefs.json"
Xvfb :99 -screen 0 1280x900x24 >/out/xvfb.log 2>&1 &
display_pid=$!
trap 'kill "$display_pid" 2>/dev/null || true' EXIT
sleep 1
dbus-run-session -- python3 /src/build/test/flatpak-native-smoke.py
