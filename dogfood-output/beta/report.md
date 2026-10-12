# Validação de usabilidade — Meu Ubuntu 0.6.0-beta.1

Data: 12 de outubro de 2026. Alvo: Ubuntu 26.04 LTS, GNOME, amd64.

## Método e limites

Avaliação especializada por tarefas, com reprodução dos problemas da [versão 0.5.0](../v0.5.0/report.md), correções e repetição dos percursos. Não houve participantes humanos: não foram medidos tempos de aprendizagem, satisfação ou taxas de sucesso de utilizadores.

O binário Wails/WebKitGTK foi exercitado como utilizador normal num container Ubuntu 26.04 com Xvfb, em `--simular`, com dados isolados. Assim, revisão, execução, eventos, deteção simulada, histórico e persistência passam pelo backend real da aplicação. O navegador complementou os testes a 900 × 600, incluindo temas e teclado, usando uma ponte de teste; não é evidência de instalação administrativa real. A simulação não altera o sistema anfitrião.

## Correções dos seis problemas iniciais

| Problema da 0.5.0 | Comportamento da beta | Evidência |
|---|---|---|
| Perfil confundido com seleção ativa | **Mostrar catálogo de** filtra a vista. **Minhas escolhas** reúne a seleção completa; perfis têm revisão antes de adicionar ou substituir. | [Escolhas na aplicação nativa](screenshots/native-04-minhas-escolhas.png), [confirmação de substituição](screenshots/12-modal-perfil-corrigido-900x600.png) |
| Começar vazio escondia o catálogo | **Escolher do zero** abre diretamente o catálogo completo; a receita vazia deixou de ser um filtro. | [Assistente](screenshots/native-01-arranque.png), [catálogo](screenshots/native-02-catalogo.png) |
| VSCode não encontrava o editor | Pesquisa reconhece nomes correntes, IDs, pacotes, acentos e espaços. Correspondências no nome têm prioridade. | [VSCode com um resultado](screenshots/native-03-pesquisa-vscode.png), testes `catalogSearch.test.mjs` |
| Limpeza apagava trabalho imediatamente | Confirmação, cancelamento e **Desfazer limpeza**, incluindo opções, origens e pedidos de reversão. | [Confirmação](screenshots/native-05-confirmar-limpeza.png), [cancelamento](screenshots/native-06-limpeza-cancelada.png), [limpeza](screenshots/native-07-escolhas-limpas.png), [recuperação](screenshots/native-08-escolhas-recuperadas.png) |
| Formulário escondia os perfis disponíveis | Perfis aparecem primeiro, separados entre incluídos, pessoais e importados. Guardar abre um diálogo dedicado. | [Perfis nativos](screenshots/native-09-perfis.png), [900 × 600 e tema escuro](screenshots/08-perfis-900x600-escuro.png) |
| Filtros ocupavam várias linhas e contavam o catálogo errado | Categoria e estado usam seletores compactos; contagens acompanham os resultados atuais. | [Pesquisa nativa](screenshots/native-03-pesquisa-vscode.png) |

O cancelamento e a recuperação foram repetidos no navegador com duas escolhas, uma delas com opção personalizada. Foram conferidos os dados persistidos de escolhas, opções e origens antes e depois. No binário nativo, a limpeza e recuperação foram repetidas com o editor escolhido.

## Bloqueios adicionais corrigidos

**Aplicar planos com dependências.** O teste nativo encontrou uma [recusa de execução](screenshots/native-14-resultado.png) mesmo sem mudança de sistema: a lista informativa de dependências adicionadas alterava o identificador de revisão quando o backend reconstruía a seleção expandida. O identificador passa a vincular a proposta efetiva, preservando a verificação dos passos, seleção e estado. Um teste de regressão reconstrói o plano e confirma também que uma instalação real da dependência invalida a revisão antiga.

O plano com Visual Studio Code e Codex adicionou Node.js, executou seis passos e terminou com três itens confirmados. [Resultado nativo](screenshots/native-23-plano-final-concluido.png), [gravação do percurso corrigido](native-aplicar-plano-corrigido.webm).

**Diálogos em janelas pequenas.** Na revisão de um perfil a 900 × 600, o [rodapé ficava abaixo da janela](screenshots/11-modal-perfil-900x600.png). Agora só o corpo do diálogo desloca; título e ações permanecem visíveis. A confirmação de substituição foi medida dentro da janela: rodapé termina em 575 px numa janela de 600 px. [Resultado](screenshots/12-modal-perfil-corrigido-900x600.png).

