#!/usr/bin/env bash
# Visual smoke test: opens the real binary in simulation mode under Xvfb and captures pages.
set -u
OUT=/out
export DISPLAY=:99 LANG=pt_PT.UTF-8 XDG_SESSION_TYPE=x11
Xvfb :99 -screen 0 1280x840x24 >/dev/null 2>&1 &
sleep 1
eval "$(dbus-launch --sh-syntax)"
THEME="${THEME:-light}"
MEU_UBUNTU_SIMULATE=1 GTK_THEME="Adwaita$([ "$THEME" = dark ] && echo :dark)" meu-ubuntu ${APPARGS:-} > $OUT/app-$THEME.log 2>&1 &
sleep 7
shot() { import -window root "$OUT/$THEME-$1.png"; }
click() { xdotool mousemove "$1" "$2" click 1; sleep "${3:-1.5}"; }
scroll() { xdotool mousemove 700 500; for i in $(seq 1 "$1"); do xdotool click 5; done; sleep 1; }
shot 01-inicio
for spec in "$@"; do
  name="${spec%%@*}"; xy="${spec#*@}"
  IFS=',' read -r x y w <<<"$xy"
  if [ "$x" = scroll ]; then scroll "$y"; else click "$x" "$y" "${w:-1.5}"; fi
  shot "$name"
done
pkill meu-ubuntu; true
