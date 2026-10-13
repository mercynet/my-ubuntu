# Contribuir para o Meu Ubuntu

O alvo da 1.0 é Ubuntu 26.04 LTS, GNOME e amd64. Leia [AGENTS.md](AGENTS.md) e [desenvolvimento](docs/desenvolvimento.md) antes de alterar receitas ou execução.

## Preparar e verificar

Instale Docker e Make. Execute `make check`, `make audit`, `make build` e `make verify`. O ambiente de build usa Ubuntu 26.04; Go e Node não são necessários no host. `make sim` abre o executor simulado. Nunca execute testes administrativos reais no computador principal; use uma VM/container descartável.

Formate Go com `gofmt`. Siga o estilo TypeScript existente. Edite o manual em `docs/manual.json` e regenere com `python3 build/scripts/generate-docs.py`.

## Propor alterações

Relate problemas com versão, ambiente, passos, resultado esperado/observado e evidência sem dados pessoais. Inclua regressões para mudanças em seleção, revisão, execução, verificação ou reversão. Receitas precisam de fonte oficial, deteção, verificação, dependências e riscos concretos.

Preserve IDs de receitas/perfis, a revisão obrigatória e as listas de operações permitidas do helper. Perfis nunca transportam comandos livres nem segredos. PRs devem explicar comportamento, validação e limitações; inclua capturas para alterações de interface.

O código usa GPL-3.0-only. Contribuições ao código devem ser compatíveis; assets e dependências mantêm os seus avisos. Use o [repositório público](https://github.com/mercynet/my-ubuntu) para PRs e issues; vulnerabilidades seguem o canal privado de SECURITY.md. A pasta `private/` presente no histórico é pública: o nome e o `.gitignore` não restringem acesso. Nunca coloque aí segredos ou novos materiais internos.

## Organização e versões

Use development para trabalho integrado e subjects de commit curtos no imperativo.
Adicione uma entrada ao changelog para alterações visíveis. Tags e publicação
seguem [RELEASING.md](RELEASING.md); VERSION e manifests precisam corresponder.
Binários, frontend compilado e caches são gerados, não versionados. Preserve
relatórios/capturas revistos em dogfood-output; deixe estados XDG e credenciais locais.
O [exemplo OfficeGest](examples/integrations/officegest/README.md) é declarativo
público e sanitizado. Não inclua o checkout do produto, dumps ou configurações privadas.