**Teclado após erro.** Uma falha injetada ao guardar mantém nome, escolhas e mensagem no diálogo. Tab recupera o foco dentro dele quando um botão desativado o perdeu; Shift+Tab circula entre as ações; Escape fecha e devolve o foco à ação de origem. [Erro preservado](screenshots/09-falha-ao-guardar.png). A injeção pertence à ponte de teste, não ao produto.

## Percursos completos

- Guardar um perfil pela interface nativa, encerrar, reabrir e encontrá-lo em **Meus perfis**: [persistência](screenshots/native-19-perfil-persistente.png).
- Abrir um perfil incluído na aplicação nativa, pedir substituição e cancelar: as duas escolhas e as origens manuais foram conservadas. [Confirmação](screenshots/native-25-confirmar-substituicao.png).
- Adicionar **Base essencial** a três escolhas existentes no navegador: seleção passou a 16 itens e conservou as três escolhas iniciais. Cancelar a substituição conservou a seleção.
- Falha controlada em Codex com `MEU_UBUNTU_SIM_FAIL=@openai/codex`: cinco passos concluídos e um falhado, sem declarar sucesso total. **Rever o que falta** gerou um passo de instalação do agente, mantendo editor e Node.js como já instalados. [Falha](screenshots/native-17-falha-controlada.png), [recuperação](screenshots/native-18-rever-o-que-falta.png).
- Diagnóstico e ajuda abriram pela ponte nativa: [diagnóstico](screenshots/native-20-diagnostico.png), [manual](screenshots/native-21-ajuda.png).

## Agentes de IA

A secção reúne **Claude Code, Codex, Gemini CLI, OpenCode, Pi e GitHub Copilot CLI**. Instalações usam os pacotes oficiais com versões explícitas em `~/.local`, sem sudo. Node.js/npm aparecem como dependência quando necessários; agentes existentes que arrancam são preservados. Autenticação ocorre dentro de cada ferramenta.

As seis receitas foram instaladas realmente num Ubuntu 26.04 isolado com Node.js 22.22.1, UID 1000 e prefixo pessoal. Todos os comandos `--version` terminaram com código zero: Claude 2.1.296, Codex 0.162.1, Gemini 0.63.0, OpenCode 1.19.0, Pi 1.1.0 e Copilot 1.0.95. Pi foi instalado com `--ignore-scripts`, conforme a sua documentação. [Resultados registados](agent-startup.json), [secção nativa final](screenshots/native-22-agentes-final.png).

O Codex avisou que não criaria aliases auxiliares num HOME temporário em `/tmp`; o comando arrancou. Não foram testados login, contas, chamadas a modelos, consumo de créditos ou trabalho dos agentes sobre projetos.

## Verificação técnica e entrega

Passaram os testes do frontend, TypeScript/Vite, `go vet` e a suite Go com `-race`, incluindo instalação dos agentes, confirmação obrigatória e regressão dos planos com dependências. O compilador foi atualizado de Go 1.26.8 para **1.26.9**, com arquivo oficial e checksum fixos, após alertas da auditoria. Com a versão corrigida, `govulncheck` e `npm audit` não reportaram vulnerabilidades conhecidas na verificação.

O `.deb` foi instalado com APT num container Ubuntu 26.04 novo, resolvendo as dependências declaradas. Versão, ajuda, manual incorporado, documentos do pacote e checksums foram conferidos. A imagem Ubuntu mínima exclui documentação por configuração do dpkg; a reinstalação do pacote com inclusão explícita confirmou também os documentos em disco. O pacote usa versão Debian `0.6.0~beta1`, anterior a uma futura `0.6.0` final.

## Próxima validação com pessoas

A beta está preparada para testes exploratórios em **Ubuntu 26.04 LTS com GNOME e amd64**. Use o [roteiro de tarefas e feedback](../../docs/beta-usabilidade.md), começando em simulação. Ainda faltam confirmação em sessões GNOME/Wayland reais, polkit, instalações administrativas reais, atualização no computador do autor e sessões com participantes. Não há promessa de suporte a Ubuntu anterior, outras distribuições ou ARM. Exportação PDF/CSV/Excel e planeamento de produto ficaram para depois desta fase.
