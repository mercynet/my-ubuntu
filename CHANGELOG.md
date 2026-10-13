# Changelog

Alterações relevantes para utilizadores e contribuidores. Candidatas não equivalem
à versão estável; alcance e limites de validação acompanham cada relatório.
Versões anteriores à RC.1 estão documentadas no [contexto histórico](docs/contexto-projeto.md)
e nos relatórios de QA, sem reconstruir notas de release não verificadas.

## Unreleased

## 1.0.0-rc.4 — 2026-10-13

### Alterado

- Prefere Flatpak em 18 aplicações com publicação ativa confirmada no Flathub.
- Preserva aplicações APT e IDs Flatpak anteriores; migração de dados não é automática.
- Atualiza Thunderbird para o ID ativo org.mozilla.thunderbird.
- Distingue cliente virt-manager Flatpak e motor QEMU/libvirt do host.
- Organiza README, índice técnico, processo de releases e exemplo público OfficeGest.
- Centraliza versão em VERSION e confere os manifests frontend no build.
- Retira artefactos gerados e estado local do índice Git, preservando-os no disco.

### Adicionado

- Aviso Flatpak nos cartões e no plano, com dependência incluída quando falta.
- Validação do executável, origem, assinaturas, estado do remoto e ID antes de instalar.
- Templates de bugs/melhorias/PRs, changelog de candidatas e artefactos identificados na CI.
- Guias para integrações do host exigidas por alguns Flatpaks.

### Corrigido

- Recusa origens com o nome Flathub que apontam para outro repositório.
- Falha de preparação Flathub impede instalar aplicações dependentes.
- Novas instalações Metadata Cleaner ficam indisponíveis porque foi descontinuado.

Candidata para Ubuntu 26.04, GNOME, amd64. Instalações Flatpak reais e desktop final
continuam sujeitos à matriz de aceitação; consulta ao catálogo não comprova execução.

## 1.0.0-rc.3 — 2026-10-13

- Acrescenta Navicat Premium Lite por referência Flatpak oficial e confirmação dos termos.
- Amplia 67 descrições e mostra a utilidade nos cartões.
- Mostra etapa, contagem e tempo durante execução, mesmo ao navegar.
- Permite pausa/retoma entre operações e recusa retoma se o estado mudou.
- Acrescenta preparação declarativa revista de ambientes Compose, histórico e controlos.
- OfficeGest 0.2.0 prepara configuração isolada, schemas, Redis e Composer; dumps e root manuais.
- Limita CPU/RAM dos containers de build; preserva artefactos RC.2 enviados.

## 1.0.0-rc.2 — 2026-10-13

- Composição WebKit por software para evitar janelas vazias em desktops virtuais sem contexto EGL utilizável.
- Planos já concluídos mantêm uma lista de passos vazia válida, evitando a perda da interface ao repetir a revisão.
- Cabeçalho PDF posicionado explicitamente também nas páginas criadas por quebras automáticas.
- Suporte opcional a pacotes de integração com diagnóstico local e guias manuais, sem instalar ou iniciar projetos automaticamente.
- Regressões para importação ZIP, caminhos inseguros, alteração após revisão, expiração, persistência e isolamento da simulação.
- Instalador reduzido gerado e verificado pelo build, com o .deb, checksum próprio e instruções; fontes continuam disponíveis separadamente.
- Repositório público e canal privado de vulnerabilidades anunciados na documentação.

Candidata para completar testes de atualização, GNOME/Wayland, polkit, recuperação e beta antes da 1.0 estável. O pacote específico OfficeGest continua fora da distribuição da aplicação.

## 1.0.0-rc.1 — 2026-10-12

- Exportação do diagnóstico em PDF, CSV e Excel (.xlsx), com seleção de secções, data e versão da fotografia apresentada.
- Recuperação de reversões parciais sem repetir chaves restauradas; planos vinculados aos valores atuais e alvos, incluindo snapshots administrativos.
- Nova tentativa de configuração/verificação sem reinstalar pacotes já presentes; falhas de verificação persistem como intervenção no histórico.
- Correção da deteção de grupos Docker, preferências Git ausentes, Xdebug, journald e drivers recomendados.
- Fecho normal protegido durante execução, com cancelamento no limite seguro.
- GPLv3, avisos de terceiros, guias de contribuição/segurança e arquivo do código-fonte correspondente.

Candidata para Ubuntu 26.04 LTS, GNOME e amd64. Testes GNOME/Wayland, polkit, atualização e feedback externo ainda precedem a 1.0 estável. Veja a [revisão adversarial](dogfood-output/rc.1/review-astra-high.md).
