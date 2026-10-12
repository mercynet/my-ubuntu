# Revisão adversarial — 1.0.0-rc.1

Data: 2026-10-12  
Revisor: Astra / high  
Estado: revisão final concluída; não foram identificados bloqueadores de código remanescentes nos cenários revistos. QA nativo continua pendente.

## Parecer final

As alterações corrigem os casos iniciais de reversão parcial GSettings, revisão de valores de GSettings, ponte legada sem token, grupos Docker e tratamento do código de saída Git. A recuperação de pós-operações, a persistência de intervenções e a barreira de fecho normal também foram revistas. As variantes e regressões descobertas durante esta revisão foram corrigidas e reavaliadas. Não identifiquei bloqueadores de código remanescentes nestes cenários; isto permite avançar para o QA nativo do candidato, sem equivaler a validação real de instalações ou a aprovação de uma publicação externa.

Os achados foram inicialmente confirmados por inspeção dos caminhos de código. Uma tentativa inicial de preparar reproduções por overlay não foi executada: o ambiente recusou `docker exec -i` com `permission denied ... docker.sock`. Posteriormente, o responsável pela implementação incorporou regressões dirigidas nas fontes; essas regressões participaram nos checks finais descritos abaixo. Nenhuma fonte foi alterada pelo revisor.

## Resultado das correções

| Achado desta revisão | Resultado final e evidência |
| --- | --- |
| 1. Recibo antigo de journald | Resolvido: a prova é invalidada na aplicação/reversão e só reposta após recarga bem-sucedida. `TestJournaldReapplyAfterRevertCannotReuseOldReceipt` cobre aplicar, reverter e falhar ao reaplicar. |
| 2. Verificação de drivers | Resolvido no código: compara pacotes recomendados com pacotes instalados e distingue erro de consulta. `TestDriverStateRequiresRecommendedPackageAndReportsQueryFailure` cobre recomendado ausente/presente e consulta falhada; não valida hardware real. |
| 3. Retry cancelado e ordem dos dependentes | Resolvido: estados não executados não limpam a intervenção anterior e a verificação é inserida antes dos dependentes. `TestCancelledVerificationRemainsPendingAndPrecedesDependants` cobre o cenário através de um registo cancelado guardado no store. |
| 4. Caminho de modelo revisto | Resolvido: a execução usa o diretório revisto e respeita o registo mais recente, incluindo quando já está revertido. `TestTemplateDirectoryChangeCannotDeleteUnreviewedOldTemplate` e `TestRevertedLatestTemplateJournalCannotFallBackToOlderActiveChange` cobrem as duas variantes. |
| 5. Extensão acrescentada após confirmação | Resolvido por inspeção: um nome sem a extensão esperada é recusado; o código não troca silenciosamente o destino escolhido no diálogo. O diálogo GTK real ainda não foi exercitado pelo revisor. |
| 6. Unicode astral no PDF | Resolvido: PDF representa esses runes como `[U+...]`; CSV/XLSX preservam o original. Fixture com 🚀 passou e os leitores independentes confirmaram os resultados. |
| 7. Reversão administrativa em lote | Resolvido para o cenário revisto: [admin_review.go](../../internal/engine/admin_review.go) liga IDs do journal, estado atual e alvo ao token da cadeia completa de snapshots root, UID e modo batch. O engine revalida e o helper confere o token antes da primeira mutação. `TestAdministrativeBatchReviewBindsPrivateSnapshotTarget` altera o alvo privado mantendo o ficheiro atual e confirma a recusa. |

As regressões encontram-se em [rc_regression_test.go](../../internal/engine/rc_regression_test.go), [rc_followup_test.go](../../internal/engine/rc_followup_test.go), [app_rc_test.go](../../app_rc_test.go) e [export_test.go](../../internal/diag/export_test.go).

Durante a conferência final do achado 4 surgiu uma segunda variante: criar um modelo (c1), apagá-lo externamente, recriá-lo pela aplicação (c2), reverter c2 no histórico e recriar manualmente o mesmo conteúdo. A proposta respeitava c2 como já revertido, mas a execução ainda podia procurar c1 ativo. A correção final deixa de saltar registos revertidos ao escolher o último registo do caminho; a regressão específica acima passou.

O problema original de reversão parcial de wallpapers tem cobertura em `TestPartialWallpaperUndoDoesNotResetAlreadyRestoredKey`. Xdebug e a recarga inicial de journald têm regressões de falha/retoma. A intervenção de verificação é guardada no run/log e pode ser repetida sem reinstalar. A chamada legada `App.Revert` recusa a mutação. `FinishReset` guarda conclusão e marcação de reversão numa só atualização do journal. A barreira de fecho está ligada a `OnBeforeClose`, partilha o mutex de início de execução e tem teste unitário; falta a validação do comportamento de desktop indicada nos limites.

