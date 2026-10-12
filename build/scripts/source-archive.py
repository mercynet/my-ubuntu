#!/usr/bin/env python3
"""Build an explicit source archive without local state or private directories."""
import gzip
import os
from pathlib import Path
import tarfile

ROOT = Path(__file__).resolve().parents[2]
version = os.environ.get("VERSION", "1.0.0-rc.1")
epoch = int(os.environ.get("SOURCE_DATE_EPOCH", "1790000000"))
prefix = "meu-ubuntu-" + version + "-source"
roots = ["LICENSE", "NOTICE.md", "README.md", "AGENTS.md", ".gitignore", "CONTRIBUTING.md",
         "SECURITY.md", "CHANGELOG.md", "Makefile", "go.mod", "go.sum", "wails.json",
         "internal", "frontend", "build", "docs", ".github"]
files = list(ROOT.glob("*.go"))
for name in roots:
    path = ROOT / name
    if path.is_dir():
        files.extend(p for p in path.rglob("*") if p.is_file() and not p.is_symlink())
    elif path.is_file():
        files.append(path)
files.extend(p for p in (ROOT / "dogfood-output").rglob("*.md") if not p.is_symlink())
excluded = {"node_modules", "dist", "__pycache__", "harness", ".cache"}
files = sorted(set(p for p in files if not excluded.intersection(p.relative_to(ROOT).parts)))
output = ROOT / "dist" / (prefix + ".tar.gz")
with output.open("wb") as dest, gzip.GzipFile(filename="", mode="wb", fileobj=dest, mtime=epoch) as zipped:
    with tarfile.open(fileobj=zipped, mode="w") as archive:
        for path in files:
            info = archive.gettarinfo(str(path), prefix + "/" + str(path.relative_to(ROOT)))
            info.uid = info.gid = 0
            info.uname = info.gname = "root"
            info.mtime = epoch
            with path.open("rb") as source:
                archive.addfile(info, source)
print(f"Source archive: {output.name} ({len(files)} files).")
