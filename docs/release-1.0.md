# Aceitação da versão 1.0

Estado atual: **1.0 estável ainda não aprovada**. Candidata em preparação:
**1.0.0-rc.4**, Ubuntu 26.04 LTS, GNOME, amd64. Inclui política Flatpak e
organização de releases, além do progresso e integrações da RC.3.
As provas RC.2/RC.3 abaixo são históricas; não validam os bytes RC.4.

## O que falta para a estável

| Trabalho | Estado atual / critério |
|---|---|
| Checks, auditoria, build e integridade RC.4 | PASS local; [relatório](../dogfood-output/rc.4/report.md). Conferir CI do commit publicado |
| Aviso e plano Flatpak; dependência ausente | PASS nativo GTK/WebKit em X11/Xvfb, executor simulado; não instala apps reais |
| Pacote e atualização 0.3.2 → RC.4 | PASS APT/CLI em container, dados XDG sintéticos preservados; interface ainda pendente |
| Instalação/atualização num desktop GNOME/Wayland | Conferir o pacote final e preferências/perfis/histórico pela interface |
| Polkit, revisão, aplicação e reversão | Exercitar autorização recusada/aceite, estado alterado e resultado real nos bytes finais |
| Falha de rede, pausa/cancelamento/fecho e retoma | Comprovar limite seguro e estado parcial persistido num desktop descartável |
| Perfis, exportações, teclado, foco, temas e janela pequena | Concluir os percursos e conferir resultados nas aplicações leitoras |
| Novas receitas Flatpak e integrações do host | Instalar realmente apps representativas; validar SDKs, dispositivos, KVM e captura onde aplicável |
| Receitas de hardware/sessão | Prova adequada ou manter manuais/indisponíveis, com limites explícitos |
| Beta com utilizadores | Uso diário por pelo menos dois utilizadores; cinco participantes é meta proposta |
| Publicação | Preparar versão 1.0.0, tag e artefactos da revisão aprovada; seguir RELEASING.md |

Consulta à API Flathub confirma publicação/ID e fonte, não instalação, licença
aceite, funcionamento ou permissões no host. Testes simulados também não
comprovam GNOME/Wayland, polkit nem a preparação real do OfficeGest.

O addon OfficeGest é opcional, distribuído separadamente e tem aceitação própria:
checkout sanitizado, build, containers, schemas, Redis, Composer e dumps da equipa.
Dumps e operações administrativas continuam manuais. A sua prontidão não deve
ser apresentada como concluída quando só existe prova simulada.

Candidatas podem ser disponibilizadas como prerelease com esses limites. A
promoção estável depende da matriz concluída, não apenas de commit/push.

---

# Provas históricas RC.2 e RC.3

Estado em 2026-10-13: **1.0 estável ainda não aprovada**. Candidata preparada: **1.0.0-rc.2**, Ubuntu 26.04 LTS, GNOME, amd64. Os resultados abaixo pertencem à árvore de trabalho local sobre `4c652f8`; a CI publicada desse commit precede estas alterações.

| Critério | Estado | Evidência / trabalho restante |
|---|---|---|
| Documentação, frontend, Go vet e testes Go com race detection | PASS | `dogfood-output/rc.2/check-final.log`; nove testes Node e regressões de importação/revisão/persistência/plano vazio |
| Auditoria Go/npm | PASS | `dogfood-output/rc.2/audit.log`; zero vulnerabilidades reportadas nesta execução |
| Binário, .deb, fontes, ZIP completo e instalador reduzido | PASS | Build e verificação locais; o instalador tem apenas .deb, checksum próprio e instruções |
| Instalação sem Go/Node/npm e atualização 0.3.2 → RC.2 | PASS em container | CLI, bibliotecas, avisos instalados e dados XDG sintéticos preservados; não comprova migração pela interface |
| APT recusa remoção implícita de outro pacote | PASS real em container | `TestRealAPTRefusesImplicitPackageRemoval`, com fixtures locais e `MEU_UBUNTU_DISPOSABLE=1` |
| Importação e diagnóstico de integração na janela GTK | PASS em simulação | ZIP, revisão, pasta, diagnóstico, cópia do comando e persistência após reabrir; X11/Xvfb |
| Janela GNOME/Wayland e repetição de plano concluído | PASS em VM | Ubuntu 26.04.1, GNOME 50.1, sessão Wayland; corrigidas janela vazia por EGL e lista nula de passos |
| Polkit e instalação real pela interface | PASS em VM | Cancelamento registado como falha; autorização instalou git-flow, conferido por dpkg; ensaio no build anterior à correção do plano vazio |
| Preferência, revisão desatualizada e reversão | PASS em VM | Valor real de animações `true → false → true`; plano alterado externamente recusado sem criar execução |
| Falha de rede e cancelamento/fecho durante operação; retoma | NOT RUN no desktop | Existem regressões Go; ainda falta testemunho pela interface com estado parcial persistido |
| Exportação PDF/CSV/Excel e cancelamento | PASS em VM | CSV/Excel com 148 linhas iguais, sem fórmulas; diálogo cancelado preservou hashes; PDF final com cinco páginas e cabeçalhos conferidos |
| Interação desktop completa | PARCIAL | Sobre/janela 900 × 600 em GTK/Xvfb; teclado/foco e temas ainda exigem QA completo |
| Addon OfficeGest: importação e diagnóstico | PASS em VM | Pacote 0.1.0 instalado pela interface, separador e diagnóstico reais; pasta sintética, sem serviços OfficeGest |
| Addon OfficeGest: preparar ambiente em poucos cliques | NÃO IMPLEMENTADO | Novo objetivo: containers, bases vazias, Redis e dependências; dump manual. O pacote 0.1.0 inclui apenas diagnóstico/guias |
| Hardware e receitas dependentes de sessão/dispositivos | NOT RUN | Extensões, energia, Firefox e drivers/TRIM exigem prova adequada; receitas sem prova devem ficar manuais/indisponíveis |
| Beta com utilizadores | NOT RUN | Uso diário por pelo menos dois utilizadores; cinco participantes é meta proposta |
| Repositório e canais de contribuição/segurança | PASS | Repositório público; relatos privados de vulnerabilidades ativados e documentados |
| GitHub Release 1.0.0 | NOT RUN | Publicar após fechar os critérios de aceitação |