## Histórico dos achados encontrados e corrigidos

As severidades e sequências seguintes documentam os defeitos da primeira passagem, não bloqueadores atualmente abertos.

1. **Alta — recibo antigo de journald pode esconder uma nova falha de recarga.**

   [helper/dropin.go](../../internal/helper/dropin.go), `applyDropin` e `revertDropin`; [detect.go](../../internal/detect/detect.go), `adminState`.

   Na implementação inicialmente revista, `.applied` só era escrito após sucesso, mas nunca invalidado antes de nova aplicação nem removido/atualizado na reversão. Sequência: aplicar valor A com sucesso; reverter A; voltar a aplicar A com falha de `systemctl restart`. O ficheiro novo tem A e o recibo antigo ainda contém o hash de A. A deteção passa a `Installed`, apesar da falha, e a retoma omite a recarga. A regressão inicial cobria apenas a primeira aplicação, sem recibo anterior.

   Correção mínima: invalidar a prova antes da alteração/recarga e atualizar ou remover a prova durante a reversão; só publicar nova prova depois da recarga bem-sucedida. Testar aplicar → reverter → reaplicar com falha.

2. **Média — operação de drivers bem-sucedida passa a ser comunicada como falha.**

   [engine.go](../../internal/engine/engine.go), `verify`; [detect.go](../../internal/detect/detect.go), ramo `ubuntu-drivers`.

   O fallback que aceitava exit code zero para operações administrativas foi retirado, mas o detetor de `ubuntu-drivers` continuava a devolver sempre `Available`. Assim, depois de `ubuntu-drivers install` terminar com sucesso, a verificação devolve `Failed`; planos seguintes voltam a propor a operação. É uma regressão introduzida pela alteração da verificação.

   Correção mínima: verificar os drivers recomendados e os pacotes presentes, distinguir indisponibilidade da consulta de trabalho pendente e conservar o aviso de reinício quando aplicável.

3. **Média — uma verificação cancelada apaga a intervenção anterior; verificações recuperadas ficam depois dos dependentes.**

   [recovery.go](../../internal/engine/recovery.go), `ReviewRecovery`.

   O código inicialmente revisto marcava `seen[item]` antes de filtrar o estado do passo. Depois de uma falha de verificação PHP, um retry cancelado antes do passo `verify` torna-se o registo mais recente. O registo `cancelled` impede a consulta da intervenção anterior e o plano seguinte fica sem verificação. `pending`/`skipped` têm o mesmo problema.

   Além disso, `verify` era acrescentado no fim do plano: após falha de verificação PHP com Composer pendente, o retry pode instalar Composer antes de voltar a verificar PHP, apesar da dependência declarada.

   Correção mínima: passos não executados não resolvem intervenções anteriores; ordenar a verificação recuperada antes dos dependentes e testar cancelamento/retoma.

4. **Média — a reversão de um modelo pode apagar um caminho diferente do revisto.**

   [undo_review.go](../../internal/engine/undo_review.go), `undoValues`; [engine.go](../../internal/engine/engine.go), ramo `MethodTemplate` de `undo`.

   Criar `tpl-md`, alterar `XDG_TEMPLATES_DIR` e pedir reversão. A revisão lê o diretório de modelos atual e pode mostrar “sem histórico; manter o valor atual”. A execução inicial continuava a procurar o primeiro registo ativo do item, encontrando e removendo o ficheiro no diretório anterior. O ficheiro tem de continuar a corresponder ao conteúdo original para a remoção ser permitida, mas o alvo não foi o mostrado na revisão.

   Correção mínima: a execução deve consumir exatamente o caminho/ChangeID revisto ou incluir explicitamente os alvos históricos na proposta. Não voltar a procurar outro registo do item depois da autorização.

5. **Alta — exportação pode sobrescrever um ficheiro que o diálogo não confirmou.**

   [app_diagnostics.go](../../app_diagnostics.go), `ExportDiagnostics`; [safefile.go](../../internal/safefile/safefile.go), `Write`.

   O código inicialmente revisto acrescentava a extensão depois de `SaveFileDialog`. Se existe `relatorio.pdf` e o utilizador escolhe `relatorio` inexistente, o diálogo confirma esse nome, mas a aplicação escreve `relatorio.pdf` por rename atómico, substituindo-o. A implementação Linux do Wails v2.16.0 foi inspecionada localmente: ativa confirmação de overwrite GTK em `window.c:793` e devolve diretamente `gtk_file_chooser_get_filenames`, sem acrescentar extensão.

   Correção mínima: confirmar o destino final efetivo ou exigir a extensão no nome devolvido pelo diálogo. A confirmação tem de referir o mesmo ficheiro que será substituído.

