#!/usr/bin/env python3
"""Check the visible Flatpak prerequisite and reviewed plan in native simulation."""
import hashlib
import json
from pathlib import Path
import subprocess
import time

import pyatspi

OUT = Path("/out")
BINARY = Path("/src/dist/meu-ubuntu")
previous_runs = set((OUT / "state/meu-ubuntu-sim/runs").glob("*.json"))
log = (OUT / "app.log").open("w")
app = subprocess.Popen([str(BINARY), "--simular"], stdout=log, stderr=log)


def nodes(obj=None, depth=0):
    if obj is None:
        obj = pyatspi.Registry.getDesktop(0)
    if depth > 32:
        return
    yield obj
    try:
        for child in obj:
            yield from nodes(child, depth + 1)
    except Exception:
        pass


def wait(check, seconds=25):
    deadline = time.monotonic() + seconds
    while time.monotonic() < deadline:
        result = check()
        if result:
            return result
        time.sleep(0.1)
    raise AssertionError("Expected native UI state was not reached")


def button(name):
    for obj in nodes():
        try:
            if obj.name == name and obj.getRoleName() in ("button", "push button", "check box"):
                if obj.getState().contains(pyatspi.STATE_ENABLED):
                    return obj
        except Exception:
            pass


def click(name):
    assert wait(lambda: button(name)).queryAction().doAction(0), name
    time.sleep(0.2)


def contains_text(text):
    for obj in nodes():
        try:
            if text in (obj.name or "") or text in obj.queryText().getText(0, -1):
                return True
        except Exception:
            pass
    return False


def search_entry():
    for obj in nodes():
        try:
            if obj.name == "Pesquisar em todo o catálogo" and obj.getState().contains(pyatspi.STATE_EDITABLE):
                return obj
        except Exception:
            pass


def shot(name):
    subprocess.run(["import", "-window", "root", str(OUT / (name + ".png"))], check=True)


def completed_run():
    for path in (OUT / "state/meu-ubuntu-sim/runs").glob("*.json"):
        if path in previous_runs:
            continue
        try:
            run = json.loads(path.read_text())
            if run.get("finished"):
                return run
        except (FileNotFoundError, ValueError):
            pass


try:
    entry = wait(search_entry, 40)
    # WebKit exposes HTML inputs as editable without the AT-SPI EditableText
    # interface. Type into their visible bounds through X11 instead.
    windows = subprocess.check_output(["xdotool", "search", "--onlyvisible", "--name", "Meu Ubuntu"], text=True).splitlines()
    subprocess.run(["xdotool", "windowfocus", "--sync", windows[-1]], check=True, timeout=5)
    rect = entry.queryComponent().getExtents(pyatspi.DESKTOP_COORDS)
    subprocess.run(["xdotool", "mousemove", str(rect.x + rect.width // 2), str(rect.y + rect.height // 2), "click", "1"], check=True)
    subprocess.run(["xdotool", "key", "ctrl+a"], check=True)
    subprocess.run(["xdotool", "type", "--clearmodifiers", "VLC"], check=True)
    wait(lambda: button("Incluir VLC nas escolhas"))
    assert contains_text("Requer Flatpak."), "missing visible Flatpak prerequisite"
    assert contains_text("Flatpak (utilizador)"), "incorrect installation method"
    shot("01-flatpak-card")
    click("Incluir VLC nas escolhas")
    click("Rever plano")
    wait(lambda: contains_text("Estas aplicações requerem Flatpak."))
    assert contains_text("Flatpak + Flathub"), "missing visible dependency"
    assert contains_text("VLC"), "missing reviewed application"
    shot("02-flatpak-review")
    click("Aplicar plano")
    run = wait(completed_run)
    assert all(step["status"] == "done" for step in run["steps"]), run["steps"]
    plan_steps = run["steps"]
    vlc_steps = [step for step in plan_steps if step.get("item") == "vlc"]
    assert len(vlc_steps) == 1 and vlc_steps[0]["op"] == "flatpak-install", vlc_steps
    assert any(step["op"] == "apt-install" and step["item"] == "flathub" for step in plan_steps), plan_steps
    assert any(step["op"] == "flatpak-remote" for step in plan_steps), plan_steps
    wait(lambda: contains_text("Concluído"))
    shot("03-flatpak-completed")
    result = {
        "binary_sha256": hashlib.sha256(BINARY.read_bytes()).hexdigest(),
        "version": subprocess.check_output([str(BINARY), "--version"], text=True).strip(),
        "scope": "Native GTK/WebKit X11/Xvfb; simulation executor; no real app installation, GNOME/Wayland or polkit",
        "result": "PASS",
        "checks": ["visible requirement and origin", "reviewed Flatpak/Flathub dependency", "simulated completion and persisted steps"],
        "steps": [{key: step.get(key) for key in ("item", "op", "title")} for step in plan_steps],
    }
    (OUT / "result.json").write_text(json.dumps(result, ensure_ascii=False, indent=2) + "\n")
    print("PASS: visible Flatpak requirement, dependency review and simulated completion in native GTK/WebKit.")
except Exception:
    shot("failure")
    print("Visible buttons:", sorted({obj.name for obj in nodes() if obj.getRoleName() in ("button", "push button", "check box")}))
    raise
finally:
    app.terminate()
    try:
        app.wait(timeout=5)
    except subprocess.TimeoutExpired:
        app.kill()
        app.wait()
    log.close()
