# Pesquisa e decisões (2026-10-04)

Alvo: Ubuntu 26.04 LTS «resolute», GNOME Shell 50, amd64. Nenhum script de referência foi executado nem copiado.

## Projetos de referência

| Projeto | Licença / estado | Aproveitado (ideias) | Rejeitado |
|---|---|---|---|
| tprasadtp/ubuntu-post-install | GPL-3.0, **arquivado** (2023), focal–mantic | Configuração declarativa (pacotes/PPAs/tarefas) separada do código | Código (GPL, arquivado, sem 26.04) |
| franckferman/ubuntu-post-install | AGPL-3.0, ativo, só 22.04/24.04 | Passos idempotentes «saltar se feito», revisão antes de transferir, grupos de opções | Código AGPL; hardening/privacidade fora de âmbito; ciclo `sudo -v` em background |
| kronflux/ubuntu-post-install | MIT, ativo, foco Server 26.04 | Perfis + catálogo por release, remoção só por deteção (sem wildcard), backup antes de alterar, dry-run, documento de reversão | Purga de snap, hardening, DNS-over-TLS |
| kostiantyn-petlia/ubuntu-dev-env-setup | GPL-3.0, parado (2024), cheat sheet 24.04 | Checklist de workstation PHP, identidades Git com `includeIf`, permissões SSH, pós-instalação Docker | Apache/PhpStorm/Chrome específicos |

Nenhum código foi reutilizado; só ideias de arquitetura (licenças respeitadas).

## Receitas verificadas

Fingerprints obtidos com `gpg --show-keys --with-colons` sobre as chaves transferidas; `dists/<suite>/Release` testado por HTTP.

| Repositório | Suite | Fingerprints primários aceites | Validação real 26.04 |
|---|---|---|---|
| Docker | `resolute` | 9DC858229FC7DD38854AE2D88D81803C0EBFCD88 | ✅ docker-ce 29.8.2 instalado (container) |
| VS Code | `stable` | BC528686B50D79E339D3721CEB3E94ADBE1229CF | ✅ repositório + `apt update` |
| Brave | `stable` | DBF1A116…0257, 47D32A74…CD96, B2A3DCA3…DCA0 (keyring com 3 chaves) | ✅ repositório + `apt update` |
| GitHub CLI | `stable` | 7F38BBB5…3325 (atual), 2C610620…6059 (expirada, ainda no keyring) | ✅ gh 2.102 instalado |
| Stripe CLI | `stable` | 6681D7C3D103DAC65D79C25EDEEBD57F917C83E3 | ✅ stripe 1.53 instalado |
| Warp | `stable` (só Release, sem InRelease) | 0913165C78D5B7A41B42AC657FF7AB39D60F803F | ✅ repositório + `apt update` aceita |
| PPA ondrej/php | `resolute` | B8DC7E53946656EFBCE4C1DD71DAEAAB4AD4CAB6 (Launchpad API) | ❌ **não publica para 26.04**: o passo falha sem alterar nada |

Regra do helper: todas as chaves primárias do ficheiro transferido têm de estar na allowlist do catálogo; caso contrário nada é escrito.