## Matriz desktop

Registar versão do executável, Ubuntu/GNOME, sessão, snapshot da VM, passos, esperado, observado e evidência de cada cenário:

1. Instalar a candidata num desktop limpo e abrir como utilizador normal.
2. Atualizar a 0.3.2 com escolhas, perfil, histórico e progresso de guias existentes; conferir na interface após reabrir.
3. **PASS em VM:** cancelar polkit sem sucesso indevido; repetir autorizando, conferir git-flow por dpkg. Repetir estes cenários com os bytes finais da candidata antes de promover.
4. **Repetição PASS:** plano concluído apresenta zero passos e botão de execução desativado. Exportação/importação de perfil com opções ainda pendente.
5. Provocar falha de rede, pedir cancelamento/fecho durante operação, reabrir e retomar; histórico deve explicar o resultado parcial.
6. **PASS em VM:** aplicar e reverter animações, conferindo por `gsettings`; alterar externamente depois de rever e confirmar recusa sem nova execução. O ensaio externo usou CLI, não a janela das Definições.
7. **Exportações PASS:** PDF/CSV/XLSX reais, equivalência CSV/Excel e cancelamento do diálogo. PDF final renderizado e conferido. Teclado, foco e temas ainda pendentes; janela 900 × 600 aprovada em GTK/Xvfb.

## Promover a candidata

Corrigir e repetir os cenários que falharem; manter a prova vinculada aos mesmos bytes do candidato. Depois atualizar versão/notas para `1.0.0`, executar checks, auditoria, build e verificação, conferir instalação/atualização do pacote final e publicar tag, fontes, binário, .deb e checksums. Os ficheiros versionados em `private/` são públicos; o pacote OfficeGest específico não acompanha a distribuição da aplicação.

## Candidata RC.3 (2026-10-13)

RC.3 acrescenta Navicat Premium Lite, descrições, progresso e pausa/retoma,
mais a preparação declarativa de ambientes. RC.2 enviada continua preservada;
a prova desktop da RC.2 não é aceitação dos bytes RC.3.

- PASS: checks completos, testes Go com race detection, frontend e seis
  regressões Python sintéticas; dogfood-output/rc.3/check.log.
- PASS: auditorias Go/npm, zero vulnerabilidades reportadas;
  dogfood-output/rc.3/audit.log.
- As regressões conferem pausa/retoma/cancelamento, recusa de alterações,
  revisão/concorrência, persistência, referências Flatpak e simulação sem
  alterar a pasta escolhida.
- A preparação OfficeGest exige validação completa numa máquina descartável
  com checkout sanitizado, autenticação privada e dumps da equipa. Não foi
  executada sobre o ambiente real de desenvolvimento existente.
- Build, verificação e ensaio nativo RC.3 devem ser consultados no relatório
  dogfood-output/rc.3/report.md. X11/Xvfb simulado não é prova GNOME/Wayland,
  instalação real Navicat ou ambiente OfficeGest pronto para trabalhar.

A versão estável 1.0 continua pendente dos critérios desktop/beta acima.
