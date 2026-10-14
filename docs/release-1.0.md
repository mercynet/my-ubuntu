# Aceitação da versão 1.0

Estado atual: **1.0.0 aprovada para o escopo definido abaixo**, Ubuntu 26.04
LTS, GNOME, amd64. A decisão considera a utilização sem erros em dois
computadores reais e os checks/builds das candidatas RC.1–RC.4. Os cenários
fora do escopo continuam explicitamente limitados; não são promessas da 1.0.
As provas RC.2/RC.3 abaixo são históricas e não substituem a verificação dos
artefactos finais 1.0.0.

## Escopo da versão 1.0.0

| Trabalho | Estado atual / critério |
|---|---|
| Checks, auditoria, build e integridade | PASS nas RC; repetir para os bytes 1.0.0 antes da publicação |
| Uso real no ambiente suportado | PASS informal: utilização sem erros em dois computadores reais |
| APT, Flatpak, revisão, execução e histórico | Incluídos no escopo; receitas sem prova suficiente permanecem manuais/indisponíveis |
| Perfis, diagnóstico, exportações e reversão | Incluídos no escopo; devem ser reportados como bugs se falharem no ambiente suportado |
| Outras distribuições e arquiteturas | Fora do escopo: apenas Ubuntu 26.04 LTS, GNOME e amd64 |
| Hardware, dispositivos e permissões externas | Fora da garantia geral; dependem de máquina e sessão e podem exigir guia manual |
| OfficeGest completo | Fora da distribuição e com aceitação própria; não é preparado automaticamente |
| Publicação | Criar tag `v1.0.0`, gerar artefactos finais e publicar a release |

Consulta à API Flathub confirma publicação/ID e fonte, não instalação, licença
aceite, funcionamento ou permissões no host. Testes simulados também não
comprovam GNOME/Wayland, polkit nem a preparação real do OfficeGest.

O addon OfficeGest é opcional, distribuído separadamente e tem aceitação própria:
checkout sanitizado, build, containers, schemas, Redis, Composer e dumps da equipa.
Dumps e operações administrativas continuam manuais. A sua prontidão não deve
ser apresentada como concluída quando só existe prova simulada.

Esta é uma versão estável de escopo controlado. A promoção não implica suporte
a sistemas, hardware ou integrações que não estejam explicitamente incluídos.

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

## Publicar a versão estável

Manter a prova vinculada aos mesmos bytes da versão final. Executar checks,
auditoria, build e verificação, conferir o pacote final e publicar tag, fontes,
binário, `.deb` e checksums. Os ficheiros versionados em `private/` são
públicos; o pacote OfficeGest específico não acompanha a distribuição.

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
