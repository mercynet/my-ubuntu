# Meu Ubuntu 1.0.0-rc.2: validação

Data: 2026-10-13. Checkout: branch `development`, base `4c652f857ca77283edea5acfe07e692746b4f5ac`, com alterações locais. Alvo: Ubuntu 26.04 LTS, GNOME, amd64. Não houve publicação de release nem promoção à 1.0 estável.

## Resultados conferidos

- `make check`: documentação, nove testes Node, TypeScript/Vite, Go vet e suíte Go com `-race` aprovados. As regressões novas exercitam ZIP inválido, traversal, symlinks, limites de tamanho, tipo de verificação proibido, revisão alterada/expirada/consumida, persistência e simulação sem Docker nem execução de guias.
- `make audit`: Go/npm não reportaram vulnerabilidades nessa execução. Compilador de build: Go 1.26.9, arquivo oficial com checksum conferido, em Ubuntu 26.04.
- `make build` e `make verify`: binário, versão Debian `1.0.0~rc2`, fontes, ZIP completo e instalador reduzido conferidos. O instalador tem apenas .deb, checksum desse pacote e instruções.
- `TestRealAPTRefusesImplicitPackageRemoval`, num container descartável sem rede: instalação conflitante recusada pelo APT real, pacote preexistente preservado. Resultado: PASS, 6.52 s, exit 0.
- Instalação limpa e atualização 0.3.2 → RC.2 em dois containers Ubuntu 26.04 sem Go/Node/npm: CLI, bibliotecas, manuais e avisos presentes; hashes de preferências, perfil, histórico, log e dados sintéticos preservados. Os primeiros ensaios falharam no roteiro porque `/etc/os-release` substituiu `VERSION`; o roteiro agora preserva `RELEASE_VERSION`. Falhas originais mantidas em `first-attempt-failed.log`, ensaios corrigidos aprovados.
- Os cinco artefactos específicos da RC.1 permanecem byte a byte iguais ao manifesto `rc1-preservation.json`.
- Repositório público confirmado; relatos privados de vulnerabilidades ativados e documentados. O nome `private/` não restringe acesso aos ficheiros já publicados.

## Janela nativa em simulação

Container próprio, UID 1000, sem rede nem socket Docker, Wails/GTK/WebKit, X11/Xvfb e pastas XDG sintéticas. Foi importado um ZIP de teste por diálogo GTK real; manifest/versão/fornecedor revistos e confirmação exigida. Seleção de pasta, diagnóstico marcado como exemplo e cópia exata do comando para a área de transferência foram conferidos. Pacote e pasta persistiram após reabrir e permitiram novo diagnóstico. Sobre e versão ficaram acessíveis a 900 × 600. Capturas em `screenshots/`.

Isso comprova o fluxo nativo simulado. Uma tentativa de reiniciar o harness reutilizou um display Xvfb indisponível; foi corrigida com um novo display e repetida. Os ensaios GNOME/Wayland reais estão descritos abaixo.

## Desktop e aceitação restante

VM descartável, a partir da [imagem Ubuntu oficial](https://cloud-images.ubuntu.com/resolute/current/), SHA256 `8800651811af9a85465ad1d552add729947bb16488dddb4a9b5305a3d97332b2`. Ambiente conferido: Ubuntu 26.04.1, GNOME 50.1, sessão Wayland ativa, UID 1000, sem Go/Node/npm. A preparação inicial atingiu o limite de memória porque o disco estava em tmpfs; pacotes interrompidos e conflito GDM foram resolvidos como problemas do harness. A VM foi desligada ao terminar e os discos preservados em `/tmp/meu-ubuntu-rc2-disks` para liberar espaço no host.

Foram encontradas e corrigidas três falhas: janela vazia por composição EGL, revisão de plano concluído com `steps: null` e posição do cabeçalho PDF em quebras automáticas. O teste do plano vazio reproduziu ambas as variantes antes da correção. Checks completos e build foram repetidos depois das alterações.

- Atualização 0.3.2 → RC.2 na VM: perfil, histórico e log sintéticos mantiveram hashes; escolhas antigas continuam visíveis. Preferências mudaram legitimamente por deteção e escolhas durante os testes. Progresso de guias ainda não foi validado com dados significativos.
- Polkit: cancelamento terminou como falha e git-flow permaneceu ausente; repetição autorizada terminou com três passos concluídos e git-flow 1.12.3-3build1 presente por dpkg. Este ensaio precedeu as correções de plano vazio/PDF; repetir com os bytes finais antes da promoção.
- Plano concluído: zero passos e botão Aplicar desativado; interface íntegra. Preferência de animações aplicada e revertida, com `gsettings` independente confirmando `true → false → true`. Alteração externa depois da revisão foi recusada sem criar execução. Provas pertencem ao binário SHA256 `73641b079a4da9f18083d24800954ed1869472f64b9b88eb284b6dc2c4f2c688`, anterior apenas ao ajuste PDF.
- PDF/CSV/Excel: diálogos GTK reais. CSV/Excel têm 148 linhas iguais, sem fórmulas, e cancelamento manteve os hashes. PDF final renderizado em cinco páginas, cabeçalhos/versão e paginação conferidos; binário final SHA256 `25970a0e1ab16c5fe99cdb5b3fcbbf5c8361d7255ae341037ac086469afa76b1`.
- OfficeGest 0.1.0 importado pela interface do binário final; separador, pasta e diagnóstico reais conferidos com fixture sem Docker. Isso não comprova instalação dos serviços OfficeGest. O addon atual é diagnóstico/guia manual; preparação automática é um novo objetivo ainda não implementado.

Recibos e capturas em `wayland/`; logs dos checks finais acompanham esta pasta. Hardware, falha de rede/cancelamento durante operação/retoma no desktop e beta externo permanecem NOT RUN. A [matriz da 1.0](../../docs/release-1.0.md) mantém os limites da aceitação; a candidata pode ser distribuída para testes.
