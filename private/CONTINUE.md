# Retoma: pacote privado OfficeGest

Pausa pedida pelo utilizador em 2026-10-12, ao terminar o expediente.

- RC.1 preservada; instalador pequeno pronto em `dist/meu-ubuntu-1.0.0-rc.1-instalar.zip` (4,3 MB). Só .deb, checksum e instruções. Fontes disponíveis separadamente.
- Novo suporte genérico em `internal/integrations/packages.go`, `app_integrations.go`, `frontend/src/Integrations.tsx`, com alterações em app.go/api.ts/App.tsx.
- Pacotes são dados: ZIP com apenas manifest.json; verificações têm tipos e argv compilados. Os comandos dos guias só são copiados, nunca executados pelo pacote/helper.
- Manifesto privado em `private/officegest/manifest.json`; preparar ZIP depois dos checks. Não contém código/segredos do ogdevel.
- Implementação ainda é rascunho: testar importação, ZIP inválido/traversal, pacote alterado após revisão, persistência e simulação sem chamadas Docker; conferir janela nativa e seleção do projeto.
- Não foram modificados ficheiros nem iniciados serviços em `../ogdevel`.
- Compose local reportou v2.3.3; o diagnóstico sinaliza requisito 2.20+. Automação e configuração isolada de containers/volumes ficam para a próxima etapa, após validar o diagnóstico/guia.
- Antes de gerar novos artefactos, mudar a versão para 1.0.0-rc.2; manter intactos os ficheiros RC.1 já enviados.
- O diretório `.git` da raiz está vazio e é somente de leitura. Metadados dos commits públicos estão em `/tmp/meu-ubuntu-public-commits.git`, usando esta pasta como worktree explícita. Os privados ficam em `/tmp/meu-ubuntu-private-commits.git`, com worktree `private/`.
- Bundles duráveis: `dist/meu-ubuntu-public.git.bundle` e `private/officegest-private.git.bundle`. Não misturar o bundle privado com a distribuição pública.
- Ambiente temporário `meu-ubuntu-packs-build`: bootstrap Go 1.26.9, Node/GTK/WebKit em Ubuntu 26.04, cache Go reutilizada. Não confundir instalação de ferramentas de build com validação da aplicação. Pode ser descartado depois de guardar os resultados.

Próximo passo: completar testes e revisão da implementação, gerar RC.2 e o pacote privado, validar numa máquina Ubuntu de teste e depois ampliar a preparação automática.

Checks rápidos antes da pausa: Go formatado; `go test -tags webkit2_41 ./internal/integrations .` passou (o pacote novo ainda não tem testes próprios); TypeScript `tsc --noEmit` passou; JSON privado válido. Importação real, diagnóstico e interface nativa ainda não foram exercitados.

O instalador reduzido foi gerado separadamente; integrar essa variante no script de release ao preparar a RC.2. O ZIP completo da RC.1 continua disponível com as fontes, e o ZIP `-instalar` contém só o necessário para instalar.
