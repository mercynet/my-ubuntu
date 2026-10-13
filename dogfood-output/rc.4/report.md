# RC.4 — Flatpak e organização de releases

Data: 2026-10-13. Versão: `1.0.0-rc.4`; Debian: `1.0.0~rc4`.
Árvore testada: alterações locais na branch `development`, sobre
`4c652f857ca77283edea5acfe07e692746b4f5ac` antes do commit de entrega.
O checkout limpo foi criado a partir do índice Git, sem `dist/`, frontend
compilado, dependências locais ou `private/`. Os resultados locais precedem
a CI do commit publicado; a origem dos artefactos está em `release-manifest.json`.

Executável testado, SHA-256:
`08ebb7784c5f92f1358f390f94c99b0ea8f1dab65bf2bfe3160de58b4f3854b3`.

## Resultados

| Verificação | Resultado e alcance | Evidência |
|---|---|---|
| Checkout limpo: versão/documentação, frontend, Go vet e testes com race detection | PASS; 9 testes Node, 6 Python e todas as packages Go | [check-clean.log](check-clean.log) |
| Auditorias Go/npm | PASS; zero vulnerabilidades reportadas nesta consulta | [audit.log](audit.log) |
| Build, binário, Debian, fontes e ZIPs | PASS em Docker Ubuntu 26.04; Go 1.26.9 | [build.log](build.log) |
| Versões, checksums, documentação/licenças, fontes e conteúdo dos ZIPs | PASS | [verify.log](verify.log) |
| Instalação e atualização 0.3.2 → RC.4 | PASS APT/CLI em container sem Go/Node/npm; preferências, perfis, histórico e dados XDG sintéticos preservados por hash | [package-upgrade.log](package/package-upgrade.log) |
| Requisito Flatpak e origem no cartão | PASS na janela GTK/WebKit, X11/Xvfb, simulação | [cartão](screenshots/01-flatpak-card.png) |
| Plano com aviso e dependência ausente | PASS nativo simulado; revisão inclui Flatpak + Flathub | [revisão](screenshots/02-flatpak-review.png) |
| Execução e persistência dos cinco passos | PASS nativo simulado; instalação Flatpak de VLC após preparação do Flathub | [resultado](screenshots/result.json), [captura](screenshots/03-flatpak-completed.png), [log](native-flatpak.log) |
| Recusa antes de instalar | PASS regressões: Flatpak ausente/quebrado; remoto ausente, incorreto, sem GPG ou desativado; ID errado ou consulta falhada | `internal/engine/flatpak_test.go` |
| Instalações existentes e indisponibilidade | PASS regressões: APT preservado, Thunderbird antigo/atual e Metadata Cleaner descontinuado | `internal/detect/flatpak_test.go` |
| Publicação das apps consultadas | PASS consulta à API; 18 novas receitas Flatpak com referência x86_64/stable ativa | [snapshot Flathub](flathub-inventory.json) |
| Artefactos RC.3 anteriormente preparados | Preservados por hash; RC.2 enviada também permanece local | [manifesto RC.3](rc3-preservation.json), [manifesto RC.2](../rc.2/rc2-sent-preservation.json) |

O controlador nativo teve duas falhas de automação antes do PASS: o input HTML
não expunha EditableText no AT-SPI, e os títulos técnicos estavam numa secção
fechada. O script passou a usar foco/teclado X11 e a conferir dependência visível
mais os passos persistidos. A [captura da tentativa](screenshots/harness-review-failure.png)
mostra o aviso e a revisão; não é apresentada como falha do produto.

## Entrega e organização

- `VERSION` centraliza a versão; manifests frontend são sincronizados e verificados.
- README resumido, índice técnico, changelog, `RELEASING.md` e templates de contribuição.
- Exemplo OfficeGest sanitizado em `examples/integrations/officegest/`.
- Artefactos, frontend compilado e estado local retirados do índice, mantidos no disco.
  O histórico anterior não foi reescrito; materiais antigos continuam públicos.
- Logs brutos preservam CR/progresso do terminal; `.gitattributes` limita a exceção
  de whitespace aos recibos `.log` de QA.
- A CI prepara artefactos do checkout com nomes ligados à versão. Commit/push
  não cria automaticamente uma tag ou GitHub Release.

## Limites e trabalho restante

**1.0 estável: ainda não aprovada.** Consulta ao Flathub e executor simulado não
comprovam instalação real, permissões ou integração com ferramentas do host.
Restam GNOME/Wayland, polkit nos bytes finais, falhas/interrupção/retoma reais,
percursos completos de perfis/exportações/teclado/temas, receitas de hardware e beta.
Consulte a [matriz de aceitação](../../docs/release-1.0.md).

OfficeGest completo com checkout, builds, schemas, Redis, Composer, autenticação
e dumps reais: **NOT RUN** nesta candidata. Navicat real e novas apps Flathub:
**NOT RUN**. Solaar, VS Code, virt-manager e Wireshark expõem os limites e guias
de integração com o host; não concedem permissões automaticamente.
