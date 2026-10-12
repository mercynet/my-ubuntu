# Escolhas e aplicações instaladas — 0.6.0-beta.3

Data: 12 de outubro de 2026. Alvo: Ubuntu 26.04 LTS, GNOME, amd64.

## Dúvida reproduzida

Um cartão com **Já instalado** podia estar desmarcado e indicar apenas **Por escolher**, sugerindo que aplicar o perfil retiraria essa aplicação. A caixa representa inclusão nas escolhas/perfil; o selo representa o estado detetado no computador. Desmarcar não cria uma operação de desinstalação ou reversão.

## Alterações

- Cartões passam a distinguir **Incluído nas escolhas** de **Fora das escolhas**. Um item instalado e desmarcado indica **Será mantido como está, mesmo sem marca**.
- O significado da caixa aparece no catálogo, inclusive a 900 × 600, e na revisão do plano. Tooltips e nomes acessíveis dizem incluir/retirar **das escolhas**.
- APT recebe `--no-remove` nas instalações do catálogo. Se resolver dependências exigir retirar outro pacote, a instalação é bloqueada; o erro apresentado explica a preservação dos pacotes. A [documentação APT](https://manpages.ubuntu.com/manpages/noble/man8/apt-get.8.html) descreve esse comportamento, confirmado também com APT no Ubuntu 26.04.
- Rever escolhas já aplicadas mantém a interface utilizável. O teste nativo encontrou uma janela vazia quando o plano tinha zero passos: a resposta serializava `steps: null`. O backend agora devolve `steps: []`; o plano mostra zero alterações e o botão Aplicar fica desativado.

Os controlos de reversão continuam a ser ações explícitas, revistas separadamente. Instalações detetadas não são incluídas automaticamente num perfil: **Mais ações do catálogo → Adicionar os já instalados nesta vista** permite fazê-lo por escolha do utilizador.

## Testes

Passaram testes frontend, TypeScript/Vite, `go vet` e a suite Go com `-race`. A regressão adicional prepara Déjà Dup, OpenVPN, rsync, SMB e FreeFileSync instalados, deixa os dois primeiros fora das escolhas e aplica apenas o trabalho em falta. Verifica que todos permanecem instalados, que repetir as escolhas gera zero passos e que retirar um ajuste das escolhas não o reverte.

O teste `TestRealAPTRefusesImplicitPackageRemoval`, com `realsystem` e `MEU_UBUNTU_DISPOSABLE=1`, correu como root apenas num container Ubuntu 26.04 descartável. Instalou um pacote fictício e preparou outro que conflita com ele. A simulação APT sem proteção propôs remover o primeiro; a execução através do helper com proteção abortou, preservou o pacote instalado e deixou o candidato ausente. Não foram usados pacotes nem configurações do computador do utilizador.

Um teste de contrato confirma que a revisão de um item instalado devolve uma lista vazia de passos. A classificação da recusa APT também tem cobertura de regressão.

## Percurso nativo

Binário Wails/WebKitGTK, UID 1000, Ubuntu 26.04 com Xvfb, modo de simulação e dados isolados:

1. Instalar Déjà Dup, OpenVPN e SMB no executor simulado: [plano inicial](screenshots/01-plano-inicial.png).
2. Ver os três itens instalados e incluídos: [escolhas](screenshots/03-escolhas-instaladas.png).
3. Desmarcar Déjà Dup e OpenVPN: [continuam instalados, com mensagem de preservação](screenshots/04-instalados-fora-das-escolhas.png).
4. Rever apenas SMB, já instalado: [zero instalações, configurações e reversões](screenshots/05-reaplicar-sem-alteracoes.png).
5. Conferir o cartão desmarcado de OpenVPN a 900 × 600: [mensagem visível](screenshots/06-mantido-900x600.png).
6. Abrir um plano só com Git já instalado e conferir o resultado explícito: [Não é necessário aplicar novamente](screenshots/07-tudo-pronto.png).

Este percurso valida a interface e o backend simulado. O teste APT usa o gestor real de pacotes em isolamento. Não substituem sessões completas com participantes em GNOME/Wayland.
