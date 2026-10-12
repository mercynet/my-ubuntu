# Avaliação de usabilidade do Meu Ubuntu 0.5.0

Data: 12 de outubro de 2026.

Interface de produção em Chromium a 1280 × 840 e 900 × 600 com ponte Wails de teste; adicionalmente, binário nativo original 0.5.0, Wails/WebKitGTK, em Ubuntu 26.04 amd64 com Xvfb a 1280 × 840 e --simular. Avaliação especializada por tarefas, sem participantes humanos. A ponte do navegador não prova execução desktop. Instalações administrativas reais, GNOME/Wayland, polkit e autenticação não foram validados.

## Resultados

Seis problemas reproduzidos: um de gravidade alta e cinco de gravidade média. A interface confunde a vista do catálogo com a seleção que será aplicada; também permite apagar escolhas sem confirmação e exige demasiada navegação para tarefas simples.

Esta fase preserva a interface da v0.5.0. As correções e a sua validação estão em [avaliação do beta](../beta/report.md). Os problemas foram observados na interface, antes de consultar a implementação para os corrigir.

## Problemas reproduzidos

### ISSUE 001 O seletor Perfil muda a vista mas mantém outro perfil nas escolhas

**Gravidade:** média. **Categoria:** usabilidade. **Prioridade:** alta.

Após começar com «Uso pessoal / escritório», escolher «Base essencial» no seletor «Perfil» muda o título e os cartões apresentados, mas mantém 17 escolhas e 12 escolhas fora da vista. Os cartões continuam a indicar «Origem: perfil Uso pessoal / escritório». O comportamento protege as escolhas, mas o nome do controlo sugere uma troca do perfil ativo. «Usar perfil…» é uma segunda ação separada.

**Reprodução confirmada duas vezes:**

1. No assistente, escolher «Uso pessoal / escritório» e continuar. [Estado inicial](screenshots/02-inicio-perfil.png).
2. Selecionar «Base essencial» no controlo «Perfil». [Resultado](screenshots/issue-001-corrected-result.png).
3. Observar o título «Base essencial», o contador «17 no perfil» e o aviso de 12 escolhas fora da vista.

**Consequência provável:** interpretar o conjunto visível como o conjunto que será aplicado. Esta consequência é uma inferência da avaliação especializada; não foi medida com participantes.

**Correção proposta:** chamar ao seletor «Mostrar catálogo de» e ao botão «Adicionar escolhas deste perfil…»; mostrar a seleção atual numa área com nome próprio e contagem acessível. Se houver um perfil ativo, identificá-lo explicitamente.

**Critério de aceitação:** antes de rever o plano, uma pessoa distingue o filtro do catálogo das escolhas que serão aplicadas sem consultar o manual.

### ISSUE 002 Começar vazio esconde o catálogo necessário para começar

**Gravidade:** média. **Categoria:** usabilidade. **Prioridade:** alta.

Selecionar «Personalizado (começar vazio)» mostra «Sem itens: escolha tudo nas áreas», mas as áreas também ficam vazias devido ao filtro. O contador pode continuar a mostrar escolhas anteriores. O botão «Explorar Todas…» permite recuperar, mas é necessário descobrir que o ponto de partida bloqueou a exploração que a própria mensagem recomenda.

**Reprodução:** selecionar «Personalizado (começar vazio)» ([primeiro estado](screenshots/04-personalizado-vazio.png)); abrir Aplicações ([resultado](screenshots/issue-002-areas-vazias.png)); repetir em Desenvolvimento ([segunda reprodução](screenshots/issue-002-desenvolvimento-vazio.png)). «Explorar Todas…» recupera o catálogo ([recuperação](screenshots/issue-002-recuperacao.png)).

**Correção proposta:** oferecer «Escolher do zero» como ação de seleção, com acesso imediato ao catálogo completo. Remover a receita vazia das opções de filtro.

**Critério de aceitação:** começar do zero permite escolher uma aplicação imediatamente, sem alterar outro filtro.

### ISSUE 003 A pesquisa não reconhece VS Code

**Gravidade:** média. **Categoria:** pesquisa. **Prioridade:** alta.

«vscode» e «VS Code» devolvem zero resultados com as opções avançadas desligadas. Com opções avançadas, «VS Code» encontra «Mais watchers inotify», pela descrição, mas não o editor. «Code» encontra «Visual Studio Code». A expressão VS Code aparece no texto do perfil PHP/Laravel, pelo que é uma pesquisa plausível sugerida pela própria aplicação.

**Reprodução:** pesquisar «VS Code» ([zero resultados](screenshots/issue-003-step-1.png)); ativar opções avançadas ([resultado indireto](screenshots/issue-003-step-2.png)); desligar e repetir ([reprodução](screenshots/issue-003-repeat.png)); pesquisar «Code» ([editor encontrado](screenshots/issue-003-code-encontrado.png)).

