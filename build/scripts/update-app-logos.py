#!/usr/bin/env python3
"""Refresh the curated, bundled app logos. Not run during normal builds."""
import hashlib
import json
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from urllib.request import urlopen
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "frontend/src/assets/apps"
SIMPLE = "https://raw.githubusercontent.com/simple-icons/simple-icons/16.13.0"
DEVICON = "https://raw.githubusercontent.com/devicons/devicon/v2.16.0"

# Explicit mapping prevents similarly named brands from receiving the wrong logo.
DEV = {
    "git": ["git", "git-flow", "git-defaults", "git-identity"],
    "github": ["gh"], "python": ["python-dev"], "go": ["golang"],
    "android": ["android-tools"], "sqlite": ["sqlite"],
    "mariadb": ["mariadb-client", "mariadb-server"], "dbeaver": ["dbeaver"],
    "php": ["php", "php-ppa", "php-herd-lite"], "composer": ["composer"],
    "erlang": ["erlang"], "java": ["java8"], "nodejs": ["nodejs", "nvm"],
    "bun": ["bun"], "docker": ["docker-engine", "docker-ubuntu", "docker-desktop"],
    "vscode": ["vscode"], "phpstorm": ["phpstorm"], "ansible": ["ansible"],
    "firefox": ["firefox-flatpak"],
}
SI = {
    "curl": ("curl", ["curl"]), "rclone": ("Rclone", ["rclone"]),
    "openvpn": ("OpenVPN", ["openvpn"]), "filezilla": ("FileZilla", ["filezilla"]),
    "brave": ("Brave", ["brave"]), "vlcmediaplayer": ("VLC media player", ["vlc"]),
    "qbittorrent": ("qbittorrent", ["qbittorrent"]),
    "flathub": ("Flathub", ["flathub"]), "zsh": ("Zsh", ["zsh"]),
    "stripe": ("Stripe", ["stripe-cli"]), "bruno": ("Bruno", ["bruno"]),
    "hoppscotch": ("Hoppscotch", ["hoppscotch"]), "warp": ("Warp", ["warp", "warp-style"]),
    "libreoffice": ("LibreOffice", ["libreoffice"]), "thunderbird": ("Thunderbird", ["thunderbird"]),
    "ente": ("Ente", ["ente-photos", "ente-auth"]), "bitwarden": ("Bitwarden", ["bitwarden"]),
    "keepassxc": ("KeePassXC", ["keepassxc"]), "yubico": ("Yubico", ["yubico-authenticator", "yubikey-manager"]),
    "cryptomator": ("Cryptomator", ["cryptomator"]), "torbrowser": ("Tor Browser", ["torbrowser-launcher"]),
    "wireshark": ("Wireshark", ["wireshark"]), "wireguard": ("WireGuard", ["wireguard-tools"]),
    "tmux": ("tmux", ["tmux"]), "virtualbox": ("VirtualBox", ["virtualbox"]),
}


def fetch(url):
    with urlopen(url, timeout=30) as response:
        return response.read()


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    simple = {i["title"]: i for i in json.loads(fetch(SIMPLE + "/data/simple-icons.json"))}
    devicon = {i["name"]: i for i in json.loads(fetch(DEVICON + "/devicon.json"))}
    jobs = []
    for name, items in DEV.items():
        assert "original" in devicon[name]["versions"]["svg"], name
        jobs.append(("dev-" + name, DEVICON + f"/icons/{name}/{name}-original.svg", items, None, {}))
    for slug, (title, items) in SI.items():
        meta = simple[title]
        jobs.append(("si-" + slug, SIMPLE + f"/icons/{slug}.svg", items, meta["hex"], meta))

    def download(job):
        name, url, items, color, meta = job
        raw = fetch(url)
        root = ET.fromstring(raw)
        assert root.tag == "{http://www.w3.org/2000/svg}svg", name
        for element in root.iter():
            tag = element.tag.rsplit("}", 1)[-1]
            assert tag not in {"script", "foreignObject", "image", "a"}, name
            for key, value in element.attrib.items():
                assert not key.lower().startswith("on"), name
                assert "href" not in key.lower() or value.startswith("#"), name
                assert "url(" not in value.lower() or value.strip().startswith("url(#"), name
        if color:
            # Simple Icons provides a monochrome shape; use its published brand color.
            raw = raw.replace(b"<svg ", f'<svg fill="#{color}" '.encode(), 1)
        (OUT / (name + ".svg")).write_bytes(raw)
        return dict(file=name + ".svg", items=items, download=url, upstream=meta,
                    sha256=hashlib.sha256(raw).hexdigest())

    records = list(ThreadPoolExecutor(max_workers=6).map(download, jobs))
    (OUT / "sources.json").write_text(json.dumps(records, ensure_ascii=False, indent=2) + "\n")
    mapping = {item: record["file"] for record in records for item in record["items"]}
    (ROOT / "frontend/src/appLogos.json").write_text(json.dumps(mapping, indent=2) + "\n")
    (OUT / "DEVICON-LICENSE.txt").write_bytes(fetch(DEVICON + "/LICENSE"))
    (OUT / "SIMPLE-ICONS-LICENSE.txt").write_bytes(fetch(SIMPLE + "/LICENSE.md"))
    print(f"Bundled {len(records)} logos covering {len(mapping)} catalog items.")


if __name__ == "__main__":
    main()
