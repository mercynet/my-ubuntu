<!-- Gerado do catálogo e dos perfis; não editar diretamente. -->

# Referência do catálogo e dos perfis

191 itens e 10 perfis incluídos. Esta referência descreve as receitas embutidas; o estado e a compatibilidade efetivos são detetados no computador. Os comandos de verificação são evidência usada pelo software, não instruções para aplicar a receita manualmente.

## Índice por área

### Aplicações

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Git](#item-git) | Essenciais | apt | Não |
| [git-flow](#item-git-flow) | Essenciais | apt | Não |
| [GitHub CLI](#item-gh) | Essenciais | apt-repo | Não |
| [build-essential](#item-build-essential) | Essenciais | apt | Não |
| [curl](#item-curl) | Essenciais | apt | Não |
| [Cliente SSH](#item-ssh-client) | Essenciais | apt | Não |
| [Servidor SSH](#item-ssh-server) | Rede e ficheiros | apt | Sim |
| [htop](#item-htop) | Utilitários | apt | Não |
| [ncdu](#item-ncdu) | Utilitários | apt | Não |
| [rsync](#item-rsync) | Rede e ficheiros | apt | Não |
| [rclone](#item-rclone) | Rede e ficheiros | apt | Não |
| [Partilhas Windows (SMB)](#item-smb) | Rede e ficheiros | apt | Não |
| [OpenVPN no NetworkManager](#item-openvpn) | Rede e ficheiros | apt | Não |
| [Déjà Dup (Cópias de segurança)](#item-deja-dup) | Rede e ficheiros | apt | Não |
| [FileZilla](#item-filezilla) | Rede e ficheiros | apt | Não |
| [FreeFileSync](#item-freefilesync) | Rede e ficheiros | flatpak | Não |
| [Firefox (Flatpak)](#item-firefox-flatpak) | Navegadores | flatpak | Não |
| [Brave](#item-brave) | Navegadores | apt-repo | Não |
| [Remover Firefox Snap](#item-firefox-snap-remove) | Navegadores | admin | Sim |
| [VLC](#item-vlc) | Multimédia | apt | Não |
| [qBittorrent](#item-qbittorrent) | Multimédia | apt | Não |
| [Flameshot](#item-flameshot) | Utilitários | apt | Não |
| [Solaar](#item-solaar) | Utilitários | apt | Não |
| [Pandoc + LaTeX](#item-pandoc) | Documentos | apt | Não |
| [Português (PT/BR)](#item-lang-pt) | Idiomas | apt | Não |
| [Inglês](#item-lang-en) | Idiomas | apt | Não |
| [Nala](#item-nala) | Sistema | apt | Não |
| [Flatpak + Flathub](#item-flathub) | Sistema | apt | Não |
| [GNOME Tweaks (Ajustes)](#item-gnome-tweaks) | Sistema | apt | Não |
| [Extension Manager](#item-extension-manager) | Sistema | apt | Não |
| [dconf Editor](#item-dconf-editor) | Sistema | apt | Sim |
| [LibreOffice](#item-libreoffice) | Escritório | apt | Não |
| [Thunderbird (Flatpak)](#item-thunderbird) | Escritório | flatpak | Não |
| [Ente Photos](#item-ente-photos) | Multimédia | flatpak | Não |

### Desenvolvimento

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Python (pip, venv, pipx)](#item-python-dev) | Linguagens | apt | Não |
| [Go](#item-golang) | Linguagens | apt | Não |
| [adb + scrcpy](#item-android-tools) | Linguagens | apt | Não |
| [SQLite](#item-sqlite) | Bases de dados | apt | Não |
| [Cliente MariaDB/MySQL](#item-mariadb-client) | Bases de dados | apt | Não |
| [Servidor MariaDB (local)](#item-mariadb-server) | Serviços locais · Dados e cache | apt | Não |
| [Supervisor](#item-supervisor) | Bases de dados | apt | Sim |
| [DBeaver Community](#item-dbeaver) | Bases de dados | flatpak | Não |
| [Zsh](#item-zsh) | Terminal | apt | Não |
| [Zsh: sugestões pelo histórico](#item-zsh-autosuggestions) | Plugins e complementos | apt | Não |
| [Zsh: destaque de sintaxe](#item-zsh-syntax-highlighting) | Plugins e complementos | apt | Não |
| [fzf — pesquisa no terminal](#item-fzf) | Plugins e complementos | apt | Não |
| [zoxide — navegar entre pastas](#item-zoxide) | Plugins e complementos | apt | Não |
| [tldr — exemplos de comandos](#item-tealdeer) | Terminal | apt | Não |
| [Oh My Zsh — configuração e plugins](#item-oh-my-zsh) | Plugins e complementos | manual | Não |
| [PHP (sistema, APT)](#item-php) | PHP | apt | Não |
| [PHP outras versões (PPA Ondřej)](#item-php-ppa) | PHP | apt-repo | Sim |
| [PHP via php.new (Herd Lite)](#item-php-herd-lite) | PHP | user-script | Não |
| [Composer](#item-composer) | PHP | apt | Não |
| [Xdebug (desativado por padrão)](#item-xdebug) | PHP | apt | Não |
| [Dependências de compilação PHP](#item-php-build-deps) | Avançado | apt | Sim |
| [Erlang](#item-erlang) | Avançado | apt | Sim |
| [Java 8 (OpenJDK)](#item-java8) | Avançado | apt | Sim |
| [wkhtmltopdf](#item-wkhtmltopdf) | Avançado | manual | Sim |
| [Mainline (kernels manuais)](#item-mainline) | Avançado | manual | Sim |
| [Node.js do sistema (APT)](#item-nodejs) | Node.js | apt | Não |
| [NVM (várias versões Node)](#item-nvm) | Node.js | manual | Não |
| [Bun](#item-bun) | Node.js | manual | Não |
| [Docker Engine (oficial)](#item-docker-engine) | Contentores | apt-repo | Não |
| [Docker (pacotes Ubuntu)](#item-docker-ubuntu) | Contentores | apt | Sim |
| [Docker Desktop](#item-docker-desktop) | Contentores | manual | Não |
| [Docker sem sudo (grupo docker)](#item-docker-group) | Contentores | admin | Não |
| [Visual Studio Code](#item-vscode) | Editores e APIs | apt-repo | Não |
| [PhpStorm](#item-phpstorm) | Editores e APIs | manual | Não |
| [Stripe CLI](#item-stripe-cli) | Editores e APIs | apt-repo | Não |
| [Bruno](#item-bruno) | Editores e APIs | flatpak | Não |
| [Hoppscotch Desktop](#item-hoppscotch) | Editores e APIs | manual | Não |
| [Warp](#item-warp) | Terminal | apt-repo | Não |
| [Devin Desktop](#item-devin) | Editores e APIs | manual | Não |
| [Git: predefinições úteis](#item-git-defaults) | Git | gitconfig | Não |
| [Git: identidade por diretório](#item-git-identity) | Git | manual | Não |
| [Mais watchers inotify](#item-inotify-watches) | Sistema | admin | Sim |
| [Claude Code](#item-claude-code) | Agentes de IA | npm-user | Não |
| [Codex](#item-codex-cli) | Agentes de IA | npm-user | Não |
| [Gemini CLI](#item-gemini-cli) | Agentes de IA | npm-user | Não |
| [OpenCode](#item-opencode) | Agentes de IA | npm-user | Não |
| [Pi](#item-pi-coding-agent) | Agentes de IA | npm-user | Não |
| [GitHub Copilot CLI](#item-github-copilot-cli) | Agentes de IA | npm-user | Não |
| [Claude Code: gerir plugins](#item-claude-code-plugins-guide) | Plugins e complementos | manual | Não |
| [Codex: skills e plugins](#item-codex-skills-guide) | Plugins e complementos | manual | Não |
| [OpenCode: plugins e configuração](#item-opencode-plugins-guide) | Plugins e complementos | manual | Não |
| [Pi: pacotes, extensões e skills](#item-pi-packages-guide) | Plugins e complementos | manual | Não |
| [Gemini CLI: gerir extensões](#item-gemini-extensions-guide) | Plugins e complementos | manual | Não |
| [Apache HTTP Server (local)](#item-apache2) | Serviços locais · Web | apt | Não |
| [Nginx (local)](#item-nginx) | Serviços locais · Web | apt | Não |
| [Redis (servidor local)](#item-redis-server) | Serviços locais · Dados e cache | apt | Não |
| [RabbitMQ (local)](#item-rabbitmq-server) | Serviços locais · Mensageria | apt | Não |
| [PostgreSQL (servidor local)](#item-postgresql-server) | Serviços locais · Dados e cache | apt | Não |
| [Memcached (local)](#item-memcached) | Serviços locais · Dados e cache | apt | Não |
| [Apache Kafka — laboratório Docker](#item-kafka-local-guide) | Serviços locais · Mensageria | manual | Não |

### Infraestrutura

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Remmina (RDP/VNC)](#item-remmina) | Acesso remoto | apt | Não |
| [PuTTY tools](#item-putty-tools) | Acesso remoto | apt | Não |
| [SSHFS](#item-sshfs) | Acesso remoto | apt | Não |
| [tmux](#item-tmux) | Acesso remoto | apt | Não |
| [Diagnóstico de rede](#item-net-diag) | Rede e diagnóstico | apt | Não |
| [Nmap + Ncat](#item-nmap) | Rede e diagnóstico | apt | Não |
| [iperf3](#item-iperf3) | Rede e diagnóstico | apt | Não |
| [Monitores de tráfego](#item-traffic-monitor) | Rede e diagnóstico | apt | Não |
| [tcpdump](#item-tcpdump) | Rede e diagnóstico | apt | Não |
| [Wireshark](#item-wireshark) | Rede e diagnóstico | apt | Não |
| [Virtualização (virt-manager/KVM)](#item-virt-manager) | Virtualização | apt | Não |
| [GNOME Boxes](#item-gnome-boxes) | Virtualização | apt | Não |
| [Ansible](#item-ansible) | Automação | apt | Não |
| [WireGuard (ferramentas)](#item-wireguard-tools) | Automação | apt | Não |
| [Discos e recuperação](#item-disk-tools) | Discos | apt | Não |
| [lnav](#item-lnav) | Discos | apt | Não |
| [minicom](#item-minicom) | Acesso remoto | apt | Sim |
| [VirtualBox](#item-virtualbox) | Virtualização | manual | Não |
| [Terraform / kubectl / Helm](#item-terraform-kubectl) | Automação | manual | Não |
| [AnyDesk / TeamViewer / RustDesk](#item-remote-support) | Acesso remoto | manual | Não |

### Segurança

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Bitwarden](#item-bitwarden) | Gestores de senhas | flatpak | Não |
| [KeePassXC](#item-keepassxc) | Gestores de senhas | apt | Não |
| [Secrets (GNOME)](#item-gnome-secrets) | Gestores de senhas | flatpak | Não |
| [Senhas e chaves](#item-seahorse) | Gestores de senhas | apt | Não |
| [Ente Auth](#item-ente-auth) | Autenticação 2FA | flatpak | Não |
| [Authenticator (GNOME)](#item-gnome-authenticator) | Autenticação 2FA | flatpak | Não |
| [Yubico Authenticator](#item-yubico-authenticator) | Autenticação 2FA | flatpak | Não |
| [YubiKey Manager (ykman)](#item-yubikey-manager) | Autenticação 2FA | apt | Sim |
| [Cryptomator](#item-cryptomator) | Cifragem | flatpak | Não |
| [Metadata Cleaner](#item-metadata-cleaner) | Privacidade | flatpak | Não |
| [Flatseal](#item-flatseal) | Privacidade | flatpak | Não |
| [Tor Browser (launcher)](#item-torbrowser-launcher) | Privacidade | apt | Não |
| [Gufw (interface da firewall)](#item-gufw) | Firewall e auditoria | apt | Não |
| [Lynis](#item-lynis) | Firewall e auditoria | apt | Não |
| [ClamAV + ClamTk](#item-clamtk) | Firewall e auditoria | apt | Sim |
| [OpenSnitch](#item-opensnitch) | Firewall e auditoria | apt | Sim |
| [Gobuster](#item-gobuster) | Auditoria web | apt | Não |
| [sqlmap](#item-sqlmap) | Auditoria web | apt | Não |
| [SSLScan](#item-sslscan) | TLS e rede | apt | Não |
| [testssl.sh](#item-testssl) | TLS e rede | apt | Não |
| [arp-scan](#item-arp-scan) | TLS e rede | apt | Não |

### Estilo

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Tema claro/escuro](#item-color-scheme) | Cores | gsettings | Não |
| [Fonte monoespaçada](#item-mono-font) | Fontes | gsettings | Não |
| [Cor de destaque](#item-accent-color) | Cores | gsettings | Não |
| [Coleção de wallpapers Ubuntu](#item-wallpapers-pack) | Wallpapers | apt | Não |
| [Wallpaper](#item-wallpaper) | Wallpapers | gsettings | Não |
| [Terminal: transparência (Ptyxis)](#item-ptyxis-opacity) | Terminal | gsettings | Não |
| [Terminal: transparência (GNOME Terminal)](#item-gnome-terminal-transparency) | Terminal | gsettings | Sim |
| [Warp: transparência e temas](#item-warp-style) | Terminal | manual | Não |
| [Dock: transparente](#item-dock-transparency) | Painel e dock | gsettings | Não |
| [Tamanho do cursor](#item-cursor-size) | Cursor e fontes | gsettings | Não |
| [Fonte da interface](#item-ui-font) | Cursor e fontes | gsettings | Não |

### Desktop

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Tamanho do texto](#item-text-scaling) | Acessibilidade | gsettings | Não |
| [Reduzir animações](#item-reduce-animations) | Acessibilidade | gsettings | Não |
| [Luz noturna](#item-night-light) | Acessibilidade | gsettings | Não |
| [Percentagem da bateria](#item-battery-pct) | Barra superior | gsettings | Não |
| [Dia da semana no relógio](#item-clock-weekday) | Barra superior | gsettings | Não |
| [Dock: posição](#item-dock-position) | Dock | gsettings | Não |
| [Dock: tamanho dos ícones](#item-dock-size) | Dock | gsettings | Não |
| [Dock: ocultar automaticamente](#item-dock-autohide) | Dock | gsettings | Não |
| [Dock: clique minimiza](#item-dock-click) | Dock | gsettings | Não |
| [Dock: em todos os monitores](#item-dock-multimonitor) | Dock | gsettings | Não |
| [Botões minimizar/maximizar](#item-window-buttons) | Janelas | gsettings | Não |
| [Centrar novas janelas](#item-center-windows) | Janelas | gsettings | Não |
| [Áreas de trabalho em todos os monitores](#item-workspaces-all) | Janelas | gsettings | Não |
| [Super+E abre a pasta pessoal](#item-super-e) | Atalhos | gsettings | Não |
| [Tocar para clicar (touchpad)](#item-tap-click) | Atalhos | gsettings | Não |
| [Ficheiros: vista em lista](#item-nautilus-list) | Ficheiros | gsettings | Não |
| [Ficheiros: pastas primeiro](#item-folders-first) | Ficheiros | gsettings | Não |
| [Notificações fora do ecrã de bloqueio](#item-notif-lock) | Notificações | gsettings | Não |
| [Modelo: Documento de texto.txt](#item-tpl-txt) | Novo documento | template | Não |
| [Modelo: Nota.md](#item-tpl-md) | Novo documento | template | Não |
| [Modelo: Script.sh](#item-tpl-sh) | Novo documento | template | Não |
| [Modelo: Classe.php](#item-tpl-php) | Novo documento | template | Não |
| [Aplicações de arranque](#item-startup-apps) | Arranque | manual | Não |

### Extensões

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Dash to Dock](#item-ext-dash-to-dock) | Já incluídas no Ubuntu | extension | Não |
| [AppIndicator](#item-ext-appindicator) | Já incluídas no Ubuntu | extension | Não |
| [Tiling Assistant](#item-ext-tiling) | Já incluídas no Ubuntu | extension | Não |
| [Clipboard Indicator](#item-ext-clipboard) | Produtividade | extension | Não |
| [Caffeine](#item-ext-caffeine) | Produtividade | extension | Não |
| [Just Perfection](#item-ext-just-perfection) | Interface | extension | Sim |
| [Vitals](#item-ext-vitals) | Monitorização | extension | Não |
| [System Monitor (GNOME oficial)](#item-ext-system-monitor) | Monitorização | apt | Não |
| [GSConnect](#item-ext-gsconnect) | Integração | apt | Não |
| [Dash to Panel](#item-ext-dash-to-panel) | Interface | manual | Sim |

### Energia e hardware

| Item | Categoria | Método | Avançado |
| --- | --- | --- | --- |
| [Energia: equilibrado](#item-power-balanced) | Perfil de energia | power-profile | Não |
| [Energia: poupança](#item-power-saver) | Perfil de energia | power-profile | Não |
| [Energia: desempenho](#item-power-performance) | Perfil de energia | power-profile | Não |
| [Ecrã: apagar após inatividade](#item-screen-blank) | Ecrã e suspensão | gsettings | Não |
| [Não suspender ligado à corrente](#item-suspend-ac) | Ecrã e suspensão | gsettings | Não |
| [Suspender na bateria após 15 min](#item-suspend-battery) | Ecrã e suspensão | gsettings | Não |
| [Comportamento da tampa](#item-lid-close) | Ecrã e suspensão | manual | Não |
| [TRIM semanal de SSD](#item-fstrim-timer) | Disco | admin | Não |
| [Limitar logs do sistema](#item-journald-limit) | Disco | admin | Não |
| [Drivers adicionais recomendados](#item-ubuntu-drivers) | Hardware | admin | Sim |
| [TLP](#item-tlp) | Avançado | manual | Sim |
| [zram / swappiness](#item-zram) | Avançado | manual | Sim |

## Perfis incluídos

### Uso pessoal / escritório

Escritório, email, navegadores, multimédia, backups e um gestor de senhas.

ID: `uso-pessoal`. Tipo: `usage`. Para quem usa o PC para documentos, email, navegação, fotos e chamadas.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [curl](#item-curl) | {} |
| [Flatpak + Flathub](#item-flathub) | {} |
| [Português (PT/BR)](#item-lang-pt) | {} |
| [Déjà Dup (Cópias de segurança)](#item-deja-dup) | {} |
| [GNOME Tweaks (Ajustes)](#item-gnome-tweaks) | {} |
| [LibreOffice](#item-libreoffice) | {} |
| [Thunderbird (Flatpak)](#item-thunderbird) | {} |
| [Firefox (Flatpak)](#item-firefox-flatpak) | {} |
| [VLC](#item-vlc) | {} |
| [Flameshot](#item-flameshot) | {} |
| [Bitwarden](#item-bitwarden) | {} |
| [Ente Auth](#item-ente-auth) | {} |
| [Partilhas Windows (SMB)](#item-smb) | {} |
| [Tema claro/escuro](#item-color-scheme) | {} |
| [Luz noturna](#item-night-light) | {} |
| [Dia da semana no relógio](#item-clock-weekday) | {} |
| [Modelo: Documento de texto.txt](#item-tpl-txt) | {} |

### Desenvolvimento PHP/Laravel

Docker oficial, PHP 8.5 do sistema com Xdebug desligado, Composer, Node, VS Code, Git e clientes de bases de dados.

ID: `dev-php-laravel`. Tipo: `usage`. Para quem desenvolve em PHP/Laravel com Docker, Node e Git.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [curl](#item-curl) | {} |
| [Flatpak + Flathub](#item-flathub) | {} |
| [Português (PT/BR)](#item-lang-pt) | {} |
| [Déjà Dup (Cópias de segurança)](#item-deja-dup) | {} |
| [GNOME Tweaks (Ajustes)](#item-gnome-tweaks) | {} |
| [Git](#item-git) | {} |
| [git-flow](#item-git-flow) | {} |
| [GitHub CLI](#item-gh) | {} |
| [build-essential](#item-build-essential) | {} |
| [Cliente SSH](#item-ssh-client) | {} |
| [htop](#item-htop) | {} |
| [ncdu](#item-ncdu) | {} |
| [rsync](#item-rsync) | {} |
| [Docker Engine (oficial)](#item-docker-engine) | {} |
| [PHP (sistema, APT)](#item-php) | {} |
| [Composer](#item-composer) | {} |
| [Xdebug (desativado por padrão)](#item-xdebug) | {} |
| [Node.js do sistema (APT)](#item-nodejs) | {} |
| [SQLite](#item-sqlite) | {} |
| [Cliente MariaDB/MySQL](#item-mariadb-client) | {} |
| [Visual Studio Code](#item-vscode) | {} |
| [DBeaver Community](#item-dbeaver) | {} |
| [Git: predefinições úteis](#item-git-defaults) | {} |
| [Modelo: Classe.php](#item-tpl-php) | {} |
| [Modelo: Script.sh](#item-tpl-sh) | {} |
| [Python (pip, venv, pipx)](#item-python-dev) | {} |
| [KeePassXC](#item-keepassxc) | {} |

### Infraestrutura — estação de trabalho

Ferramentas de cliente apenas: não configura servidores, firewall, DNS, VPN nem serviços. Validado em simulação/container; não validado ainda em máquinas de clientes.

ID: `infra-estacao`. Tipo: `usage`. Para técnicos de infraestrutura: acesso remoto, rede, diagnóstico e virtualização num Ubuntu Desktop.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [curl](#item-curl) | {} |
| [Flatpak + Flathub](#item-flathub) | {} |
| [Português (PT/BR)](#item-lang-pt) | {} |
| [Déjà Dup (Cópias de segurança)](#item-deja-dup) | {} |
| [GNOME Tweaks (Ajustes)](#item-gnome-tweaks) | {} |
| [Git](#item-git) | {} |
| [Cliente SSH](#item-ssh-client) | {} |
| [htop](#item-htop) | {} |
| [ncdu](#item-ncdu) | {} |
| [rsync](#item-rsync) | {} |
| [rclone](#item-rclone) | {} |
| [Remmina (RDP/VNC)](#item-remmina) | {} |
| [PuTTY tools](#item-putty-tools) | {} |
| [SSHFS](#item-sshfs) | {} |
| [tmux](#item-tmux) | {} |
| [Diagnóstico de rede](#item-net-diag) | {} |
| [Nmap + Ncat](#item-nmap) | {} |
| [iperf3](#item-iperf3) | {} |
| [Monitores de tráfego](#item-traffic-monitor) | {} |
| [tcpdump](#item-tcpdump) | {} |
| [Wireshark](#item-wireshark) | {} |
| [Virtualização (virt-manager/KVM)](#item-virt-manager) | {} |
| [Ansible](#item-ansible) | {} |
| [WireGuard (ferramentas)](#item-wireguard-tools) | {} |
| [Discos e recuperação](#item-disk-tools) | {} |
| [lnav](#item-lnav) | {} |
| [KeePassXC](#item-keepassxc) | {} |
| [OpenVPN no NetworkManager](#item-openvpn) | {} |
| [FileZilla](#item-filezilla) | {} |

### Personalizado (começar vazio)

Sem itens: escolha tudo nas áreas.

ID: `personalizado`. Tipo: `usage`. Para quem quer escolher item a item.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

Começa com escolhas vazias.

### Base essencial

Ferramentas base, Flathub, idiomas e backups.

ID: `essencial`. Tipo: `complement`. Complemento útil para qualquer uso.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [curl](#item-curl) | {} |
| [Flatpak + Flathub](#item-flathub) | {} |
| [Português (PT/BR)](#item-lang-pt) | {} |
| [Déjà Dup (Cópias de segurança)](#item-deja-dup) | {} |
| [GNOME Tweaks (Ajustes)](#item-gnome-tweaks) | {} |
| [Git](#item-git) | {} |
| [build-essential](#item-build-essential) | {} |
| [Cliente SSH](#item-ssh-client) | {} |
| [htop](#item-htop) | {} |
| [ncdu](#item-ncdu) | {} |
| [rsync](#item-rsync) | {} |
| [Extension Manager](#item-extension-manager) | {} |
| [Inglês](#item-lang-en) | {} |

### Desktop confortável

Dock, janelas, ficheiros, notificações e algumas extensões úteis.

ID: `desktop-confortavel`. Tipo: `complement`. Complemento de aparência e produtividade do GNOME.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [Tema claro/escuro](#item-color-scheme) | {} |
| [Cor de destaque](#item-accent-color) | {} |
| [Dock: posição](#item-dock-position) | {} |
| [Dock: tamanho dos ícones](#item-dock-size) | {} |
| [Dock: ocultar automaticamente](#item-dock-autohide) | {} |
| [Dock: clique minimiza](#item-dock-click) | {} |
| [Dock: em todos os monitores](#item-dock-multimonitor) | {} |
| [Centrar novas janelas](#item-center-windows) | {} |
| [Áreas de trabalho em todos os monitores](#item-workspaces-all) | {} |
| [Super+E abre a pasta pessoal](#item-super-e) | {} |
| [Ficheiros: vista em lista](#item-nautilus-list) | {} |
| [Ficheiros: pastas primeiro](#item-folders-first) | {} |
| [Notificações fora do ecrã de bloqueio](#item-notif-lock) | {} |
| [Luz noturna](#item-night-light) | {} |
| [Dia da semana no relógio](#item-clock-weekday) | {} |
| [Modelo: Documento de texto.txt](#item-tpl-txt) | {} |
| [Modelo: Nota.md](#item-tpl-md) | {} |
| [Caffeine](#item-ext-caffeine) | {} |
| [Clipboard Indicator](#item-ext-clipboard) | {} |

### Terminal produtivo

Nala, Zsh/Oh My Zsh com sugestões e destaque de sintaxe, pesquisa com fzf, navegação com zoxide e exemplos com tldr. A instalação dos pacotes é automática; a integração com a configuração pessoal tem passos guiados.

ID: `terminal-produtivo`. Tipo: `complement`. Para quem usa o terminal e quer reduzir tarefas repetitivas.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [Nala](#item-nala) | {} |
| [Zsh: sugestões pelo histórico](#item-zsh-autosuggestions) | {} |
| [Zsh: destaque de sintaxe](#item-zsh-syntax-highlighting) | {} |
| [fzf — pesquisa no terminal](#item-fzf) | {} |
| [zoxide — navegar entre pastas](#item-zoxide) | {} |
| [tldr — exemplos de comandos](#item-tealdeer) | {} |
| [Oh My Zsh — configuração e plugins](#item-oh-my-zsh) | {} |

### Segurança — laboratório e auditoria

Ferramentas open source para análise de rede, auditoria web, TLS e verificação local. A aplicação apenas instala as ferramentas; os testes são executados manualmente.

ID: `seguranca-laboratorio`. Tipo: `complement`. Para quem desenvolve aplicações, estuda segurança ou audita sistemas e redes.

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [Nmap + Ncat](#item-nmap) | {} |
| [Wireshark](#item-wireshark) | {} |
| [tcpdump](#item-tcpdump) | {} |
| [Lynis](#item-lynis) | {} |
| [Gobuster](#item-gobuster) | {} |
| [sqlmap](#item-sqlmap) | {} |
| [SSLScan](#item-sslscan) | {} |
| [testssl.sh](#item-testssl) | {} |
| [arp-scan](#item-arp-scan) | {} |

### Portátil equilibrado

Perfil equilibrado, suspensão em bateria e indicadores úteis.

ID: `portatil-equilibrado`. Tipo: `complement`. Complemento de energia para portáteis (alternativo ao Desempenho).

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [Energia: equilibrado](#item-power-balanced) | {} |
| [Suspender na bateria após 15 min](#item-suspend-battery) | {} |
| [Ecrã: apagar após inatividade](#item-screen-blank) | {} |
| [Percentagem da bateria](#item-battery-pct) | {} |
| [Tocar para clicar (touchpad)](#item-tap-click) | {} |
| [TRIM semanal de SSD](#item-fstrim-timer) | {} |

### Desempenho ligado à corrente

Perfil de desempenho e sem suspensão com carregador. Ganho não medido.

ID: `desempenho-corrente`. Tipo: `complement`. Complemento de energia para desktops sempre ligados (alternativo ao Portátil).

Requisitos: {"minRelease": "26.04", "arch": ["amd64"]}

| Item | Opções explícitas |
| --- | --- |
| [Energia: desempenho](#item-power-performance) | {} |
| [Não suspender ligado à corrente](#item-suspend-ac) | {} |

## Repositórios APT

### Docker

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "docker",
  "name": "Docker",
  "keyUrl": "https://download.docker.com/linux/ubuntu/gpg",
  "fingerprints": [
    "9DC858229FC7DD38854AE2D88D81803C0EBFCD88"
  ],
  "url": "https://download.docker.com/linux/ubuntu",
  "suite": "{codename}",
  "components": [
    "stable"
  ],
  "arch": "amd64",
  "docs": "https://docs.docker.com/engine/install/ubuntu/"
}
```

### Visual Studio Code

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "vscode",
  "name": "Visual Studio Code",
  "keyUrl": "https://packages.microsoft.com/keys/microsoft.asc",
  "fingerprints": [
    "BC528686B50D79E339D3721CEB3E94ADBE1229CF"
  ],
  "url": "https://packages.microsoft.com/repos/code",
  "suite": "stable",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://code.visualstudio.com/docs/setup/linux"
}
```

### Brave

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "brave",
  "name": "Brave",
  "keyUrl": "https://brave-browser-apt-release.s3.brave.com/brave-browser-archive-keyring.gpg",
  "fingerprints": [
    "DBF1A116C220B8C7164F98230686B78420038257",
    "47D32A74E9A9E013A4B4926C68D513D36A73CD96",
    "B2A3DCA350E67256740DF904DE4EC67BE4B0DCA0"
  ],
  "url": "https://brave-browser-apt-release.s3.brave.com",
  "suite": "stable",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://brave.com/linux/"
}
```

### GitHub CLI

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "github-cli",
  "name": "GitHub CLI",
  "keyUrl": "https://cli.github.com/packages/githubcli-archive-keyring.gpg",
  "fingerprints": [
    "7F38BBB59D064DBCB3D84D725612B36462313325",
    "2C6106201985B60E6C7AC87323F3D4EA75716059"
  ],
  "url": "https://cli.github.com/packages",
  "suite": "stable",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://github.com/cli/cli/blob/trunk/docs/install_linux.md"
}
```

### Stripe CLI

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "stripe",
  "name": "Stripe CLI",
  "keyUrl": "https://packages.stripe.dev/api/security/keypair/stripe-cli-gpg/public",
  "fingerprints": [
    "6681D7C3D103DAC65D79C25EDEEBD57F917C83E3"
  ],
  "url": "https://packages.stripe.dev/stripe-cli-debian-local",
  "suite": "stable",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://docs.stripe.com/stripe-cli"
}
```

### Warp

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "warp",
  "name": "Warp",
  "keyUrl": "https://releases.warp.dev/linux/keys/warp.asc",
  "fingerprints": [
    "0913165C78D5B7A41B42AC657FF7AB39D60F803F"
  ],
  "url": "https://releases.warp.dev/linux/deb",
  "suite": "stable",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://docs.warp.dev/getting-started/getting-started-with-warp"
}
```

### PPA ondrej/php

Receita completa, incluindo URL da chave, fingerprints permitidos e documentação do fornecedor:

```json
{
  "id": "ondrej-php",
  "name": "PPA ondrej/php",
  "keyUrl": "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0xB8DC7E53946656EFBCE4C1DD71DAEAAB4AD4CAB6",
  "fingerprints": [
    "B8DC7E53946656EFBCE4C1DD71DAEAAB4AD4CAB6"
  ],
  "url": "https://ppa.launchpadcontent.net/ondrej/php/ubuntu",
  "suite": "{codename}",
  "components": [
    "main"
  ],
  "arch": "amd64",
  "docs": "https://launchpad.net/~ondrej/+archive/ubuntu/php"
}
```

## Grupos exclusivos

### Perfil de energia

Perfil do power-profiles-daemon, o gestor de energia padrão do Ubuntu. Escolha um: equilibrado (padrão), poupança (bateria) ou desempenho (ligado à corrente; ganho não medido).

[Energia: equilibrado](#item-power-balanced), [Energia: poupança](#item-power-saver), [Energia: desempenho](#item-power-performance)

### Monitor de recursos no painel

CPU, memória, rede e temperaturas no painel superior. System Monitor é o oficial do GNOME (APT); Vitals tem mais sensores (extensions.gnome.org). Escolha um.

[Vitals](#item-ext-vitals), [System Monitor (GNOME oficial)](#item-ext-system-monitor)

## Receitas dos itens

<a id="item-git"></a>

### Git

Controlo de versões distribuído.

Base de qualquer projeto; necessário para Composer, NVM e clonagem de repositórios.

| Campo | Valor |
| --- | --- |
| ID | git |
| Área / categoria | Aplicações / Essenciais |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | git |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "git",
  "name": "Git",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Controlo de versões distribuído.",
  "why": "Base de qualquer projeto; necessário para Composer, NVM e clonagem de repositórios.",
  "method": "apt",
  "verify": {
    "cmd": [
      "git",
      "--version"
    ]
  },
  "packages": [
    "git"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-git-flow"></a>

### git-flow

Extensão Git para o modelo de ramos git-flow.

Útil se os seus projetos usam feature/release/hotfix branches.

| Campo | Valor |
| --- | --- |
| ID | git-flow |
| Área / categoria | Aplicações / Essenciais |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | git-flow |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Dependências | [Git](#item-git) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "git-flow",
  "name": "git-flow",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Extensão Git para o modelo de ramos git-flow.",
  "why": "Útil se os seus projetos usam feature/release/hotfix branches.",
  "method": "apt",
  "packages": [
    "git-flow"
  ],
  "requires": [
    "git"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gh"></a>

### GitHub CLI

Cliente oficial do GitHub na linha de comandos.

PRs, issues e autenticação Git sem sair do terminal. O repositório oficial tem versões recentes (o arquivo tem 2.46, antiga).

| Campo | Valor |
| --- | --- |
| ID | gh |
| Área / categoria | Aplicações / Essenciais |
| Método | apt-repo |
| Fonte | Repositório oficial GitHub CLI |
| Pacotes | gh |
| Repositório | github-cli |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gh",
  "name": "GitHub CLI",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Cliente oficial do GitHub na linha de comandos.",
  "why": "PRs, issues e autenticação Git sem sair do terminal. O repositório oficial tem versões recentes (o arquivo tem 2.46, antiga).",
  "method": "apt-repo",
  "verify": {
    "cmd": [
      "gh",
      "--version"
    ]
  },
  "packages": [
    "gh"
  ],
  "repo": "github-cli",
  "source": "Repositório oficial GitHub CLI",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-build-essential"></a>

### build-essential

Compiladores GCC, make e cabeçalhos.

Necessário para compilar extensões nativas (npm, PECL, pip).

| Campo | Valor |
| --- | --- |
| ID | build-essential |
| Área / categoria | Aplicações / Essenciais |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | build-essential |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "build-essential",
  "name": "build-essential",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Compiladores GCC, make e cabeçalhos.",
  "why": "Necessário para compilar extensões nativas (npm, PECL, pip).",
  "method": "apt",
  "packages": [
    "build-essential"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-curl"></a>

### curl

Transferências HTTP na linha de comandos.

Testar APIs e transferir ficheiros.

| Campo | Valor |
| --- | --- |
| ID | curl |
| Área / categoria | Aplicações / Essenciais |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | curl |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "curl",
  "name": "curl",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Transferências HTTP na linha de comandos.",
  "why": "Testar APIs e transferir ficheiros.",
  "method": "apt",
  "packages": [
    "curl"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ssh-client"></a>

### Cliente SSH

Ligações SSH e chaves para Git/servidores.

Deploy, servidores remotos e autenticação Git por chave.

| Campo | Valor |
| --- | --- |
| ID | ssh-client |
| Área / categoria | Aplicações / Essenciais |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | openssh-client |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ssh-client",
  "name": "Cliente SSH",
  "area": "apps",
  "category": "Essenciais",
  "summary": "Ligações SSH e chaves para Git/servidores.",
  "why": "Deploy, servidores remotos e autenticação Git por chave.",
  "method": "apt",
  "packages": [
    "openssh-client"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ssh-server"></a>

### Servidor SSH

Permite aceder a esta máquina por SSH.

Só se precisar de entrar remotamente neste PC. Abre a porta 22.

| Campo | Valor |
| --- | --- |
| ID | ssh-server |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | openssh-server |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Atenção | Abre a porta 22 (acesso remoto a esta máquina) |
| Avançado | Sim |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ssh-server",
  "name": "Servidor SSH",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Permite aceder a esta máquina por SSH.",
  "why": "Só se precisar de entrar remotamente neste PC. Abre a porta 22.",
  "method": "apt",
  "packages": [
    "openssh-server"
  ],
  "service": true,
  "advanced": true,
  "sensitive": "Abre a porta 22 (acesso remoto a esta máquina)",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-htop"></a>

### htop

Monitor de processos interativo.

Ver rapidamente o que está a consumir CPU/RAM.

| Campo | Valor |
| --- | --- |
| ID | htop |
| Área / categoria | Aplicações / Utilitários |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | htop |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "htop",
  "name": "htop",
  "area": "apps",
  "category": "Utilitários",
  "summary": "Monitor de processos interativo.",
  "why": "Ver rapidamente o que está a consumir CPU/RAM.",
  "method": "apt",
  "packages": [
    "htop"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ncdu"></a>

### ncdu

Analisador de espaço em disco no terminal.

Descobrir o que ocupa espaço antes de apagar algo.

| Campo | Valor |
| --- | --- |
| ID | ncdu |
| Área / categoria | Aplicações / Utilitários |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | ncdu |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ncdu",
  "name": "ncdu",
  "area": "apps",
  "category": "Utilitários",
  "summary": "Analisador de espaço em disco no terminal.",
  "why": "Descobrir o que ocupa espaço antes de apagar algo.",
  "method": "apt",
  "packages": [
    "ncdu"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-rsync"></a>

### rsync

Cópia/sincronização incremental de ficheiros.

Backups e deploys eficientes.

| Campo | Valor |
| --- | --- |
| ID | rsync |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | rsync |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "rsync",
  "name": "rsync",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Cópia/sincronização incremental de ficheiros.",
  "why": "Backups e deploys eficientes.",
  "method": "apt",
  "packages": [
    "rsync"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-rclone"></a>

### rclone

Sincronização com armazenamento cloud (S3, Drive…).

Copiar backups para cloud sem clientes proprietários.

| Campo | Valor |
| --- | --- |
| ID | rclone |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | rclone |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "rclone",
  "name": "rclone",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Sincronização com armazenamento cloud (S3, Drive…).",
  "why": "Copiar backups para cloud sem clientes proprietários.",
  "method": "apt",
  "packages": [
    "rclone"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-smb"></a>

### Partilhas Windows (SMB)

Montar e aceder a partilhas SMB/CIFS.

NAS e partilhas de rede da empresa.

| Campo | Valor |
| --- | --- |
| ID | smb |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | cifs-utils, smbclient |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "smb",
  "name": "Partilhas Windows (SMB)",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Montar e aceder a partilhas SMB/CIFS.",
  "why": "NAS e partilhas de rede da empresa.",
  "method": "apt",
  "packages": [
    "cifs-utils",
    "smbclient"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-openvpn"></a>

### OpenVPN no NetworkManager

Importar perfis .ovpn nas definições de rede.

VPN da empresa pelo painel do GNOME. Não importa perfis nem certificados automaticamente.

| Campo | Valor |
| --- | --- |
| ID | openvpn |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | network-manager-openvpn-gnome |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "openvpn",
  "name": "OpenVPN no NetworkManager",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Importar perfis .ovpn nas definições de rede.",
  "why": "VPN da empresa pelo painel do GNOME. Não importa perfis nem certificados automaticamente.",
  "method": "apt",
  "packages": [
    "network-manager-openvpn-gnome"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-deja-dup"></a>

### Déjà Dup (Cópias de segurança)

Backups agendados e simples da pasta pessoal.

Ferramenta oficial do GNOME para backups regulares e restauro de ficheiros.

| Campo | Valor |
| --- | --- |
| ID | deja-dup |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | deja-dup |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "deja-dup",
  "name": "Déjà Dup (Cópias de segurança)",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Backups agendados e simples da pasta pessoal.",
  "why": "Ferramenta oficial do GNOME para backups regulares e restauro de ficheiros.",
  "method": "apt",
  "packages": [
    "deja-dup"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-filezilla"></a>

### FileZilla

Cliente FTP/SFTP gráfico.

Transferências para alojamentos que só oferecem FTP/SFTP.

| Campo | Valor |
| --- | --- |
| ID | filezilla |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | filezilla |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "filezilla",
  "name": "FileZilla",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Cliente FTP/SFTP gráfico.",
  "why": "Transferências para alojamentos que só oferecem FTP/SFTP.",
  "method": "apt",
  "packages": [
    "filezilla"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-freefilesync"></a>

### FreeFileSync

Comparação e sincronização de pastas.

Espelhar pastas/discos com pré-visualização das diferenças.

| Campo | Valor |
| --- | --- |
| ID | freefilesync |
| Área / categoria | Aplicações / Rede e ficheiros |
| Método | flatpak |
| Fonte | Flathub (editor não verificado) |
| Flatpak | org.freefilesync.FreeFileSync |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "freefilesync",
  "name": "FreeFileSync",
  "area": "apps",
  "category": "Rede e ficheiros",
  "summary": "Comparação e sincronização de pastas.",
  "why": "Espelhar pastas/discos com pré-visualização das diferenças.",
  "method": "flatpak",
  "flatpak": "org.freefilesync.FreeFileSync",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor não verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-firefox-flatpak"></a>

### Firefox (Flatpak)

Firefox empacotado pela Mozilla no Flathub.

Atualizações diretas da Mozilla. Pode coexistir com o Snap até validar o novo perfil (ver Migração).

| Campo | Valor |
| --- | --- |
| ID | firefox-flatpak |
| Área / categoria | Aplicações / Navegadores |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | org.mozilla.firefox |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "firefox-flatpak",
  "name": "Firefox (Flatpak)",
  "area": "apps",
  "category": "Navegadores",
  "summary": "Firefox empacotado pela Mozilla no Flathub.",
  "why": "Atualizações diretas da Mozilla. Pode coexistir com o Snap até validar o novo perfil (ver Migração).",
  "method": "flatpak",
  "flatpak": "org.mozilla.firefox",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "tags": [
    "snap:firefox"
  ],
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-brave"></a>

### Brave

Navegador Chromium com bloqueio de anúncios.

Repositório oficial Brave. Não importa automaticamente perfis do Firefox Flatpak/Snap.

| Campo | Valor |
| --- | --- |
| ID | brave |
| Área / categoria | Aplicações / Navegadores |
| Método | apt-repo |
| Fonte | Repositório oficial Brave |
| Pacotes | brave-browser |
| Repositório | brave |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "brave",
  "name": "Brave",
  "area": "apps",
  "category": "Navegadores",
  "summary": "Navegador Chromium com bloqueio de anúncios.",
  "why": "Repositório oficial Brave. Não importa automaticamente perfis do Firefox Flatpak/Snap.",
  "method": "apt-repo",
  "verify": {
    "cmd": [
      "brave-browser",
      "--version"
    ]
  },
  "packages": [
    "brave-browser"
  ],
  "repo": "brave",
  "source": "Repositório oficial Brave",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "tags": [
    "snap:brave"
  ],
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-firefox-snap-remove"></a>

### Remover Firefox Snap

Remove apenas o snap firefox (snapd fica).

Faça isto só depois de validar o perfil no Firefox Flatpak. Os dados em ~/snap/firefox são apagados pelo snapd (ele guarda um snapshot automático).

| Campo | Valor |
| --- | --- |
| ID | firefox-snap-remove |
| Área / categoria | Aplicações / Navegadores |
| Método | admin |
| Fonte | snap remove firefox |
| Atenção | Remove o Firefox Snap e os seus dados em ~/snap/firefox |
| Confirmação | Confirmo que abri o Firefox Flatpak e validei o perfil migrado (marcadores, senhas, extensões) e que os dados do Firefox Snap podem ser removidos. |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "firefox-snap-remove",
  "name": "Remover Firefox Snap",
  "area": "apps",
  "category": "Navegadores",
  "summary": "Remove apenas o snap firefox (snapd fica).",
  "why": "Faça isto só depois de validar o perfil no Firefox Flatpak. Os dados em ~/snap/firefox são apagados pelo snapd (ele guarda um snapshot automático).",
  "method": "admin",
  "adminOp": "snap-remove-firefox",
  "source": "snap remove firefox",
  "advanced": true,
  "conflicts": [],
  "confirm": "Confirmo que abri o Firefox Flatpak e validei o perfil migrado (marcadores, senhas, extensões) e que os dados do Firefox Snap podem ser removidos.",
  "sensitive": "Remove o Firefox Snap e os seus dados em ~/snap/firefox"
}
```

</details>

<a id="item-vlc"></a>

### VLC

Leitor multimédia que lê quase tudo.

Vídeos e áudio em qualquer formato.

| Campo | Valor |
| --- | --- |
| ID | vlc |
| Área / categoria | Aplicações / Multimédia |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | vlc |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "vlc",
  "name": "VLC",
  "area": "apps",
  "category": "Multimédia",
  "summary": "Leitor multimédia que lê quase tudo.",
  "why": "Vídeos e áudio em qualquer formato.",
  "method": "apt",
  "packages": [
    "vlc"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-qbittorrent"></a>

### qBittorrent

Cliente BitTorrent sem publicidade.

ISOs e ficheiros grandes distribuídos por torrent.

| Campo | Valor |
| --- | --- |
| ID | qbittorrent |
| Área / categoria | Aplicações / Multimédia |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | qbittorrent |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "qbittorrent",
  "name": "qBittorrent",
  "area": "apps",
  "category": "Multimédia",
  "summary": "Cliente BitTorrent sem publicidade.",
  "why": "ISOs e ficheiros grandes distribuídos por torrent.",
  "method": "apt",
  "packages": [
    "qbittorrent"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-flameshot"></a>

### Flameshot

Capturas de ecrã com anotações.

Anotar screenshots para tickets e documentação. Em Wayland usa o portal de captura.

| Campo | Valor |
| --- | --- |
| ID | flameshot |
| Área / categoria | Aplicações / Utilitários |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | flameshot |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "flameshot",
  "name": "Flameshot",
  "area": "apps",
  "category": "Utilitários",
  "summary": "Capturas de ecrã com anotações.",
  "why": "Anotar screenshots para tickets e documentação. Em Wayland usa o portal de captura.",
  "method": "apt",
  "packages": [
    "flameshot"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-solaar"></a>

### Solaar

Gestão de dispositivos Logitech Unifying/Bolt.

Emparelhar ratos/teclados Logitech e ver bateria.

| Campo | Valor |
| --- | --- |
| ID | solaar |
| Área / categoria | Aplicações / Utilitários |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | solaar |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "solaar",
  "name": "Solaar",
  "area": "apps",
  "category": "Utilitários",
  "summary": "Gestão de dispositivos Logitech Unifying/Bolt.",
  "why": "Emparelhar ratos/teclados Logitech e ver bateria.",
  "method": "apt",
  "packages": [
    "solaar"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-pandoc"></a>

### Pandoc + LaTeX

Conversão de documentos (Markdown→PDF/DOCX).

Gerar PDFs a partir de Markdown. O LaTeX ocupa ~500 MB.

| Campo | Valor |
| --- | --- |
| ID | pandoc |
| Área / categoria | Aplicações / Documentos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | pandoc, texlive-latex-recommended, texlive-xetex |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "pandoc",
  "name": "Pandoc + LaTeX",
  "area": "apps",
  "category": "Documentos",
  "summary": "Conversão de documentos (Markdown→PDF/DOCX).",
  "why": "Gerar PDFs a partir de Markdown. O LaTeX ocupa ~500 MB.",
  "method": "apt",
  "packages": [
    "pandoc",
    "texlive-latex-recommended",
    "texlive-xetex"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-lang-pt"></a>

### Português (PT/BR)

Traduções e dicionários de português.

Interface e corretor ortográfico em português.

| Campo | Valor |
| --- | --- |
| ID | lang-pt |
| Área / categoria | Aplicações / Idiomas |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | language-pack-pt, language-pack-gnome-pt, hunspell-pt-pt, hunspell-pt-br |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "lang-pt",
  "name": "Português (PT/BR)",
  "area": "apps",
  "category": "Idiomas",
  "summary": "Traduções e dicionários de português.",
  "why": "Interface e corretor ortográfico em português.",
  "method": "apt",
  "packages": [
    "language-pack-pt",
    "language-pack-gnome-pt",
    "hunspell-pt-pt",
    "hunspell-pt-br"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-lang-en"></a>

### Inglês

Traduções e dicionário de inglês.

Corretor ortográfico para documentação técnica em inglês.

| Campo | Valor |
| --- | --- |
| ID | lang-en |
| Área / categoria | Aplicações / Idiomas |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | language-pack-en, language-pack-gnome-en, hunspell-en-us |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "lang-en",
  "name": "Inglês",
  "area": "apps",
  "category": "Idiomas",
  "summary": "Traduções e dicionário de inglês.",
  "why": "Corretor ortográfico para documentação técnica em inglês.",
  "method": "apt",
  "packages": [
    "language-pack-en",
    "language-pack-gnome-en",
    "hunspell-en-us"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-nala"></a>

### Nala

Gestor de pacotes no terminal com resumos claros e histórico.

Facilita rever instalações e atualizações: nala install, nala upgrade e nala history.

| Campo | Valor |
| --- | --- |
| ID | nala |
| Área / categoria | Aplicações / Sistema |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | nala |
| Evidência da receita | Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04) |

Exemplos: sudo nala update; sudo nala install nome-do-pacote; nala history. O nala fetch altera fontes de espelhos: use-o apenas se quiser escolher esses espelhos.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nala",
  "name": "Nala",
  "area": "apps",
  "category": "Sistema",
  "summary": "Gestor de pacotes no terminal com resumos claros e histórico.",
  "why": "Facilita rever instalações e atualizações: nala install, nala upgrade e nala history.",
  "method": "apt",
  "packages": [
    "nala"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "Exemplos: sudo nala update; sudo nala install nome-do-pacote; nala history. O nala fetch altera fontes de espelhos: use-o apenas se quiser escolher esses espelhos.",
  "verified": "Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "verify": {
    "cmd": [
      "nala",
      "--version"
    ]
  }
}
```

</details>

<a id="item-flathub"></a>

### Flatpak + Flathub

Suporte a aplicações Flatpak (o remoto Flathub é adicionado por utilizador).

Necessário para Firefox Flatpak, FreeFileSync, DBeaver, etc.

| Campo | Valor |
| --- | --- |
| ID | flathub |
| Área / categoria | Aplicações / Sistema |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | flatpak, gnome-software-plugin-flatpak |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "flathub",
  "name": "Flatpak + Flathub",
  "area": "apps",
  "category": "Sistema",
  "summary": "Suporte a aplicações Flatpak (o remoto Flathub é adicionado por utilizador).",
  "why": "Necessário para Firefox Flatpak, FreeFileSync, DBeaver, etc.",
  "method": "apt",
  "packages": [
    "flatpak",
    "gnome-software-plugin-flatpak"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gnome-tweaks"></a>

### GNOME Tweaks (Ajustes)

Fontes, aplicações de arranque, botões de janela.

Ajustes que as Definições não expõem, como aplicações de arranque.

| Campo | Valor |
| --- | --- |
| ID | gnome-tweaks |
| Área / categoria | Aplicações / Sistema |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gnome-tweaks |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gnome-tweaks",
  "name": "GNOME Tweaks (Ajustes)",
  "area": "apps",
  "category": "Sistema",
  "summary": "Fontes, aplicações de arranque, botões de janela.",
  "why": "Ajustes que as Definições não expõem, como aplicações de arranque.",
  "method": "apt",
  "packages": [
    "gnome-tweaks"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-extension-manager"></a>

### Extension Manager

Procura, instala e configura extensões GNOME.

Forma segura e gráfica de gerir extensões e as suas preferências.

| Campo | Valor |
| --- | --- |
| ID | extension-manager |
| Área / categoria | Aplicações / Sistema |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gnome-shell-extension-manager |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "extension-manager",
  "name": "Extension Manager",
  "area": "apps",
  "category": "Sistema",
  "summary": "Procura, instala e configura extensões GNOME.",
  "why": "Forma segura e gráfica de gerir extensões e as suas preferências.",
  "method": "apt",
  "packages": [
    "gnome-shell-extension-manager"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-dconf-editor"></a>

### dconf Editor

Editor gráfico de todas as preferências GNOME.

Para inspecionar/alterar chaves avançadas. Use com cuidado.

| Campo | Valor |
| --- | --- |
| ID | dconf-editor |
| Área / categoria | Aplicações / Sistema |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | dconf-editor |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dconf-editor",
  "name": "dconf Editor",
  "area": "apps",
  "category": "Sistema",
  "summary": "Editor gráfico de todas as preferências GNOME.",
  "why": "Para inspecionar/alterar chaves avançadas. Use com cuidado.",
  "method": "apt",
  "packages": [
    "dconf-editor"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-python-dev"></a>

### Python (pip, venv, pipx)

Ferramentas Python para ambientes isolados.

pipx instala CLIs Python sem tocar no Python do sistema.

| Campo | Valor |
| --- | --- |
| ID | python-dev |
| Área / categoria | Desenvolvimento / Linguagens |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | python3-pip, python3-venv, pipx |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "python-dev",
  "name": "Python (pip, venv, pipx)",
  "area": "dev",
  "category": "Linguagens",
  "summary": "Ferramentas Python para ambientes isolados.",
  "why": "pipx instala CLIs Python sem tocar no Python do sistema.",
  "method": "apt",
  "packages": [
    "python3-pip",
    "python3-venv",
    "pipx"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-golang"></a>

### Go

Compilador e ferramentas Go do arquivo (1.26).

Ferramentas CLI modernas e projetos Go.

| Campo | Valor |
| --- | --- |
| ID | golang |
| Área / categoria | Desenvolvimento / Linguagens |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | golang-go |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "golang",
  "name": "Go",
  "area": "dev",
  "category": "Linguagens",
  "summary": "Compilador e ferramentas Go do arquivo (1.26).",
  "why": "Ferramentas CLI modernas e projetos Go.",
  "method": "apt",
  "verify": {
    "cmd": [
      "go",
      "version"
    ]
  },
  "packages": [
    "golang-go"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-android-tools"></a>

### adb + scrcpy

Depuração Android e espelho do ecrã do telemóvel.

Testar apps/PWAs no telemóvel e controlá-lo a partir do PC.

| Campo | Valor |
| --- | --- |
| ID | android-tools |
| Área / categoria | Desenvolvimento / Linguagens |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | adb, scrcpy |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "android-tools",
  "name": "adb + scrcpy",
  "area": "dev",
  "category": "Linguagens",
  "summary": "Depuração Android e espelho do ecrã do telemóvel.",
  "why": "Testar apps/PWAs no telemóvel e controlá-lo a partir do PC.",
  "method": "apt",
  "packages": [
    "adb",
    "scrcpy"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-sqlite"></a>

### SQLite

Cliente SQLite.

Testes Laravel com SQLite e inspeção de ficheiros .sqlite.

| Campo | Valor |
| --- | --- |
| ID | sqlite |
| Área / categoria | Desenvolvimento / Bases de dados |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | sqlite3 |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "sqlite",
  "name": "SQLite",
  "area": "dev",
  "category": "Bases de dados",
  "summary": "Cliente SQLite.",
  "why": "Testes Laravel com SQLite e inspeção de ficheiros .sqlite.",
  "method": "apt",
  "packages": [
    "sqlite3"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-mariadb-client"></a>

### Cliente MariaDB/MySQL

Cliente mysql/mariadb na linha de comandos.

Ligar a bases de dados em Docker ou remotas, dumps e imports.

| Campo | Valor |
| --- | --- |
| ID | mariadb-client |
| Área / categoria | Desenvolvimento / Bases de dados |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | mariadb-client |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "mariadb-client",
  "name": "Cliente MariaDB/MySQL",
  "area": "dev",
  "category": "Bases de dados",
  "summary": "Cliente mysql/mariadb na linha de comandos.",
  "why": "Ligar a bases de dados em Docker ou remotas, dumps e imports.",
  "method": "apt",
  "packages": [
    "mariadb-client"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-mariadb-server"></a>

### Servidor MariaDB (local)

Servidor de base de dados a correr como serviço.

Só se não usar Docker para bases de dados. Fica a arrancar com o sistema.

| Campo | Valor |
| --- | --- |
| ID | mariadb-server |
| Área / categoria | Desenvolvimento / Serviços locais · Dados e cache |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | mariadb-server |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "mariadb-server",
  "name": "Servidor MariaDB (local)",
  "area": "dev",
  "category": "Serviços locais · Dados e cache",
  "summary": "Servidor de base de dados a correr como serviço.",
  "why": "Só se não usar Docker para bases de dados. Fica a arrancar com o sistema.",
  "method": "apt",
  "packages": [
    "mariadb-server"
  ],
  "service": true,
  "advanced": false,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-supervisor"></a>

### Supervisor

Gestor de processos (queues/workers).

Só para workers locais fora do Docker. Arranca como serviço.

| Campo | Valor |
| --- | --- |
| ID | supervisor |
| Área / categoria | Desenvolvimento / Bases de dados |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | supervisor |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "supervisor",
  "name": "Supervisor",
  "area": "dev",
  "category": "Bases de dados",
  "summary": "Gestor de processos (queues/workers).",
  "why": "Só para workers locais fora do Docker. Arranca como serviço.",
  "method": "apt",
  "packages": [
    "supervisor"
  ],
  "service": true,
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-dbeaver"></a>

### DBeaver Community

Cliente gráfico universal de bases de dados.

Explorar MySQL/MariaDB/Postgres/SQLite com uma só ferramenta.

| Campo | Valor |
| --- | --- |
| ID | dbeaver |
| Área / categoria | Desenvolvimento / Bases de dados |
| Método | flatpak |
| Fonte | Flathub (editor não verificado) |
| Flatpak | io.dbeaver.DBeaverCommunity |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dbeaver",
  "name": "DBeaver Community",
  "area": "dev",
  "category": "Bases de dados",
  "summary": "Cliente gráfico universal de bases de dados.",
  "why": "Explorar MySQL/MariaDB/Postgres/SQLite com uma só ferramenta.",
  "method": "flatpak",
  "flatpak": "io.dbeaver.DBeaverCommunity",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor não verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-zsh"></a>

### Zsh

Shell alternativa ao Bash.

Melhor completação; a shell padrão NÃO é alterada automaticamente (chsh pede a sua senha; ver notas).

| Campo | Valor |
| --- | --- |
| ID | zsh |
| Área / categoria | Desenvolvimento / Terminal |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | zsh |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

Para tornar padrão: chsh -s /usr/bin/zsh e termine sessão. Os seus ~/.bashrc/~/.zshrc existentes não são tocados.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "zsh",
  "name": "Zsh",
  "area": "dev",
  "category": "Terminal",
  "summary": "Shell alternativa ao Bash.",
  "why": "Melhor completação; a shell padrão NÃO é alterada automaticamente (chsh pede a sua senha; ver notas).",
  "method": "apt",
  "packages": [
    "zsh"
  ],
  "notes": "Para tornar padrão: chsh -s /usr/bin/zsh e termine sessão. Os seus ~/.bashrc/~/.zshrc existentes não são tocados.",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-zsh-autosuggestions"></a>

### Zsh: sugestões pelo histórico

Sugere comandos à medida que escreve; aceite com a seta para a direita.

Poupa repetição no terminal, reutilizando o seu próprio histórico.

| Campo | Valor |
| --- | --- |
| ID | zsh-autosuggestions |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | zsh-autosuggestions |
| Evidência da receita | Instalação real e carregamento com Zsh verificados no Ubuntu 26.04 (container descartável, 2026-10-04) |
| Dependências | [Zsh](#item-zsh) |

Depois de instalar, acrescente ao seu .zshrc: source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh. Se já usa outra instalação do plugin, mantenha uma só. O Meu Ubuntu instala o pacote; a ativação é manual.

#### Guia manual

1. Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.
2. Para ativar no Zsh, acrescente ao ~/.zshrc:
```zsh
source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh
```
Se já carrega o plugin por Oh My Zsh ou outra instalação, mantenha esse método e não acrescente uma segunda cópia.
3. Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual.

[Documentação do fornecedor](https://github.com/zsh-users/zsh-autosuggestions)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "zsh-autosuggestions",
  "name": "Zsh: sugestões pelo histórico",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Sugere comandos à medida que escreve; aceite com a seta para a direita.",
  "why": "Poupa repetição no terminal, reutilizando o seu próprio histórico.",
  "method": "apt",
  "packages": [
    "zsh-autosuggestions"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "Depois de instalar, acrescente ao seu .zshrc: source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh. Se já usa outra instalação do plugin, mantenha uma só. O Meu Ubuntu instala o pacote; a ativação é manual.",
  "verified": "Instalação real e carregamento com Zsh verificados no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "requires": [
    "zsh"
  ],
  "addonsFor": [
    "zsh"
  ],
  "manual": {
    "url": "https://github.com/zsh-users/zsh-autosuggestions",
    "steps": [
      "Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.",
      "Para ativar no Zsh, acrescente ao ~/.zshrc:\n```zsh\nsource /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh\n```\nSe já carrega o plugin por Oh My Zsh ou outra instalação, mantenha esse método e não acrescente uma segunda cópia.",
      "Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual."
    ]
  }
}
```

</details>

<a id="item-zsh-syntax-highlighting"></a>

### Zsh: destaque de sintaxe

Destaca comandos, opções e caminhos enquanto escreve.

Ajuda a perceber erros de escrita antes de executar um comando.

| Campo | Valor |
| --- | --- |
| ID | zsh-syntax-highlighting |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | zsh-syntax-highlighting |
| Evidência da receita | Instalação real e carregamento com Zsh verificados no Ubuntu 26.04 (container descartável, 2026-10-04) |
| Dependências | [Zsh](#item-zsh) |

No fim do .zshrc, após os outros plugins: source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh. Não carregue duas cópias do plugin. O Meu Ubuntu instala o pacote; a ativação é manual.

#### Guia manual

1. Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.
2. Para ativar no Zsh, acrescente ao ~/.zshrc:
```zsh
source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
```
Coloque esta linha no fim do .zshrc, depois dos outros plugins. Não carregue duas cópias.
3. Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual.

[Documentação do fornecedor](https://github.com/zsh-users/zsh-syntax-highlighting)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "zsh-syntax-highlighting",
  "name": "Zsh: destaque de sintaxe",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Destaca comandos, opções e caminhos enquanto escreve.",
  "why": "Ajuda a perceber erros de escrita antes de executar um comando.",
  "method": "apt",
  "packages": [
    "zsh-syntax-highlighting"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "No fim do .zshrc, após os outros plugins: source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh. Não carregue duas cópias do plugin. O Meu Ubuntu instala o pacote; a ativação é manual.",
  "verified": "Instalação real e carregamento com Zsh verificados no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "requires": [
    "zsh"
  ],
  "addonsFor": [
    "zsh"
  ],
  "manual": {
    "url": "https://github.com/zsh-users/zsh-syntax-highlighting",
    "steps": [
      "Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.",
      "Para ativar no Zsh, acrescente ao ~/.zshrc:\n```zsh\nsource /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n```\nColoque esta linha no fim do .zshrc, depois dos outros plugins. Não carregue duas cópias.",
      "Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual."
    ]
  }
}
```

</details>

<a id="item-fzf"></a>

### fzf — pesquisa no terminal

Pesquisa interativa e aproximada para histórico, ficheiros e listas.

Encontre um comando antigo com Ctrl+R, ou um ficheiro com Ctrl+T, depois de ativar a integração com a shell.

| Campo | Valor |
| --- | --- |
| ID | fzf |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | fzf |
| Evidência da receita | Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04) |

Zsh: acrescente source <(fzf --zsh) ao .zshrc, ou ative o plugin fzf do Oh My Zsh. Bash: acrescente eval "$(fzf --bash)" ao .bashrc. Use um só método de integração.

#### Guia manual

1. Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.
2. Para ativar no Zsh, acrescente ao ~/.zshrc:
```zsh
source <(fzf --zsh)
```
Se já usa o plugin fzf do Oh My Zsh, mantenha apenas esse método. Depois teste Ctrl+R para pesquisar o histórico.
3. Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual.

[Documentação do fornecedor](https://github.com/junegunn/fzf)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "fzf",
  "name": "fzf — pesquisa no terminal",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Pesquisa interativa e aproximada para histórico, ficheiros e listas.",
  "why": "Encontre um comando antigo com Ctrl+R, ou um ficheiro com Ctrl+T, depois de ativar a integração com a shell.",
  "method": "apt",
  "packages": [
    "fzf"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "Zsh: acrescente source <(fzf --zsh) ao .zshrc, ou ative o plugin fzf do Oh My Zsh. Bash: acrescente eval \"$(fzf --bash)\" ao .bashrc. Use um só método de integração.",
  "verified": "Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "verify": {
    "cmd": [
      "fzf",
      "--version"
    ]
  },
  "addonsFor": [
    "zsh"
  ],
  "manual": {
    "url": "https://github.com/junegunn/fzf",
    "steps": [
      "Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.",
      "Para ativar no Zsh, acrescente ao ~/.zshrc:\n```zsh\nsource <(fzf --zsh)\n```\nSe já usa o plugin fzf do Oh My Zsh, mantenha apenas esse método. Depois teste Ctrl+R para pesquisar o histórico.",
      "Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual."
    ]
  }
}
```

</details>

<a id="item-zoxide"></a>

### zoxide — navegar entre pastas

Aprende as pastas que usa e permite saltar para elas com z.

Após visitar uma pasta, z nome evita escrever o caminho inteiro; zi abre uma escolha interativa com fzf.

| Campo | Valor |
| --- | --- |
| ID | zoxide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | zoxide |
| Evidência da receita | Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04) |

Zsh: acrescente eval "$(zoxide init zsh)" ao .zshrc. Bash: acrescente eval "$(zoxide init bash)" ao .bashrc. Não substitui o comando cd. A integração é manual.

#### Guia manual

1. Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.
2. Para ativar no Zsh, acrescente ao ~/.zshrc:
```zsh
eval "$(zoxide init zsh)"
```
Abra um novo terminal e use z para visitar pastas conhecidas. A integração mantém o comando cd.
3. Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual.

[Documentação do fornecedor](https://github.com/ajeetdsouza/zoxide)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "zoxide",
  "name": "zoxide — navegar entre pastas",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Aprende as pastas que usa e permite saltar para elas com z.",
  "why": "Após visitar uma pasta, z nome evita escrever o caminho inteiro; zi abre uma escolha interativa com fzf.",
  "method": "apt",
  "packages": [
    "zoxide"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "Zsh: acrescente eval \"$(zoxide init zsh)\" ao .zshrc. Bash: acrescente eval \"$(zoxide init bash)\" ao .bashrc. Não substitui o comando cd. A integração é manual.",
  "verified": "Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "verify": {
    "cmd": [
      "zoxide",
      "--version"
    ]
  },
  "addonsFor": [
    "zsh"
  ],
  "manual": {
    "url": "https://github.com/ajeetdsouza/zoxide",
    "steps": [
      "Instale o pacote pelo plano. Faça uma cópia do seu ~/.zshrc antes de editar a configuração.",
      "Para ativar no Zsh, acrescente ao ~/.zshrc:\n```zsh\neval \"$(zoxide init zsh)\"\n```\nAbra um novo terminal e use z para visitar pastas conhecidas. A integração mantém o comando cd.",
      "Abra um novo terminal Zsh e confirme o funcionamento. O Meu Ubuntu verifica o pacote instalado; a ativação na shell é um passo manual."
    ]
  }
}
```

</details>

<a id="item-tealdeer"></a>

### tldr — exemplos de comandos

Páginas curtas com exemplos práticos, através do cliente tealdeer.

Consulte tldr tar, tldr git ou tldr rsync quando precisa de um exemplo rápido.

| Campo | Valor |
| --- | --- |
| ID | tealdeer |
| Área / categoria | Desenvolvimento / Terminal |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | tealdeer |
| Evidência da receita | Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04) |

Na primeira utilização: tldr --update, depois tldr tar. A atualização descarrega as páginas para a cache do utilizador; para informação completa, consulte também man.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tealdeer",
  "name": "tldr — exemplos de comandos",
  "area": "dev",
  "category": "Terminal",
  "summary": "Páginas curtas com exemplos práticos, através do cliente tealdeer.",
  "why": "Consulte tldr tar, tldr git ou tldr rsync quando precisa de um exemplo rápido.",
  "method": "apt",
  "packages": [
    "tealdeer"
  ],
  "source": "Arquivo Ubuntu",
  "notes": "Na primeira utilização: tldr --update, depois tldr tar. A atualização descarrega as páginas para a cache do utilizador; para informação completa, consulte também man.",
  "verified": "Instalação real e consulta de versão no Ubuntu 26.04 (container descartável, 2026-10-04)",
  "verify": {
    "cmd": [
      "tldr",
      "--version"
    ]
  }
}
```

</details>

<a id="item-oh-my-zsh"></a>

### Oh My Zsh — configuração e plugins

Configuração guiada de Zsh com atalhos Git, sugestões e destaque de sintaxe.

Reúne uma configuração de terminal útil, com os plugins do arquivo Ubuntu e o tema robbyrussell, sem exigir fontes especiais.

| Campo | Valor |
| --- | --- |
| ID | oh-my-zsh |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Oh My Zsh oficial; configuração guiada por utilizador |
| Evidência da receita | Instalador oficial revisto em 2026-10-04; commit 4d4cfc287e9d887b81242c0e431b5f49f9cec5c1. Integração e preservação de .zshrc verificadas em container. |
| Dependências | [Zsh](#item-zsh), [Git](#item-git), [curl](#item-curl), [Zsh: sugestões pelo histórico](#item-zsh-autosuggestions), [Zsh: destaque de sintaxe](#item-zsh-syntax-highlighting) |

O plano instala as dependências e apresenta estes passos. A aplicação não altera a shell de login nem o .zshrc. Se já usa Oh My Zsh, ajuste a configuração existente em vez de o carregar duas vezes. Se usa ZDOTDIR ou uma pasta ZSH personalizada, adapte os caminhos.

#### Guia manual

1. Depois de aplicar as dependências do plano, faça uma cópia do seu .zshrc antes de editar. Se já usa Oh My Zsh, salte a instalação e ajuste apenas os plugins.
2. Transfira o instalador oficial para uma pasta temporária:
```sh
omz_tmp=$(mktemp -d)
curl --fail --show-error --location https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/4d4cfc287e9d887b81242c0e431b5f49f9cec5c1/tools/install.sh --output "$omz_tmp/install.sh"
```
Leia o ficheiro antes de o executar. Use a mesma sessão do terminal para o próximo passo.
3. Execute como o seu utilizador, na mesma sessão do terminal:
```sh
sh "$omz_tmp/install.sh" --unattended --keep-zshrc
```
Estas opções preservam um .zshrc existente e deixam a shell de login como está.
4. Integre estas linhas no .zshrc, preservando as suas definições e sem duplicar carregamentos. Este bloco é configuração para o ficheiro:
```zsh
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME="robbyrussell"
plugins=(git colored-man-pages command-not-found)
if command -v fzf >/dev/null 2>&1; then plugins+=(fzf); fi
source "$ZSH/oh-my-zsh.sh"
if command -v zoxide >/dev/null 2>&1; then eval "$(zoxide init zsh)"; fi
source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh
source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
```
5. Abra zsh para testar. Experimente atalhos Git (gst), sugestões com a seta direita e, se incluiu fzf, Ctrl+R para o histórico. Apenas depois, se quiser tornar Zsh a shell de login, execute chsh -s /usr/bin/zsh e termine sessão.

[Documentação do fornecedor](https://github.com/ohmyzsh/ohmyzsh)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "oh-my-zsh",
  "name": "Oh My Zsh — configuração e plugins",
  "area": "dev",
  "category": "Plugins e complementos",
  "method": "manual",
  "summary": "Configuração guiada de Zsh com atalhos Git, sugestões e destaque de sintaxe.",
  "why": "Reúne uma configuração de terminal útil, com os plugins do arquivo Ubuntu e o tema robbyrussell, sem exigir fontes especiais.",
  "source": "Oh My Zsh oficial; configuração guiada por utilizador",
  "requires": [
    "zsh",
    "git",
    "curl",
    "zsh-autosuggestions",
    "zsh-syntax-highlighting"
  ],
  "detect": {
    "files": [
      "~/.oh-my-zsh/oh-my-zsh.sh"
    ]
  },
  "notes": "O plano instala as dependências e apresenta estes passos. A aplicação não altera a shell de login nem o .zshrc. Se já usa Oh My Zsh, ajuste a configuração existente em vez de o carregar duas vezes. Se usa ZDOTDIR ou uma pasta ZSH personalizada, adapte os caminhos.",
  "verified": "Instalador oficial revisto em 2026-10-04; commit 4d4cfc287e9d887b81242c0e431b5f49f9cec5c1. Integração e preservação de .zshrc verificadas em container.",
  "manual": {
    "url": "https://github.com/ohmyzsh/ohmyzsh",
    "steps": [
      "Depois de aplicar as dependências do plano, faça uma cópia do seu .zshrc antes de editar. Se já usa Oh My Zsh, salte a instalação e ajuste apenas os plugins.",
      "Transfira o instalador oficial para uma pasta temporária:\n```sh\nomz_tmp=$(mktemp -d)\ncurl --fail --show-error --location https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/4d4cfc287e9d887b81242c0e431b5f49f9cec5c1/tools/install.sh --output \"$omz_tmp/install.sh\"\n```\nLeia o ficheiro antes de o executar. Use a mesma sessão do terminal para o próximo passo.",
      "Execute como o seu utilizador, na mesma sessão do terminal:\n```sh\nsh \"$omz_tmp/install.sh\" --unattended --keep-zshrc\n```\nEstas opções preservam um .zshrc existente e deixam a shell de login como está.",
      "Integre estas linhas no .zshrc, preservando as suas definições e sem duplicar carregamentos. Este bloco é configuração para o ficheiro:\n```zsh\nexport ZSH=\"$HOME/.oh-my-zsh\"\nZSH_THEME=\"robbyrussell\"\nplugins=(git colored-man-pages command-not-found)\nif command -v fzf >/dev/null 2>&1; then plugins+=(fzf); fi\nsource \"$ZSH/oh-my-zsh.sh\"\nif command -v zoxide >/dev/null 2>&1; then eval \"$(zoxide init zsh)\"; fi\nsource /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh\nsource /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n```",
      "Abra zsh para testar. Experimente atalhos Git (gst), sugestões com a seta direita e, se incluiu fzf, Ctrl+R para o histórico. Apenas depois, se quiser tornar Zsh a shell de login, execute chsh -s /usr/bin/zsh e termine sessão."
    ]
  },
  "addonsFor": [
    "zsh"
  ]
}
```

</details>

<a id="item-php"></a>

### PHP (sistema, APT)

PHP CLI com as extensões que o Laravel pede.

Versão fixa escolhida por si; nunca é atualizada implicitamente para outra versão principal. O Ubuntu 26.04 só traz 8.5.

| Campo | Valor |
| --- | --- |
| ID | php |
| Área / categoria | Desenvolvimento / PHP |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | php{opt:versions}-cli, php{opt:versions}-common, php{opt:versions}-mbstring, php{opt:versions}-xml, php{opt:versions}-curl, php{opt:versions}-zip, php{opt:versions}-sqlite3, php{opt:versions}-mysql, php{opt:versions}-intl, php{opt:versions}-bcmath, php{opt:versions}-gd |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Conflitos | [PHP via php.new (Herd Lite)](#item-php-herd-lite) (aviso): Dois PHP no PATH: o herd-lite fica à frente do PHP do sistema. |

Opção **Versões** (`versions`): múltiplos valores. Padrão: 8.5. 

| Valor | Descrição |
| --- | --- |
| 8.5 | PHP 8.5 (arquivo Ubuntu 26.04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "php",
  "name": "PHP (sistema, APT)",
  "area": "dev",
  "category": "PHP",
  "summary": "PHP CLI com as extensões que o Laravel pede.",
  "why": "Versão fixa escolhida por si; nunca é atualizada implicitamente para outra versão principal. O Ubuntu 26.04 só traz 8.5.",
  "method": "apt",
  "verify": {
    "cmd": [
      "php",
      "-v"
    ]
  },
  "packages": [
    "php{opt:versions}-cli",
    "php{opt:versions}-common",
    "php{opt:versions}-mbstring",
    "php{opt:versions}-xml",
    "php{opt:versions}-curl",
    "php{opt:versions}-zip",
    "php{opt:versions}-sqlite3",
    "php{opt:versions}-mysql",
    "php{opt:versions}-intl",
    "php{opt:versions}-bcmath",
    "php{opt:versions}-gd"
  ],
  "options": [
    {
      "key": "versions",
      "label": "Versões",
      "values": [
        {
          "value": "8.5",
          "label": "PHP 8.5 (arquivo Ubuntu 26.04)"
        }
      ],
      "default": [
        "8.5"
      ],
      "multi": true
    }
  ],
  "conflicts": [
    {
      "id": "php-herd-lite",
      "hard": false,
      "reason": "Dois PHP no PATH: o herd-lite fica à frente do PHP do sistema."
    }
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-php-ppa"></a>

### PHP outras versões (PPA Ondřej)

PHP 8.3/8.4 lado a lado com o do sistema.

Para projetos que exigem versão anterior. Em 2026-10-04 o PPA NÃO publica para 26.04: o passo falha sem alterar nada até ser publicado.

| Campo | Valor |
| --- | --- |
| ID | php-ppa |
| Área / categoria | Desenvolvimento / PHP |
| Método | apt-repo |
| Fonte | PPA ppa:ondrej/php (Launchpad) |
| Pacotes | php{opt:versions}-cli, php{opt:versions}-common, php{opt:versions}-mbstring, php{opt:versions}-xml, php{opt:versions}-curl, php{opt:versions}-zip, php{opt:versions}-sqlite3, php{opt:versions}-mysql, php{opt:versions}-intl, php{opt:versions}-bcmath, php{opt:versions}-gd |
| Repositório | ondrej-php |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |
| Avançado | Sim |

Opção **Versões** (`versions`): múltiplos valores. Padrão: 8.4. 

| Valor | Descrição |
| --- | --- |
| 8.3 | PHP 8.3 |
| 8.4 | PHP 8.4 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "php-ppa",
  "name": "PHP outras versões (PPA Ondřej)",
  "area": "dev",
  "category": "PHP",
  "summary": "PHP 8.3/8.4 lado a lado com o do sistema.",
  "why": "Para projetos que exigem versão anterior. Em 2026-10-04 o PPA NÃO publica para 26.04: o passo falha sem alterar nada até ser publicado.",
  "method": "apt-repo",
  "packages": [
    "php{opt:versions}-cli",
    "php{opt:versions}-common",
    "php{opt:versions}-mbstring",
    "php{opt:versions}-xml",
    "php{opt:versions}-curl",
    "php{opt:versions}-zip",
    "php{opt:versions}-sqlite3",
    "php{opt:versions}-mysql",
    "php{opt:versions}-intl",
    "php{opt:versions}-bcmath",
    "php{opt:versions}-gd"
  ],
  "repo": "ondrej-php",
  "source": "PPA ppa:ondrej/php (Launchpad)",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "advanced": true,
  "options": [
    {
      "key": "versions",
      "label": "Versões",
      "values": [
        {
          "value": "8.3",
          "label": "PHP 8.3"
        },
        {
          "value": "8.4",
          "label": "PHP 8.4"
        }
      ],
      "default": [
        "8.4"
      ],
      "multi": true
    }
  ],
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-php-herd-lite"></a>

### PHP via php.new (Herd Lite)

PHP estático + Composer + instalador Laravel em ~/.config/herd-lite/bin, sem root.

Instalação rápida de uma versão à escolha (8.3/8.4/8.5) sem mexer no sistema. Substitui o binário ao mudar de versão.

| Campo | Valor |
| --- | --- |
| ID | php-herd-lite |
| Área / categoria | Desenvolvimento / PHP |
| Método | user-script |
| Fonte | php.new (Laravel Herd Lite) — script do fornecedor inspecionado antes de correr |
| Evidência da receita | Script inspecionado em 2026-10-04: sem sudo; hosts download.herdphp.com e curl.se |
| Confirmação | Revi o script (SHA-256 indicado) e aceito executar este instalador do fornecedor como o meu utilizador. |
| Conflitos | [PHP (sistema, APT)](#item-php) (aviso): Dois PHP no PATH: o herd-lite fica à frente do PHP do sistema. |

Opção **Versão** (`version`): um valor. Padrão: 8.5. 

| Valor | Descrição |
| --- | --- |
| 8.5 | PHP 8.5 |
| 8.4 | PHP 8.4 |
| 8.3 | PHP 8.3 |

O instalador acrescenta o PATH ao .bashrc/.zshrc (backup feito antes) e não verifica checksums dos binários. Remoção: ~/.config/herd-lite/bin/uninstall_herd_lite.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "php-herd-lite",
  "name": "PHP via php.new (Herd Lite)",
  "area": "dev",
  "category": "PHP",
  "summary": "PHP estático + Composer + instalador Laravel em ~/.config/herd-lite/bin, sem root.",
  "why": "Instalação rápida de uma versão à escolha (8.3/8.4/8.5) sem mexer no sistema. Substitui o binário ao mudar de versão.",
  "method": "user-script",
  "source": "php.new (Laravel Herd Lite) — script do fornecedor inspecionado antes de correr",
  "script": {
    "url": "https://php.new/install/linux/{opt:version}",
    "allowHosts": [
      "download.herdphp.com",
      "curl.se",
      "herd.laravel.com"
    ],
    "backup": [
      "~/.bashrc",
      "~/.bash_profile",
      "~/.profile",
      "~/.zshrc",
      "~/.zprofile"
    ],
    "binary": "~/.config/herd-lite/bin/php",
    "versionArgs": [
      "-r",
      "echo PHP_MAJOR_VERSION.\".\".PHP_MINOR_VERSION;"
    ],
    "option": "version"
  },
  "options": [
    {
      "key": "version",
      "label": "Versão",
      "values": [
        {
          "value": "8.5",
          "label": "PHP 8.5"
        },
        {
          "value": "8.4",
          "label": "PHP 8.4"
        },
        {
          "value": "8.3",
          "label": "PHP 8.3"
        }
      ],
      "default": [
        "8.5"
      ]
    }
  ],
  "notes": "O instalador acrescenta o PATH ao .bashrc/.zshrc (backup feito antes) e não verifica checksums dos binários. Remoção: ~/.config/herd-lite/bin/uninstall_herd_lite.",
  "verified": "Script inspecionado em 2026-10-04: sem sudo; hosts download.herdphp.com e curl.se",
  "confirm": "Revi o script (SHA-256 indicado) e aceito executar este instalador do fornecedor como o meu utilizador.",
  "conflicts": [
    {
      "id": "php",
      "hard": false,
      "reason": "Dois PHP no PATH: o herd-lite fica à frente do PHP do sistema."
    }
  ]
}
```

</details>

<a id="item-composer"></a>

### Composer

Gestor de dependências PHP (arquivo, 2.9).

Obrigatório para Laravel. O php.new já inclui Composer próprio.

| Campo | Valor |
| --- | --- |
| ID | composer |
| Área / categoria | Desenvolvimento / PHP |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | composer |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Dependências | [PHP (sistema, APT)](#item-php) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "composer",
  "name": "Composer",
  "area": "dev",
  "category": "PHP",
  "summary": "Gestor de dependências PHP (arquivo, 2.9).",
  "why": "Obrigatório para Laravel. O php.new já inclui Composer próprio.",
  "method": "apt",
  "verify": {
    "cmd": [
      "composer",
      "--version"
    ]
  },
  "packages": [
    "composer"
  ],
  "requires": [
    "php"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-xdebug"></a>

### Xdebug (desativado por padrão)

Depurador/profiler PHP.

Instalado mas desligado (phpdismod) para não abrandar tudo; ative quando for depurar: sudo phpenmod xdebug.

| Campo | Valor |
| --- | --- |
| ID | xdebug |
| Área / categoria | Desenvolvimento / PHP |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | php{opt:versions}-xdebug |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Dependências | [PHP (sistema, APT)](#item-php) |

Opção **Versões** (`versions`): múltiplos valores. Padrão: 8.5. 

| Valor | Descrição |
| --- | --- |
| 8.5 | PHP 8.5 (arquivo Ubuntu 26.04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "xdebug",
  "name": "Xdebug (desativado por padrão)",
  "area": "dev",
  "category": "PHP",
  "summary": "Depurador/profiler PHP.",
  "why": "Instalado mas desligado (phpdismod) para não abrandar tudo; ative quando for depurar: sudo phpenmod xdebug.",
  "method": "apt",
  "packages": [
    "php{opt:versions}-xdebug"
  ],
  "requires": [
    "php"
  ],
  "post": [
    "xdebug-off"
  ],
  "options": [
    {
      "key": "versions",
      "label": "Versões",
      "values": [
        {
          "value": "8.5",
          "label": "PHP 8.5 (arquivo Ubuntu 26.04)"
        }
      ],
      "default": [
        "8.5"
      ],
      "multi": true
    }
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-php-build-deps"></a>

### Dependências de compilação PHP

Bibliotecas para compilar PHP/PECL a partir do código.

Só se compilar PHP (phpbrew/asdf) ou extensões PECL. Com APT/php.new não é necessário.

| Campo | Valor |
| --- | --- |
| ID | php-build-deps |
| Área / categoria | Desenvolvimento / Avançado |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | autoconf, bison, re2c, pkg-config, libxml2-dev, libsqlite3-dev, libssl-dev, libcurl4-openssl-dev, libonig-dev, libzip-dev |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "php-build-deps",
  "name": "Dependências de compilação PHP",
  "area": "dev",
  "category": "Avançado",
  "summary": "Bibliotecas para compilar PHP/PECL a partir do código.",
  "why": "Só se compilar PHP (phpbrew/asdf) ou extensões PECL. Com APT/php.new não é necessário.",
  "method": "apt",
  "packages": [
    "autoconf",
    "bison",
    "re2c",
    "pkg-config",
    "libxml2-dev",
    "libsqlite3-dev",
    "libssl-dev",
    "libcurl4-openssl-dev",
    "libonig-dev",
    "libzip-dev"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-erlang"></a>

### Erlang

Runtime Erlang/OTP.

Só para RabbitMQ local ou projetos Elixir. Prefira RabbitMQ em Docker.

| Campo | Valor |
| --- | --- |
| ID | erlang |
| Área / categoria | Desenvolvimento / Avançado |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | erlang |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "erlang",
  "name": "Erlang",
  "area": "dev",
  "category": "Avançado",
  "summary": "Runtime Erlang/OTP.",
  "why": "Só para RabbitMQ local ou projetos Elixir. Prefira RabbitMQ em Docker.",
  "method": "apt",
  "packages": [
    "erlang"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-java8"></a>

### Java 8 (OpenJDK)

JDK antigo.

Só para ferramentas legadas que exigem Java 8. Ainda existe no arquivo 26.04.

| Campo | Valor |
| --- | --- |
| ID | java8 |
| Área / categoria | Desenvolvimento / Avançado |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | openjdk-8-jdk |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "java8",
  "name": "Java 8 (OpenJDK)",
  "area": "dev",
  "category": "Avançado",
  "summary": "JDK antigo.",
  "why": "Só para ferramentas legadas que exigem Java 8. Ainda existe no arquivo 26.04.",
  "method": "apt",
  "packages": [
    "openjdk-8-jdk"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-wkhtmltopdf"></a>

### wkhtmltopdf

Conversor HTML→PDF descontinuado.

Projeto arquivado e sem pacote no 26.04. Prefira Chromium headless (Browsershot) ou Gotenberg em Docker.

| Campo | Valor |
| --- | --- |
| ID | wkhtmltopdf |
| Área / categoria | Desenvolvimento / Avançado |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Avançado | Sim |

#### Guia manual

1. Não suportado: o pacote não existe no Ubuntu 26.04 e o projeto foi arquivado.
2. Alternativa: spatie/browsershot ou o contentor gotenberg/gotenberg.

[Documentação do fornecedor](https://wkhtmltopdf.org/status.html)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "wkhtmltopdf",
  "name": "wkhtmltopdf",
  "area": "dev",
  "category": "Avançado",
  "summary": "Conversor HTML→PDF descontinuado.",
  "why": "Projeto arquivado e sem pacote no 26.04. Prefira Chromium headless (Browsershot) ou Gotenberg em Docker.",
  "method": "manual",
  "manual": {
    "steps": [
      "Não suportado: o pacote não existe no Ubuntu 26.04 e o projeto foi arquivado.",
      "Alternativa: spatie/browsershot ou o contentor gotenberg/gotenberg."
    ],
    "url": "https://wkhtmltopdf.org/status.html"
  },
  "source": "Manual / não suportado automaticamente",
  "advanced": true
}
```

</details>

<a id="item-mainline"></a>

### Mainline (kernels manuais)

Instalador de kernels não suportados.

Rejeitado: alterar o kernel está fora do âmbito; use o kernel HWE/GA do Ubuntu.

| Campo | Valor |
| --- | --- |
| ID | mainline |
| Área / categoria | Desenvolvimento / Avançado |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Avançado | Sim |

#### Guia manual

1. Não suportado: o Meu Ubuntu não altera kernel nem boot.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "mainline",
  "name": "Mainline (kernels manuais)",
  "area": "dev",
  "category": "Avançado",
  "summary": "Instalador de kernels não suportados.",
  "why": "Rejeitado: alterar o kernel está fora do âmbito; use o kernel HWE/GA do Ubuntu.",
  "method": "manual",
  "manual": {
    "steps": [
      "Não suportado: o Meu Ubuntu não altera kernel nem boot."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente",
  "advanced": true
}
```

</details>

<a id="item-nodejs"></a>

### Node.js do sistema (APT)

Node 22 LTS do arquivo Ubuntu.

Simples e atualizado pelo sistema; versão fixa. Escolha isto OU NVM.

| Campo | Valor |
| --- | --- |
| ID | nodejs |
| Área / categoria | Desenvolvimento / Node.js |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | nodejs, npm |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Conflitos | [NVM (várias versões Node)](#item-nvm) (aviso): Node do sistema + NVM: o NVM fica à frente no PATH; escolha um como principal. |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nodejs",
  "name": "Node.js do sistema (APT)",
  "area": "dev",
  "category": "Node.js",
  "summary": "Node 22 LTS do arquivo Ubuntu.",
  "why": "Simples e atualizado pelo sistema; versão fixa. Escolha isto OU NVM.",
  "method": "apt",
  "verify": {
    "cmd": [
      "node",
      "--version"
    ]
  },
  "packages": [
    "nodejs",
    "npm"
  ],
  "conflicts": [
    {
      "id": "nvm",
      "hard": false,
      "reason": "Node do sistema + NVM: o NVM fica à frente no PATH; escolha um como principal."
    }
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-nvm"></a>

### NVM (várias versões Node)

Gestor de versões Node por utilizador.

Projetos com versões Node diferentes (.nvmrc). Recomendado para desenvolvimento.

| Campo | Valor |
| --- | --- |
| ID | nvm |
| Área / categoria | Desenvolvimento / Node.js |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Conflitos | [Node.js do sistema (APT)](#item-nodejs) (aviso): Node do sistema + NVM: escolha um como principal. |

#### Guia manual

1. Leia a versão atual em https://github.com/nvm-sh/nvm#installing-and-updating
2. Transfira o install.sh da tag (ex.: v0.40.x) para ficheiro, reveja e execute: bash install.sh
3. Abra novo terminal e: nvm install --lts

[Documentação do fornecedor](https://github.com/nvm-sh/nvm)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nvm",
  "name": "NVM (várias versões Node)",
  "area": "dev",
  "category": "Node.js",
  "summary": "Gestor de versões Node por utilizador.",
  "why": "Projetos com versões Node diferentes (.nvmrc). Recomendado para desenvolvimento.",
  "method": "manual",
  "detect": {
    "files": [
      "~/.nvm/nvm.sh"
    ]
  },
  "manual": {
    "steps": [
      "Leia a versão atual em https://github.com/nvm-sh/nvm#installing-and-updating",
      "Transfira o install.sh da tag (ex.: v0.40.x) para ficheiro, reveja e execute: bash install.sh",
      "Abra novo terminal e: nvm install --lts"
    ],
    "url": "https://github.com/nvm-sh/nvm"
  },
  "source": "Manual / não suportado automaticamente",
  "conflicts": [
    {
      "id": "nodejs",
      "hard": false,
      "reason": "Node do sistema + NVM: escolha um como principal."
    }
  ]
}
```

</details>

<a id="item-bun"></a>

### Bun

Runtime/gestor de pacotes JS rápido.

Alternativa rápida ao npm para scripts e testes.

| Campo | Valor |
| --- | --- |
| ID | bun |
| Área / categoria | Desenvolvimento / Node.js |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Instale com o método oficial (script em https://bun.sh/install): transfira para ficheiro, reveja e corra com bash.
2. Ou: npm install -g bun

[Documentação do fornecedor](https://bun.sh/docs/installation)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "bun",
  "name": "Bun",
  "area": "dev",
  "category": "Node.js",
  "summary": "Runtime/gestor de pacotes JS rápido.",
  "why": "Alternativa rápida ao npm para scripts e testes.",
  "method": "manual",
  "detect": {
    "commands": [
      "bun"
    ],
    "files": [
      "~/.bun/bin/bun"
    ]
  },
  "manual": {
    "steps": [
      "Instale com o método oficial (script em https://bun.sh/install): transfira para ficheiro, reveja e corra com bash.",
      "Ou: npm install -g bun"
    ],
    "url": "https://bun.sh/docs/installation"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-docker-engine"></a>

### Docker Engine (oficial)

Docker Engine + Compose v2 (docker compose) do repositório oficial.

Recomendado para desenvolvimento Laravel com Sail/Compose. Corre nativamente, sem VM.

| Campo | Valor |
| --- | --- |
| ID | docker-engine |
| Área / categoria | Desenvolvimento / Contentores |
| Método | apt-repo |
| Fonte | Repositório oficial Docker (docs.docker.com) |
| Pacotes | docker-ce, docker-ce-cli, containerd.io, docker-buildx-plugin, docker-compose-plugin |
| Repositório | docker |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |
| Conflitos | [Docker (pacotes Ubuntu)](#item-docker-ubuntu) (bloqueia): docker.io (Ubuntu) e docker-ce (Docker) são pacotes concorrentes; instale apenas um. |
| Serviço local | Sim |

O docker-compose legado (v1) não é instalado; use 'docker compose'.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "docker-engine",
  "name": "Docker Engine (oficial)",
  "area": "dev",
  "category": "Contentores",
  "summary": "Docker Engine + Compose v2 (docker compose) do repositório oficial.",
  "why": "Recomendado para desenvolvimento Laravel com Sail/Compose. Corre nativamente, sem VM.",
  "method": "apt-repo",
  "verify": {
    "cmd": [
      "docker",
      "--version"
    ]
  },
  "packages": [
    "docker-ce",
    "docker-ce-cli",
    "containerd.io",
    "docker-buildx-plugin",
    "docker-compose-plugin"
  ],
  "repo": "docker",
  "source": "Repositório oficial Docker (docs.docker.com)",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "conflicts": [
    {
      "id": "docker-ubuntu",
      "hard": true,
      "reason": "docker.io (Ubuntu) e docker-ce (Docker) são pacotes concorrentes; instale apenas um."
    }
  ],
  "notes": "O docker-compose legado (v1) não é instalado; use 'docker compose'.",
  "compat": {
    "arch": [
      "amd64"
    ]
  },
  "service": true
}
```

</details>

<a id="item-docker-ubuntu"></a>

### Docker (pacotes Ubuntu)

Docker e Compose v2 mantidos pelo Ubuntu.

Alternativa sem repositório externo; versões um pouco mais antigas.

| Campo | Valor |
| --- | --- |
| ID | docker-ubuntu |
| Área / categoria | Desenvolvimento / Contentores |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | docker.io, docker-compose-v2, docker-buildx |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Conflitos | [Docker Engine (oficial)](#item-docker-engine) (bloqueia): Concorrente do Docker Engine oficial. |
| Avançado | Sim |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "docker-ubuntu",
  "name": "Docker (pacotes Ubuntu)",
  "area": "dev",
  "category": "Contentores",
  "summary": "Docker e Compose v2 mantidos pelo Ubuntu.",
  "why": "Alternativa sem repositório externo; versões um pouco mais antigas.",
  "method": "apt",
  "packages": [
    "docker.io",
    "docker-compose-v2",
    "docker-buildx"
  ],
  "conflicts": [
    {
      "id": "docker-engine",
      "hard": true,
      "reason": "Concorrente do Docker Engine oficial."
    }
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)",
  "service": true
}
```

</details>

<a id="item-docker-desktop"></a>

### Docker Desktop

Docker numa VM com interface gráfica (licença comercial para empresas >250 pessoas).

Diferenças: corre numa VM própria (KVM), socket e contexto próprios (desktop-linux); pode coexistir com o Engine mas os contentores/volumes não são partilhados.

| Campo | Valor |
| --- | --- |
| ID | docker-desktop |
| Área / categoria | Desenvolvimento / Contentores |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Requisitos: KVM ativo, 4 GB RAM, repositório oficial Docker configurado (ative 'Docker Engine' primeiro), 'pass' inicializado para credenciais.
2. Transfira https://desktop.docker.com/linux/main/amd64/docker-desktop-amd64.deb (≈465 MB) e verifique a checksum indicada na página.
3. Depois de conferir a origem e checksum, instale o pacote transferido:
```sh
sudo apt install ./docker-desktop-amd64.deb
```
4. Mude de contexto com: docker context use desktop-linux / default

[Documentação do fornecedor](https://docs.docker.com/desktop/setup/install/linux/ubuntu/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "docker-desktop",
  "name": "Docker Desktop",
  "area": "dev",
  "category": "Contentores",
  "summary": "Docker numa VM com interface gráfica (licença comercial para empresas >250 pessoas).",
  "why": "Diferenças: corre numa VM própria (KVM), socket e contexto próprios (desktop-linux); pode coexistir com o Engine mas os contentores/volumes não são partilhados.",
  "method": "manual",
  "detect": {
    "files": [
      "/opt/docker-desktop/bin/docker-desktop"
    ]
  },
  "manual": {
    "steps": [
      "Requisitos: KVM ativo, 4 GB RAM, repositório oficial Docker configurado (ative 'Docker Engine' primeiro), 'pass' inicializado para credenciais.",
      "Transfira https://desktop.docker.com/linux/main/amd64/docker-desktop-amd64.deb (≈465 MB) e verifique a checksum indicada na página.",
      "Depois de conferir a origem e checksum, instale o pacote transferido:\n```sh\nsudo apt install ./docker-desktop-amd64.deb\n```",
      "Mude de contexto com: docker context use desktop-linux / default"
    ],
    "url": "https://docs.docker.com/desktop/setup/install/linux/ubuntu/"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-docker-group"></a>

### Docker sem sudo (grupo docker)

Adiciona o seu utilizador ao grupo docker.

Conveniente, mas o grupo docker equivale a acesso root (qualquer contentor pode montar /). Efetivo após terminar sessão.

| Campo | Valor |
| --- | --- |
| ID | docker-group |
| Área / categoria | Desenvolvimento / Contentores |
| Método | admin |
| Fonte | usermod -aG docker (documentação pós-instalação Docker) |
| Atenção | Acesso equivalente a root para o utilizador (grupo docker) |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "docker-group",
  "name": "Docker sem sudo (grupo docker)",
  "area": "dev",
  "category": "Contentores",
  "summary": "Adiciona o seu utilizador ao grupo docker.",
  "why": "Conveniente, mas o grupo docker equivale a acesso root (qualquer contentor pode montar /). Efetivo após terminar sessão.",
  "method": "admin",
  "adminOp": "docker-group",
  "sensitive": "Acesso equivalente a root para o utilizador (grupo docker)",
  "source": "usermod -aG docker (documentação pós-instalação Docker)",
  "relogin": true
}
```

</details>

<a id="item-vscode"></a>

### Visual Studio Code

Editor de código da Microsoft.

Extensões para PHP/Laravel, Docker e Git. Repositório oficial atualiza com o sistema.

| Campo | Valor |
| --- | --- |
| ID | vscode |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | apt-repo |
| Fonte | Repositório oficial Microsoft |
| Pacotes | code |
| Repositório | vscode |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "vscode",
  "name": "Visual Studio Code",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "Editor de código da Microsoft.",
  "why": "Extensões para PHP/Laravel, Docker e Git. Repositório oficial atualiza com o sistema.",
  "method": "apt-repo",
  "verify": {
    "cmd": [
      "code",
      "--version"
    ]
  },
  "packages": [
    "code"
  ],
  "repo": "vscode",
  "source": "Repositório oficial Microsoft",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "tags": [
    "snap:code"
  ],
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-phpstorm"></a>

### PhpStorm

IDE PHP da JetBrains (pago).

O mais completo para Laravel. Instale pelo JetBrains Toolbox ou Flatpak (com.jetbrains.PhpStorm, não oficial).

| Campo | Valor |
| --- | --- |
| ID | phpstorm |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Recomendado: JetBrains Toolbox (https://www.jetbrains.com/toolbox-app/), extraia o .tar.gz e execute.
2. Alternativa Flathub: com.jetbrains.PhpStorm (empacotamento comunitário).

[Documentação do fornecedor](https://www.jetbrains.com/help/phpstorm/installation-guide.html)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "phpstorm",
  "name": "PhpStorm",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "IDE PHP da JetBrains (pago).",
  "why": "O mais completo para Laravel. Instale pelo JetBrains Toolbox ou Flatpak (com.jetbrains.PhpStorm, não oficial).",
  "method": "manual",
  "detect": {
    "commands": [
      "phpstorm"
    ]
  },
  "manual": {
    "steps": [
      "Recomendado: JetBrains Toolbox (https://www.jetbrains.com/toolbox-app/), extraia o .tar.gz e execute.",
      "Alternativa Flathub: com.jetbrains.PhpStorm (empacotamento comunitário)."
    ],
    "url": "https://www.jetbrains.com/help/phpstorm/installation-guide.html"
  },
  "source": "Manual / não suportado automaticamente",
  "tags": [
    "snap:phpstorm"
  ]
}
```

</details>

<a id="item-stripe-cli"></a>

### Stripe CLI

Testar webhooks e eventos Stripe localmente.

stripe listen encaminha webhooks para o seu Laravel local.

| Campo | Valor |
| --- | --- |
| ID | stripe-cli |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | apt-repo |
| Fonte | Repositório oficial Stripe |
| Pacotes | stripe |
| Repositório | stripe |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "stripe-cli",
  "name": "Stripe CLI",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "Testar webhooks e eventos Stripe localmente.",
  "why": "stripe listen encaminha webhooks para o seu Laravel local.",
  "method": "apt-repo",
  "verify": {
    "cmd": [
      "stripe",
      "version"
    ]
  },
  "packages": [
    "stripe"
  ],
  "repo": "stripe",
  "source": "Repositório oficial Stripe",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-bruno"></a>

### Bruno

Cliente de APIs offline, coleções em ficheiros.

Coleções versionáveis no Git, sem conta cloud.

| Campo | Valor |
| --- | --- |
| ID | bruno |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | flatpak |
| Fonte | Flathub (editor não verificado) |
| Flatpak | com.usebruno.Bruno |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "bruno",
  "name": "Bruno",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "Cliente de APIs offline, coleções em ficheiros.",
  "why": "Coleções versionáveis no Git, sem conta cloud.",
  "method": "flatpak",
  "flatpak": "com.usebruno.Bruno",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor não verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-hoppscotch"></a>

### Hoppscotch Desktop

Cliente de APIs open source.

Não está no Flathub; o pacote oficial é distribuído no GitHub.

| Campo | Valor |
| --- | --- |
| ID | hoppscotch |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Transfira o .deb/AppImage oficial em https://hoppscotch.com/download e verifique a origem.
2. Alternativa sem instalação: https://hoppscotch.io

[Documentação do fornecedor](https://hoppscotch.com/download)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "hoppscotch",
  "name": "Hoppscotch Desktop",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "Cliente de APIs open source.",
  "why": "Não está no Flathub; o pacote oficial é distribuído no GitHub.",
  "method": "manual",
  "detect": {
    "commands": [
      "hoppscotch"
    ]
  },
  "manual": {
    "steps": [
      "Transfira o .deb/AppImage oficial em https://hoppscotch.com/download e verifique a origem.",
      "Alternativa sem instalação: https://hoppscotch.io"
    ],
    "url": "https://hoppscotch.com/download"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-warp"></a>

### Warp

Terminal moderno com IA (requer conta).

Blocos de comandos e IA integrada. Telemetria e conta são da Warp — avalie.

| Campo | Valor |
| --- | --- |
| ID | warp |
| Área / categoria | Desenvolvimento / Terminal |
| Método | apt-repo |
| Fonte | Repositório oficial Warp |
| Pacotes | warp-terminal |
| Repositório | warp |
| Evidência da receita | Chave/fingerprint e dists verificados em 2026-10-04 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "warp",
  "name": "Warp",
  "area": "dev",
  "category": "Terminal",
  "summary": "Terminal moderno com IA (requer conta).",
  "why": "Blocos de comandos e IA integrada. Telemetria e conta são da Warp — avalie.",
  "method": "apt-repo",
  "packages": [
    "warp-terminal"
  ],
  "repo": "warp",
  "source": "Repositório oficial Warp",
  "verified": "Chave/fingerprint e dists verificados em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-devin"></a>

### Devin Desktop

Aplicação desktop do Devin.

Sem receita Linux oficial verificada.

| Campo | Valor |
| --- | --- |
| ID | devin |
| Área / categoria | Desenvolvimento / Editores e APIs |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Não suportado: verifique em https://devin.ai se existe instalador Linux oficial.

[Documentação do fornecedor](https://devin.ai)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "devin",
  "name": "Devin Desktop",
  "area": "dev",
  "category": "Editores e APIs",
  "summary": "Aplicação desktop do Devin.",
  "why": "Sem receita Linux oficial verificada.",
  "method": "manual",
  "manual": {
    "steps": [
      "Não suportado: verifique em https://devin.ai se existe instalador Linux oficial."
    ],
    "url": "https://devin.ai"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-git-defaults"></a>

### Git: predefinições úteis

Ramo inicial, pull e limpeza de ramos remotos.

Evita avisos e mantém o repositório local arrumado. Valores anteriores ficam registados para reverter.

| Campo | Valor |
| --- | --- |
| ID | git-defaults |
| Área / categoria | Desenvolvimento / Git |
| Método | gitconfig |
| Fonte | git config --global |
| Dependências | [Git](#item-git) |

Opção **Ramo inicial** (`branch`): um valor. Padrão: main. 

| Valor | Descrição |
| --- | --- |
| main | main |
| master | master |

Opção **git pull** (`rebase`): um valor. Padrão: false. 

| Valor | Descrição |
| --- | --- |
| false | merge (padrão) |
| true | rebase |

Opção **Editor** (`editor`): um valor. Padrão: nano. 

| Valor | Descrição |
| --- | --- |
| nano | nano |
| vim | vim |
| code --wait | VS Code |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "git-defaults",
  "name": "Git: predefinições úteis",
  "area": "dev",
  "category": "Git",
  "summary": "Ramo inicial, pull e limpeza de ramos remotos.",
  "why": "Evita avisos e mantém o repositório local arrumado. Valores anteriores ficam registados para reverter.",
  "method": "gitconfig",
  "requires": [
    "git"
  ],
  "source": "git config --global",
  "git": [
    {
      "key": "init.defaultBranch",
      "option": "branch"
    },
    {
      "key": "pull.rebase",
      "option": "rebase"
    },
    {
      "key": "fetch.prune",
      "value": "true"
    },
    {
      "key": "core.editor",
      "option": "editor"
    }
  ],
  "options": [
    {
      "key": "branch",
      "label": "Ramo inicial",
      "values": [
        {
          "value": "main",
          "label": "main"
        },
        {
          "value": "master",
          "label": "master"
        }
      ],
      "default": [
        "main"
      ]
    },
    {
      "key": "rebase",
      "label": "git pull",
      "values": [
        {
          "value": "false",
          "label": "merge (padrão)"
        },
        {
          "value": "true",
          "label": "rebase"
        }
      ],
      "default": [
        "false"
      ]
    },
    {
      "key": "editor",
      "label": "Editor",
      "values": [
        {
          "value": "nano",
          "label": "nano"
        },
        {
          "value": "vim",
          "label": "vim"
        },
        {
          "value": "code --wait",
          "label": "VS Code"
        }
      ],
      "default": [
        "nano"
      ]
    }
  ]
}
```

</details>

<a id="item-git-identity"></a>

### Git: identidade por diretório

Identidade pessoal e profissional automáticas por pasta.

Commits do trabalho com o email da empresa sem misturar com projetos pessoais.

| Campo | Valor |
| --- | --- |
| ID | git-identity |
| Área / categoria | Desenvolvimento / Git |
| Método | manual |
| Fonte | Manual (dados pessoais não são guardados em perfis) |

#### Guia manual

1. Crie ~/.gitconfig-trabalho com [user] name/email profissional.
2. Acrescente ao ~/.gitconfig:
```ini
[includeIf "gitdir:~/www/"]
    path = ~/.gitconfig-trabalho
```
3. Confirme dentro de um repositório: git config user.email

[Documentação do fornecedor](https://git-scm.com/docs/git-config#_conditional_includes)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "git-identity",
  "name": "Git: identidade por diretório",
  "area": "dev",
  "category": "Git",
  "summary": "Identidade pessoal e profissional automáticas por pasta.",
  "why": "Commits do trabalho com o email da empresa sem misturar com projetos pessoais.",
  "method": "manual",
  "manual": {
    "steps": [
      "Crie ~/.gitconfig-trabalho com [user] name/email profissional.",
      "Acrescente ao ~/.gitconfig:\n```ini\n[includeIf \"gitdir:~/www/\"]\n    path = ~/.gitconfig-trabalho\n```",
      "Confirme dentro de um repositório: git config user.email"
    ],
    "url": "https://git-scm.com/docs/git-config#_conditional_includes"
  },
  "source": "Manual (dados pessoais não são guardados em perfis)"
}
```

</details>

<a id="item-inotify-watches"></a>

### Mais watchers inotify

Aumenta fs.inotify.max_user_watches.

Só se o VS Code/Vite avisar 'unable to watch'. O padrão atual do kernel já é proporcional à RAM.

| Campo | Valor |
| --- | --- |
| ID | inotify-watches |
| Área / categoria | Desenvolvimento / Sistema |
| Método | admin |
| Fonte | /etc/sysctl.d (drop-in próprio, reversível) |
| Avançado | Sim |

Opção **Valor** (`value`): um valor. Padrão: 262144. 

| Valor | Descrição |
| --- | --- |
| 262144 | 262144 |
| 524288 | 524288 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "inotify-watches",
  "name": "Mais watchers inotify",
  "area": "dev",
  "category": "Sistema",
  "summary": "Aumenta fs.inotify.max_user_watches.",
  "why": "Só se o VS Code/Vite avisar 'unable to watch'. O padrão atual do kernel já é proporcional à RAM.",
  "method": "admin",
  "adminOp": "inotify",
  "source": "/etc/sysctl.d (drop-in próprio, reversível)",
  "advanced": true,
  "options": [
    {
      "key": "value",
      "label": "Valor",
      "values": [
        {
          "value": "262144",
          "label": "262144"
        },
        {
          "value": "524288",
          "label": "524288"
        }
      ],
      "default": [
        "262144"
      ]
    }
  ]
}
```

</details>

<a id="item-color-scheme"></a>

### Tema claro/escuro

Estilo preferido das aplicações.

Tema escuro reduz brilho à noite; claro é mais legível em ambientes iluminados.

| Campo | Valor |
| --- | --- |
| ID | color-scheme |
| Área / categoria | Estilo / Cores |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Estilo** (`scheme`): um valor. Padrão: 'prefer-dark'. 

| Valor | Descrição |
| --- | --- |
| 'prefer-dark' | Escuro |
| 'default' | Claro (padrão) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "color-scheme",
  "name": "Tema claro/escuro",
  "area": "style",
  "category": "Cores",
  "summary": "Estilo preferido das aplicações.",
  "why": "Tema escuro reduz brilho à noite; claro é mais legível em ambientes iluminados.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "color-scheme",
      "option": "scheme"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "scheme",
      "label": "Estilo",
      "values": [
        {
          "value": "'prefer-dark'",
          "label": "Escuro"
        },
        {
          "value": "'default'",
          "label": "Claro (padrão)"
        }
      ],
      "default": [
        "'prefer-dark'"
      ]
    }
  ]
}
```

</details>

<a id="item-text-scaling"></a>

### Tamanho do texto

Fator de escala do texto.

Ecrãs com alta densidade ou conforto visual sem mudar a escala inteira.

| Campo | Valor |
| --- | --- |
| ID | text-scaling |
| Área / categoria | Desktop / Acessibilidade |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Escala** (`scale`): um valor. Padrão: 1.0. 

| Valor | Descrição |
| --- | --- |
| 1.0 | 100% |
| 1.1 | 110% |
| 1.25 | 125% |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "text-scaling",
  "name": "Tamanho do texto",
  "area": "desktop",
  "category": "Acessibilidade",
  "summary": "Fator de escala do texto.",
  "why": "Ecrãs com alta densidade ou conforto visual sem mudar a escala inteira.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "text-scaling-factor",
      "option": "scale"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "scale",
      "label": "Escala",
      "values": [
        {
          "value": "1.0",
          "label": "100%"
        },
        {
          "value": "1.1",
          "label": "110%"
        },
        {
          "value": "1.25",
          "label": "125%"
        }
      ],
      "default": [
        "1.0"
      ]
    }
  ]
}
```

</details>

<a id="item-mono-font"></a>

### Fonte monoespaçada

Fonte usada em terminais e editores que seguem o sistema.

Tamanho confortável para longas sessões de código.

| Campo | Valor |
| --- | --- |
| ID | mono-font |
| Área / categoria | Estilo / Fontes |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Fonte** (`font`): um valor. Padrão: 'Ubuntu Sans Mono 13'. 

| Valor | Descrição |
| --- | --- |
| 'Ubuntu Sans Mono 13' | Ubuntu Sans Mono 13 (padrão) |
| 'Ubuntu Sans Mono 12' | Ubuntu Sans Mono 12 |
| 'Ubuntu Sans Mono 14' | Ubuntu Sans Mono 14 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "mono-font",
  "name": "Fonte monoespaçada",
  "area": "style",
  "category": "Fontes",
  "summary": "Fonte usada em terminais e editores que seguem o sistema.",
  "why": "Tamanho confortável para longas sessões de código.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "monospace-font-name",
      "option": "font"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "font",
      "label": "Fonte",
      "values": [
        {
          "value": "'Ubuntu Sans Mono 13'",
          "label": "Ubuntu Sans Mono 13 (padrão)"
        },
        {
          "value": "'Ubuntu Sans Mono 12'",
          "label": "Ubuntu Sans Mono 12"
        },
        {
          "value": "'Ubuntu Sans Mono 14'",
          "label": "Ubuntu Sans Mono 14"
        }
      ],
      "default": [
        "'Ubuntu Sans Mono 13'"
      ]
    }
  ]
}
```

</details>

<a id="item-reduce-animations"></a>

### Reduzir animações

Desliga animações do shell.

Interface mais imediata e menos movimento (acessibilidade).

| Campo | Valor |
| --- | --- |
| ID | reduce-animations |
| Área / categoria | Desktop / Acessibilidade |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "reduce-animations",
  "name": "Reduzir animações",
  "area": "desktop",
  "category": "Acessibilidade",
  "summary": "Desliga animações do shell.",
  "why": "Interface mais imediata e menos movimento (acessibilidade).",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "enable-animations",
      "value": "false"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-night-light"></a>

### Luz noturna

Tons quentes ao fim do dia.

Menos luz azul à noite.

| Campo | Valor |
| --- | --- |
| ID | night-light |
| Área / categoria | Desktop / Acessibilidade |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "night-light",
  "name": "Luz noturna",
  "area": "desktop",
  "category": "Acessibilidade",
  "summary": "Tons quentes ao fim do dia.",
  "why": "Menos luz azul à noite.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.settings-daemon.plugins.color",
      "key": "night-light-enabled",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-battery-pct"></a>

### Percentagem da bateria

Mostra a percentagem da bateria.

Saber quanto falta sem abrir o menu.

| Campo | Valor |
| --- | --- |
| ID | battery-pct |
| Área / categoria | Desktop / Barra superior |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "battery-pct",
  "name": "Percentagem da bateria",
  "area": "desktop",
  "category": "Barra superior",
  "summary": "Mostra a percentagem da bateria.",
  "why": "Saber quanto falta sem abrir o menu.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "show-battery-percentage",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-clock-weekday"></a>

### Dia da semana no relógio

Mostra o dia da semana.

Contexto rápido de data.

| Campo | Valor |
| --- | --- |
| ID | clock-weekday |
| Área / categoria | Desktop / Barra superior |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "clock-weekday",
  "name": "Dia da semana no relógio",
  "area": "desktop",
  "category": "Barra superior",
  "summary": "Mostra o dia da semana.",
  "why": "Contexto rápido de data.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "clock-show-weekday",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-dock-position"></a>

### Dock: posição

Onde fica a Ubuntu Dock.

Em ecrãs largos, à esquerda poupa altura; em baixo é mais familiar.

| Campo | Valor |
| --- | --- |
| ID | dock-position |
| Área / categoria | Desktop / Dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Posição** (`pos`): um valor. Padrão: 'BOTTOM'. 

| Valor | Descrição |
| --- | --- |
| 'LEFT' | Esquerda (padrão) |
| 'BOTTOM' | Em baixo |
| 'RIGHT' | Direita |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-position",
  "name": "Dock: posição",
  "area": "desktop",
  "category": "Dock",
  "summary": "Onde fica a Ubuntu Dock.",
  "why": "Em ecrãs largos, à esquerda poupa altura; em baixo é mais familiar.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "dock-position",
      "option": "pos"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "pos",
      "label": "Posição",
      "values": [
        {
          "value": "'LEFT'",
          "label": "Esquerda (padrão)"
        },
        {
          "value": "'BOTTOM'",
          "label": "Em baixo"
        },
        {
          "value": "'RIGHT'",
          "label": "Direita"
        }
      ],
      "default": [
        "'BOTTOM'"
      ]
    }
  ]
}
```

</details>

<a id="item-dock-size"></a>

### Dock: tamanho dos ícones

Tamanho máximo dos ícones.

Ícones menores libertam espaço em portáteis.

| Campo | Valor |
| --- | --- |
| ID | dock-size |
| Área / categoria | Desktop / Dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Tamanho** (`size`): um valor. Padrão: 40. 

| Valor | Descrição |
| --- | --- |
| 32 | 32 px |
| 40 | 40 px |
| 48 | 48 px (padrão) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-size",
  "name": "Dock: tamanho dos ícones",
  "area": "desktop",
  "category": "Dock",
  "summary": "Tamanho máximo dos ícones.",
  "why": "Ícones menores libertam espaço em portáteis.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "dash-max-icon-size",
      "option": "size"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "size",
      "label": "Tamanho",
      "values": [
        {
          "value": "32",
          "label": "32 px"
        },
        {
          "value": "40",
          "label": "40 px"
        },
        {
          "value": "48",
          "label": "48 px (padrão)"
        }
      ],
      "default": [
        "40"
      ]
    }
  ]
}
```

</details>

<a id="item-dock-autohide"></a>

### Dock: ocultar automaticamente

Esconde a dock quando uma janela a tapa.

Mais espaço útil para o editor e o navegador.

| Campo | Valor |
| --- | --- |
| ID | dock-autohide |
| Área / categoria | Desktop / Dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-autohide",
  "name": "Dock: ocultar automaticamente",
  "area": "desktop",
  "category": "Dock",
  "summary": "Esconde a dock quando uma janela a tapa.",
  "why": "Mais espaço útil para o editor e o navegador.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "dock-fixed",
      "value": "false"
    },
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "intellihide",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-dock-click"></a>

### Dock: clique minimiza

O que acontece ao clicar numa app aberta.

Clicar para minimizar/mostrar pré-visualizações é mais rápido com várias janelas.

| Campo | Valor |
| --- | --- |
| ID | dock-click |
| Área / categoria | Desktop / Dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Ação** (`action`): um valor. Padrão: 'minimize-or-previews'. 

| Valor | Descrição |
| --- | --- |
| 'minimize-or-previews' | Minimizar ou pré-visualizar |
| 'focus-or-previews' | Focar ou pré-visualizar |
| 'previews' | Pré-visualizar |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-click",
  "name": "Dock: clique minimiza",
  "area": "desktop",
  "category": "Dock",
  "summary": "O que acontece ao clicar numa app aberta.",
  "why": "Clicar para minimizar/mostrar pré-visualizações é mais rápido com várias janelas.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "click-action",
      "option": "action"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "action",
      "label": "Ação",
      "values": [
        {
          "value": "'minimize-or-previews'",
          "label": "Minimizar ou pré-visualizar"
        },
        {
          "value": "'focus-or-previews'",
          "label": "Focar ou pré-visualizar"
        },
        {
          "value": "'previews'",
          "label": "Pré-visualizar"
        }
      ],
      "default": [
        "'minimize-or-previews'"
      ]
    }
  ]
}
```

</details>

<a id="item-dock-multimonitor"></a>

### Dock: em todos os monitores

Mostra a dock em cada monitor.

Útil com monitor externo.

| Campo | Valor |
| --- | --- |
| ID | dock-multimonitor |
| Área / categoria | Desktop / Dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-multimonitor",
  "name": "Dock: em todos os monitores",
  "area": "desktop",
  "category": "Dock",
  "summary": "Mostra a dock em cada monitor.",
  "why": "Útil com monitor externo.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "multi-monitor",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-window-buttons"></a>

### Botões minimizar/maximizar

Garante os três botões nas janelas.

Padrão do Ubuntu; repõe se outra configuração o removeu.

| Campo | Valor |
| --- | --- |
| ID | window-buttons |
| Área / categoria | Desktop / Janelas |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "window-buttons",
  "name": "Botões minimizar/maximizar",
  "area": "desktop",
  "category": "Janelas",
  "summary": "Garante os três botões nas janelas.",
  "why": "Padrão do Ubuntu; repõe se outra configuração o removeu.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.wm.preferences",
      "key": "button-layout",
      "value": "'appmenu:minimize,maximize,close'"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-center-windows"></a>

### Centrar novas janelas

Janelas novas abrem ao centro.

Previsível em monitores grandes.

| Campo | Valor |
| --- | --- |
| ID | center-windows |
| Área / categoria | Desktop / Janelas |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "center-windows",
  "name": "Centrar novas janelas",
  "area": "desktop",
  "category": "Janelas",
  "summary": "Janelas novas abrem ao centro.",
  "why": "Previsível em monitores grandes.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.mutter",
      "key": "center-new-windows",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-workspaces-all"></a>

### Áreas de trabalho em todos os monitores

Trocar de área de trabalho muda todos os ecrãs.

Organizar projetos por área de trabalho com dois monitores.

| Campo | Valor |
| --- | --- |
| ID | workspaces-all |
| Área / categoria | Desktop / Janelas |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "workspaces-all",
  "name": "Áreas de trabalho em todos os monitores",
  "area": "desktop",
  "category": "Janelas",
  "summary": "Trocar de área de trabalho muda todos os ecrãs.",
  "why": "Organizar projetos por área de trabalho com dois monitores.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.mutter",
      "key": "workspaces-only-on-primary",
      "value": "false"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-super-e"></a>

### Super+E abre a pasta pessoal

Atalho para o gestor de ficheiros.

Hábito vindo do Windows; rápido.

| Campo | Valor |
| --- | --- |
| ID | super-e |
| Área / categoria | Desktop / Atalhos |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "super-e",
  "name": "Super+E abre a pasta pessoal",
  "area": "desktop",
  "category": "Atalhos",
  "summary": "Atalho para o gestor de ficheiros.",
  "why": "Hábito vindo do Windows; rápido.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.settings-daemon.plugins.media-keys",
      "key": "home",
      "value": "['<Super>e']"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-tap-click"></a>

### Tocar para clicar (touchpad)

Toque leve no touchpad conta como clique.

Mais confortável em portáteis.

| Campo | Valor |
| --- | --- |
| ID | tap-click |
| Área / categoria | Desktop / Atalhos |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tap-click",
  "name": "Tocar para clicar (touchpad)",
  "area": "desktop",
  "category": "Atalhos",
  "summary": "Toque leve no touchpad conta como clique.",
  "why": "Mais confortável em portáteis.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.peripherals.touchpad",
      "key": "tap-to-click",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-nautilus-list"></a>

### Ficheiros: vista em lista

Vista em lista com colunas por omissão.

Ver tamanhos e datas de imediato em pastas de projetos.

| Campo | Valor |
| --- | --- |
| ID | nautilus-list |
| Área / categoria | Desktop / Ficheiros |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nautilus-list",
  "name": "Ficheiros: vista em lista",
  "area": "desktop",
  "category": "Ficheiros",
  "summary": "Vista em lista com colunas por omissão.",
  "why": "Ver tamanhos e datas de imediato em pastas de projetos.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.nautilus.preferences",
      "key": "default-folder-viewer",
      "value": "'list-view'"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-folders-first"></a>

### Ficheiros: pastas primeiro

Pastas antes dos ficheiros.

Navegação mais rápida.

| Campo | Valor |
| --- | --- |
| ID | folders-first |
| Área / categoria | Desktop / Ficheiros |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "folders-first",
  "name": "Ficheiros: pastas primeiro",
  "area": "desktop",
  "category": "Ficheiros",
  "summary": "Pastas antes dos ficheiros.",
  "why": "Navegação mais rápida.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gtk.gtk4.Settings.FileChooser",
      "key": "sort-directories-first",
      "value": "true"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-notif-lock"></a>

### Notificações fora do ecrã de bloqueio

Não mostra notificações com o ecrã bloqueado.

Privacidade: mensagens não ficam visíveis a quem passa.

| Campo | Valor |
| --- | --- |
| ID | notif-lock |
| Área / categoria | Desktop / Notificações |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "notif-lock",
  "name": "Notificações fora do ecrã de bloqueio",
  "area": "desktop",
  "category": "Notificações",
  "summary": "Não mostra notificações com o ecrã bloqueado.",
  "why": "Privacidade: mensagens não ficam visíveis a quem passa.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.notifications",
      "key": "show-in-lock-screen",
      "value": "false"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-tpl-txt"></a>

### Modelo: Documento de texto.txt

Ficheiro de texto vazio no menu «Novo documento».

Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.

| Campo | Valor |
| --- | --- |
| ID | tpl-txt |
| Área / categoria | Desktop / Novo documento |
| Método | template |
| Fonte | ~/Modelos (XDG_TEMPLATES_DIR) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tpl-txt",
  "name": "Modelo: Documento de texto.txt",
  "area": "desktop",
  "category": "Novo documento",
  "summary": "Ficheiro de texto vazio no menu «Novo documento».",
  "why": "Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.",
  "method": "template",
  "template": {
    "name": "Documento de texto.txt",
    "content": ""
  },
  "source": "~/Modelos (XDG_TEMPLATES_DIR)"
}
```

</details>

<a id="item-tpl-md"></a>

### Modelo: Nota.md

Nota Markdown no menu «Novo documento».

Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.

| Campo | Valor |
| --- | --- |
| ID | tpl-md |
| Área / categoria | Desktop / Novo documento |
| Método | template |
| Fonte | ~/Modelos (XDG_TEMPLATES_DIR) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tpl-md",
  "name": "Modelo: Nota.md",
  "area": "desktop",
  "category": "Novo documento",
  "summary": "Nota Markdown no menu «Novo documento».",
  "why": "Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.",
  "method": "template",
  "template": {
    "name": "Nota.md",
    "content": "# Título\n\n"
  },
  "source": "~/Modelos (XDG_TEMPLATES_DIR)"
}
```

</details>

<a id="item-tpl-sh"></a>

### Modelo: Script.sh

Script Bash com modo estrito no menu «Novo documento».

Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.

| Campo | Valor |
| --- | --- |
| ID | tpl-sh |
| Área / categoria | Desktop / Novo documento |
| Método | template |
| Fonte | ~/Modelos (XDG_TEMPLATES_DIR) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tpl-sh",
  "name": "Modelo: Script.sh",
  "area": "desktop",
  "category": "Novo documento",
  "summary": "Script Bash com modo estrito no menu «Novo documento».",
  "why": "Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.",
  "method": "template",
  "template": {
    "name": "Script.sh",
    "content": "#!/usr/bin/env bash\nset -euo pipefail\n\n"
  },
  "source": "~/Modelos (XDG_TEMPLATES_DIR)"
}
```

</details>

<a id="item-tpl-php"></a>

### Modelo: Classe.php

Ficheiro PHP com strict_types no menu «Novo documento».

Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.

| Campo | Valor |
| --- | --- |
| ID | tpl-php |
| Área / categoria | Desktop / Novo documento |
| Método | template |
| Fonte | ~/Modelos (XDG_TEMPLATES_DIR) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tpl-php",
  "name": "Modelo: Classe.php",
  "area": "desktop",
  "category": "Novo documento",
  "summary": "Ficheiro PHP com strict_types no menu «Novo documento».",
  "why": "Clique direito → Novo documento no gestor de ficheiros. Nunca sobrescreve modelos existentes.",
  "method": "template",
  "template": {
    "name": "Classe.php",
    "content": "<?php\n\ndeclare(strict_types=1);\n\n"
  },
  "source": "~/Modelos (XDG_TEMPLATES_DIR)"
}
```

</details>

<a id="item-startup-apps"></a>

### Aplicações de arranque

Rever o que arranca com a sessão.

A página Diagnóstico lista as entradas; desative só as que reconhece.

| Campo | Valor |
| --- | --- |
| ID | startup-apps |
| Área / categoria | Desktop / Arranque |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Abra GNOME Tweaks → Aplicações de arranque, ou Definições → Aplicações → (app) → Arrancar com a sessão.
2. Não desative serviços desconhecidos do sistema.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "startup-apps",
  "name": "Aplicações de arranque",
  "area": "desktop",
  "category": "Arranque",
  "summary": "Rever o que arranca com a sessão.",
  "why": "A página Diagnóstico lista as entradas; desative só as que reconhece.",
  "method": "manual",
  "manual": {
    "steps": [
      "Abra GNOME Tweaks → Aplicações de arranque, ou Definições → Aplicações → (app) → Arrancar com a sessão.",
      "Não desative serviços desconhecidos do sistema."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-accent-color"></a>

### Cor de destaque

Cor de botões, seleções e destaques (GNOME 47+).

Personalizar sem temas de terceiros; o Ubuntu usa laranja por omissão e respeita esta chave.

| Campo | Valor |
| --- | --- |
| ID | accent-color |
| Área / categoria | Estilo / Cores |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |
| Evidência da receita | Chave e valores lidos do schema no Ubuntu 26.04; override Ubuntu = 'orange' |

Opção **Cor** (`accent`): um valor. Padrão: 'blue'. 

| Valor | Descrição |
| --- | --- |
| 'orange' | Laranja (padrão Ubuntu) |
| 'blue' | Azul |
| 'teal' | Verde-azulado |
| 'green' | Verde |
| 'yellow' | Amarelo |
| 'red' | Vermelho |
| 'pink' | Rosa |
| 'purple' | Roxo |
| 'slate' | Ardósia |
| 'brown' | Castanho |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "accent-color",
  "name": "Cor de destaque",
  "area": "style",
  "category": "Cores",
  "summary": "Cor de botões, seleções e destaques (GNOME 47+).",
  "why": "Personalizar sem temas de terceiros; o Ubuntu usa laranja por omissão e respeita esta chave.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "accent-color",
      "option": "accent"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "accent",
      "label": "Cor",
      "values": [
        {
          "value": "'orange'",
          "label": "Laranja (padrão Ubuntu)"
        },
        {
          "value": "'blue'",
          "label": "Azul"
        },
        {
          "value": "'teal'",
          "label": "Verde-azulado"
        },
        {
          "value": "'green'",
          "label": "Verde"
        },
        {
          "value": "'yellow'",
          "label": "Amarelo"
        },
        {
          "value": "'red'",
          "label": "Vermelho"
        },
        {
          "value": "'pink'",
          "label": "Rosa"
        },
        {
          "value": "'purple'",
          "label": "Roxo"
        },
        {
          "value": "'slate'",
          "label": "Ardósia"
        },
        {
          "value": "'brown'",
          "label": "Castanho"
        }
      ],
      "default": [
        "'blue'"
      ]
    }
  ],
  "verified": "Chave e valores lidos do schema no Ubuntu 26.04; override Ubuntu = 'orange'"
}
```

</details>

<a id="item-wallpapers-pack"></a>

### Coleção de wallpapers Ubuntu

Wallpapers oficiais das versões 24.04 a 26.04.

Mais escolhas de fundo, todas do arquivo Ubuntu.

| Campo | Valor |
| --- | --- |
| ID | wallpapers-pack |
| Área / categoria | Estilo / Wallpapers |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | ubuntu-wallpapers, ubuntu-wallpapers-resolute, ubuntu-wallpapers-questing, ubuntu-wallpapers-noble |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "wallpapers-pack",
  "name": "Coleção de wallpapers Ubuntu",
  "area": "style",
  "category": "Wallpapers",
  "summary": "Wallpapers oficiais das versões 24.04 a 26.04.",
  "why": "Mais escolhas de fundo, todas do arquivo Ubuntu.",
  "method": "apt",
  "packages": [
    "ubuntu-wallpapers",
    "ubuntu-wallpapers-resolute",
    "ubuntu-wallpapers-questing",
    "ubuntu-wallpapers-noble"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-wallpaper"></a>

### Wallpaper

Fundo do ambiente de trabalho, com versão clara e escura.

Cada opção define a imagem para o tema claro e para o escuro.

| Campo | Valor |
| --- | --- |
| ID | wallpaper |
| Área / categoria | Estilo / Wallpapers |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |
| Evidência da receita | Ficheiros listados nos pacotes ubuntu-wallpapers-* e gnome-backgrounds do 26.04 |
| Dependências | [Coleção de wallpapers Ubuntu](#item-wallpapers-pack) |

Opção **Imagem** (`wall`): um valor. Padrão: raccoon. 

| Valor | Descrição |
| --- | --- |
| raccoon | Resolute Raccoon (26.04) |
| simple-raccoon | Simple Raccoon |
| boxes | Flying Boxes |
| mug | Ubuntu Coffee Mug |
| quokka | Questing Quokka (25.10) |
| numbat | Noble Numbat (24.04) |
| adwaita | GNOME Adwaita |
| blobs | GNOME Blobs |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "wallpaper",
  "name": "Wallpaper",
  "area": "style",
  "category": "Wallpapers",
  "summary": "Fundo do ambiente de trabalho, com versão clara e escura.",
  "why": "Cada opção define a imagem para o tema claro e para o escuro.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.background",
      "key": "picture-uri",
      "option": "wall",
      "map": {
        "raccoon": "'file:///usr/share/backgrounds/Resolute_Raccoon_Wallpaper_Light_3840x2160.png'",
        "simple-raccoon": "'file:///usr/share/backgrounds/cnusr25-Simple_Raccoon_Light.png'",
        "boxes": "'file:///usr/share/backgrounds/endycal-Flying_Boxes_Light.png'",
        "mug": "'file:///usr/share/backgrounds/ezspain-Ubuntu_Coffee_Mug_Light.png'",
        "quokka": "'file:///usr/share/backgrounds/Questing_Quokka_Wallpaper_Light_3840x2160.png'",
        "numbat": "'file:///usr/share/backgrounds/Numbat_wallpaper_light_3480x2160.png'",
        "adwaita": "'file:///usr/share/backgrounds/gnome/adwaita-l.jxl'",
        "blobs": "'file:///usr/share/backgrounds/gnome/blobs-l.svg'"
      }
    },
    {
      "schema": "org.gnome.desktop.background",
      "key": "picture-uri-dark",
      "option": "wall",
      "map": {
        "raccoon": "'file:///usr/share/backgrounds/Resolute_Raccoon_Wallpaper_Dimmed_3840x2160.png'",
        "simple-raccoon": "'file:///usr/share/backgrounds/cnusr25-Simple_Raccoon_Dark.png'",
        "boxes": "'file:///usr/share/backgrounds/endycal-Flying_Boxes_Dark.png'",
        "mug": "'file:///usr/share/backgrounds/ezspain-Ubuntu_Coffee_Mug_Dark.png'",
        "quokka": "'file:///usr/share/backgrounds/Questing_Quokka_Wallpaper_Dimmed_3840x2160.png'",
        "numbat": "'file:///usr/share/backgrounds/Numbat_wallpaper_dimmed_3480x2160.png'",
        "adwaita": "'file:///usr/share/backgrounds/gnome/adwaita-d.jxl'",
        "blobs": "'file:///usr/share/backgrounds/gnome/blobs-d.svg'"
      }
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "requires": [
    "wallpapers-pack"
  ],
  "options": [
    {
      "key": "wall",
      "label": "Imagem",
      "values": [
        {
          "value": "raccoon",
          "label": "Resolute Raccoon (26.04)"
        },
        {
          "value": "simple-raccoon",
          "label": "Simple Raccoon"
        },
        {
          "value": "boxes",
          "label": "Flying Boxes"
        },
        {
          "value": "mug",
          "label": "Ubuntu Coffee Mug"
        },
        {
          "value": "quokka",
          "label": "Questing Quokka (25.10)"
        },
        {
          "value": "numbat",
          "label": "Noble Numbat (24.04)"
        },
        {
          "value": "adwaita",
          "label": "GNOME Adwaita"
        },
        {
          "value": "blobs",
          "label": "GNOME Blobs"
        }
      ],
      "default": [
        "raccoon"
      ]
    }
  ],
  "verified": "Ficheiros listados nos pacotes ubuntu-wallpapers-* e gnome-backgrounds do 26.04"
}
```

</details>

<a id="item-ptyxis-opacity"></a>

### Terminal: transparência (Ptyxis)

Opacidade do perfil predefinido do Terminal (Ptyxis, o padrão do Ubuntu 26.04).

Ver o editor ou a documentação por trás do terminal.

| Campo | Valor |
| --- | --- |
| ID | ptyxis-opacity |
| Área / categoria | Estilo / Terminal |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |
| Evidência da receita | Chave opacity (0.0–1.0) do schema org.gnome.Ptyxis.Profile no 26.04 |

Opção **Opacidade** (`opacity`): um valor. Padrão: 0.9. 

| Valor | Descrição |
| --- | --- |
| 1.0 | 100% (sem transparência) |
| 0.95 | 95% |
| 0.9 | 90% |
| 0.85 | 85% |
| 0.8 | 80% |

Requer que o Terminal tenha sido aberto pelo menos uma vez (cria o perfil predefinido).

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ptyxis-opacity",
  "name": "Terminal: transparência (Ptyxis)",
  "area": "style",
  "category": "Terminal",
  "summary": "Opacidade do perfil predefinido do Terminal (Ptyxis, o padrão do Ubuntu 26.04).",
  "why": "Ver o editor ou a documentação por trás do terminal.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.Ptyxis.Profile:/org/gnome/Ptyxis/Profiles/{ptyxis-profile}/",
      "key": "opacity",
      "option": "opacity"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "opacity",
      "label": "Opacidade",
      "values": [
        {
          "value": "1.0",
          "label": "100% (sem transparência)"
        },
        {
          "value": "0.95",
          "label": "95%"
        },
        {
          "value": "0.9",
          "label": "90%"
        },
        {
          "value": "0.85",
          "label": "85%"
        },
        {
          "value": "0.8",
          "label": "80%"
        }
      ],
      "default": [
        "0.9"
      ]
    }
  ],
  "notes": "Requer que o Terminal tenha sido aberto pelo menos uma vez (cria o perfil predefinido).",
  "verified": "Chave opacity (0.0–1.0) do schema org.gnome.Ptyxis.Profile no 26.04"
}
```

</details>

<a id="item-gnome-terminal-transparency"></a>

### Terminal: transparência (GNOME Terminal)

Transparência do perfil predefinido do GNOME Terminal.

Só se usar o GNOME Terminal em vez do Ptyxis.

| Campo | Valor |
| --- | --- |
| ID | gnome-terminal-transparency |
| Área / categoria | Estilo / Terminal |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |
| Avançado | Sim |

Opção **Transparência** (`pct`): um valor. Padrão: 15. 

| Valor | Descrição |
| --- | --- |
| 10 | 10% |
| 15 | 15% |
| 20 | 20% |
| 30 | 30% |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gnome-terminal-transparency",
  "name": "Terminal: transparência (GNOME Terminal)",
  "area": "style",
  "category": "Terminal",
  "summary": "Transparência do perfil predefinido do GNOME Terminal.",
  "why": "Só se usar o GNOME Terminal em vez do Ptyxis.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.Terminal.Legacy.Profile:/org/gnome/terminal/legacy/profiles:/:{gnome-terminal-profile}/",
      "key": "use-transparent-background",
      "value": "true"
    },
    {
      "schema": "org.gnome.Terminal.Legacy.Profile:/org/gnome/terminal/legacy/profiles:/:{gnome-terminal-profile}/",
      "key": "background-transparency-percent",
      "option": "pct"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "advanced": true,
  "options": [
    {
      "key": "pct",
      "label": "Transparência",
      "values": [
        {
          "value": "10",
          "label": "10%"
        },
        {
          "value": "15",
          "label": "15%"
        },
        {
          "value": "20",
          "label": "20%"
        },
        {
          "value": "30",
          "label": "30%"
        }
      ],
      "default": [
        "15"
      ]
    }
  ]
}
```

</details>

<a id="item-warp-style"></a>

### Warp: transparência e temas

Opacidade da janela e temas do Warp.

O Warp guarda o aspeto nas suas próprias definições; não há chave de sistema fiável para automatizar.

| Campo | Valor |
| --- | --- |
| ID | warp-style |
| Área / categoria | Estilo / Terminal |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Warp → Settings → Appearance → Window opacity (e desfoque).
2. Temas personalizados: ficheiros .yaml em ~/.local/share/warp-terminal/themes (ver documentação).
3. Escolha o tema em Settings → Appearance → Themes.

[Documentação do fornecedor](https://docs.warp.dev/terminal/appearance/custom-themes)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "warp-style",
  "name": "Warp: transparência e temas",
  "area": "style",
  "category": "Terminal",
  "summary": "Opacidade da janela e temas do Warp.",
  "why": "O Warp guarda o aspeto nas suas próprias definições; não há chave de sistema fiável para automatizar.",
  "method": "manual",
  "manual": {
    "steps": [
      "Warp → Settings → Appearance → Window opacity (e desfoque).",
      "Temas personalizados: ficheiros .yaml em ~/.local/share/warp-terminal/themes (ver documentação).",
      "Escolha o tema em Settings → Appearance → Themes."
    ],
    "url": "https://docs.warp.dev/terminal/appearance/custom-themes"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-dock-transparency"></a>

### Dock: transparente

Fundo da Ubuntu Dock com opacidade fixa.

Dock mais leve visualmente sobre o wallpaper.

| Campo | Valor |
| --- | --- |
| ID | dock-transparency |
| Área / categoria | Estilo / Painel e dock |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Opacidade** (`op`): um valor. Padrão: 0.6. 

| Valor | Descrição |
| --- | --- |
| 0.8 | 80% |
| 0.6 | 60% |
| 0.4 | 40% |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "dock-transparency",
  "name": "Dock: transparente",
  "area": "style",
  "category": "Painel e dock",
  "summary": "Fundo da Ubuntu Dock com opacidade fixa.",
  "why": "Dock mais leve visualmente sobre o wallpaper.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "transparency-mode",
      "value": "'FIXED'"
    },
    {
      "schema": "org.gnome.shell.extensions.dash-to-dock",
      "key": "background-opacity",
      "option": "op"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "op",
      "label": "Opacidade",
      "values": [
        {
          "value": "0.8",
          "label": "80%"
        },
        {
          "value": "0.6",
          "label": "60%"
        },
        {
          "value": "0.4",
          "label": "40%"
        }
      ],
      "default": [
        "0.6"
      ]
    }
  ]
}
```

</details>

<a id="item-cursor-size"></a>

### Tamanho do cursor

Tamanho do ponteiro do rato.

Ecrãs grandes ou de alta densidade.

| Campo | Valor |
| --- | --- |
| ID | cursor-size |
| Área / categoria | Estilo / Cursor e fontes |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Tamanho** (`size`): um valor. Padrão: 32. 

| Valor | Descrição |
| --- | --- |
| 24 | Normal (24) |
| 32 | Médio (32) |
| 48 | Grande (48) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "cursor-size",
  "name": "Tamanho do cursor",
  "area": "style",
  "category": "Cursor e fontes",
  "summary": "Tamanho do ponteiro do rato.",
  "why": "Ecrãs grandes ou de alta densidade.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "cursor-size",
      "option": "size"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "size",
      "label": "Tamanho",
      "values": [
        {
          "value": "24",
          "label": "Normal (24)"
        },
        {
          "value": "32",
          "label": "Médio (32)"
        },
        {
          "value": "48",
          "label": "Grande (48)"
        }
      ],
      "default": [
        "32"
      ]
    }
  ]
}
```

</details>

<a id="item-ui-font"></a>

### Fonte da interface

Fonte usada em menus e janelas.

Ligeiramente maior para leitura confortável.

| Campo | Valor |
| --- | --- |
| ID | ui-font |
| Área / categoria | Estilo / Cursor e fontes |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Fonte** (`font`): um valor. Padrão: 'Ubuntu Sans 11'. 

| Valor | Descrição |
| --- | --- |
| 'Ubuntu Sans 11' | Ubuntu Sans 11 (padrão) |
| 'Ubuntu Sans 12' | Ubuntu Sans 12 |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ui-font",
  "name": "Fonte da interface",
  "area": "style",
  "category": "Cursor e fontes",
  "summary": "Fonte usada em menus e janelas.",
  "why": "Ligeiramente maior para leitura confortável.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.interface",
      "key": "font-name",
      "option": "font"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "font",
      "label": "Fonte",
      "values": [
        {
          "value": "'Ubuntu Sans 11'",
          "label": "Ubuntu Sans 11 (padrão)"
        },
        {
          "value": "'Ubuntu Sans 12'",
          "label": "Ubuntu Sans 12"
        }
      ],
      "default": [
        "'Ubuntu Sans 11'"
      ]
    }
  ]
}
```

</details>

<a id="item-libreoffice"></a>

### LibreOffice

Processador de texto, folha de cálculo e apresentações.

Abrir e editar documentos Office/ODF; interface em português.

| Campo | Valor |
| --- | --- |
| ID | libreoffice |
| Área / categoria | Aplicações / Escritório |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | libreoffice-writer, libreoffice-calc, libreoffice-impress, libreoffice-gnome, libreoffice-l10n-pt, libreoffice-l10n-pt-br |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "libreoffice",
  "name": "LibreOffice",
  "area": "apps",
  "category": "Escritório",
  "summary": "Processador de texto, folha de cálculo e apresentações.",
  "why": "Abrir e editar documentos Office/ODF; interface em português.",
  "method": "apt",
  "packages": [
    "libreoffice-writer",
    "libreoffice-calc",
    "libreoffice-impress",
    "libreoffice-gnome",
    "libreoffice-l10n-pt",
    "libreoffice-l10n-pt-br"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-thunderbird"></a>

### Thunderbird (Flatpak)

Cliente de email e calendário da Mozilla.

No 26.04 o pacote APT é só uma transição para o Snap; o Flatpak vem da Mozilla.

| Campo | Valor |
| --- | --- |
| ID | thunderbird |
| Área / categoria | Aplicações / Escritório |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | org.mozilla.Thunderbird |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "thunderbird",
  "name": "Thunderbird (Flatpak)",
  "area": "apps",
  "category": "Escritório",
  "summary": "Cliente de email e calendário da Mozilla.",
  "why": "No 26.04 o pacote APT é só uma transição para o Snap; o Flatpak vem da Mozilla.",
  "method": "flatpak",
  "flatpak": "org.mozilla.Thunderbird",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "tags": [
    "snap:thunderbird"
  ],
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-ente-photos"></a>

### Ente Photos

Fotos com cifragem ponta a ponta (conta Ente).

Alternativa privada ao Google Fotos. Requer conta.

| Campo | Valor |
| --- | --- |
| ID | ente-photos |
| Área / categoria | Aplicações / Multimédia |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | io.ente.photos |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ente-photos",
  "name": "Ente Photos",
  "area": "apps",
  "category": "Multimédia",
  "summary": "Fotos com cifragem ponta a ponta (conta Ente).",
  "why": "Alternativa privada ao Google Fotos. Requer conta.",
  "method": "flatpak",
  "flatpak": "io.ente.photos",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-bitwarden"></a>

### Bitwarden

Gestor de senhas open source, sincronizado entre dispositivos.

Senhas fortes e únicas em todo o lado; extensão para navegadores. Requer conta (plano gratuito).

| Campo | Valor |
| --- | --- |
| ID | bitwarden |
| Área / categoria | Segurança / Gestores de senhas |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | com.bitwarden.desktop |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "bitwarden",
  "name": "Bitwarden",
  "area": "security",
  "category": "Gestores de senhas",
  "summary": "Gestor de senhas open source, sincronizado entre dispositivos.",
  "why": "Senhas fortes e únicas em todo o lado; extensão para navegadores. Requer conta (plano gratuito).",
  "method": "flatpak",
  "flatpak": "com.bitwarden.desktop",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-keepassxc"></a>

### KeePassXC

Gestor de senhas offline (ficheiro .kdbx).

Cofre local sem conta nem cloud; sincronize o ficheiro como quiser. Também guarda TOTP.

| Campo | Valor |
| --- | --- |
| ID | keepassxc |
| Área / categoria | Segurança / Gestores de senhas |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | keepassxc |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "keepassxc",
  "name": "KeePassXC",
  "area": "security",
  "category": "Gestores de senhas",
  "summary": "Gestor de senhas offline (ficheiro .kdbx).",
  "why": "Cofre local sem conta nem cloud; sincronize o ficheiro como quiser. Também guarda TOTP.",
  "method": "apt",
  "packages": [
    "keepassxc"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gnome-secrets"></a>

### Secrets (GNOME)

Gestor de senhas KeePass simples, ao estilo GNOME.

Mais simples que o KeePassXC, mesmo formato .kdbx.

| Campo | Valor |
| --- | --- |
| ID | gnome-secrets |
| Área / categoria | Segurança / Gestores de senhas |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | org.gnome.World.Secrets |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gnome-secrets",
  "name": "Secrets (GNOME)",
  "area": "security",
  "category": "Gestores de senhas",
  "summary": "Gestor de senhas KeePass simples, ao estilo GNOME.",
  "why": "Mais simples que o KeePassXC, mesmo formato .kdbx.",
  "method": "flatpak",
  "flatpak": "org.gnome.World.Secrets",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-seahorse"></a>

### Senhas e chaves

Gere o chaveiro GNOME, chaves SSH e GPG.

Ver/exportar o chaveiro na migração e gerir chaves GPG graficamente.

| Campo | Valor |
| --- | --- |
| ID | seahorse |
| Área / categoria | Segurança / Gestores de senhas |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | seahorse |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "seahorse",
  "name": "Senhas e chaves",
  "area": "security",
  "category": "Gestores de senhas",
  "summary": "Gere o chaveiro GNOME, chaves SSH e GPG.",
  "why": "Ver/exportar o chaveiro na migração e gerir chaves GPG graficamente.",
  "method": "apt",
  "packages": [
    "seahorse"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ente-auth"></a>

### Ente Auth

Códigos 2FA (TOTP) com backup cifrado e sincronização.

Não perder os códigos 2FA ao trocar de telemóvel/PC. Funciona offline; conta opcional para sincronizar.

| Campo | Valor |
| --- | --- |
| ID | ente-auth |
| Área / categoria | Segurança / Autenticação 2FA |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | io.ente.auth |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ente-auth",
  "name": "Ente Auth",
  "area": "security",
  "category": "Autenticação 2FA",
  "summary": "Códigos 2FA (TOTP) com backup cifrado e sincronização.",
  "why": "Não perder os códigos 2FA ao trocar de telemóvel/PC. Funciona offline; conta opcional para sincronizar.",
  "method": "flatpak",
  "flatpak": "io.ente.auth",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-gnome-authenticator"></a>

### Authenticator (GNOME)

Códigos 2FA locais no desktop.

Alternativa simples e só local ao Ente Auth.

| Campo | Valor |
| --- | --- |
| ID | gnome-authenticator |
| Área / categoria | Segurança / Autenticação 2FA |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | com.belmoussaoui.Authenticator |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gnome-authenticator",
  "name": "Authenticator (GNOME)",
  "area": "security",
  "category": "Autenticação 2FA",
  "summary": "Códigos 2FA locais no desktop.",
  "why": "Alternativa simples e só local ao Ente Auth.",
  "method": "flatpak",
  "flatpak": "com.belmoussaoui.Authenticator",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-yubico-authenticator"></a>

### Yubico Authenticator

Códigos 2FA guardados numa YubiKey.

Só se tiver YubiKey. Editor não verificado no Flathub.

| Campo | Valor |
| --- | --- |
| ID | yubico-authenticator |
| Área / categoria | Segurança / Autenticação 2FA |
| Método | flatpak |
| Fonte | Flathub (editor não verificado) |
| Flatpak | com.yubico.yubioath |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "yubico-authenticator",
  "name": "Yubico Authenticator",
  "area": "security",
  "category": "Autenticação 2FA",
  "summary": "Códigos 2FA guardados numa YubiKey.",
  "why": "Só se tiver YubiKey. Editor não verificado no Flathub.",
  "method": "flatpak",
  "flatpak": "com.yubico.yubioath",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor não verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-yubikey-manager"></a>

### YubiKey Manager (ykman)

Configurar YubiKeys (OTP, FIDO2, PIV, OpenPGP) na linha de comandos.

Gerir a chave física e usá-la para GPG/SSH.

| Campo | Valor |
| --- | --- |
| ID | yubikey-manager |
| Área / categoria | Segurança / Autenticação 2FA |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | yubikey-manager, scdaemon |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "yubikey-manager",
  "name": "YubiKey Manager (ykman)",
  "area": "security",
  "category": "Autenticação 2FA",
  "summary": "Configurar YubiKeys (OTP, FIDO2, PIV, OpenPGP) na linha de comandos.",
  "why": "Gerir a chave física e usá-la para GPG/SSH.",
  "method": "apt",
  "packages": [
    "yubikey-manager",
    "scdaemon"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-cryptomator"></a>

### Cryptomator

Cofres cifrados dentro de pastas sincronizadas (Drive, Dropbox…).

Proteger documentos sensíveis guardados na cloud.

| Campo | Valor |
| --- | --- |
| ID | cryptomator |
| Área / categoria | Segurança / Cifragem |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | org.cryptomator.Cryptomator |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "cryptomator",
  "name": "Cryptomator",
  "area": "security",
  "category": "Cifragem",
  "summary": "Cofres cifrados dentro de pastas sincronizadas (Drive, Dropbox…).",
  "why": "Proteger documentos sensíveis guardados na cloud.",
  "method": "flatpak",
  "flatpak": "org.cryptomator.Cryptomator",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-metadata-cleaner"></a>

### Metadata Cleaner

Remove metadados (GPS, autor, câmara) de fotos e documentos.

Antes de partilhar ficheiros com clientes ou publicar fotos.

| Campo | Valor |
| --- | --- |
| ID | metadata-cleaner |
| Área / categoria | Segurança / Privacidade |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | fr.romainvigier.MetadataCleaner |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "metadata-cleaner",
  "name": "Metadata Cleaner",
  "area": "security",
  "category": "Privacidade",
  "summary": "Remove metadados (GPS, autor, câmara) de fotos e documentos.",
  "why": "Antes de partilhar ficheiros com clientes ou publicar fotos.",
  "method": "flatpak",
  "flatpak": "fr.romainvigier.MetadataCleaner",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-flatseal"></a>

### Flatseal

Rever e restringir permissões das aplicações Flatpak.

Saber a que pastas/dispositivos cada Flatpak acede e reduzir.

| Campo | Valor |
| --- | --- |
| ID | flatseal |
| Área / categoria | Segurança / Privacidade |
| Método | flatpak |
| Fonte | Flathub (editor verificado) |
| Flatpak | com.github.tchx84.Flatseal |
| Evidência da receita | Flathub API 200 em 2026-10-04 |
| Dependências | [Flatpak + Flathub](#item-flathub) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "flatseal",
  "name": "Flatseal",
  "area": "security",
  "category": "Privacidade",
  "summary": "Rever e restringir permissões das aplicações Flatpak.",
  "why": "Saber a que pastas/dispositivos cada Flatpak acede e reduzir.",
  "method": "flatpak",
  "flatpak": "com.github.tchx84.Flatseal",
  "requires": [
    "flathub"
  ],
  "source": "Flathub (editor verificado)",
  "verified": "Flathub API 200 em 2026-10-04",
  "compat": {
    "arch": [
      "amd64"
    ]
  }
}
```

</details>

<a id="item-torbrowser-launcher"></a>

### Tor Browser (launcher)

Transfere e mantém o Tor Browser, verificando a assinatura do Tor Project.

Navegação anónima quando necessário.

| Campo | Valor |
| --- | --- |
| ID | torbrowser-launcher |
| Área / categoria | Segurança / Privacidade |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | torbrowser-launcher |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "torbrowser-launcher",
  "name": "Tor Browser (launcher)",
  "area": "security",
  "category": "Privacidade",
  "summary": "Transfere e mantém o Tor Browser, verificando a assinatura do Tor Project.",
  "why": "Navegação anónima quando necessário.",
  "method": "apt",
  "packages": [
    "torbrowser-launcher"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gufw"></a>

### Gufw (interface da firewall)

Interface gráfica para o ufw.

Só instala — a firewall NÃO é ativada nem configurada. Atenção: o Docker publica portas contornando as regras do ufw.

| Campo | Valor |
| --- | --- |
| ID | gufw |
| Área / categoria | Segurança / Firewall e auditoria |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gufw |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

Ative e configure manualmente, considerando portas de serviços, Docker e VPN.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gufw",
  "name": "Gufw (interface da firewall)",
  "area": "security",
  "category": "Firewall e auditoria",
  "summary": "Interface gráfica para o ufw.",
  "why": "Só instala — a firewall NÃO é ativada nem configurada. Atenção: o Docker publica portas contornando as regras do ufw.",
  "method": "apt",
  "packages": [
    "gufw"
  ],
  "notes": "Ative e configure manualmente, considerando portas de serviços, Docker e VPN.",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-lynis"></a>

### Lynis

Auditoria de segurança do sistema (só leitura).

Relatório de sugestões de hardening: sudo lynis audit system. Não altera nada.

| Campo | Valor |
| --- | --- |
| ID | lynis |
| Área / categoria | Segurança / Firewall e auditoria |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | lynis |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "lynis",
  "name": "Lynis",
  "area": "security",
  "category": "Firewall e auditoria",
  "summary": "Auditoria de segurança do sistema (só leitura).",
  "why": "Relatório de sugestões de hardening: sudo lynis audit system. Não altera nada.",
  "method": "apt",
  "packages": [
    "lynis"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-clamtk"></a>

### ClamAV + ClamTk

Antivírus por análise a pedido.

Útil para verificar ficheiros/anexos que vão para Windows. Não protege o Linux em tempo real; o atualizador freshclam fica como serviço.

| Campo | Valor |
| --- | --- |
| ID | clamtk |
| Área / categoria | Segurança / Firewall e auditoria |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | clamtk |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "clamtk",
  "name": "ClamAV + ClamTk",
  "area": "security",
  "category": "Firewall e auditoria",
  "summary": "Antivírus por análise a pedido.",
  "why": "Útil para verificar ficheiros/anexos que vão para Windows. Não protege o Linux em tempo real; o atualizador freshclam fica como serviço.",
  "method": "apt",
  "packages": [
    "clamtk"
  ],
  "advanced": true,
  "service": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-opensnitch"></a>

### OpenSnitch

Firewall de aplicações interativa (pergunta que app pode ligar-se à rede).

Avançado: gera muitos pedidos no início e pode bloquear atualizações se responder mal.

| Campo | Valor |
| --- | --- |
| ID | opensnitch |
| Área / categoria | Segurança / Firewall e auditoria |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | opensnitch |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "opensnitch",
  "name": "OpenSnitch",
  "area": "security",
  "category": "Firewall e auditoria",
  "summary": "Firewall de aplicações interativa (pergunta que app pode ligar-se à rede).",
  "why": "Avançado: gera muitos pedidos no início e pode bloquear atualizações se responder mal.",
  "method": "apt",
  "packages": [
    "opensnitch"
  ],
  "advanced": true,
  "service": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gobuster"></a>

### Gobuster

Descoberta de diretórios, ficheiros e subdomínios na linha de comandos.

Verificar caminhos expostos e a superfície de aplicações web num laboratório ou projeto próprio.

| Campo | Valor |
| --- | --- |
| ID | gobuster |
| Área / categoria | Segurança / Auditoria web |
| Método | apt |
| Fonte | Arquivo Ubuntu (universe) |
| Pacotes | gobuster |
| Evidência da receita | Instalação e gobuster --version verificados em container Ubuntu 26.04: 3.8.2, 2026-10-04 |
| Atenção | Enumeração ativa de aplicações e DNS |

Requer uma wordlist fornecida por si. Instalar não executa testes nem transfere listas automaticamente.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gobuster",
  "name": "Gobuster",
  "area": "security",
  "category": "Auditoria web",
  "summary": "Descoberta de diretórios, ficheiros e subdomínios na linha de comandos.",
  "why": "Verificar caminhos expostos e a superfície de aplicações web num laboratório ou projeto próprio.",
  "method": "apt",
  "packages": [
    "gobuster"
  ],
  "verify": {
    "cmd": [
      "gobuster",
      "--version"
    ]
  },
  "sensitive": "Enumeração ativa de aplicações e DNS",
  "notes": "Requer uma wordlist fornecida por si. Instalar não executa testes nem transfere listas automaticamente.",
  "source": "Arquivo Ubuntu (universe)",
  "verified": "Instalação e gobuster --version verificados em container Ubuntu 26.04: 3.8.2, 2026-10-04"
}
```

</details>

<a id="item-sqlmap"></a>

### sqlmap

Ferramenta de testes de injeção SQL em aplicações web.

Investigar e reproduzir falhas de SQL injection em aplicações de laboratório ou sob auditoria.

| Campo | Valor |
| --- | --- |
| ID | sqlmap |
| Área / categoria | Segurança / Auditoria web |
| Método | apt |
| Fonte | Arquivo Ubuntu (universe) |
| Pacotes | sqlmap |
| Evidência da receita | Instalação e sqlmap --version verificados em container Ubuntu 26.04: 1.10.4, 2026-10-04 |
| Atenção | Testes ativos de injeção SQL podem alterar dados do alvo |

A aplicação apenas instala a ferramenta. Os alvos e os testes são escolhidos por si depois.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "sqlmap",
  "name": "sqlmap",
  "area": "security",
  "category": "Auditoria web",
  "summary": "Ferramenta de testes de injeção SQL em aplicações web.",
  "why": "Investigar e reproduzir falhas de SQL injection em aplicações de laboratório ou sob auditoria.",
  "method": "apt",
  "packages": [
    "sqlmap"
  ],
  "verify": {
    "cmd": [
      "sqlmap",
      "--version"
    ]
  },
  "sensitive": "Testes ativos de injeção SQL podem alterar dados do alvo",
  "notes": "A aplicação apenas instala a ferramenta. Os alvos e os testes são escolhidos por si depois.",
  "source": "Arquivo Ubuntu (universe)",
  "verified": "Instalação e sqlmap --version verificados em container Ubuntu 26.04: 1.10.4, 2026-10-04"
}
```

</details>

<a id="item-sslscan"></a>

### SSLScan

Inspeção rápida dos protocolos e cifras TLS aceites por um serviço.

Rever a configuração HTTPS e identificar protocolos ou cifras que já não pretende aceitar.

| Campo | Valor |
| --- | --- |
| ID | sslscan |
| Área / categoria | Segurança / TLS e rede |
| Método | apt |
| Fonte | Arquivo Ubuntu (universe) |
| Pacotes | sslscan |
| Evidência da receita | Instalação e sslscan --version verificados em container Ubuntu 26.04: 2.1.5, 2026-10-04 |

Consulta serviços quando o executar manualmente; não altera certificados nem configurações.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "sslscan",
  "name": "SSLScan",
  "area": "security",
  "category": "TLS e rede",
  "summary": "Inspeção rápida dos protocolos e cifras TLS aceites por um serviço.",
  "why": "Rever a configuração HTTPS e identificar protocolos ou cifras que já não pretende aceitar.",
  "method": "apt",
  "packages": [
    "sslscan"
  ],
  "verify": {
    "cmd": [
      "sslscan",
      "--version"
    ]
  },
  "notes": "Consulta serviços quando o executar manualmente; não altera certificados nem configurações.",
  "source": "Arquivo Ubuntu (universe)",
  "verified": "Instalação e sslscan --version verificados em container Ubuntu 26.04: 2.1.5, 2026-10-04"
}
```

</details>

<a id="item-testssl"></a>

### testssl.sh

Auditoria detalhada de TLS, certificados e falhas criptográficas conhecidas.

Complementa o SSLScan com uma análise mais extensa de HTTPS e outros serviços TLS/STARTTLS.

| Campo | Valor |
| --- | --- |
| ID | testssl |
| Área / categoria | Segurança / TLS e rede |
| Método | apt |
| Fonte | Arquivo Ubuntu (universe) |
| Pacotes | testssl.sh |
| Evidência da receita | Instalação e testssl --version verificados em container Ubuntu 26.04: 3.2.2, 2026-10-04 |

O comando no Ubuntu é testssl. A instalação não executa auditorias.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "testssl",
  "name": "testssl.sh",
  "area": "security",
  "category": "TLS e rede",
  "summary": "Auditoria detalhada de TLS, certificados e falhas criptográficas conhecidas.",
  "why": "Complementa o SSLScan com uma análise mais extensa de HTTPS e outros serviços TLS/STARTTLS.",
  "method": "apt",
  "packages": [
    "testssl.sh"
  ],
  "verify": {
    "cmd": [
      "testssl",
      "--version"
    ]
  },
  "notes": "O comando no Ubuntu é testssl. A instalação não executa auditorias.",
  "source": "Arquivo Ubuntu (universe)",
  "verified": "Instalação e testssl --version verificados em container Ubuntu 26.04: 3.2.2, 2026-10-04"
}
```

</details>

<a id="item-arp-scan"></a>

### arp-scan

Descoberta de dispositivos IPv4 na rede local através de ARP.

Inventariar dispositivos na mesma rede local, incluindo os que não respondem a ping.

| Campo | Valor |
| --- | --- |
| ID | arp-scan |
| Área / categoria | Segurança / TLS e rede |
| Método | apt |
| Fonte | Arquivo Ubuntu (universe) |
| Pacotes | arp-scan |
| Evidência da receita | Instalação e arp-scan --version verificados em container Ubuntu 26.04: 1.10.0, 2026-10-04 |
| Atenção | Envia pedidos ARP na rede local |

O varrimento pode exigir privilégios. A aplicação não altera permissões nem executa capturas ou varrimentos.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "arp-scan",
  "name": "arp-scan",
  "area": "security",
  "category": "TLS e rede",
  "summary": "Descoberta de dispositivos IPv4 na rede local através de ARP.",
  "why": "Inventariar dispositivos na mesma rede local, incluindo os que não respondem a ping.",
  "method": "apt",
  "packages": [
    "arp-scan"
  ],
  "verify": {
    "cmd": [
      "arp-scan",
      "--version"
    ]
  },
  "sensitive": "Envia pedidos ARP na rede local",
  "notes": "O varrimento pode exigir privilégios. A aplicação não altera permissões nem executa capturas ou varrimentos.",
  "source": "Arquivo Ubuntu (universe)",
  "verified": "Instalação e arp-scan --version verificados em container Ubuntu 26.04: 1.10.0, 2026-10-04"
}
```

</details>

<a id="item-remmina"></a>

### Remmina (RDP/VNC)

Cliente de ambiente de trabalho remoto RDP, VNC e SSH.

Aceder a servidores Windows e máquinas de clientes.

| Campo | Valor |
| --- | --- |
| ID | remmina |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | remmina, remmina-plugin-rdp, remmina-plugin-vnc |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "remmina",
  "name": "Remmina (RDP/VNC)",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "Cliente de ambiente de trabalho remoto RDP, VNC e SSH.",
  "why": "Aceder a servidores Windows e máquinas de clientes.",
  "method": "apt",
  "packages": [
    "remmina",
    "remmina-plugin-rdp",
    "remmina-plugin-vnc"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-putty-tools"></a>

### PuTTY tools

puttygen/pscp/plink: converter chaves .ppk e copiar ficheiros.

Clientes que só entregam chaves PuTTY.

| Campo | Valor |
| --- | --- |
| ID | putty-tools |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | putty-tools |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "putty-tools",
  "name": "PuTTY tools",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "puttygen/pscp/plink: converter chaves .ppk e copiar ficheiros.",
  "why": "Clientes que só entregam chaves PuTTY.",
  "method": "apt",
  "packages": [
    "putty-tools"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-sshfs"></a>

### SSHFS

Montar pastas remotas por SSH.

Editar ficheiros de um servidor como se fossem locais.

| Campo | Valor |
| --- | --- |
| ID | sshfs |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | sshfs |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "sshfs",
  "name": "SSHFS",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "Montar pastas remotas por SSH.",
  "why": "Editar ficheiros de um servidor como se fossem locais.",
  "method": "apt",
  "packages": [
    "sshfs"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-tmux"></a>

### tmux

Multiplexador de terminal com sessões persistentes.

Sessões que sobrevivem a quedas de ligação.

| Campo | Valor |
| --- | --- |
| ID | tmux |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | tmux |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tmux",
  "name": "tmux",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "Multiplexador de terminal com sessões persistentes.",
  "why": "Sessões que sobrevivem a quedas de ligação.",
  "method": "apt",
  "packages": [
    "tmux"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-net-diag"></a>

### Diagnóstico de rede

mtr, traceroute, dig/nslookup, whois, ipcalc, ethtool.

Kit básico para diagnosticar latência, DNS, rotas e interfaces.

| Campo | Valor |
| --- | --- |
| ID | net-diag |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | mtr-tiny, traceroute, bind9-dnsutils, whois, ipcalc, ethtool |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "net-diag",
  "name": "Diagnóstico de rede",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "mtr, traceroute, dig/nslookup, whois, ipcalc, ethtool.",
  "why": "Kit básico para diagnosticar latência, DNS, rotas e interfaces.",
  "method": "apt",
  "packages": [
    "mtr-tiny",
    "traceroute",
    "bind9-dnsutils",
    "whois",
    "ipcalc",
    "ethtool"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-nmap"></a>

### Nmap + Ncat

Descoberta de hosts e portas.

Inventário de rede. Use só em redes onde tem autorização.

| Campo | Valor |
| --- | --- |
| ID | nmap |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | nmap, ncat |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Atenção | Varrimento de rede: use apenas com autorização |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nmap",
  "name": "Nmap + Ncat",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "Descoberta de hosts e portas.",
  "why": "Inventário de rede. Use só em redes onde tem autorização.",
  "method": "apt",
  "packages": [
    "nmap",
    "ncat"
  ],
  "sensitive": "Varrimento de rede: use apenas com autorização",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-iperf3"></a>

### iperf3

Medir largura de banda entre duas máquinas.

Validar links e Wi-Fi. Não fica como serviço (opção padrão do pacote).

| Campo | Valor |
| --- | --- |
| ID | iperf3 |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | iperf3 |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "iperf3",
  "name": "iperf3",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "Medir largura de banda entre duas máquinas.",
  "why": "Validar links e Wi-Fi. Não fica como serviço (opção padrão do pacote).",
  "method": "apt",
  "packages": [
    "iperf3"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-traffic-monitor"></a>

### Monitores de tráfego

Tráfego por ligação, por processo e por interface.

Descobrir quem está a consumir a rede.

| Campo | Valor |
| --- | --- |
| ID | traffic-monitor |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | iftop, nethogs, nload |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "traffic-monitor",
  "name": "Monitores de tráfego",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "Tráfego por ligação, por processo e por interface.",
  "why": "Descobrir quem está a consumir a rede.",
  "method": "apt",
  "packages": [
    "iftop",
    "nethogs",
    "nload"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-tcpdump"></a>

### tcpdump

Captura de pacotes na linha de comandos.

Diagnóstico de protocolos; requer sudo para capturar.

| Campo | Valor |
| --- | --- |
| ID | tcpdump |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | tcpdump |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tcpdump",
  "name": "tcpdump",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "Captura de pacotes na linha de comandos.",
  "why": "Diagnóstico de protocolos; requer sudo para capturar.",
  "method": "apt",
  "packages": [
    "tcpdump"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-wireshark"></a>

### Wireshark

Análise gráfica de pacotes.

Instalado sem permitir captura a utilizadores comuns (padrão seguro). Para capturar sem sudo: sudo dpkg-reconfigure wireshark-common e adicionar-se ao grupo wireshark (manual).

| Campo | Valor |
| --- | --- |
| ID | wireshark |
| Área / categoria | Infraestrutura / Rede e diagnóstico |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | wireshark |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Atenção | Ferramenta de captura de tráfego |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "wireshark",
  "name": "Wireshark",
  "area": "infra",
  "category": "Rede e diagnóstico",
  "summary": "Análise gráfica de pacotes.",
  "why": "Instalado sem permitir captura a utilizadores comuns (padrão seguro). Para capturar sem sudo: sudo dpkg-reconfigure wireshark-common e adicionar-se ao grupo wireshark (manual).",
  "method": "apt",
  "packages": [
    "wireshark"
  ],
  "sensitive": "Ferramenta de captura de tráfego",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-virt-manager"></a>

### Virtualização (virt-manager/KVM)

Máquinas virtuais KVM com interface gráfica.

Laboratórios e VMs de clientes. Instala o serviço libvirtd; adicionar-se ao grupo libvirt é manual.

| Campo | Valor |
| --- | --- |
| ID | virt-manager |
| Área / categoria | Infraestrutura / Virtualização |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | virt-manager, qemu-system-x86, libvirt-daemon-system |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Serviço local | Sim |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "virt-manager",
  "name": "Virtualização (virt-manager/KVM)",
  "area": "infra",
  "category": "Virtualização",
  "summary": "Máquinas virtuais KVM com interface gráfica.",
  "why": "Laboratórios e VMs de clientes. Instala o serviço libvirtd; adicionar-se ao grupo libvirt é manual.",
  "method": "apt",
  "packages": [
    "virt-manager",
    "qemu-system-x86",
    "libvirt-daemon-system"
  ],
  "service": true,
  "relogin": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-gnome-boxes"></a>

### GNOME Boxes

VMs simples por utilizador.

Testar distribuições rapidamente sem configurar libvirt.

| Campo | Valor |
| --- | --- |
| ID | gnome-boxes |
| Área / categoria | Infraestrutura / Virtualização |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gnome-boxes |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gnome-boxes",
  "name": "GNOME Boxes",
  "area": "infra",
  "category": "Virtualização",
  "summary": "VMs simples por utilizador.",
  "why": "Testar distribuições rapidamente sem configurar libvirt.",
  "method": "apt",
  "packages": [
    "gnome-boxes"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ansible"></a>

### Ansible

Automação de configuração por SSH.

Preparar vários servidores/postos de forma repetível.

| Campo | Valor |
| --- | --- |
| ID | ansible |
| Área / categoria | Infraestrutura / Automação |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | ansible |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ansible",
  "name": "Ansible",
  "area": "infra",
  "category": "Automação",
  "summary": "Automação de configuração por SSH.",
  "why": "Preparar vários servidores/postos de forma repetível.",
  "method": "apt",
  "packages": [
    "ansible"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-wireguard-tools"></a>

### WireGuard (ferramentas)

wg/wg-quick para túneis WireGuard.

Só instala as ferramentas; nenhuma VPN é configurada.

| Campo | Valor |
| --- | --- |
| ID | wireguard-tools |
| Área / categoria | Infraestrutura / Automação |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | wireguard-tools |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "wireguard-tools",
  "name": "WireGuard (ferramentas)",
  "area": "infra",
  "category": "Automação",
  "summary": "wg/wg-quick para túneis WireGuard.",
  "why": "Só instala as ferramentas; nenhuma VPN é configurada.",
  "method": "apt",
  "packages": [
    "wireguard-tools"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-disk-tools"></a>

### Discos e recuperação

SMART, partições, recuperação de partições/ficheiros e cópia de discos com erros.

Diagnóstico e salvamento de discos de clientes. smartd fica como serviço.

| Campo | Valor |
| --- | --- |
| ID | disk-tools |
| Área / categoria | Infraestrutura / Discos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | smartmontools, gparted, testdisk, gddrescue |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Serviço local | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "disk-tools",
  "name": "Discos e recuperação",
  "area": "infra",
  "category": "Discos",
  "summary": "SMART, partições, recuperação de partições/ficheiros e cópia de discos com erros.",
  "why": "Diagnóstico e salvamento de discos de clientes. smartd fica como serviço.",
  "method": "apt",
  "packages": [
    "smartmontools",
    "gparted",
    "testdisk",
    "gddrescue"
  ],
  "service": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-lnav"></a>

### lnav

Visualizador de logs com pesquisa e realce.

Ler logs de servidores rapidamente.

| Campo | Valor |
| --- | --- |
| ID | lnav |
| Área / categoria | Infraestrutura / Discos |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | lnav |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "lnav",
  "name": "lnav",
  "area": "infra",
  "category": "Discos",
  "summary": "Visualizador de logs com pesquisa e realce.",
  "why": "Ler logs de servidores rapidamente.",
  "method": "apt",
  "packages": [
    "lnav"
  ],
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-minicom"></a>

### minicom

Consola série.

Equipamentos de rede com porta de consola.

| Campo | Valor |
| --- | --- |
| ID | minicom |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | minicom |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Avançado | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "minicom",
  "name": "minicom",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "Consola série.",
  "why": "Equipamentos de rede com porta de consola.",
  "method": "apt",
  "packages": [
    "minicom"
  ],
  "advanced": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-virtualbox"></a>

### VirtualBox

Hipervisor da Oracle.

Manual: módulos do kernel e Secure Boot exigem assinatura (MOK); prefira KVM/virt-manager.

| Campo | Valor |
| --- | --- |
| ID | virtualbox |
| Área / categoria | Infraestrutura / Virtualização |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Recomendado: usar virt-manager (KVM).
2. Se precisar mesmo: siga https://www.virtualbox.org/wiki/Linux_Downloads e registe a chave MOK no arranque.

[Documentação do fornecedor](https://www.virtualbox.org/wiki/Linux_Downloads)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "virtualbox",
  "name": "VirtualBox",
  "area": "infra",
  "category": "Virtualização",
  "summary": "Hipervisor da Oracle.",
  "why": "Manual: módulos do kernel e Secure Boot exigem assinatura (MOK); prefira KVM/virt-manager.",
  "method": "manual",
  "manual": {
    "steps": [
      "Recomendado: usar virt-manager (KVM).",
      "Se precisar mesmo: siga https://www.virtualbox.org/wiki/Linux_Downloads e registe a chave MOK no arranque."
    ],
    "url": "https://www.virtualbox.org/wiki/Linux_Downloads"
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-terraform-kubectl"></a>

### Terraform / kubectl / Helm

Ferramentas de infraestrutura como código e Kubernetes.

Repositórios de terceiros ainda não verificados nesta versão; siga a documentação oficial.

| Campo | Valor |
| --- | --- |
| ID | terraform-kubectl |
| Área / categoria | Infraestrutura / Automação |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Terraform: https://developer.hashicorp.com/terraform/install (repositório APT HashiCorp, verifique a fingerprint).
2. kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
3. Helm: https://helm.sh/docs/intro/install/

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "terraform-kubectl",
  "name": "Terraform / kubectl / Helm",
  "area": "infra",
  "category": "Automação",
  "summary": "Ferramentas de infraestrutura como código e Kubernetes.",
  "why": "Repositórios de terceiros ainda não verificados nesta versão; siga a documentação oficial.",
  "method": "manual",
  "manual": {
    "steps": [
      "Terraform: https://developer.hashicorp.com/terraform/install (repositório APT HashiCorp, verifique a fingerprint).",
      "kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/",
      "Helm: https://helm.sh/docs/intro/install/"
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-remote-support"></a>

### AnyDesk / TeamViewer / RustDesk

Suporte remoto a clientes.

Software de terceiros não verificado nesta versão; RustDesk é open source.

| Campo | Valor |
| --- | --- |
| ID | remote-support |
| Área / categoria | Infraestrutura / Acesso remoto |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Instale a partir do site oficial do fornecedor e verifique a assinatura do pacote.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "remote-support",
  "name": "AnyDesk / TeamViewer / RustDesk",
  "area": "infra",
  "category": "Acesso remoto",
  "summary": "Suporte remoto a clientes.",
  "why": "Software de terceiros não verificado nesta versão; RustDesk é open source.",
  "method": "manual",
  "manual": {
    "steps": [
      "Instale a partir do site oficial do fornecedor e verifique a assinatura do pacote."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-ext-dash-to-dock"></a>

### Dash to Dock

Dock configurável.

O Ubuntu já inclui a Ubuntu Dock (fork do Dash to Dock) — não instale ambos.

| Campo | Valor |
| --- | --- |
| ID | ext-dash-to-dock |
| Área / categoria | Extensões / Já incluídas no Ubuntu |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | dash-to-dock@micxgx.gmail.com |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-dash-to-dock",
  "name": "Dash to Dock",
  "area": "extensions",
  "category": "Já incluídas no Ubuntu",
  "summary": "Dock configurável.",
  "why": "O Ubuntu já inclui a Ubuntu Dock (fork do Dash to Dock) — não instale ambos.",
  "method": "extension",
  "extension": "dash-to-dock@micxgx.gmail.com",
  "providedBy": [
    "ubuntu-dock@ubuntu.com"
  ],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04"
}
```

</details>

<a id="item-ext-appindicator"></a>

### AppIndicator

Ícones de bandeja (Slack, Docker Desktop, Solaar).

Já fornecida pelo Ubuntu (ubuntu-appindicators).

| Campo | Valor |
| --- | --- |
| ID | ext-appindicator |
| Área / categoria | Extensões / Já incluídas no Ubuntu |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | appindicatorsupport@rgcjonas.gmail.com |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-appindicator",
  "name": "AppIndicator",
  "area": "extensions",
  "category": "Já incluídas no Ubuntu",
  "summary": "Ícones de bandeja (Slack, Docker Desktop, Solaar).",
  "why": "Já fornecida pelo Ubuntu (ubuntu-appindicators).",
  "method": "extension",
  "extension": "appindicatorsupport@rgcjonas.gmail.com",
  "providedBy": [
    "ubuntu-appindicators@ubuntu.com"
  ],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04"
}
```

</details>

<a id="item-ext-tiling"></a>

### Tiling Assistant

Encaixar janelas em metades/quartos.

Já fornecida pelo Ubuntu (Ubuntu Tiling Assistant).

| Campo | Valor |
| --- | --- |
| ID | ext-tiling |
| Área / categoria | Extensões / Já incluídas no Ubuntu |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | tiling-assistant@leleat-on-github |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-tiling",
  "name": "Tiling Assistant",
  "area": "extensions",
  "category": "Já incluídas no Ubuntu",
  "summary": "Encaixar janelas em metades/quartos.",
  "why": "Já fornecida pelo Ubuntu (Ubuntu Tiling Assistant).",
  "method": "extension",
  "extension": "tiling-assistant@leleat-on-github",
  "providedBy": [
    "tiling-assistant@ubuntu.com"
  ],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04"
}
```

</details>

<a id="item-ext-clipboard"></a>

### Clipboard Indicator

Histórico da área de transferência.

Recuperar algo copiado há minutos. Atenção: guarda texto copiado (inclusive senhas) em memória/disco — ative 'privado' nas preferências.

| Campo | Valor |
| --- | --- |
| ID | ext-clipboard |
| Área / categoria | Extensões / Produtividade |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | clipboard-indicator@tudmotu.com |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-clipboard",
  "name": "Clipboard Indicator",
  "area": "extensions",
  "category": "Produtividade",
  "summary": "Histórico da área de transferência.",
  "why": "Recuperar algo copiado há minutos. Atenção: guarda texto copiado (inclusive senhas) em memória/disco — ative 'privado' nas preferências.",
  "method": "extension",
  "extension": "clipboard-indicator@tudmotu.com",
  "providedBy": [],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04"
}
```

</details>

<a id="item-ext-caffeine"></a>

### Caffeine

Impede suspensão/ecrã a apagar com um clique.

Builds longos, apresentações e downloads.

| Campo | Valor |
| --- | --- |
| ID | ext-caffeine |
| Área / categoria | Extensões / Produtividade |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | caffeine@patapon.info |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-caffeine",
  "name": "Caffeine",
  "area": "extensions",
  "category": "Produtividade",
  "summary": "Impede suspensão/ecrã a apagar com um clique.",
  "why": "Builds longos, apresentações e downloads.",
  "method": "extension",
  "extension": "caffeine@patapon.info",
  "providedBy": [],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04"
}
```

</details>

<a id="item-ext-just-perfection"></a>

### Just Perfection

Ajuste fino da interface do GNOME Shell.

Esconder/mover elementos do painel. Poderoso; opcional.

| Campo | Valor |
| --- | --- |
| ID | ext-just-perfection |
| Área / categoria | Extensões / Interface |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | just-perfection-desktop@just-perfection |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Avançado | Sim |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-just-perfection",
  "name": "Just Perfection",
  "area": "extensions",
  "category": "Interface",
  "summary": "Ajuste fino da interface do GNOME Shell.",
  "why": "Esconder/mover elementos do painel. Poderoso; opcional.",
  "method": "extension",
  "extension": "just-perfection-desktop@just-perfection",
  "providedBy": [],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04",
  "advanced": true
}
```

</details>

<a id="item-ext-vitals"></a>

### Vitals

Temperaturas, CPU, RAM e rede no painel.

Ver carga e temperatura durante builds/Docker.

| Campo | Valor |
| --- | --- |
| ID | ext-vitals |
| Área / categoria | Extensões / Monitorização |
| Método | extension |
| Fonte | extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME) |
| UUID | Vitals@CoreCoding.com |
| Evidência da receita | API extension-info com shell_version=50: HTTP 200 em 2026-10-04 |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-vitals",
  "name": "Vitals",
  "area": "extensions",
  "category": "Monitorização",
  "summary": "Temperaturas, CPU, RAM e rede no painel.",
  "why": "Ver carga e temperatura durante builds/Docker.",
  "method": "extension",
  "extension": "Vitals@CoreCoding.com",
  "providedBy": [],
  "source": "extensions.gnome.org (API pública; compatibilidade verificada contra a versão real do GNOME)",
  "relogin": true,
  "verified": "API extension-info com shell_version=50: HTTP 200 em 2026-10-04",
  "group": "monitor"
}
```

</details>

<a id="item-ext-system-monitor"></a>

### System Monitor (GNOME oficial)

Monitor de recursos oficial do projeto GNOME.

Leve e mantido com o próprio GNOME; alternativa ao Vitals (escolha um).

| Campo | Valor |
| --- | --- |
| ID | ext-system-monitor |
| Área / categoria | Extensões / Monitorização |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gnome-shell-extension-system-monitor |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Terminar sessão | Sim |

Ative em Extensões após terminar sessão.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-system-monitor",
  "name": "System Monitor (GNOME oficial)",
  "area": "extensions",
  "category": "Monitorização",
  "summary": "Monitor de recursos oficial do projeto GNOME.",
  "why": "Leve e mantido com o próprio GNOME; alternativa ao Vitals (escolha um).",
  "method": "apt",
  "packages": [
    "gnome-shell-extension-system-monitor"
  ],
  "group": "monitor",
  "relogin": true,
  "notes": "Ative em Extensões após terminar sessão.",
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ext-gsconnect"></a>

### GSConnect

Integração com o telemóvel Android (KDE Connect).

Notificações, ficheiros e área de transferência entre telemóvel e PC. Se usar firewall, abra 1714-1764 TCP/UDP só na rede local.

| Campo | Valor |
| --- | --- |
| ID | ext-gsconnect |
| Área / categoria | Extensões / Integração |
| Método | apt |
| Fonte | Arquivo Ubuntu |
| Pacotes | gnome-shell-extension-gsconnect |
| Evidência da receita | Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04) |
| Terminar sessão | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-gsconnect",
  "name": "GSConnect",
  "area": "extensions",
  "category": "Integração",
  "summary": "Integração com o telemóvel Android (KDE Connect).",
  "why": "Notificações, ficheiros e área de transferência entre telemóvel e PC. Se usar firewall, abra 1714-1764 TCP/UDP só na rede local.",
  "method": "apt",
  "packages": [
    "gnome-shell-extension-gsconnect"
  ],
  "relogin": true,
  "source": "Arquivo Ubuntu",
  "verified": "Verificado: arquivo Ubuntu 26.04 (apt-cache policy, 2026-10-04)"
}
```

</details>

<a id="item-ext-dash-to-panel"></a>

### Dash to Panel

Barra única estilo Windows.

Substitui a dock: exige desativar a Ubuntu Dock para não haver conflito; por isso é manual.

| Campo | Valor |
| --- | --- |
| ID | ext-dash-to-panel |
| Área / categoria | Extensões / Interface |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Avançado | Sim |

#### Guia manual

1. Desative a Ubuntu Dock em Extensões.
2. Instale Dash to Panel pelo Extension Manager (verifica a compatibilidade com o seu GNOME).
3. Para voltar: desative Dash to Panel e reative Ubuntu Dock.

[Documentação do fornecedor](https://extensions.gnome.org/extension/1160/dash-to-panel/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ext-dash-to-panel",
  "name": "Dash to Panel",
  "area": "extensions",
  "category": "Interface",
  "summary": "Barra única estilo Windows.",
  "why": "Substitui a dock: exige desativar a Ubuntu Dock para não haver conflito; por isso é manual.",
  "method": "manual",
  "manual": {
    "steps": [
      "Desative a Ubuntu Dock em Extensões.",
      "Instale Dash to Panel pelo Extension Manager (verifica a compatibilidade com o seu GNOME).",
      "Para voltar: desative Dash to Panel e reative Ubuntu Dock."
    ],
    "url": "https://extensions.gnome.org/extension/1160/dash-to-panel/"
  },
  "source": "Manual / não suportado automaticamente",
  "advanced": true
}
```

</details>

<a id="item-power-balanced"></a>

### Energia: equilibrado

Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).

Padrão recomendado para portáteis e uso geral.

| Campo | Valor |
| --- | --- |
| ID | power-balanced |
| Área / categoria | Energia e hardware / Perfil de energia |
| Método | power-profile |
| Fonte | powerprofilesctl (power-profiles-daemon) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "power-balanced",
  "name": "Energia: equilibrado",
  "area": "energy",
  "category": "Perfil de energia",
  "summary": "Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).",
  "why": "Padrão recomendado para portáteis e uso geral.",
  "method": "power-profile",
  "powerProfile": "balanced",
  "group": "power-profile",
  "source": "powerprofilesctl (power-profiles-daemon)"
}
```

</details>

<a id="item-power-saver"></a>

### Energia: poupança

Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).

Prolonga a bateria reduzindo desempenho.

| Campo | Valor |
| --- | --- |
| ID | power-saver |
| Área / categoria | Energia e hardware / Perfil de energia |
| Método | power-profile |
| Fonte | powerprofilesctl (power-profiles-daemon) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "power-saver",
  "name": "Energia: poupança",
  "area": "energy",
  "category": "Perfil de energia",
  "summary": "Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).",
  "why": "Prolonga a bateria reduzindo desempenho.",
  "method": "power-profile",
  "powerProfile": "power-saver",
  "group": "power-profile",
  "source": "powerprofilesctl (power-profiles-daemon)"
}
```

</details>

<a id="item-power-performance"></a>

### Energia: desempenho

Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).

Ligado à corrente, para compilações/containers. Ganho depende do CPU (não medido); mais calor e ruído.

| Campo | Valor |
| --- | --- |
| ID | power-performance |
| Área / categoria | Energia e hardware / Perfil de energia |
| Método | power-profile |
| Fonte | powerprofilesctl (power-profiles-daemon) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "power-performance",
  "name": "Energia: desempenho",
  "area": "energy",
  "category": "Perfil de energia",
  "summary": "Define o perfil do power-profiles-daemon (o gestor padrão do Ubuntu).",
  "why": "Ligado à corrente, para compilações/containers. Ganho depende do CPU (não medido); mais calor e ruído.",
  "method": "power-profile",
  "powerProfile": "performance",
  "group": "power-profile",
  "source": "powerprofilesctl (power-profiles-daemon)"
}
```

</details>

<a id="item-screen-blank"></a>

### Ecrã: apagar após inatividade

Tempo até o ecrã apagar.

Equilíbrio entre privacidade/energia e não apagar durante leitura.

| Campo | Valor |
| --- | --- |
| ID | screen-blank |
| Área / categoria | Energia e hardware / Ecrã e suspensão |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

Opção **Tempo** (`delay`): um valor. Padrão: uint32 600. 

| Valor | Descrição |
| --- | --- |
| uint32 300 | 5 minutos (padrão) |
| uint32 600 | 10 minutos |
| uint32 900 | 15 minutos |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "screen-blank",
  "name": "Ecrã: apagar após inatividade",
  "area": "energy",
  "category": "Ecrã e suspensão",
  "summary": "Tempo até o ecrã apagar.",
  "why": "Equilíbrio entre privacidade/energia e não apagar durante leitura.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.desktop.session",
      "key": "idle-delay",
      "option": "delay"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)",
  "options": [
    {
      "key": "delay",
      "label": "Tempo",
      "values": [
        {
          "value": "uint32 300",
          "label": "5 minutos (padrão)"
        },
        {
          "value": "uint32 600",
          "label": "10 minutos"
        },
        {
          "value": "uint32 900",
          "label": "15 minutos"
        }
      ],
      "default": [
        "uint32 600"
      ]
    }
  ]
}
```

</details>

<a id="item-suspend-ac"></a>

### Não suspender ligado à corrente

Evita suspensão automática com carregador.

Downloads, builds e acessos remotos não são interrompidos.

| Campo | Valor |
| --- | --- |
| ID | suspend-ac |
| Área / categoria | Energia e hardware / Ecrã e suspensão |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "suspend-ac",
  "name": "Não suspender ligado à corrente",
  "area": "energy",
  "category": "Ecrã e suspensão",
  "summary": "Evita suspensão automática com carregador.",
  "why": "Downloads, builds e acessos remotos não são interrompidos.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.settings-daemon.plugins.power",
      "key": "sleep-inactive-ac-type",
      "value": "'nothing'"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-suspend-battery"></a>

### Suspender na bateria após 15 min

Suspende em bateria após 15 minutos inativo.

Protege a bateria quando se esquece do portátil.

| Campo | Valor |
| --- | --- |
| ID | suspend-battery |
| Área / categoria | Energia e hardware / Ecrã e suspensão |
| Método | gsettings |
| Fonte | Preferência GNOME (gsettings, utilizador) |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "suspend-battery",
  "name": "Suspender na bateria após 15 min",
  "area": "energy",
  "category": "Ecrã e suspensão",
  "summary": "Suspende em bateria após 15 minutos inativo.",
  "why": "Protege a bateria quando se esquece do portátil.",
  "method": "gsettings",
  "settings": [
    {
      "schema": "org.gnome.settings-daemon.plugins.power",
      "key": "sleep-inactive-battery-type",
      "value": "'suspend'"
    },
    {
      "schema": "org.gnome.settings-daemon.plugins.power",
      "key": "sleep-inactive-battery-timeout",
      "value": "900"
    }
  ],
  "source": "Preferência GNOME (gsettings, utilizador)"
}
```

</details>

<a id="item-lid-close"></a>

### Comportamento da tampa

Ação ao fechar a tampa.

Gerido pelo systemd-logind (HandleLidSwitch); alterar exige editar /etc/systemd/logind.conf.d — não automatizado.

| Campo | Valor |
| --- | --- |
| ID | lid-close |
| Área / categoria | Energia e hardware / Ecrã e suspensão |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |

#### Guia manual

1. Padrão: suspender. Com monitor externo ligado, o GNOME não suspende.
2. Para mudar: crie /etc/systemd/logind.conf.d/tampa.conf com [Login] HandleLidSwitchExternalPower=ignore e reinicie.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "lid-close",
  "name": "Comportamento da tampa",
  "area": "energy",
  "category": "Ecrã e suspensão",
  "summary": "Ação ao fechar a tampa.",
  "why": "Gerido pelo systemd-logind (HandleLidSwitch); alterar exige editar /etc/systemd/logind.conf.d — não automatizado.",
  "method": "manual",
  "manual": {
    "steps": [
      "Padrão: suspender. Com monitor externo ligado, o GNOME não suspende.",
      "Para mudar: crie /etc/systemd/logind.conf.d/tampa.conf com [Login] HandleLidSwitchExternalPower=ignore e reinicie."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente"
}
```

</details>

<a id="item-fstrim-timer"></a>

### TRIM semanal de SSD

Ativa fstrim.timer (padrão do Ubuntu).

Mantém o desempenho do SSD. Só ative se o Diagnóstico mostrar SSD com discard e o timer desativado.

| Campo | Valor |
| --- | --- |
| ID | fstrim-timer |
| Área / categoria | Energia e hardware / Disco |
| Método | admin |
| Fonte | systemctl enable fstrim.timer |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "fstrim-timer",
  "name": "TRIM semanal de SSD",
  "area": "energy",
  "category": "Disco",
  "summary": "Ativa fstrim.timer (padrão do Ubuntu).",
  "why": "Mantém o desempenho do SSD. Só ative se o Diagnóstico mostrar SSD com discard e o timer desativado.",
  "method": "admin",
  "adminOp": "fstrim-timer",
  "source": "systemctl enable fstrim.timer"
}
```

</details>

<a id="item-journald-limit"></a>

### Limitar logs do sistema

Define SystemMaxUse no journald.

Logs antigos podem ocupar GBs; o limite remove os mais antigos sem desligar logs. Reversível.

| Campo | Valor |
| --- | --- |
| ID | journald-limit |
| Área / categoria | Energia e hardware / Disco |
| Método | admin |
| Fonte | /etc/systemd/journald.conf.d (drop-in próprio, reversível) |

Opção **Limite** (`value`): um valor. Padrão: 1G. 

| Valor | Descrição |
| --- | --- |
| 500M | 500 MB |
| 1G | 1 GB |
| 2G | 2 GB |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "journald-limit",
  "name": "Limitar logs do sistema",
  "area": "energy",
  "category": "Disco",
  "summary": "Define SystemMaxUse no journald.",
  "why": "Logs antigos podem ocupar GBs; o limite remove os mais antigos sem desligar logs. Reversível.",
  "method": "admin",
  "adminOp": "journald",
  "source": "/etc/systemd/journald.conf.d (drop-in próprio, reversível)",
  "options": [
    {
      "key": "value",
      "label": "Limite",
      "values": [
        {
          "value": "500M",
          "label": "500 MB"
        },
        {
          "value": "1G",
          "label": "1 GB"
        },
        {
          "value": "2G",
          "label": "2 GB"
        }
      ],
      "default": [
        "1G"
      ]
    }
  ]
}
```

</details>

<a id="item-ubuntu-drivers"></a>

### Drivers adicionais recomendados

Instala os drivers que o ubuntu-drivers recomenda (ex.: NVIDIA).

Só se o Diagnóstico listar drivers. Requer reiniciar.

| Campo | Valor |
| --- | --- |
| ID | ubuntu-drivers |
| Área / categoria | Energia e hardware / Hardware |
| Método | admin |
| Fonte | ubuntu-drivers install (ferramenta oficial) |
| Avançado | Sim |
| Reiniciar | Sim |

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "ubuntu-drivers",
  "name": "Drivers adicionais recomendados",
  "area": "energy",
  "category": "Hardware",
  "summary": "Instala os drivers que o ubuntu-drivers recomenda (ex.: NVIDIA).",
  "why": "Só se o Diagnóstico listar drivers. Requer reiniciar.",
  "method": "admin",
  "adminOp": "ubuntu-drivers",
  "source": "ubuntu-drivers install (ferramenta oficial)",
  "reboot": true,
  "advanced": true
}
```

</details>

<a id="item-tlp"></a>

### TLP

Gestor de energia alternativo.

Não recomendado: conflitua com o power-profiles-daemon padrão. Só remova o PPD se souber porquê.

| Campo | Valor |
| --- | --- |
| ID | tlp |
| Área / categoria | Energia e hardware / Avançado |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Avançado | Sim |

#### Guia manual

1. Não suportado automaticamente para evitar dois gestores de energia.

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "tlp",
  "name": "TLP",
  "area": "energy",
  "category": "Avançado",
  "summary": "Gestor de energia alternativo.",
  "why": "Não recomendado: conflitua com o power-profiles-daemon padrão. Só remova o PPD se souber porquê.",
  "method": "manual",
  "manual": {
    "steps": [
      "Não suportado automaticamente para evitar dois gestores de energia."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente",
  "advanced": true
}
```

</details>

<a id="item-zram"></a>

### zram / swappiness

Compressão de memória e agressividade do swap.

Não alterado sem evidência: só considere com pouca RAM e swap usado frequentemente (ver Diagnóstico).

| Campo | Valor |
| --- | --- |
| ID | zram |
| Área / categoria | Energia e hardware / Avançado |
| Método | manual |
| Fonte | Manual / não suportado automaticamente |
| Avançado | Sim |

#### Guia manual

1. Meça primeiro: free -h e vmstat 5 durante o seu trabalho normal.
2. Se necessário: sudo apt install systemd-zram-generator (documentação Ubuntu).

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "zram",
  "name": "zram / swappiness",
  "area": "energy",
  "category": "Avançado",
  "summary": "Compressão de memória e agressividade do swap.",
  "why": "Não alterado sem evidência: só considere com pouca RAM e swap usado frequentemente (ver Diagnóstico).",
  "method": "manual",
  "manual": {
    "steps": [
      "Meça primeiro: free -h e vmstat 5 durante o seu trabalho normal.",
      "Se necessário: sudo apt install systemd-zram-generator (documentação Ubuntu)."
    ],
    "url": ""
  },
  "source": "Manual / não suportado automaticamente",
  "advanced": true
}
```

</details>

<a id="item-claude-code"></a>

### Claude Code

Agente da Anthropic para trabalhar com código no terminal.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | claude-code |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://code.claude.com/docs/en/setup |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | @anthropic-ai/claude-code@2.1.296 |
| Comando | ~/.local/bin/claude |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/claude
2. Entre na sua conta seguindo o assistente do Claude Code.
3. Para usar o comando claude sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://code.claude.com/docs/en/setup)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "claude-code",
  "name": "Claude Code",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Agente da Anthropic para trabalhar com código no terminal.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://code.claude.com/docs/en/setup",
  "npm": {
    "package": "@anthropic-ai/claude-code",
    "version": "2.1.296",
    "binary": "claude"
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/claude",
      "Entre na sua conta seguindo o assistente do Claude Code.",
      "Para usar o comando claude sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://code.claude.com/docs/en/setup"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-codex-cli"></a>

### Codex

Agente da OpenAI para editar código, executar tarefas e verificar resultados.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | codex-cli |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://learn.chatgpt.com/docs/codex/cli |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | @openai/codex@0.162.1 |
| Comando | ~/.local/bin/codex |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/codex
2. Entre com a sua conta ChatGPT ou configure o acesso indicado pelo Codex.
3. Para usar o comando codex sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://learn.chatgpt.com/docs/codex/cli)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "codex-cli",
  "name": "Codex",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Agente da OpenAI para editar código, executar tarefas e verificar resultados.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://learn.chatgpt.com/docs/codex/cli",
  "npm": {
    "package": "@openai/codex",
    "version": "0.162.1",
    "binary": "codex"
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/codex",
      "Entre com a sua conta ChatGPT ou configure o acesso indicado pelo Codex.",
      "Para usar o comando codex sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://learn.chatgpt.com/docs/codex/cli"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-gemini-cli"></a>

### Gemini CLI

Agente da Google para explorar e modificar projetos no terminal.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | gemini-cli |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://geminicli.com/docs/get-started/installation/ |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | @google/gemini-cli@0.63.0 |
| Comando | ~/.local/bin/gemini |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/gemini
2. Siga o assistente de autenticação do Gemini CLI e escolha o seu método de acesso.
3. Para usar o comando gemini sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://geminicli.com/docs/get-started/installation/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gemini-cli",
  "name": "Gemini CLI",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Agente da Google para explorar e modificar projetos no terminal.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://geminicli.com/docs/get-started/installation/",
  "npm": {
    "package": "@google/gemini-cli",
    "version": "0.63.0",
    "binary": "gemini"
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/gemini",
      "Siga o assistente de autenticação do Gemini CLI e escolha o seu método de acesso.",
      "Para usar o comando gemini sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://geminicli.com/docs/get-started/installation/"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-opencode"></a>

### OpenCode

Agente para programar com diferentes fornecedores de modelos.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | opencode |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://opencode.ai/docs/ |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | opencode-ai@1.19.0 |
| Comando | ~/.local/bin/opencode |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/opencode
2. Dentro do OpenCode, use /connect para ligar o fornecedor que pretende usar.
3. Para usar o comando opencode sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://opencode.ai/docs/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "opencode",
  "name": "OpenCode",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Agente para programar com diferentes fornecedores de modelos.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://opencode.ai/docs/",
  "npm": {
    "package": "opencode-ai",
    "version": "1.19.0",
    "binary": "opencode"
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/opencode",
      "Dentro do OpenCode, use /connect para ligar o fornecedor que pretende usar.",
      "Para usar o comando opencode sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://opencode.ai/docs/"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-pi-coding-agent"></a>

### Pi

Agente de programação extensível com suporte a vários fornecedores.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | pi-coding-agent |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://github.com/earendil-works/pi/tree/main/packages/coding-agent |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | @earendil-works/pi-coding-agent@1.1.0 |
| Comando | ~/.local/bin/pi |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/pi
2. Dentro do Pi, use /login para escolher e ligar o fornecedor de modelos.
3. Para usar o comando pi sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://github.com/earendil-works/pi/tree/main/packages/coding-agent)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "pi-coding-agent",
  "name": "Pi",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Agente de programação extensível com suporte a vários fornecedores.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://github.com/earendil-works/pi/tree/main/packages/coding-agent",
  "npm": {
    "package": "@earendil-works/pi-coding-agent",
    "version": "1.1.0",
    "binary": "pi",
    "ignoreScripts": true
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/pi",
      "Dentro do Pi, use /login para escolher e ligar o fornecedor de modelos.",
      "Para usar o comando pi sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://github.com/earendil-works/pi/tree/main/packages/coding-agent"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-github-copilot-cli"></a>

### GitHub Copilot CLI

Assistente de programação no terminal integrado com o GitHub Copilot.

Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.

| Campo | Valor |
| --- | --- |
| ID | github-copilot-cli |
| Área / categoria | Desenvolvimento / Agentes de IA |
| Método | npm-user |
| Fonte | https://docs.github.com/en/copilot/how-tos/copilot-cli/set-up-copilot-cli/install-copilot-cli |
| Evidência da receita | Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar. |
| Confirmação | Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador. |
| Dependências | [Node.js do sistema (APT)](#item-nodejs) |
| Pacote npm fixo | @github/copilot@1.0.95 |
| Comando | ~/.local/bin/copilot |
| Destino | ~/.local (utilizador, sem sudo) |

A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.

#### Guia manual

1. Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:
~/.local/bin/copilot
2. Dentro do Copilot CLI, use /login e siga as instruções da sua conta GitHub.
3. Para usar o comando copilot sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH.

[Documentação do fornecedor](https://docs.github.com/en/copilot/how-tos/copilot-cli/set-up-copilot-cli/install-copilot-cli)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "github-copilot-cli",
  "name": "GitHub Copilot CLI",
  "area": "dev",
  "category": "Agentes de IA",
  "summary": "Assistente de programação no terminal integrado com o GitHub Copilot.",
  "why": "Instalação para o seu utilizador, com versão conhecida e verificação antes de usar.",
  "method": "npm-user",
  "source": "https://docs.github.com/en/copilot/how-tos/copilot-cli/set-up-copilot-cli/install-copilot-cli",
  "npm": {
    "package": "@github/copilot",
    "version": "1.0.95",
    "binary": "copilot"
  },
  "requires": [
    "nodejs"
  ],
  "detect": {},
  "verify": {},
  "confirm": "Aceito instalar este pacote oficial e o código das suas dependências como o meu utilizador.",
  "notes": "A instalação existente é preservada. O acesso aos modelos e eventuais custos são tratados pelo fornecedor, dentro do agente. O Meu Ubuntu não recolhe credenciais.",
  "manual": {
    "steps": [
      "Depois de aplicar o plano, abra um terminal na pasta do projeto e execute:\n~/.local/bin/copilot",
      "Dentro do Copilot CLI, use /login e siga as instruções da sua conta GitHub.",
      "Para usar o comando copilot sem o caminho completo, termine e volte a iniciar a sessão do Ubuntu caso ~/.local/bin ainda não faça parte do PATH."
    ],
    "url": "https://docs.github.com/en/copilot/how-tos/copilot-cli/set-up-copilot-cli/install-copilot-cli"
  },
  "verified": "Instalação por utilizador e comando --version verificados em container Ubuntu 26.04, em 2026-10-12. Autenticação e utilização de modelos por testar."
}
```

</details>

<a id="item-claude-code-plugins-guide"></a>

### Claude Code: gerir plugins

Descobrir, instalar e consultar plugins no catálogo do Claude Code.

Centraliza o guia de complementos junto da ferramenta correspondente.

| Campo | Valor |
| --- | --- |
| ID | claude-code-plugins-guide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Documentação oficial da ferramenta |
| Evidência da receita | Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente. |
| Dependências | [Claude Code](#item-claude-code) |

#### Guia manual

1. Abra claude no seu projeto. Dentro da sessão, execute /plugin para explorar o catálogo e consultar os detalhes dos plugins.
2. Escolha o plugin e reveja o que acrescenta. Prefira o âmbito de utilizador para uso pessoal; o âmbito de projeto altera a configuração partilhada do repositório.
3. Para consultar os plugins instalados no terminal:
```sh
claude plugin list
```
Confirme o âmbito e o estado. A instalação e ativação são feitas pelo Claude Code.

[Documentação do fornecedor](https://code.claude.com/docs/en/discover-plugins)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "claude-code-plugins-guide",
  "name": "Claude Code: gerir plugins",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Descobrir, instalar e consultar plugins no catálogo do Claude Code.",
  "why": "Centraliza o guia de complementos junto da ferramenta correspondente.",
  "method": "manual",
  "source": "Documentação oficial da ferramenta",
  "addonsFor": [
    "claude-code"
  ],
  "requires": [
    "claude-code"
  ],
  "manual": {
    "steps": [
      "Abra claude no seu projeto. Dentro da sessão, execute /plugin para explorar o catálogo e consultar os detalhes dos plugins.",
      "Escolha o plugin e reveja o que acrescenta. Prefira o âmbito de utilizador para uso pessoal; o âmbito de projeto altera a configuração partilhada do repositório.",
      "Para consultar os plugins instalados no terminal:\n```sh\nclaude plugin list\n```\nConfirme o âmbito e o estado. A instalação e ativação são feitas pelo Claude Code."
    ],
    "url": "https://code.claude.com/docs/en/discover-plugins"
  },
  "verified": "Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente."
}
```

</details>

<a id="item-codex-skills-guide"></a>

### Codex: skills e plugins

Acrescentar fluxos reutilizáveis através das skills do Codex.

Centraliza o guia de complementos junto da ferramenta correspondente.

| Campo | Valor |
| --- | --- |
| ID | codex-skills-guide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Documentação oficial da ferramenta |
| Evidência da receita | Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente. |
| Dependências | [Codex](#item-codex-cli) |

#### Guia manual

1. Abra codex no projeto. Dentro da sessão, use /skills para consultar as skills disponíveis.
2. Peça ao Codex: Use $skill-installer para listar as skills disponíveis. Escolha as que servem o seu trabalho e peça a instalação pelo nome.
3. Skills pessoais podem ficar em ~/.agents/skills; skills partilhadas com a equipa ficam em .agents/skills no repositório. Cada skill contém um SKILL.md. Se uma nova skill não aparecer, reinicie o Codex. A documentação explica também a distribuição de skills em plugins.

[Documentação do fornecedor](https://learn.chatgpt.com/docs/build-skills)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "codex-skills-guide",
  "name": "Codex: skills e plugins",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Acrescentar fluxos reutilizáveis através das skills do Codex.",
  "why": "Centraliza o guia de complementos junto da ferramenta correspondente.",
  "method": "manual",
  "source": "Documentação oficial da ferramenta",
  "addonsFor": [
    "codex-cli"
  ],
  "requires": [
    "codex-cli"
  ],
  "manual": {
    "steps": [
      "Abra codex no projeto. Dentro da sessão, use /skills para consultar as skills disponíveis.",
      "Peça ao Codex: Use $skill-installer para listar as skills disponíveis. Escolha as que servem o seu trabalho e peça a instalação pelo nome.",
      "Skills pessoais podem ficar em ~/.agents/skills; skills partilhadas com a equipa ficam em .agents/skills no repositório. Cada skill contém um SKILL.md. Se uma nova skill não aparecer, reinicie o Codex. A documentação explica também a distribuição de skills em plugins."
    ],
    "url": "https://learn.chatgpt.com/docs/build-skills"
  },
  "verified": "Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente."
}
```

</details>

<a id="item-opencode-plugins-guide"></a>

### OpenCode: plugins e configuração

Organizar plugins npm e plugins locais sem substituir a configuração existente.

Centraliza o guia de complementos junto da ferramenta correspondente.

| Campo | Valor |
| --- | --- |
| ID | opencode-plugins-guide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Documentação oficial da ferramenta |
| Evidência da receita | Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente. |
| Dependências | [OpenCode](#item-opencode) |

#### Guia manual

1. Consulte os plugins no ecossistema ligado pela documentação oficial. Reveja o pacote e escolha uma versão antes de o adicionar.
2. Faça backup de ~/.config/opencode/opencode.json. Acrescente o pacote escolhido à lista plugin existente; preserve as restantes opções. O OpenCode instala os pacotes npm configurados ao arrancar.
3. Plugins locais pessoais ficam em ~/.config/opencode/plugins/; os do projeto ficam em .opencode/plugins/. Abra opencode e confirme o carregamento. Um plugin local e um plugin npm semelhantes podem carregar em simultâneo.

[Documentação do fornecedor](https://opencode.ai/docs/plugins/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "opencode-plugins-guide",
  "name": "OpenCode: plugins e configuração",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Organizar plugins npm e plugins locais sem substituir a configuração existente.",
  "why": "Centraliza o guia de complementos junto da ferramenta correspondente.",
  "method": "manual",
  "source": "Documentação oficial da ferramenta",
  "addonsFor": [
    "opencode"
  ],
  "requires": [
    "opencode"
  ],
  "manual": {
    "steps": [
      "Consulte os plugins no ecossistema ligado pela documentação oficial. Reveja o pacote e escolha uma versão antes de o adicionar.",
      "Faça backup de ~/.config/opencode/opencode.json. Acrescente o pacote escolhido à lista plugin existente; preserve as restantes opções. O OpenCode instala os pacotes npm configurados ao arrancar.",
      "Plugins locais pessoais ficam em ~/.config/opencode/plugins/; os do projeto ficam em .opencode/plugins/. Abra opencode e confirme o carregamento. Um plugin local e um plugin npm semelhantes podem carregar em simultâneo."
    ],
    "url": "https://opencode.ai/docs/plugins/"
  },
  "verified": "Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente."
}
```

</details>

<a id="item-pi-packages-guide"></a>

### Pi: pacotes, extensões e skills

Gerir pacotes que podem incluir extensões, skills, prompts e temas.

Centraliza o guia de complementos junto da ferramenta correspondente.

| Campo | Valor |
| --- | --- |
| ID | pi-packages-guide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Documentação oficial da ferramenta |
| Evidência da receita | Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente. |
| Dependências | [Pi](#item-pi-coding-agent) |

#### Guia manual

1. Consulte os pacotes já configurados:
```sh
pi list
```
2. Depois de escolher e rever uma fonte, use pi install com o pacote npm, repositório Git ou caminho local. A documentação explica como fixar uma versão e escolher o âmbito pessoal ou do projeto.
3. Escolha os recursos que ficam ativos:
```sh
pi config
```
Pode ativar ou desativar extensões, skills, prompts e temas sem carregar tudo o que um pacote oferece.

[Documentação do fornecedor](https://github.com/earendil-works/pi/blob/main/packages/coding-agent/docs/packages.md)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "pi-packages-guide",
  "name": "Pi: pacotes, extensões e skills",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Gerir pacotes que podem incluir extensões, skills, prompts e temas.",
  "why": "Centraliza o guia de complementos junto da ferramenta correspondente.",
  "method": "manual",
  "source": "Documentação oficial da ferramenta",
  "addonsFor": [
    "pi-coding-agent"
  ],
  "requires": [
    "pi-coding-agent"
  ],
  "manual": {
    "steps": [
      "Consulte os pacotes já configurados:\n```sh\npi list\n```",
      "Depois de escolher e rever uma fonte, use pi install com o pacote npm, repositório Git ou caminho local. A documentação explica como fixar uma versão e escolher o âmbito pessoal ou do projeto.",
      "Escolha os recursos que ficam ativos:\n```sh\npi config\n```\nPode ativar ou desativar extensões, skills, prompts e temas sem carregar tudo o que um pacote oferece."
    ],
    "url": "https://github.com/earendil-works/pi/blob/main/packages/coding-agent/docs/packages.md"
  },
  "verified": "Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente."
}
```

</details>

<a id="item-gemini-extensions-guide"></a>

### Gemini CLI: gerir extensões

Explorar extensões e consultar o seu estado no Gemini CLI.

Centraliza o guia de complementos junto da ferramenta correspondente.

| Campo | Valor |
| --- | --- |
| ID | gemini-extensions-guide |
| Área / categoria | Desenvolvimento / Plugins e complementos |
| Método | manual |
| Fonte | Documentação oficial da ferramenta |
| Evidência da receita | Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente. |
| Dependências | [Gemini CLI](#item-gemini-cli) |

#### Guia manual

1. Explore a galeria de extensões ligada pela documentação oficial. Reveja a origem e os recursos que cada extensão acrescenta.
2. Consulte as extensões instaladas no terminal:
```sh
gemini extensions list
```
Dentro de uma sessão Gemini, também pode usar /extensions list.
3. Instale a extensão escolhida com gemini extensions install seguido da URL do repositório. Confirme o estado na lista antes de a utilizar. A instalação é feita pelo Gemini CLI.

[Documentação do fornecedor](https://geminicli.com/docs/extensions/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "gemini-extensions-guide",
  "name": "Gemini CLI: gerir extensões",
  "area": "dev",
  "category": "Plugins e complementos",
  "summary": "Explorar extensões e consultar o seu estado no Gemini CLI.",
  "why": "Centraliza o guia de complementos junto da ferramenta correspondente.",
  "method": "manual",
  "source": "Documentação oficial da ferramenta",
  "addonsFor": [
    "gemini-cli"
  ],
  "requires": [
    "gemini-cli"
  ],
  "manual": {
    "steps": [
      "Explore a galeria de extensões ligada pela documentação oficial. Reveja a origem e os recursos que cada extensão acrescenta.",
      "Consulte as extensões instaladas no terminal:\n```sh\ngemini extensions list\n```\nDentro de uma sessão Gemini, também pode usar /extensions list.",
      "Instale a extensão escolhida com gemini extensions install seguido da URL do repositório. Confirme o estado na lista antes de a utilizar. A instalação é feita pelo Gemini CLI."
    ],
    "url": "https://geminicli.com/docs/extensions/"
  },
  "verified": "Guia conferido na documentação oficial em 2026-10-12; não instala nem deteta plugins automaticamente."
}
```

</details>

<a id="item-apache2"></a>

### Apache HTTP Server (local)

Servidor web para desenvolvimento e testes de aplicações.

Servir projetos locais, incluindo aplicações PHP, ou testar configurações Apache.

| Campo | Valor |
| --- | --- |
| ID | apache2 |
| Área / categoria | Desenvolvimento / Serviços locais · Web |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | apache2 |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Conflitos | [Nginx (local)](#item-nginx) (aviso): Por omissão, Apache e Nginx usam a porta 80. Configure portas ou interfaces distintas para os usar em simultâneo. |
| Serviço local | Sim |

Pode escutar na porta 80 em todas as interfaces. Instalar não configura os seus projetos nem restringe o acesso à rede. Configuração em /etc/apache2/.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status apache2
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. Pode escutar na porta 80 em todas as interfaces. Instalar não configura os seus projetos nem restringe o acesso à rede. Configuração em /etc/apache2/.

[Documentação do fornecedor](https://ubuntu.com/server/docs/how-to/web-services/install-apache2/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "apache2",
  "name": "Apache HTTP Server (local)",
  "area": "dev",
  "category": "Serviços locais · Web",
  "summary": "Servidor web para desenvolvimento e testes de aplicações.",
  "why": "Servir projetos locais, incluindo aplicações PHP, ou testar configurações Apache.",
  "method": "apt",
  "packages": [
    "apache2"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "Pode escutar na porta 80 em todas as interfaces. Instalar não configura os seus projetos nem restringe o acesso à rede. Configuração em /etc/apache2/.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://ubuntu.com/server/docs/how-to/web-services/install-apache2/",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status apache2\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "Pode escutar na porta 80 em todas as interfaces. Instalar não configura os seus projetos nem restringe o acesso à rede. Configuração em /etc/apache2/."
    ]
  },
  "verify": {
    "cmd": [
      "/usr/sbin/apache2",
      "-v"
    ]
  },
  "conflicts": [
    {
      "id": "nginx",
      "hard": false,
      "reason": "Por omissão, Apache e Nginx usam a porta 80. Configure portas ou interfaces distintas para os usar em simultâneo."
    }
  ]
}
```

</details>

<a id="item-nginx"></a>

### Nginx (local)

Servidor web e proxy inverso para aplicações locais.

Testar aplicações servidas por Nginx e configurações de proxy inverso.

| Campo | Valor |
| --- | --- |
| ID | nginx |
| Área / categoria | Desenvolvimento / Serviços locais · Web |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | nginx |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Conflitos | [Apache HTTP Server (local)](#item-apache2) (aviso): Por omissão, Apache e Nginx usam a porta 80. Configure portas ou interfaces distintas para os usar em simultâneo. |
| Serviço local | Sim |

Pode escutar na porta 80 em todas as interfaces. Instalar não configura proxies, certificados ou acesso à rede. Configuração em /etc/nginx/.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status nginx
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. Pode escutar na porta 80 em todas as interfaces. Instalar não configura proxies, certificados ou acesso à rede. Configuração em /etc/nginx/.

[Documentação do fornecedor](https://ubuntu.com/server/docs/how-to/web-services/install-nginx/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "nginx",
  "name": "Nginx (local)",
  "area": "dev",
  "category": "Serviços locais · Web",
  "summary": "Servidor web e proxy inverso para aplicações locais.",
  "why": "Testar aplicações servidas por Nginx e configurações de proxy inverso.",
  "method": "apt",
  "packages": [
    "nginx"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "Pode escutar na porta 80 em todas as interfaces. Instalar não configura proxies, certificados ou acesso à rede. Configuração em /etc/nginx/.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://ubuntu.com/server/docs/how-to/web-services/install-nginx/",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status nginx\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "Pode escutar na porta 80 em todas as interfaces. Instalar não configura proxies, certificados ou acesso à rede. Configuração em /etc/nginx/."
    ]
  },
  "verify": {
    "cmd": [
      "/usr/sbin/nginx",
      "-v"
    ]
  },
  "conflicts": [
    {
      "id": "apache2",
      "hard": false,
      "reason": "Por omissão, Apache e Nginx usam a porta 80. Configure portas ou interfaces distintas para os usar em simultâneo."
    }
  ]
}
```

</details>

<a id="item-redis-server"></a>

### Redis (servidor local)

Servidor Redis para cache, sessões e filas de desenvolvimento.

Testar aplicações que usam Redis sem depender de um servidor remoto.

| Campo | Valor |
| --- | --- |
| ID | redis-server |
| Área / categoria | Desenvolvimento / Serviços locais · Dados e cache |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | redis-server |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Serviço local | Sim |

O serviço usa normalmente a porta 6379. Reveja /etc/redis/redis.conf, interfaces e autenticação antes de permitir acesso pela rede. Não altera a configuração das suas aplicações.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status redis-server
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. O serviço usa normalmente a porta 6379. Reveja /etc/redis/redis.conf, interfaces e autenticação antes de permitir acesso pela rede. Não altera a configuração das suas aplicações.

[Documentação do fornecedor](https://redis.io/docs/latest/operate/oss_and_stack/management/config/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "redis-server",
  "name": "Redis (servidor local)",
  "area": "dev",
  "category": "Serviços locais · Dados e cache",
  "summary": "Servidor Redis para cache, sessões e filas de desenvolvimento.",
  "why": "Testar aplicações que usam Redis sem depender de um servidor remoto.",
  "method": "apt",
  "packages": [
    "redis-server"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "O serviço usa normalmente a porta 6379. Reveja /etc/redis/redis.conf, interfaces e autenticação antes de permitir acesso pela rede. Não altera a configuração das suas aplicações.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://redis.io/docs/latest/operate/oss_and_stack/management/config/",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status redis-server\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "O serviço usa normalmente a porta 6379. Reveja /etc/redis/redis.conf, interfaces e autenticação antes de permitir acesso pela rede. Não altera a configuração das suas aplicações."
    ]
  },
  "verify": {
    "cmd": [
      "/usr/bin/redis-server",
      "--version"
    ]
  }
}
```

</details>

<a id="item-rabbitmq-server"></a>

### RabbitMQ (local)

Broker de mensagens para filas e aplicações com AMQP.

Desenvolver e testar produtores, consumidores e processamento assíncrono.

| Campo | Valor |
| --- | --- |
| ID | rabbitmq-server |
| Área / categoria | Desenvolvimento / Serviços locais · Mensageria |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | rabbitmq-server |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Serviço local | Sim |

O serviço AMQP usa normalmente a porta 5672. O painel de gestão não é ativado pelo Meu Ubuntu. Utilizadores, permissões, TLS e acesso pela rede precisam de configuração própria.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status rabbitmq-server
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. O serviço AMQP usa normalmente a porta 5672. O painel de gestão não é ativado pelo Meu Ubuntu. Utilizadores, permissões, TLS e acesso pela rede precisam de configuração própria.

[Documentação do fornecedor](https://www.rabbitmq.com/docs/install-debian)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "rabbitmq-server",
  "name": "RabbitMQ (local)",
  "area": "dev",
  "category": "Serviços locais · Mensageria",
  "summary": "Broker de mensagens para filas e aplicações com AMQP.",
  "why": "Desenvolver e testar produtores, consumidores e processamento assíncrono.",
  "method": "apt",
  "packages": [
    "rabbitmq-server"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "O serviço AMQP usa normalmente a porta 5672. O painel de gestão não é ativado pelo Meu Ubuntu. Utilizadores, permissões, TLS e acesso pela rede precisam de configuração própria.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://www.rabbitmq.com/docs/install-debian",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status rabbitmq-server\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "O serviço AMQP usa normalmente a porta 5672. O painel de gestão não é ativado pelo Meu Ubuntu. Utilizadores, permissões, TLS e acesso pela rede precisam de configuração própria."
    ]
  }
}
```

</details>

<a id="item-postgresql-server"></a>

### PostgreSQL (servidor local)

Servidor PostgreSQL mantido pelo Ubuntu.

Desenvolver e testar aplicações com PostgreSQL instalado neste computador.

| Campo | Valor |
| --- | --- |
| ID | postgresql-server |
| Área / categoria | Desenvolvimento / Serviços locais · Dados e cache |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | postgresql |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Serviço local | Sim |

Pode criar um cluster e ativar um serviço na instalação. A porta habitual é 5432. Roles, bases de dados e regras de acesso não são criadas pelo Meu Ubuntu.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status postgresql
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. Pode criar um cluster e ativar um serviço na instalação. A porta habitual é 5432. Roles, bases de dados e regras de acesso não são criadas pelo Meu Ubuntu.

[Documentação do fornecedor](https://ubuntu.com/server/docs/how-to/databases/install-postgresql/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "postgresql-server",
  "name": "PostgreSQL (servidor local)",
  "area": "dev",
  "category": "Serviços locais · Dados e cache",
  "summary": "Servidor PostgreSQL mantido pelo Ubuntu.",
  "why": "Desenvolver e testar aplicações com PostgreSQL instalado neste computador.",
  "method": "apt",
  "packages": [
    "postgresql"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "Pode criar um cluster e ativar um serviço na instalação. A porta habitual é 5432. Roles, bases de dados e regras de acesso não são criadas pelo Meu Ubuntu.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://ubuntu.com/server/docs/how-to/databases/install-postgresql/",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status postgresql\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "Pode criar um cluster e ativar um serviço na instalação. A porta habitual é 5432. Roles, bases de dados e regras de acesso não são criadas pelo Meu Ubuntu."
    ]
  },
  "verify": {
    "cmd": [
      "/usr/bin/psql",
      "--version"
    ]
  }
}
```

</details>

<a id="item-memcached"></a>

### Memcached (local)

Cache simples em memória para aplicações locais.

Testar aplicações que usam Memcached como cache temporária.

| Campo | Valor |
| --- | --- |
| ID | memcached |
| Área / categoria | Desenvolvimento / Serviços locais · Dados e cache |
| Método | apt |
| Fonte | Arquivo Ubuntu 26.04 |
| Pacotes | memcached |
| Evidência da receita | Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço. |
| Confirmação | Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim. |
| Serviço local | Sim |

O serviço usa normalmente a porta 11211. Reveja /etc/memcached.conf antes de alterar o endereço de escuta. Não oferece persistência para os dados da cache.

#### Guia manual

1. Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.
2. Depois da instalação, confira o serviço no terminal:
```sh
systemctl status memcached
```
O selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.
3. O serviço usa normalmente a porta 11211. Reveja /etc/memcached.conf antes de alterar o endereço de escuta. Não oferece persistência para os dados da cache.

[Documentação do fornecedor](https://docs.memcached.org/serverguide/configuring/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "memcached",
  "name": "Memcached (local)",
  "area": "dev",
  "category": "Serviços locais · Dados e cache",
  "summary": "Cache simples em memória para aplicações locais.",
  "why": "Testar aplicações que usam Memcached como cache temporária.",
  "method": "apt",
  "packages": [
    "memcached"
  ],
  "service": true,
  "source": "Arquivo Ubuntu 26.04",
  "notes": "O serviço usa normalmente a porta 11211. Reveja /etc/memcached.conf antes de alterar o endereço de escuta. Não oferece persistência para os dados da cache.",
  "confirm": "Confirmo a instalação deste servidor local. O pacote pode ativar um serviço e fazê-lo arrancar com o sistema; a configuração de acesso e das portas será revista por mim.",
  "verified": "Pacote e candidato conferidos no arquivo Ubuntu 26.04 em 2026-10-12; o estado deteta a instalação, não a saúde do serviço.",
  "manual": {
    "url": "https://docs.memcached.org/serverguide/configuring/",
    "steps": [
      "Instale pelo plano e confirme o comportamento de serviço local. A versão instalada é a do arquivo Ubuntu 26.04.",
      "Depois da instalação, confira o serviço no terminal:\n```sh\nsystemctl status memcached\n```\nO selo Já instalado confirma o pacote; não confirma que o serviço esteja ativo.",
      "O serviço usa normalmente a porta 11211. Reveja /etc/memcached.conf antes de alterar o endereço de escuta. Não oferece persistência para os dados da cache."
    ]
  },
  "verify": {
    "cmd": [
      "/usr/bin/memcached",
      "-V"
    ]
  }
}
```

</details>

<a id="item-kafka-local-guide"></a>

### Apache Kafka — laboratório Docker

Guia para testar Kafka num contentor local com a imagem oficial.

Testar aplicações baseadas em eventos sem preparar um cluster de produção.

| Campo | Valor |
| --- | --- |
| ID | kafka-local-guide |
| Área / categoria | Desenvolvimento / Serviços locais · Mensageria |
| Método | manual |
| Fonte | Imagem oficial apache/kafka e documentação Apache |
| Evidência da receita | Imagem e exemplo conferidos no quickstart oficial Apache em 2026-10-12; a conclusão do guia não é deteção de um broker instalado. |

Requer Docker disponível. O guia é manual: não instala nem inicia um contentor pelo plano. O exemplo usa a porta 9092 só no localhost e dados temporários.

#### Guia manual

1. Tenha Docker disponível e confirme docker version no terminal. Se necessário, escolha uma das opções Docker do catálogo; não instale os dois pacotes concorrentes.
2. Exemplo temporário para testes, acessível apenas neste computador:
```sh
docker run --rm --name meu-ubuntu-kafka-dev -p 127.0.0.1:9092:9092 apache/kafka:4.3.1
```
Mantenha este terminal aberto. O exemplo não configura persistência; Ctrl+C encerra e remove este contentor.
3. Ligue a sua aplicação a localhost:9092. Siga o quickstart oficial para criar tópicos, produzir e consumir mensagens. Para dados persistentes ou uso em equipa, prepare uma configuração própria antes de utilizar este exemplo.

[Documentação do fornecedor](https://kafka.apache.org/quickstart/)

<details>
<summary>Receita JSON integral</summary>

```json
{
  "id": "kafka-local-guide",
  "name": "Apache Kafka — laboratório Docker",
  "area": "dev",
  "category": "Serviços locais · Mensageria",
  "summary": "Guia para testar Kafka num contentor local com a imagem oficial.",
  "why": "Testar aplicações baseadas em eventos sem preparar um cluster de produção.",
  "method": "manual",
  "source": "Imagem oficial apache/kafka e documentação Apache",
  "notes": "Requer Docker disponível. O guia é manual: não instala nem inicia um contentor pelo plano. O exemplo usa a porta 9092 só no localhost e dados temporários.",
  "manual": {
    "url": "https://kafka.apache.org/quickstart/",
    "steps": [
      "Tenha Docker disponível e confirme docker version no terminal. Se necessário, escolha uma das opções Docker do catálogo; não instale os dois pacotes concorrentes.",
      "Exemplo temporário para testes, acessível apenas neste computador:\n```sh\ndocker run --rm --name meu-ubuntu-kafka-dev -p 127.0.0.1:9092:9092 apache/kafka:4.3.1\n```\nMantenha este terminal aberto. O exemplo não configura persistência; Ctrl+C encerra e remove este contentor.",
      "Ligue a sua aplicação a localhost:9092. Siga o quickstart oficial para criar tópicos, produzir e consumir mensagens. Para dados persistentes ou uso em equipa, prepare uma configuração própria antes de utilizar este exemplo."
    ]
  },
  "verified": "Imagem e exemplo conferidos no quickstart oficial Apache em 2026-10-12; a conclusão do guia não é deteção de um broker instalado."
}
```

</details>
