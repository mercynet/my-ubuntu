# Correção do diagnóstico de arranque — 0.6.0-beta.2

Data: 12 de outubro de 2026. Alvo: Ubuntu 26.04 LTS, GNOME, amd64.

## Problema confirmado

O diagnóstico mostrava «Skype for Linux (utilizador)» porque encontrava `~/.config/autostart/skypeforlinux.desktop`. A entrada referencia `Exec=/usr/bin/skypeforlinux` e está configurada como ativada. O executável referido não existe neste computador, não está no PATH e o dpkg não encontra o pacote `skypeforlinux`. A origem histórica do ficheiro não foi determinada; a sua presença não comprova uma instalação atual.

A versão anterior lia apenas o nome dos ficheiros de arranque e apresentava-os como aplicações. Não verificava o executável, ignorava a desativação específica do GNOME e não respeitava a precedência de uma entrada pessoal sobre outra do sistema com o mesmo nome.

## Comportamento corrigido

O painel passa a chamar-se **Entradas de arranque**. Cada entrada indica **Comando disponível**, **Executável não encontrado**, **Desativada nesta sessão** ou **Por verificar**. **Ver origem da entrada** mostra o ficheiro e o executável referido, sem executar comandos de arranque nem alterar esses ficheiros.

A leitura usa os diretórios XDG, respeita a precedência das entradas pessoais, interpreta apenas `[Desktop Entry]` e considera `Hidden`, `X-GNOME-Autostart-enabled`, restrições de sessão e `TryExec`. `NoDisplay` regula visibilidade nos menus, não desativação do arranque. Ativação D-Bus, lançadores indiretos e condições de arranque não avaliadas ficam explicitamente por verificar. A disponibilidade do executável não é apresentada como prova de que a aplicação esteja instalada por um gestor de pacotes ou esteja em execução.

Referências: [especificação Autostart](https://specifications.freedesktop.org/autostart/latest/), [chaves Desktop Entry](https://specifications.freedesktop.org/desktop-entry/latest/recognized-keys.html), [sintaxe Exec](https://specifications.freedesktop.org/desktop-entry/latest/exec-variables.html).

## Validação

- A lógica corrigida foi executada em modo de leitura neste computador: a entrada real «Skype for Linux» devolveu `status: missing` e identificou `/usr/bin/skypeforlinux` como alvo ausente.
- Testes Go com `-race` cobrem executável ausente, permissões, caminhos com espaços, entradas desativadas, precedência pessoal, ambiente GNOME, `TryExec`, D-Bus, lançadores indiretos, condições e secções de ações. Um executável de teste deixaria uma marca se fosse iniciado; o diagnóstico localizou-o sem o executar.
- Passaram testes frontend, TypeScript/Vite, `go vet` e a suite Go completa com `-race`.
- O binário nativo Wails/WebKitGTK foi aberto em Ubuntu 26.04 com Xvfb como UID 1000, em modo real, com configuração de teste isolada. Três ficheiros de teste demonstram os estados disponível, desativado e ausente. [Estados no diagnóstico](screenshots/02-estados-de-arranque.png), [origem da referência Skype](screenshots/03-origem-skype.png).
- Binário, pacote Debian `0.6.0~beta2`, manual e checksums foram conferidos. A beta anterior permanece em `dist/`.

O teste nativo cobre a leitura e a apresentação; não é uma sessão GNOME/Wayland completa. Não foi executado nenhum plano nem removida a entrada pessoal do utilizador.
