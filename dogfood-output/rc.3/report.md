# RC.3 — relatório de retoma, 2026-10-13

Candidata para testes: Meu Ubuntu 1.0.0-rc.3 e addon OfficeGest 0.2.0.
Checkout: /home/paulo/www/my-ubuntu, branch development, base 4c652f857ca77283edea5acfe07e692746b4f5ac,
com alterações locais sem commit. A 1.0 estável continua sem aprovação.

Binário ensaiado: SHA256 `7f9771a12d5c58345d6e27c60822cbbc487fe393c5e3f00411f8c0dfa0e4e7d7`. Artefactos RC.2 enviados foram
conferidos contra rc.2/rc2-sent-preservation.json e permanecem iguais.

## Alterações

- Navicat Premium Lite por referência Flatpak oficial incorporada; identificação,
  hash e confirmação dos termos revalidados. DBeaver continua disponível.
- 67 descrições ampliadas; utilidade visível em cada cartão.
- Progresso global e por item: etapa, contagem, tempo, mensagens e estados.
- Pausa/retoma entre operações, cancelamento seguro e recusa de estado alterado.
- Preparação declarativa revista: configuração Compose pessoal, recursos
  identificados e isolados, portas locais, schemas, Redis e Composer.
- Histórico da preparação e lock de kernel contra concorrência após crash.
- Containers de checks/build limitados por padrão a 2 CPUs e 4 GiB. Não foi
  diagnosticada a causa do bloqueio do PC; não foram reiniciados serviços da equipa.

## Verificação

| Verificação | Estado | Evidência e alcance |
|---|---|---|
| Documentação, frontend, vet e Go com race detection | PASS | check.log; 9 testes Node, 6 Python e todos os pacotes Go; Python usa fixtures sem Docker real |
| Auditoria Go/npm | PASS | audit.log; zero vulnerabilidades reportadas nesta execução |
| Build e verificação de binário/.deb/fontes/ZIPs | PASS | build.log e verify.log; versões e checksums conferidos |
| Atualização 0.3.2 → RC.3, sem Go/Node/npm | PASS em container | package-upgrade/package-upgrade.log; preferências, perfis, histórico e dados sintéticos preservados |
| Importação ZIP, revisão e preparação na janela nativa | PASS em simulação | native.log e screenshots; GTK/WebKit em X11/Xvfb, utilizador 1000 |
| Pausa, retoma e navegação durante preparação | PASS em simulação | screenshots/04-paused-other-page.png; estado real do job conferido independentemente em JSON, passos não avançam em pausa |
| Seleção da pasta no diálogo GTK | INCONCLUSIVE na RC.3 | Controlador não concluiu o seletor; pasta pré-carregada como fixture. Não atribuir esse bloqueio ao produto nem considerar esse fluxo aprovado |
| ZIP do colega | PASS | .deb, fontes, addon manifest-only, TESTAR.md e checksums internos conferidos |
| Arranque real OfficeGest completo | NOT RUN | Não foram executados scripts nem alterados ficheiros/serviços de ../guisoft/ogdevel |
| Instalação real Navicat | NOT RUN | Referência oficial e simulação conferidas; não foi instalado no host |
| GNOME/Wayland, polkit e operação real pela UI RC.3 | NOT RUN | Provas RC.2 não são aceitação destes bytes |

O ensaio nativo percorre o ZIP e confirmações reais, com uma pasta sintética
previamente configurada. A simulação não invoca Docker/Python/Composer, não
modifica esse projeto e não comprova SQL, Redis ou arranque do OfficeGest reais.
Os testes Python conferem configuração isolada, conservação de ficheiros,
SQL sem DROP, cluster sem RESET, validação, pausa e lock com subprocesses simulados.

## Utilização pelo colega

Enviar dist/meu-ubuntu-1.0.0-rc.3-officegest-teste.zip (cerca de 6,7 MB).
Instalar o .deb, importar officegest-0.2.0.zip e escolher um checkout sanitizado.
Rever ferramentas, aplicar o plano, depois rever e confirmar a preparação.

Dumps e credenciais privadas Composer ficam manuais. /etc/hosts e confiança
opcional do certificado exigem os comandos apresentados no resultado.
Configurações de tenants existentes são preservadas; devem ser conferidas antes
 dos dumps. Use o comando dock apresentado para gerir apenas este ambiente.
Não usar os scripts Docker originais para gerir os recursos isolados do addon.

A conclusão da preparação prova apenas os passos conferidos. Aceitação como
ambiente de desenvolvimento utilizável depende de teste real com os dumps e
acessos da equipa; esta entrega não é aprovação da 1.0 estável.

Referências: [Navicat Premium Lite](https://www.navicat.com/en/products/navicat-premium-lite)
e [Docker Compose](https://docs.docker.com/reference/cli/docker/compose/).
