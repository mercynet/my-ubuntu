# Revisão adversarial — Meu Ubuntu

Data: 2026-10-12  
Revisor: Astra / high

A base tem proteções úteis: catálogo incorporado, comandos por argumentos, validação de opções, `--no-remove`, journal anterior à alteração e snapshots administrativos separados. **Não encontrei uma elevação arbitrária a root demonstrável. Encontrei falhas concretas de recuperação e revisão que corrigiria antes da 1.0.**

Os achados abaixo estão confirmados pelo caminho de código, salvo indicação contrária. Não executei estas reproduções dirigidas; executei os testes existentes.

1. **Alta — repetir uma reversão parcial pode substituir configurações pessoais já restauradas.**  
   Referências: [engine.go:774](../internal/engine/engine.go#L774), [engine.go:788](../internal/engine/engine.go#L788), [engine.go:674](../internal/engine/engine.go#L674).

   Reprodução determinística pelo código: aplicar `wallpaper`, que altera `picture-uri` e `picture-uri-dark`, tendo anteriormente wallpapers pessoais. Alterar externamente apenas a segunda chave e pedir reversão do item. A primeira chave volta corretamente ao wallpaper pessoal e o seu registo fica marcado como revertido; a segunda recusa a reversão porque mudou depois. Ao repetir, `latest()` já não encontra um registo ativo para a primeira chave e entra no caminho `gsettings reset`, substituindo o wallpaper pessoal restaurado pelo padrão do sistema.

   Impacto: a recuperação de uma falha provoca uma alteração adicional indesejada. Outros itens com várias chaves têm o mesmo desenho, incluindo transparência e suspensão.

   Correção mínima: distinguir “nunca gerido” de “já revertido”; persistir a conclusão por chave e retomar apenas as reversões pendentes. O fallback para reset não pode tratar estes dois estados como equivalentes.

2. **Alta — alterações parcialmente aplicadas são consideradas concluídas, impedindo uma retoma eficaz.**  
   Referências: [plan.go:195](../internal/plan/plan.go#L195), [plan.go:238](../internal/plan/plan.go#L238), [detect.go:290](../internal/detect/detect.go#L290), [dropin.go:89](../internal/helper/dropin.go#L89), [detect.go:474](../internal/detect/detect.go#L474).

   Há dois exemplos concretos:

   - Xdebug tem instalação seguida de `phpdismod`. Cancelar depois da instalação, antes do `post`, ou falhar o `phpdismod`, deixa o pacote instalado. A nova deteção consulta apenas os pacotes; `BuildPlan` ignora o item como `Installed` antes de gerar a pós-operação. “Retomar” nunca volta a tentar a desativação prometida.
   - Inotify/journald escrevem o ficheiro antes de chamar `reloadDropin`. Se `sysctl` ou a recarga do serviço falhar, a deteção posterior declara sucesso apenas porque o ficheiro contém o texto esperado. A nova execução omite precisamente a operação que falhou.

   Não existe `defaultActionVerify` a corrigir isto; [engine.go:361](../internal/engine/engine.go#L361) exclui explicitamente `post` da verificação.

   Correção mínima: separar instalação, configuração persistente e aplicação efetiva. A deteção/retoma tem de verificar ou conservar as pós-condições pendentes, sem reinstalar automaticamente software já existente.

3. **Média — o token de revisão não protege os valores das reversões incluídas num plano.**  
   Referências: [plan.go:352](../internal/plan/plan.go#L352), [plan.go:358](../internal/plan/plan.go#L358), [app.go:365](../app.go#L365), [engine.go:780](../internal/engine/engine.go#L780).

   `undoSteps` lê o estado, mas produz texto genérico: não inclui valor atual, alvo ou journal em `Changes`/`Detail`. Exemplo: pedir reversão de `color-scheme` sem histórico, rever o plano e mudar depois essa preferência nas definições do GNOME. Ao executar, o plano reconstruído continua idêntico; o token é aceite e o reset substitui a nova preferência.

   Isto acontece pelo fluxo normal da aplicação. A proteção mais forte existente em `RevertReviewed` não é usada neste caminho.

   Correção mínima: construir uma proposta concreta por chave, incluindo valor atual e alvo; incorporar essa proposta no fingerprint e revalidá-la antes da reversão.

4. **Média — a ponte mantém uma reversão forçada sem revisão.**  
   Referências: [app.go:408](../app.go#L408), [engine.go:580](../internal/engine/engine.go#L580), [engine.go:609](../internal/engine/engine.go#L609), [main.go:91](../main.go#L91).

   O método público `App.Revert(id, force)` continua exposto pelo binding de `App`. Chamar `window.go.main.App.Revert(id, true)` usa `RevertChange` e ignora completamente o token exigido por `RevertReviewed`. Com um registo de template, também permite remover um ficheiro entretanto editado.

   **Modelo de ameaça:** exige capacidade de chamar a ponte fora dos botões atuais, por exemplo frontend comprometido ou código cliente antigo. Não é uma vulnerabilidade remota demonstrada, nem ultrapassa por si só a autorização polkit. É um bypass confirmado do contrato de revisão.

   Correção mínima: retirar esse método mutante da ponte ou fazê-lo recusar chamadas, mantendo a lógica interna necessária exclusivamente atrás da revisão.

5. **Média — a deteção do grupo Docker usa a sessão antiga e pode repetir operações administrativas.**  
   Referências: [detect.go:436](../internal/detect/detect.go#L436), [helper/preview.go:36](../internal/helper/preview.go#L36), [engine.go:387](../internal/engine/engine.go#L387).

   `id -nG` sem utilizador consulta os grupos do processo atual. Depois de `usermod -aG docker`, essa sessão continua com os grupos anteriores; a deteção diz “fora do grupo” e outro plano volta a executar `usermod` e criar journal. Após remoção, o problema inverso pode mostrar pertença já retirada da configuração persistente.

   O próprio helper de preview usa corretamente `id -nG <utilizador>`, pelo que há duas definições diferentes de estado no projeto.

   Correção mínima: consultar a pertença persistente para decidir se há trabalho e apresentar separadamente “requer nova sessão”. Não esconder esta diferença através do fallback que considera qualquer operação administrativa com exit code zero como concluída.

6. **Média — falha de verificação pode terminar como “Concluída”, perder a explicação no histórico e não ser corrigível por retoma.**  
   Referências: [engine.go:372](../internal/engine/engine.go#L372), [engine.go:242](../internal/engine/engine.go#L242), [engine.go:251](../internal/engine/engine.go#L251), [engine.go:279](../internal/engine/engine.go#L279), [PlanRun.tsx:280](../frontend/src/PlanRun.tsx#L280).

   Exemplo: APT instala PHP, mas `php -v` falha. `verify` devolve `StIntervention` com uma nota. Apenas `StFailed` altera o resultado agregado; o run termina `done`. A nota é emitida ao frontend, mas não guardada em `StepRecord`, e o erro dessa verificação também não é escrito no log neste caminho. Depois de fechar a vista, perde-se a causa concreta. Ao rever o que falta, a deteção de pacotes volta a considerar PHP instalado e omite a verificação.

   A interface apresenta a contagem de intervenções em direto; portanto não é ocultação total. O resumo e o histórico são, contudo, inconsistentes.

   Correção mínima: persistir a nota/erro de verificação, representar conclusão com intervenção no run e disponibilizar repetição da verificação/configuração sem reinstalação cega.

7. **Baixa — a reversão Git trata uma preferência ausente como erro de leitura.**  
   Referências: [review_revert.go:162](../internal/engine/review_revert.go#L162), [runner.go:47](../internal/sys/runner.go#L47).

   `gitValue` aceita ausência apenas se `err.Error() == "exit status 1"`. O executor real embrulha o erro como `git: exit status 1: ...`; essa comparação nunca reconhece o resultado real. Se alguém retirar uma preferência depois da aplicação, o preview falha em vez de mostrar “ausente” e permitir a decisão prevista.

   Correção mínima: preservar e inspecionar o código de saída através de um erro tipado, usando `errors.As`, sem comparar mensagens formatadas.

**Risco importante ainda por validar numa sessão GNOME:** [main.go:80](../main.go#L80) não configura `OnBeforeClose`/`OnShutdown`, apesar de a execução correr numa goroutine e de o cancelamento prometer terminar o passo corrente. Fechar a janela pode terminar o coordenador durante uma instalação, perder acompanhamento e deixar processos filhos a terminar de forma não controlada. A ausência da proteção está confirmada; **não demonstrei que fechar a janela mata dpkg ou corrompe pacotes**, pois isso depende do comportamento concreto dos processos e pipes. Validaria este cenário numa VM descartável e bloquearia o fecho normal enquanto se aguarda um limite seguro. O arranque seguinte marca runs como interrompidos, mas isso não substitui proteger o fecho normal.

Para a **1.0**, considero bloqueadores os dois problemas de recuperação, a consistência da revisão de reversões e a validação/proteção do fecho durante execução. Corrigiria também os estados Docker e os resultados com intervenção antes de apresentar o produto como previsível para utilizadores sem conhecimentos técnicos. O bypass legado é pequeno de remover e não tem benefício suficiente para o conservar numa superfície pública.

A instalação real da **0.3.2** reportada pelo utilizador é evidência positiva de utilização. Não comprova estes fluxos da beta atual. Nesta revisão passaram, no container `meu-ubuntu-about-build`:

`go test ./internal/plan ./internal/engine ./internal/helper ./internal/detect ./internal/sys`

Não alterei fontes nem dependências, não executei receitas no host e não executei testes `realsystem`. As reproduções específicas acima resultam de inspeção dos caminhos; os testes existentes não exercitam essas sequências de falha e retoma. Não validei a sessão GNOME, polkit interativo, fornecedores externos atuais ou todas as receitas de Ubuntu 26.04.

Sobre reutilização para premium/multiOS: **reutilizaria os contratos de catálogo, seleção, revisão, resultados e boa parte da interface; não trataria o executor atual como um agente multiOS pronto.** `sys.Runner` abstrai comandos, mas `detect`, `plan` e `engine` conhecem diretamente APT, GNOME, systemd, paths Linux e polkit. O planner também ordena por método e ID, em vez de fazer ordenação topológica — `composer` precede o seu requisito `php` no catálogo atual, por exemplo.

Para gerir várias máquinas, faltam identidade de destino ligada à revisão, estado/journal por máquina, execução durável fora da janela, reconciliação após desconexão e um modelo explícito de autorização remota. São alterações estruturais, especialmente perante os problemas de retoma encontrados. Há uma base aproveitável para evoluir, mas acrescentar SSH e trocar APT por outro gestor de pacotes não basta. As correções de segurança e recuperação devem integrar a edição gratuita existente; o premium pode acrescentar a gestão centralizada sem depender de uma base local menos segura.
