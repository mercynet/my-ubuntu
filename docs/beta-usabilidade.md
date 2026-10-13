# Testar o Meu Ubuntu 1.0.0-rc.4

## Ambiente suportado

Ubuntu **26.04 LTS**, GNOME, computador **amd64/x86_64**. Não há suporte nesta candidata para versões anteriores, outras distribuições ou ARM. Comece pelo modo de simulação; a instalação real precisa de ligação à Internet e autorização administrativa para operações de sistema.

## Instalar e abrir

Para teste, envie o ZIP da RC ou estes ficheiros de `dist/`: `meu-ubuntu_1.0.0-rc.4_amd64.deb`, `meu-ubuntu-1.0.0-rc.4-linux-amd64`, `meu-ubuntu-1.0.0-rc.4-source.tar.gz` e `SHA256SUMS`. Na pasta recebida:

```bash
sha256sum -c SHA256SUMS
sudo apt install ./meu-ubuntu_1.0.0-rc.4_amd64.deb
meu-ubuntu --simular
```

Abra a aplicação como utilizador normal. Para experimentar operações reais depois da simulação, encerre-a e execute `meu-ubuntu`, ou use o atalho do menu. O manual acompanha o pacote em `/usr/share/doc/meu-ubuntu/manual.md`; também está disponível em **Ajuda e manual**.

## Roteiro de 15–20 minutos

Faça as tarefas sem consultar o manual inicialmente. Registe onde hesitou, o que esperava e o resultado. Não há respostas certas a decorar.

1. Comece com **Escolher no catálogo**. Encontre o Visual Studio Code pesquisando `VSCode`; escolha-o e também `htop`.
2. Descubra onde ver tudo o que escolheu. Antes de abrir o plano, explique com as suas palavras o que será instalado.
3. No catálogo, use **Mostrar catálogo de** para ver outro perfil. Confira se isso mudou as suas escolhas ou apenas os itens apresentados.
4. Em **Minhas escolhas**, tente limpar e cancele. Depois confirme a limpeza e use **Desfazer limpeza**. Confira se recuperou tudo.
5. Guarde as escolhas como um perfil chamado `Meu teste`. Encontre-o novamente em **Perfis prontos → Os meus**.
6. Abra um perfil incluído e examine as escolhas antes de o adicionar. Veja se distingue adicionar de substituir.
7. Abra **Agentes de IA**, escolha um agente e consulte as instruções de arranque. Reveja o plano: a dependência Node.js deve aparecer se necessária, e a instalação do agente deve ser por utilizador.
8. Feche a simulação e abra `meu-ubuntu` sem aplicar um plano. Em **Diagnóstico → Entradas de arranque**, confira se distingue comandos disponíveis, entradas desativadas e executáveis não encontrados; consulte a origem de uma entrada que não reconhece. Feche e volte a abrir com `--simular`.
9. Ainda em simulação, reveja e aplique um plano pequeno. Observe progresso, resultado e histórico. Retire uma aplicação instalada das escolhas: deve continuar como Já instalado e indicar que permanece instalada. Atualizar estado e reabrir a aplicação devem respeitar essa exclusão explícita. Reveja as escolhas já aplicadas: o plano deve mostrar zero alterações, sem reinstalações ou reversões. Feche e reabra para verificar que o perfil guardado permanece.

10. Ao abrir, confirme que os programas já instalados ficam marcados sem carregar num botão de inclusão. Experimente atualizar o estado após instalar uma ferramenta fora do Meu Ubuntu.
11. Abra Plugins e complementos, filtre Zsh e inclua o grupo. Reveja os pacotes propostos, abra o guia de ativação e experimente os grupos de agentes de IA. Os guias destes agentes devem ser apresentados como manuais.

12. Abra Desenvolvimento → Ver serviços locais. Encontre Apache, Redis, RabbitMQ e o guia Kafka. Reveja a confirmação de instalação de um servidor, sem a executar no computador principal. Já instalado confirma o pacote; o guia explica como verificar o serviço.

13. Expanda Explorar por área e encontre Sobre no rodapé. Confirme a versão em execução, consulte o manual e copie as informações técnicas para suporte. Repita numa janela de 900 × 600: Sobre e a versão devem permanecer visíveis enquanto o menu tem deslocamento próprio.

Opcionalmente, num Ubuntu 26.04 de teste, repita uma instalação real pequena, como `htop`. Instale um agente apenas se quiser usá-lo; autenticação, subscrição e eventuais custos dependem do fornecedor. A aplicação não pede chaves nem inicia sessões por si.

## Enviar feedback

Para cada dificuldade, envie: versão do Meu Ubuntu; versão do Ubuntu e sessão Wayland/X11; resolução; tarefa; passos; resultado esperado; resultado observado; captura ou gravação, se possível. Remova credenciais e dados pessoais. Diga também se concluiu sem ajuda e avalie a facilidade da tarefa de 1 (muito difícil) a 5 (muito fácil).

## Limites da candidata

É uma candidata para recolher feedback antes da 1.0 estável. A avaliação automatizada e especializada não substitui sessões com pessoas: instalações administrativas, GNOME/Wayland, polkit e autenticação dos agentes precisam de confirmação em computadores reais. PDF/CSV/XLSX já estão disponíveis: carregue o diagnóstico, exporte cada formato, retire uma secção e confirme que o conteúdo respeita a escolha. Cancelar a janela de gravação não deve criar um ficheiro. Outras distribuições, Windows/macOS, gestão remota e novos perfis profissionais ficam para depois da 1.0.

A [revisão final Astra/high](../dogfood-output/rc.1/review-astra-high.md) conferiu as correções e regressões, sem bloqueadores de código remanescentes nos cenários revistos. Isso não substitui a matriz desktop: registe testes de atualização, autenticação/cancelamento polkit e alterações reais numa máquina ou VM com backup.