6. **Alta — caracteres fora do BMP podem provocar panic na exportação PDF.**

   [export.go](../../internal/diag/export.go), `cleanText` e `pdfReport`; dependência local `github.com/go-pdf/fpdf@v0.9.0`.

   `cleanText` aceita emoji e outros runes acima de U+FFFF. A biblioteca cria `CharWidths` com 65536 elementos (`utf8fontfile.go:839`), atualiza `LastRune` para o maior rune utilizado e percorre até esse valor indexando `font.Cw[cid]` sem limite (`fpdf.go:4452–4456`). Um nome de entrada de arranque com 🚀 pode atingir esse caminho. A fixture original usa Ω, que não testa esta fronteira.

   Correção mínima: testar `Encode` com um emoji, impedir panic e representar de forma explícita os caracteres não suportados pelo renderer PDF, preservando o texto completo nos formatos que o suportam. Este cenário foi identificado por inspeção da implementação da biblioteca; não foi executada uma reprodução adicional nesta passagem.

7. **Média — a revisão das reversões administrativas em lote continua incompleta.**

   [undo_review.go](../../internal/engine/undo_review.go), `ReviewUndo`; [plan.go](../../internal/plan/plan.go), `undoSteps`; [engine.go](../../internal/engine/engine.go), `revertAdminItem`.

   A revisão/validação nova cobre `Op == "undo"`, mas não `Op == "revert"`. O caminho administrativo usa os registos/snapshots disponíveis ao executar. Incluir apenas os valores do ficheiro atual em `Changes` melhora a deteção de mudanças, mas não liga o token ao alvo anterior e à cadeia de snapshots que serão revertidos. A revisão histórica individual já dispõe de uma proposta mais completa.

   Correção mínima: vincular também os alvos e snapshots administrativos à autorização, ou encaminhar a reversão administrativa para o fluxo individual `PreviewRevert`/`RevertReviewed` e recusar a alternativa em lote no RC.

## Evidência executada

No container `meu-ubuntu-rc-build`, com `GOFLAGS=-mod=readonly`, passaram os testes sem `-race` durante a revisão e o check final seguinte:

```text
go test -race -tags webkit2_41 ./internal/engine ./internal/plan ./internal/helper ./internal/detect ./internal/diag ./internal/sys .
```

O check final retornou sucesso para os sete pacotes, reutilizando a cache válida de testes Go para as fontes finais. Passou também `TestPDFAndIndependentReaderFixtures` com `-count=1`, com artefactos isolados finais em `/tmp/astra-diagnostics-final` no container. Foram abertos com leitores independentes já instalados:

- `openpyxl`: XLSX legível, células de texto sem fórmulas, texto `João · ação · configuração · Ω` e emoji 🚀 preservados.
- `pdftotext`: extração correta do mesmo texto português/grego no PDF e de `[U+1F680]` para o emoji.
- `csv.reader`: CSV legível, emoji preservado e 52 linhas, correspondentes às linhas do XLSX.

A exportação usa um snapshot validado antes da serialização. CSV neutraliza prefixos de fórmulas e usa quoting; XLSX usa `inlineStr` e XML escapado. A escrita é atómica com modo 0600. Estes pontos foram verificados no código e nas fixtures.

## Limites

Não foram executadas receitas privilegiadas no host, testes `realsystem`, instalação de dependências ou alterações de fontes. Não foi testado o diálogo GTK real, o fecho durante APT numa sessão GNOME nem a matriz completa de drivers/hardware. O fecho normal tem agora uma barreira e um teste unitário, mas o comportamento do desktop exige validação em VM/sessão apropriada. Não foi demonstrado que fechar a janela mata dpkg ou corrompe pacotes.

Os leitores independentes validaram a fixture portuguesa/grega e o caso de Unicode astral; não constituem prova de renderização visual completa nem de todos os caracteres Unicode. O utilizador autorizou a opção GPLv3 recomendada e a licença local foi identificada como GPL-3.0-only. Esta revisão não afirma que houve publicação externa. O próximo QA Wails/GTK/Xvfb em simulação, quando executado, deve continuar distinguido de testes GNOME/Wayland reais e de instalação numa máquina descartável.
