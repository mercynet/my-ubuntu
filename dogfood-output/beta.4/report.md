# Marcação automática, complementos e serviços locais — 0.6.0-beta.4

Data: 12 de outubro de 2026. Alvo: Ubuntu 26.04 LTS, GNOME, amd64.

## Correção da marcação

A deteção atualizava o selo **Já instalado**, mas deixava a seleção independente. A beta.3 esclarecia essa separação sem cumprir a marcação automática pedida. Esta revisão altera o comportamento: aplicações, extensões e preferências confirmadas ficam marcadas ao abrir e ao atualizar o estado, incluindo instalações feitas fora do Meu Ubuntu.

Desmarcações explícitas ficam guardadas. A deteção respeita reversões pedidas, opções entretanto alteradas, alternativas já escolhidas e perfis em edição. Ações administrativas, como remoções, continuam a exigir escolha. O plano não instala dependências de um item que já corresponde à proposta; um Oh My Zsh existente, por exemplo, não passa a instalar plugins em falta só por ter sido detetado.

## Organização dos complementos

**Plugins e complementos** agrupa por ferramenta e permite filtrar, abrir e fechar grupos. Zsh reúne sugestões pelo histórico, destaque de sintaxe, fzf e zoxide, com inclusão conjunta dos pacotes, guias de ativação e o guia Oh My Zsh. Os cartões das ferramentas abrem diretamente o grupo correspondente.

Claude Code, Codex, OpenCode, Pi e Gemini CLI têm guias manuais próprios. A interface distingue esses guias dos pacotes instaláveis. A marca da checklist acompanha passos pessoais; não deteta nem confirma um plugin instalado. As instruções foram conferidas nas documentações de [Claude Code](https://code.claude.com/docs/en/discover-plugins), [Codex](https://learn.chatgpt.com/docs/build-skills), [OpenCode](https://opencode.ai/docs/plugins/), [Pi](https://github.com/earendil-works/pi/blob/main/packages/coding-agent/docs/packages.md) e [Gemini CLI](https://geminicli.com/docs/extensions/).

## Serviços locais

**Desenvolvimento → Ver serviços locais**, também acessível em **Explorar por área**, organiza Web, dados/cache, mensageria e contentores. Foram acrescentados Apache, Nginx, Redis, RabbitMQ, PostgreSQL e Memcached. MariaDB e Docker integram a mesma vista.

Os novos pacotes de servidor pedem confirmação explícita de que a instalação pode ativar um serviço. Os guias distinguem pacote instalado de serviço ativo e configurado. Apache e Nginx apresentam um aviso sobre a porta 80, permitindo configurações deliberadas com portas diferentes. Não são criados sites, bases de dados, credenciais ou configurações de produção.

Kafka apresenta um laboratório manual com `apache/kafka:4.3.1`, conforme o [quickstart oficial](https://kafka.apache.org/quickstart/). O exemplo publica apenas no localhost e não conserva dados após remover o contentor. Não força a instalação de outro Docker.

## Validação

Passaram os nove testes frontend, TypeScript/Vite, documentação gerada, `go vet`, a suite Go com `-race` e a verificação dos artefactos. Novas regressões cobrem marcação automática, exclusões persistentes, resultados antigos de deteção, grupos exclusivos, dependências de instalações existentes, confirmação dos servidores, conclusão da instalação e reaplicação sem alterações.

| Percurso | Resultado e evidência |
| --- | --- |
| Primeira abertura | Itens instalados marcados; continuar para o catálogo conserva as marcas. [Catálogo nativo](screenshots/native-04-catalogo-confirmado.png). |
| Instalações externas | OpenVPN e SMB instalados no container depois da abertura; Atualizar estado marcou ambos. [Resultado](screenshots/native-06-rede-atualizada.png). |
| Desmarcar e reabrir | Déjà Dup permanece instalado e excluído das escolhas após atualização e reinício. Preferências verificadas no disco. |
| Rever os já prontos | Zero instalações, configurações e reversões. [Plano nativo](screenshots/native-08-plano-sem-alteracoes.png). |
| Complementos Zsh | Inclusão conjunta mantém sugestões já instaladas e propõe apenas os três pacotes em falta. [Grupo](screenshots/native-13-zsh-final.png), [plano](screenshots/native-11-plano-complementos.png). |
| Guias por agente | Os cinco filtros abrem guias de três passos; atalho do cartão Codex abre só esse grupo. [Atalho](screenshots/browser-07-atalho-codex.png). |
| Progresso pessoal | Passo do guia Codex permanece assinalado após recarregar. [Guia](screenshots/browser-05-codex-900-escuro.png). |
| Editar um perfil | Uma instalação detetada durante a edição não entra no perfil. Após descartar e atualizar, entra nas escolhas do computador. |
| Ecrã pequeno | Navegação, filtros e guias conferidos a 900 × 600, com tema escuro e sem transbordo horizontal. [Nativo](screenshots/native-16-codex-900.png). |
| Servidores Ubuntu | Seis pacotes instalados realmente em container Ubuntu 26.04; versões e configurações Apache/Nginx conferidas. [Versões](services-versions.txt). Arranque automático foi bloqueado apenas no ambiente de teste. |
| Estado dos servidores | Os seis servidores instalados fora da aplicação aparecem marcados automaticamente. [Vista nativa](screenshots/native-18-servicos-instalados.png). |
| Execução pela aplicação | Apache em simulação: confirmação obrigatória, três passos concluídos e nova revisão com zero alterações. [Execução](screenshots/native-22-servidor-executado-sim.png), [reaplicação](screenshots/native-23-reaplicacao-servidor.png). |
| Kafka real | Imagem oficial iniciada num container sem portas publicadas no anfitrião; tópico criado e uma mensagem produzida e consumida. Container de teste removido após validação. |

Os testes de navegador usam os assets de produção com uma ponte Wails de teste. Validam interação e persistência da interface; a execução real do backend é coberta separadamente pelos testes Go e pela aplicação nativa. Nenhuma instalação foi feita no sistema do utilizador.

## Ajustes encontrados nos testes

Corrigidos o ícone que ocupava a largura do cabeçalho de um grupo, a barra de preferências sem região acessível e o papel acessível dos filtros. As confirmações explícitas deixam de repetir o aviso genérico de item sensível. A análise automática de acessibilidade não encontrou violações nas vistas verificadas. A verificação inicialmente inconclusiva do contraste da navegação escura foi complementada pelo cálculo das cores efetivas: 4,74:1.

## Limites e distribuição

Esta avaliação técnica não substitui sessões com outras pessoas em GNOME/Wayland, polkit real ou configuração dos serviços. Os guias de extensões dos agentes são manuais; não existe ainda instalação automática de plugins de IA. PDF/CSV/Excel, outras distribuições e perfis profissionais ficam para a fase seguinte.

O ZIP desta revisão contém o `.deb`, o binário portátil, checksums, o roteiro de testes e o manual. É uma beta para recolher feedback em Ubuntu 26.04 LTS amd64.
