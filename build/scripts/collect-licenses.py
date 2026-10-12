#!/usr/bin/env python3
"""Collect redistributable dependency notices from the resolved build inputs."""
import json
import os
from pathlib import Path
import shutil
import subprocess

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "dist" / "third-party"
if OUT.exists():
    shutil.rmtree(OUT)
OUT.mkdir(parents=True)

raw = subprocess.check_output(["go", "list", "-m", "-json", "all"], cwd=ROOT, text=True)
decoder = json.JSONDecoder()
modules = []
while raw.strip():
    value, end = decoder.raw_decode(raw.lstrip())
    modules.append(value)
    raw = raw.lstrip()[end:]

index = []


def collect(name, version, folder):
    folder = Path(folder)
    candidates = sorted(p for p in folder.iterdir() if p.is_file() and (
        p.name.upper().startswith(("LICENSE", "LICENCE", "COPYING", "NOTICE"))))
    if not candidates:
        raise SystemExit(f"Missing license notice: {name} {version}")
    target = OUT / name.replace("/", "__").replace("@", "")
    target.mkdir(exist_ok=True)
    for source in candidates:
        shutil.copyfile(source, target / source.name)
    index.append({"name": name, "version": version,
                  "notices": [str((target / p.name).relative_to(OUT)) for p in candidates]})


for module in modules:
    if module.get("Main"):
        continue
    source = module.get("Replace", module)
    folder = source.get("Dir")
    if not folder:
        subprocess.run(["go", "mod", "download", source["Path"] + "@" + source["Version"]],
                       cwd=ROOT, check=True)
        metadata = json.loads(subprocess.check_output(
            ["go", "mod", "download", "-json", source["Path"] + "@" + source["Version"]],
            cwd=ROOT, text=True))
        folder = metadata["Dir"]
    collect(module["Path"], module.get("Version", "local"), folder)

goroot = subprocess.check_output(["go", "env", "GOROOT"], text=True).strip()
collect("Go", subprocess.check_output(["go", "version"], text=True).strip(), goroot)
for name in ("react", "react-dom", "scheduler", "lucide-react"):
    folder = ROOT / "frontend" / "node_modules" / name
    package = json.loads((folder / "package.json").read_text())
    collect("npm/" + name, package["version"], folder)

(OUT / "modules.json").write_text(json.dumps(index, ensure_ascii=False, indent=2) + "\n")
print(f"Collected license notices for {len(index)} dependencies.")