**Correção proposta:** reconhecer nomes correntes, IDs e nomes de pacotes; normalizar espaços, maiúsculas e acentos. Priorizar correspondências no nome em relação ao texto descritivo.

**Critério de aceitação:** «VS Code», «vscode» e «Visual Studio Code» encontram o mesmo editor, antes de receitas que apenas o mencionam.

### ISSUE 004 Limpar seleção apaga trabalho sem confirmação ou desfazer

**Gravidade:** alta. **Categoria:** usabilidade e perda de trabalho. **Prioridade:** alta.

O botão elimina todas as escolhas imediatamente. Não surge confirmação nem opção de desfazer. A seleção vazia é guardada nas preferências. A reprodução com 17 itens e com uma escolha manual de htop teve o mesmo resultado.

**Reprodução:** preparar escolhas e abrir Gerir perfis ([uma escolha manual](screenshots/issue-004-step-1.png)); clicar «Limpar seleção» ([zero escolhas](screenshots/issue-004-step-2.png)). A primeira reprodução está em [captura adicional](screenshots/10-limpar-selecao.png).

**Correção proposta:** confirmar a limpeza e permitir desfazer, preservando opções e origens. Não apresentar a ação junto dos botões principais de guardar e exportar.

**Critério de aceitação:** cancelar mantém as escolhas; confirmar limpa; desfazer recupera exatamente escolhas, opções e origens.

### ISSUE 005 Gerir perfis dá prioridade ao formulário de criação

**Gravidade:** média. **Categoria:** hierarquia visual. **Prioridade:** média.

A página começa com um formulário de quatro campos e ações distintas de guardar, exportar, adicionar existentes e limpar. As ações dos perfis incluídos ficam abaixo da dobra a 1280 × 840; a 900 × 600 nem os cartões aparecem. O percurso comum de escolher um perfil exige passar por um formulário que não é necessário para essa tarefa.

**Evidência:** [1280 × 840](screenshots/06-gerir-perfis.png) e [900 × 600](screenshots/07-perfis-900x600.png). É um problema de hierarquia e descoberta, não de impossibilidade de deslocar a página.

**Correção proposta:** separar usar perfis de criar/guardar um perfil. Apresentar os perfis primeiro, com uma ação principal clara e as operações secundárias num menu ou detalhe.

**Critério de aceitação:** quem chega à página para usar um perfil reconhece a ação principal sem passar pelo formulário de criação.

### ISSUE 006 Os filtros ocupam espaço e apresentam contagens alheias à pesquisa

**Gravidade:** média. **Categoria:** hierarquia visual e conteúdo. **Prioridade:** média.

A pesquisa com zero ou um resultado mantém dezenas de botões de categoria e a contagem «Todas 150». Isso ocupa várias linhas e afasta o resultado do campo de pesquisa. As contagens referem-se ao catálogo anterior aos filtros, enquanto o título informa o número de resultados atuais.

**Evidência:** [pesquisa sem resultados](screenshots/issue-003-step-1.png) e [pesquisa com o editor encontrado](screenshots/issue-003-code-encontrado.png). A mudança entre opções avançadas também altera as contagens sem tornar esse âmbito explícito.

**Correção proposta:** usar um filtro compacto de categoria; contar os resultados atuais; mostrar os filtros ativos e uma forma de os limpar.

**Critério de aceitação:** pesquisar uma aplicação apresenta o resultado antes de uma lista de categorias; as contagens correspondem ao conjunto apresentado.


## Confirmação na aplicação nativa original

O binário 0.5.0 foi aberto como utilizador normal, com dados de simulação isolados. As capturas seguintes usam a ponte Wails e o backend Go da aplicação:

- [Arranque original](screenshots/native-01-arranque.png) e [início com escolhas](screenshots/native-02-inicio.png).
- [Filtro «começar vazio» mantendo 17 escolhas](screenshots/native-03-filtro-perfil.png).
- [Gerir perfis antes de limpar](screenshots/native-07-gerir-perfis.png) e [limpeza imediata, sem confirmação](screenshots/native-08-limpeza-confirmada.png). [Gravação da perda das escolhas](videos/native-limpar-selecao.webm).
- [VSCode sem resultados](screenshots/native-09-vscode-sem-resultados.png) e [Code encontra o editor](screenshots/native-10-code-encontrado.png). [Gravação da comparação](videos/native-pesquisa-vscode.webm).

Estas reproduções reforçam os problemas de descoberta, pesquisa e limpeza. Não constituem medidas de tempo, taxa de sucesso ou satisfação de utilizadores humanos.