### Outras decisões
- **PHP**: o arquivo 26.04 só tem 8.5 (php8.5-*); 8.3/8.4 via PPA (bloqueado até publicar) ou **php.new / Herd Lite** (pedido do utilizador): o script é transferido para ficheiro, o SHA-256 é registado, é recusado se contiver `sudo` ou hosts fora de `download.herdphp.com`, `curl.se`, `herd.laravel.com`, faz backup de `.bashrc/.zshrc/.profile` e corre como utilizador. Limitação: o instalador do fornecedor não verifica checksums dos binários.
- **Xdebug**: instalado e desligado com `phpdismod` (validado: `php -m` não o carrega).
- **Docker**: Engine oficial (recomendado) e `docker.io`+`docker-compose-v2` do Ubuntu são conflito bloqueante. O `docker-compose` v1 não existe no 26.04. Docker Desktop é manual (deb de ~465 MB, KVM, `pass`; coexiste com contexto `desktop-linux`). Grupo `docker` é opcional e explica que equivale a root.
- **Node**: APT (22 LTS) ou NVM (manual), com conflito-aviso entre ambos.
- **Extensões**: o Ubuntu 26.04 inclui `gnome-shell-ubuntu-extensions` (Ubuntu Dock, AppIndicators, Tiling Assistant, DING) — detetadas como «Já incluídas», nunca duplicadas. Instalação via API pública `extensions.gnome.org/extension-info/?uuid=…&shell_version=<versão real>`; todas as candidatas responderam 200 para GNOME 50. GSConnect e System Monitor vêm por APT. Dash to Panel é manual (exige desativar a Ubuntu Dock). Vitals e System Monitor são alternativas.
- **Energia**: o gestor padrão é power-profiles-daemon; TLP é desaconselhado (conflito) e não automatizado. zram/swappiness/sysctl não são alterados sem evidência.
- **Flathub**: IDs verificados (firefox, brave, freefilesync, dbeaver, bruno…). Hoppscotch e Warp não estão no Flathub → manual / repo oficial.
- **Firefox Flatpak**: Firefox ≥147 usa XDG para perfis novos (`~/.var/app/org.mozilla.firefox/config/mozilla/firefox`), perfis antigos ficam em `.mozilla/firefox`; ambos são sondados. O Brave não importa perfis Flatpak/Snap.
- **Wails**: v2.16.0 (estável, 2026-09-14); v3 ainda beta → rejeitado. Build com tag `webkit2_41` (webkit2gtk-4.0 não existe no 26.04).

## 0.3.0 — Infraestrutura, Segurança e perfis

- **Pacotes verificados no arquivo 26.04** (`apt-cache policy`, 2026-10-04): remmina (+ plugins rdp e vnc), virt-manager, qemu-system-x86, libvirt-daemon-system, wireshark, nmap, ncat, mtr-tiny, iperf3, tcpdump, traceroute, bind9-dnsutils, whois, ipcalc, ansible, wireguard-tools, putty-tools, minicom, tmux, keepassxc, nethogs, iftop, smartmontools, lnav, gparted, ethtool, nload, testdisk, gddrescue, gnome-boxes, sshfs, gufw, lynis, clamtk, torbrowser-launcher, seahorse, yubikey-manager, opensnitch e libreoffice-* (+ l10n pt/pt-br).
- O `thunderbird` do arquivo é um pacote de transição para o Snap, por isso foi usado o Flatpak `org.mozilla.Thunderbird`.
- **Flathub (API, 2026-10-04):**
  - Editor verificado: com.bitwarden.desktop, io.ente.auth, io.ente.photos, com.belmoussaoui.Authenticator, org.cryptomator.Cryptomator, org.gnome.World.Secrets, fr.romainvigier.MetadataCleaner, com.github.tchx84.Flatseal, org.mozilla.Thunderbird.
  - Editor **não verificado**: com.yubico.yubioath.
  - O Tor Browser launcher não está no Flathub, por isso foi usado o APT.
- **Instalação não interativa validada:** wireshark-common/install-setuid fica `false` (sem captura para utilizadores comuns) e o daemon do iperf3 não é ativado.
- **Rejeitados:**
  - rkhunter e chkrootkit (falsos positivos), firejail (complexidade) e libpam-u2f (altera o PAM/login).
  - VirtualBox (módulos e Secure Boot), Terraform/kubectl/Helm e AnyDesk/TeamViewer: repositórios de terceiros ainda não verificados, por isso ficam manuais.
- **php.new:** passou a exigir confirmação do SHA-256 do script revisto. As heurísticas (sudo, hosts) mantêm-se só como deteção de surpresas óbvias, não como garantia.
