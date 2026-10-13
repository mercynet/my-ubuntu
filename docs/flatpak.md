# Política Flatpak

Aplicações desktop com publicação ativa confirmada preferem Flatpak por utilizador.
IDs de receitas/perfis permanecem estáveis. CLI, linguagens, motores e serviços
do host mantêm o método necessário ao seu funcionamento. Não existe pesquisa
automática que aceite comandos ou IDs arbitrários de um perfil.

## Antes da instalação

O cartão e o plano mostram o requisito. Flatpak + Flathub é incluído quando
falta; a autorização APT para o host segue a revisão normal. Antes de instalar,
o executor confirma executável/versão, remoto do utilizador, URL permitida e
verificação de assinaturas. O ID é consultado no Flathub; erro ou ID diferente
bloqueia a instalação. Remotos existentes incorretos não são sobrescritos.

Navicat usa uma referência oficial incorporada, ligada ao hash revisto, com
confirmação dos termos; um remoto existente desse fornecedor também é conferido.
Não há fallback silencioso para APT nem execução de comandos do manifesto.

## Apps migradas na RC.4

| Receita | ID Flatpak | Editor no Flathub |
|---|---|---|
| Déjà Dup (Cópias de segurança) | [org.gnome.DejaDup](https://flathub.org/apps/org.gnome.DejaDup) | Verificado |
| FileZilla | [org.filezillaproject.Filezilla](https://flathub.org/apps/org.filezillaproject.Filezilla) | Não verificado |
| Brave | [com.brave.Browser](https://flathub.org/apps/com.brave.Browser) | Verificado |
| VLC | [org.videolan.VLC](https://flathub.org/apps/org.videolan.VLC) | Não verificado |
| qBittorrent | [org.qbittorrent.qBittorrent](https://flathub.org/apps/org.qbittorrent.qBittorrent) | Verificado |
| Flameshot | [org.flameshot.Flameshot](https://flathub.org/apps/org.flameshot.Flameshot) | Verificado |
| Solaar | [io.github.pwr_solaar.solaar](https://flathub.org/apps/io.github.pwr_solaar.solaar) | Não verificado |
| LibreOffice | [org.libreoffice.LibreOffice](https://flathub.org/apps/org.libreoffice.LibreOffice) | Verificado |
| Extension Manager | [com.mattjakeman.ExtensionManager](https://flathub.org/apps/com.mattjakeman.ExtensionManager) | Verificado |
| dconf Editor | [ca.desrt.dconf-editor](https://flathub.org/apps/ca.desrt.dconf-editor) | Verificado |
| Visual Studio Code | [com.visualstudio.code](https://flathub.org/apps/com.visualstudio.code) | Não verificado |
| PhpStorm | [com.jetbrains.PhpStorm](https://flathub.org/apps/com.jetbrains.PhpStorm) | Não verificado |
| KeePassXC | [org.keepassxc.KeePassXC](https://flathub.org/apps/org.keepassxc.KeePassXC) | Verificado |
| Wireshark | [org.wireshark.Wireshark](https://flathub.org/apps/org.wireshark.Wireshark) | Não verificado |
| Remmina (RDP/VNC) | [org.remmina.Remmina](https://flathub.org/apps/org.remmina.Remmina) | Verificado |
| Virtualização (virt-manager) | [org.virt_manager.virt-manager](https://flathub.org/apps/org.virt_manager.virt-manager) | Verificado |
| GNOME Boxes | [org.gnome.Boxes](https://flathub.org/apps/org.gnome.Boxes) | Verificado |
| Tor Browser (launcher) | [org.torproject.torbrowser-launcher](https://flathub.org/apps/org.torproject.torbrowser-launcher) | Verificado |

Consulta à API pública em 2026-10-13: [snapshot](../dogfood-output/rc.4/flathub-inventory.json).
HTTP 200, ID e referência x86_64/stable ativa não comprovam instalação ou integração
com o desktop. O editor verificado indica controlo de publicação, não uma garantia
geral sobre comportamento ou permissões.

## Instalações e limites preservados

- Pacotes APT e IDs Flatpak anteriores reconhecidos não são removidos ou migrados.
- Thunderbird passa para org.mozilla.thunderbird; o ID anterior é conservado na deteção.
- Metadata Cleaner foi descontinuado; novas instalações ficam indisponíveis.
- Solaar exige regras udev do host, apresentadas como preparação manual.
- VS Code Flatpak é comunitário; terminal e SDKs do host precisam de configuração.
- virt-manager usa cliente Flatpak e QEMU/libvirt do host para manter KVM local.
- Wireshark não ganha permissões de captura automaticamente.
- Apps sem publicação ativa verificável, como Warp ou Hoppscotch nesta consulta,
  conservam a receita existente. Não trate um ID adivinhado ou erro de rede como
  prova de disponibilidade ou ausência universal.

Referências: [Flatpak CLI](https://docs.flatpak.org/en/latest/flatpak-command-reference.html),
[VS Code](https://github.com/flathub/com.visualstudio.code),
[Solaar](https://github.com/flathub/io.github.pwr_solaar.solaar),
[virt-manager](https://github.com/flathub/org.virt_manager.virt-manager),
[Metadata Cleaner](https://flathub.org/en/apps/fr.romainvigier.MetadataCleaner).
