# Preparar e publicar versões

O desenvolvimento acontece em `development`. Tags usam `v<versão>`:
`v1.0.0-rc.4` para uma candidata e `v1.0.0` para a estável. A versão Debian da
candidata usa `~rc`, ordenada antes da estável. `VERSION` é a fonte de versão;
os manifests frontend são sincronizados e conferidos pelo build.

## Preparar a revisão

1. Atualize a versão com `python3 build/scripts/set-version.py <versão>`.
2. Acrescente a entrada datada ao `CHANGELOG.md`, agrupando mudanças e correções.
   Preserve versões históricas e identifique limitações da candidata.
3. Atualize instruções atuais, manual e catálogo; execute `make docs`.
4. Execute `make check`, `make audit`, `make build` e `make verify`.
5. Registe os resultados em `dogfood-output/<candidata>/report.md`, com checksum
   do executável e a distinção entre simulação, container e desktop real.
6. Faça commit e push. Confira que a CI do mesmo SHA passou; publique artefactos
   construídos dessa revisão. Mudanças depois do build exigem conferir a origem
   e os checksums novamente.

O build usa Docker Ubuntu 26.04, Go com versão/checksum fixos e dependências
lockadas. Binários, `.deb`, fontes e ZIPs ficam em `dist/`, ignorados pelo Git.
A CI constrói os artefactos a partir do checkout e não cria uma Release automaticamente.
`SOURCE_DATE_EPOCH` fixa timestamps; sozinho não comprova builds idênticos.

## Critérios de aceitação

Uma candidata pode ser publicada como **prerelease**, com limites conhecidos.
Para uma versão estável, documente o [escopo](docs/release-1.0.md) nos bytes finais:
instalação/atualização, polkit, revisão/reversão, falhas, pausa/cancelamento/fecho,
retoma, perfis, exportações, teclado, temas e receitas dependentes do sistema.
Registe o beta com utilizadores. Receitas sem prova suficiente devem ficar
manuais/indisponíveis ou ter o limite exposto; zero falhas no teste Go não prova
aceitação administrativa ou desktop.

O pacote de integração OfficeGest é opcional e tem validação própria. Não
inclua dumps, configurações privadas ou o checkout do produto na Release.

## Publicar uma revisão aprovada

1. Trabalhe a partir do commit aprovado e mantenha a árvore limpa para o build.
2. Crie a tag anotada da versão e envie-a ao remoto.
3. Aguarde a CI da tag e confira SHA, versão e checksums dos artefactos.
4. Crie a GitHub Release ligada à tag; marque candidatas como prerelease.
5. Use a entrada do changelog como notas e inclua alvo, instalação, validação e limites.
6. Anexe binário, `.deb`, arquivo de fontes, ZIP completo, instalador reduzido,
   `SHA256SUMS` e `release-manifest.json`. Confira os downloads publicados.

Publicação da Release e promoção da estável são ações explícitas, separadas de
commit/push. Versões candidatas anteriores já distribuídas não são sobrescritas.

## Integrações para teste

Depois do build, gere o pacote e bundle opcionais com:

```sh
python3 build/scripts/integration-test-bundle.py \
  --manifest examples/integrations/officegest/manifest.json
```

O ZIP do addon contém apenas manifesto declarativo. O bundle de teste inclui
instalador, fontes, addon, instruções e checksums, sem dados privados.
