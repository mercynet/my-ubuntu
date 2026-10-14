<!-- Gerado de docs/manual.json; execute python3 build/scripts/generate-docs.py. -->

# Manual do Meu Ubuntu

Guia da versão 1.0.0. O Meu Ubuntu ajuda a escolher aplicações e ajustes, rever as alterações, executá-las e conferir o resultado. O catálogo, este manual e os logos estão incluídos na aplicação; a instalação de software pode precisar de internet.

## Índice

- [Primeiros passos](#primeiros-passos)
- [Instalar, atualizar e remover](#instalar-e-atualizar)
- [Explorar o catálogo e ver as escolhas](#perfil-no-topo)
- [Estado, escolhas e alterações pendentes](#estado-e-selecao)
- [Plugins e complementos por ferramenta](#plugins-complementos)
- [Serviços locais para desenvolvimento](#servicos-locais)
- [Catálogo, áreas e pesquisa](#areas-e-pesquisa)
- [Guardar, combinar e partilhar perfis](#perfis-reutilizaveis)
- [Rever o plano](#rever-plano)
- [Executar, acompanhar, pausar e cancelar](#executar-e-cancelar)
- [Recuperar de falhas](#falhas-e-retoma)
- [Histórico e reversão](#historico-e-reversao)
- [Desktop e estilo](#desktop-e-estilo)
- [Extensões GNOME](#extensoes)
- [Energia e hardware](#energia-e-hardware)
- [Desenvolvimento e terminal](#desenvolvimento-e-terminal)
- [Infraestrutura e segurança](#infraestrutura-e-seguranca)
- [Diagnóstico](#diagnostico)
- [Migração e dados pessoais](#migracao)
- [Guias manuais](#guias-manuais)
- [Permissões e ligações à rede](#permissoes-e-rede)
- [Dados locais e privacidade](#dados-e-privacidade)
- [Simulação e opções de terminal](#simulacao-e-terminal)
- [Problemas frequentes](#problemas-frequentes)
- [Pacotes de integração](#pacotes-de-integracao)
- [Contribuir, licenças e limites](#contribuir-e-licencas)
- [Instalar agentes de IA](#agentes-de-ia)
- [Flatpak e origens de instalação](#flatpak-e-origens)

<a id="primeiros-passos"></a>

## Primeiros passos

Um percurso curto da primeira abertura ao resultado.

O alvo desta versão é Ubuntu 26.04 com GNOME, arquitetura amd64 e sessão de utilizador. Fedora e outras distribuições ainda não têm suporte de instalação validado. A aplicação não deve ser aberta com sudo.

1. Abra em simulação para conhecer os fluxos: meu-ubuntu --simular.
2. Na primeira abertura, explore o catálogo, use um perfil pronto ou reveja os itens deste computador. Os já instalados ficam marcados automaticamente.
3. Explore Catálogo ou Perfis prontos. Adicionar às escolhas preserva as anteriores; Substituir escolhas troca as escolhas pendentes e conserva os itens detetados como instalados, exceto exclusões explícitas.
4. Veja tudo o que escolheu em Minhas escolhas. Ajuste as opções, retire itens ou guarde as escolhas como perfil.
5. Clique em Rever plano. Confira instalações, dependências, alterações e confirmações necessárias.
6. Aplique apenas depois da revisão. As operações administrativas pedem autorização ao sistema.
7. Confira os resultados. Os guias dos agentes de IA explicam como iniciar sessão dentro de cada ferramenta.

> Comece por um plano pequeno. Uma execução concluída pode ainda ter tarefas manuais pendentes; as duas informações aparecem separadamente.

<a id="instalar-e-atualizar"></a>

## Instalar, atualizar e remover

Pacote .deb, binário portátil e conservação dos dados.

O .deb instala o executável, o atalho no menu, o ícone e a documentação. O APT resolve as bibliotecas declaradas. Não precisa de Go, Node ou npm para usar a aplicação.

```sh
sha256sum -c SHA256SUMS --ignore-missing
sudo apt install ./meu-ubuntu_1.0.0_amd64.deb
meu-ubuntu --version
meu-ubuntu
```

O checksum confirma que o ficheiro corresponde à lista recebida; a confiança também depende de obter o pacote e a lista de uma origem de confiança. Para atualizar, instale o novo .deb pelo mesmo comando. Não existe atualizador automático integrado.

```sh
cp meu-ubuntu-1.0.0-linux-amd64 ~/meu-ubuntu
chmod +x ~/meu-ubuntu
~/meu-ubuntu --simular
```

O binário precisa das mesmas bibliotecas de execução GTK 3 e WebKitGTK 4.1. Uma pen montada com noexec impede executá-lo. Copie-o para uma pasta local do utilizador; isso também evita problemas de acesso pelo pkexec.

```sh
sudo apt remove meu-ubuntu
```

Remover o Meu Ubuntu não desinstala as aplicações que escolheu nem repõe os ajustes feitos. Os dados pessoais em pastas XDG permanecem. Antes de apagar esses dados, guarde perfis e backups de que precisa e tenha em conta que perderá o histórico para reversões.

<a id="perfil-no-topo"></a>

## Explorar o catálogo e ver as escolhas

O filtro do catálogo e as escolhas que serão aplicadas são apresentados em locais distintos.

Catálogo permite explorar aplicações e ajustes. O controlo Mostrar catálogo de filtra os itens de um perfil, com as suas dependências e alternativas. Todo o catálogo remove esse filtro. Filtrar não adiciona nem apaga escolhas.

| Ação | Resultado |
| --- | --- |
| Mostrar catálogo de | Filtra a vista sem mudar as escolhas. |
| Adicionar este perfil… | Abre uma revisão antes de adicionar ou substituir escolhas. |
| Todo o catálogo | Mostra as aplicações e os ajustes de todas as áreas. |
| Pesquisar no topo | Encontra itens em todo o catálogo, independentemente da vista anterior. |
| Minhas escolhas | Mostra as escolhas completas e os pedidos de reversão. |
| Perfis prontos | Permite usar, guardar, editar, importar ou exportar perfis. |

O contador no topo refere-se às escolhas completas, não aos itens visíveis. O catálogo indica quando existem escolhas fora da vista atual. Os itens avançados escolhidos continuam visíveis mesmo com Mostrar opções avançadas desligado.

> Limpar escolhas pede confirmação e permite Desfazer limpeza enquanto não fizer novas escolhas. Limpar escolhas não desinstala aplicações nem apaga perfis guardados.

<a id="estado-e-selecao"></a>

## Estado, escolhas e alterações pendentes

Os itens confirmados como instalados ficam marcados automaticamente.

Ao abrir a aplicação e ao atualizar o estado, aplicações, extensões e preferências confirmadas como instaladas, configuradas ou já incluídas no Ubuntu ficam marcadas automaticamente. Isto inclui programas instalados fora do Meu Ubuntu. Ações administrativas, como remoções, exigem a sua escolha. O plano volta a verificar o sistema e não reinstala o que já está pronto.

> Por exemplo: Déjà Dup e OpenVPN já instalados aparecem marcados. Se os desmarcar, permanecem no computador; a sua exclusão é guardada e respeitada na próxima abertura. Configurações com reversão pedida e alternativas de um grupo com outra escolha pendente não são marcadas automaticamente.

| Estado | Como interpretar |
| --- | --- |
| Já instalado / Já configurado / Já existe | A deteção encontrou o programa, valor proposto ou ficheiro. Confira os detalhes da deteção. |
| Disponível | Pode ser proposto para instalação. |
| Requer configuração | Existe, mas o estado atual difere do pretendido. |
| Manual | A tarefa tem instruções para concluir por si. |
| Incompatível | Não pode ser aplicado neste ambiente. |
| Desconhecido | Não houve evidência suficiente; reveja a causa antes de aplicar. |
| Já incluído no Ubuntu | Um equivalente já vem no sistema; evita duplicação quando detetável. |

Atualizar estado deteta também novas instalações e marca-as. Para recuperar itens que desmarcou, use Mais ações do catálogo → Voltar a marcar os já instalados nesta vista. Ao editar um perfil guardado, a deteção não acrescenta itens ao perfil sem a sua escolha.

O painel do início separa escolhas confirmadas, pendentes, manuais e por verificar. Essa informação não mede a saúde global do computador: refere-se aos itens escolhidos no catálogo.

<a id="plugins-complementos"></a>

## Plugins e complementos por ferramenta

Zsh, Claude Code, Codex, OpenCode, Pi e Gemini CLI organizados em grupos.

Abra Plugins e complementos na navegação. O filtro Ferramenta mostra só o grupo escolhido; cada cabeçalho abre ou fecha os complementos. Nos cartões das ferramentas, Ver plugins e complementos leva diretamente ao grupo correspondente.

O grupo Zsh reúne sugestões pelo histórico, destaque de sintaxe, fzf, zoxide e o guia Oh My Zsh. Incluir complementos do grupo marca os quatro pacotes; o plano instala o que falta e prepara as dependências. Os guias Como ativar este complemento explicam a integração no .zshrc.

> Os grupos dos agentes de IA apresentam guias manuais com documentação oficial. Skills, plugins, extensões e pacotes usam mecanismos próprios de cada ferramenta. Estes guias não são uma deteção de plugins instalados e não instalam extensões automaticamente. As marcas da checklist servem para acompanhar os seus passos.

<a id="servicos-locais"></a>

## Serviços locais para desenvolvimento

Web, dados, cache, mensageria e Docker na mesma vista.

Abra Desenvolvimento e use Ver serviços locais, ou Explorar por área → Serviços locais. A vista reúne Apache e Nginx, Redis e Memcached, PostgreSQL e MariaDB, RabbitMQ, o laboratório Kafka e as opções Docker. Estes grupos podem fazer parte de perfis de desenvolvimento ou DevOps.

> Pacotes de servidor podem ativar serviços e arrancar com o sistema. As novas instalações pedem confirmação no plano. Já instalado confirma o pacote, não que o serviço esteja ativo ou configurado para os seus projetos. Consulte Verificar e configurar o serviço e a documentação oficial de cada ferramenta.

Apache e Nginx usam a porta 80 por omissão; o plano avisa se escolher ambos. Configure portas ou interfaces distintas para os usar em simultâneo. A aplicação não cria sites, bases de dados, credenciais, regras de firewall nem configurações de produção.

Kafka é um guia manual com a imagem oficial apache/kafka:4.3.1. O exemplo só publica a porta 9092 no localhost e usa um contentor temporário sem persistência. O guia não instala outro Docker automaticamente, para preservar a opção Docker que já utiliza.

<a id="areas-e-pesquisa"></a>

## Catálogo, áreas e pesquisa

Encontre aplicações, fontes, alternativas e opções avançadas.

O catálogo está dividido em Aplicações, Desenvolvimento, Infraestrutura, Segurança, Estilo, Desktop, Extensões e Energia. Use categorias e o filtro de estado para reduzir a lista. A pesquisa aceita várias palavras e ignora diferenças de acentos e maiúsculas.

Os logos locais ajudam a reconhecer aplicações; alguns itens usam um símbolo da categoria quando não há logo incluído. O nome, a fonte e o método de instalação continuam a identificar a proposta. Logos e marcas não significam patrocínio dos fornecedores.

Abra Detalhes, origem e verificação para ver pacotes, repositório, UUID de extensão, chaves de configuração, conflitos, notas e evidência de deteção. Itens avançados incluem serviços locais, dependências de compilação e ferramentas de uso específico. Itens avançados já escolhidos continuam visíveis.

O menu mantém as entradas principais no mesmo lugar. Explorar por área dá acesso a Aplicações, Desenvolvimento, Infraestrutura, Segurança, Estilo, Desktop, Extensões e Energia. A página Agentes de IA reúne as ferramentas de programação no terminal.

Grupos exclusivos, como o perfil de energia, usam um select com Não alterar. Só uma alternativa pode ser escolhida. Ver uma alternativa atual não a inclui automaticamente no perfil.

<a id="perfis-reutilizaveis"></a>

## Guardar, combinar e partilhar perfis

Perfis de uso, complementos e revisão de importações.

Um perfil é uma receita declarativa de IDs do catálogo e opções permitidas. Os perfis incluídos sugerem um uso principal; complementos acrescentam capacidades como terminal, segurança ou ajustes de desktop. Não contêm comandos arbitrários, senhas, chaves privadas ou ficheiros pessoais.

1. Em Perfis prontos, reveja um perfil antes de o adicionar às escolhas. Pode repetir o processo para combinar perfis.
2. Juntar preserva escolhas atuais; Substituir propõe trocar o conjunto. Reveja requisitos, itens sensíveis e avisos.
3. Se dois perfis discordarem numa opção ou alternativa, resolva o conflito no diálogo. A última seleção não ganha silenciosamente.
4. Guarde as escolhas com um nome e uma descrição. Um nome já existente exige uma escolha explícita de substituição.
5. Exporte o JSON para partilhar. No outro computador, importe e reveja a compatibilidade e o estado desse sistema.

Editar um perfil abre uma sessão de edição. Guardar perfil conserva as alterações da receita; Descartar edição recupera as escolhas anteriores à edição. Eliminar um perfil remove a receita guardada e não desinstala as aplicações.

Importar abre sempre uma revisão. Pode guardar o perfil importado para usar mais tarde sem o aplicar. Abrir um JSON pelo terminal também abre uma revisão na janela existente; nunca inicia uma instalação automaticamente.

```sh
meu-ubuntu --perfil ./meu-perfil.json
```

> O formato atual é meu-ubuntu/profile, versão 2. Marcas pessoais dos guias manuais, histórico, logs e filtro do catálogo não são exportados com a receita. Um perfil não é um backup completo do computador.

<a id="rever-plano"></a>

## Rever o plano

Dependências, diferenças, conflitos e confirmações.

Rever plano calcula o trabalho necessário com base nas escolhas e na deteção atual. O plano distingue operações automáticas, itens mantidos como estão, itens ignorados ou incompatíveis, dependências acrescentadas e guias manuais. Selecionar um item já confirmado não obriga a reinstalá-lo.

- Confira valores atuais e propostos das preferências e as opções de cada item.
- Veja serviços locais, pedidos de terminar sessão ou reiniciar e alterações administrativas.
- Resolva conflitos que bloqueiam a execução. Avisos também merecem leitura, mesmo quando não bloqueiam.
- Use a edição de um item a partir do plano para ajustar escolhas e voltar à revisão.
- Itens sensíveis podem exigir uma confirmação própria; a autorização de root é uma confirmação adicional do sistema.
- Para instaladores do fornecedor, reveja o conteúdo apresentado. A aplicação confere que o conteúdo executado corresponde ao conteúdo revisto.

A execução volta a validar o plano e a evidência. Se escolhas, sistema ou conteúdo relevante mudarem entretanto, pode ser necessário gerar e rever um plano novo. Uma revisão antiga não autoriza alterações diferentes.

<a id="executar-e-cancelar"></a>

## Executar, acompanhar, pausar e cancelar

Progresso por passo e efeitos de uma interrupção.

Durante a execução, o topo mostra a etapa atual, passos terminados, tempo decorrido e última mensagem. Pode navegar entre separadores; os cartões também mostram estado e passos de cada item. O Meu Ubuntu permite uma execução de cada vez.

Pausar aguarda o fim da operação atual e impede iniciar a próxima. Retomar confere o estado antes de prosseguir; alterações externas podem exigir nova revisão. Cancelar aguarda o limite seguro e não desfaz passos concluídos. Depois de fechar ou reiniciar o computador, uma execução interrompida exige um plano novo. Evite desligar enquanto um gestor de pacotes está a trabalhar.

| Resultado | Próximo passo |
| --- | --- |
| Confirmado | O teste posterior confirmou o efeito esperado. |
| Falhou | Leia a causa e o registo; corrija o problema e reveja novamente. |
| Requer intervenção | O comando terminou, mas o efeito esperado não foi confirmado; verifique o sistema. |
| Manual por concluir | Siga as instruções do guia; a execução automática não conclui essa tarefa. |

Os avisos de reiniciar ou terminar sessão referem-se aos passos que foram executados e precisam disso. Essas ações só acontecem se as pedir. Guarde o trabalho de outras aplicações antes de usar os botões correspondentes.

Ao fechar a janela durante uma execução, pode aguardar ou pedir cancelamento e aguardar o limite seguro. O fecho normal não interrompe um gestor de pacotes a meio de uma operação. Uma falha de verificação fica como intervenção, com a causa guardada no histórico.

<a id="falhas-e-retoma"></a>

## Recuperar de falhas

Use a evidência da falha e construa um plano atualizado.

1. Identifique o primeiro passo que falhou ou não foi confirmado.
2. Leia o diagnóstico apresentado e expanda o registo original. A explicação só é específica quando há evidência reconhecida.
3. Resolva a causa no sistema: ligação de rede, espaço livre, processo APT ativo, autorização recusada ou outro erro indicado.
4. Use a nova tentativa para detetar de novo e rever um plano atualizado. Passos já satisfeitos são evitados quando a deteção os confirma.
5. Se precisar de ajuda, partilhe a versão e o erro relevante. Reveja o registo antes de o publicar.

> Não apague ficheiros de lock do APT por rotina. Espere que a atualização do sistema termine e confirme o processo responsável. Cancelar ou fechar a aplicação não reverte instalações concluídas.

O histórico conserva execuções e logs. Uma tarefa manual continua manual mesmo que tenha marcado todos os passos do guia. Essas marcas são uma memória pessoal, não uma verificação automática do sistema.

<a id="historico-e-reversao"></a>

## Histórico e reversão

Reponha ajustes individuais com comparação do estado atual.

Histórico e reversão mostra execuções e alterações registadas com o valor anterior e o valor aplicado. A reversão é individual e depende do tipo de operação. Instalações ou remoções de pacotes não têm rollback automático.

1. Localize a alteração e clique em Reverter…. Também pode pedir uma reversão no cartão de um ajuste suportado e levá-la ao plano.
2. Confira o valor atual, o valor a repor, avisos e necessidade de autorização administrativa.
3. Se o valor mudou fora da aplicação, a reversão normal é recusada. Quando disponível para uma preferência, a opção de forçar exige decisão explícita.
4. Confirme a revisão e confira o resultado. Se o estado mudar durante a revisão, volte a rever.

Ficheiros administrativos têm backups protegidos. A reversão recusa ficheiros alterados entretanto ou backups não confiáveis; não existe um botão genérico para forçar a reposição desses ficheiros. Conserve histórico e backups enquanto puder precisar de reverter.

> Retirar uma aplicação das escolhas não pede reversão. Escolha Reverter… quando pretende mesmo repor um ajuste.

A revisão inclui valores atuais e os valores a restaurar. Se mudarem entretanto, gere e confirme um plano novo. Uma reversão parcial retoma apenas o que falta. Consultar uma reversão administrativa pode pedir autorização para ler o snapshot protegido; sem histórico confiável, a configuração existente é preservada.

<a id="desktop-e-estilo"></a>

## Desktop e estilo

Preferências GNOME, modelos e personalização com valores explícitos.

Estas áreas propõem preferências de janelas, dock, relógio, ficheiros, acessibilidade, cores, fontes, cursor e outras personalizações disponíveis no catálogo. Ajustes GNOME mostram o valor atual e o proposto; schemas ausentes podem tornar um item incompatível ou indisponível.

Escolher claro, escuro ou sistema no topo muda apenas o tema do Meu Ubuntu. Para mudar o desktop, escolha o ajuste correspondente no catálogo e reveja o plano. Modelos de ficheiros são criados sem substituir ficheiros pessoais existentes.

Algumas personalizações precisam de ferramentas ou passos manuais. Use os guias do item e confira a origem nos detalhes. Os ajustes são preferências; não há ganhos de desempenho medidos ou garantidos.

<a id="extensoes"></a>

## Extensões GNOME

Compatibilidade, equivalentes incluídos e ativação.

O catálogo inclui um conjunto limitado de extensões. A deteção compara UUIDs e equivalentes conhecidos para evitar duplicar funções já incluídas no Ubuntu. A compatibilidade depende da versão do GNOME e da extensão disponível.

1. Confirme a versão do GNOME em Sobre ou no rodapé.
2. Veja o UUID, a fonte e a deteção nos detalhes do item.
3. Reveja a instalação no plano; uma extensão incompatível não deve ser forçada.
4. Depois de instalar, confira a ativação e configuração numa sessão GNOME real. Use a aplicação de gestão de extensões se forem necessários ajustes.

Não é possível validar o comportamento visual de uma extensão num container sem GNOME Shell. A verificação de instalação não garante ausência de conflitos entre extensões na sua sessão.

<a id="energia-e-hardware"></a>

## Energia e hardware

Escolha ajustes segundo a máquina e a evidência do diagnóstico.

Os perfis de energia são alternativas entre si e dependem do que o power-profiles-daemon e o hardware suportam. Um perfil chamado Desempenho não promete tornar qualquer computador mais rápido. Considere consumo, calor, ruído e autonomia.

Abra Diagnóstico antes de escolher ações de disco ou drivers. Confira SSD, suporte de discard, espaço livre, bateria e sugestões com evidência. Operações como fstrim e ubuntu-drivers têm efeitos reais e devem ser revistas no plano.

> A aplicação não faz tuning geral do kernel, não troca drivers silenciosamente nem apresenta ganhos inventados. As sugestões são contextualizadas, não um benchmark.

<a id="desenvolvimento-e-terminal"></a>

## Desenvolvimento e terminal

Ferramentas, runtimes, Docker, Nala e Oh My Zsh.

Desenvolvimento reúne editores, Git, runtimes, bases de dados e clientes. Veja versões, opções e dependências no cartão e no plano. Bases de dados locais e Docker podem instalar serviços que arrancam com o sistema; os selos assinalam essa diferença.

O complemento Terminal produtivo inclui Nala, Zsh com plugins de sugestões e destaque de sintaxe, fzf, zoxide, tealdeer e um guia para Oh My Zsh. Nala é instalado como ferramenta APT adicional; o Meu Ubuntu continua a executar as suas próprias operações de pacotes pelo APT.

A instalação de pacotes não altera automaticamente a sua configuração de shell. O guia Oh My Zsh usa uma revisão fixa do instalador e preserva o .zshrc na instalação orientada. Reveja as instruções de integração, faça backup da configuração e ative só plugins que utiliza.

> Adicionar um utilizador ao grupo docker concede acesso ao daemon com poder equivalente a root. A aplicação trata essa escolha como sensível. Uma mudança de grupo pode exigir terminar sessão.

Git e perfis transportam preferências permitidas, não credenciais. Configure autenticação SSH, tokens e segredos com as ferramentas próprias. Não coloque ficheiros .env ou chaves privadas nos perfis.

<a id="infraestrutura-e-seguranca"></a>

## Infraestrutura e segurança

Clientes, diagnóstico de rede e ferramentas de laboratório.

Infraestrutura inclui clientes de acesso remoto, rede, virtualização e discos. Instalar um cliente não cria automaticamente um servidor, uma VPN, regras de firewall ou um DNS. Itens que instalam serviços locais são identificados e revistos.

Segurança inclui gestores de senhas, 2FA, cifragem, privacidade e ferramentas de auditoria. O complemento de laboratório oferece ferramentas úteis no Ubuntu, sem tentar converter o sistema em Kali Linux. A firewall é instalada, mas não é ativada automaticamente.

Ferramentas como Nmap e Wireshark são instaladas para uso consciente em sistemas e redes onde tem autorização. Configuração de capturas, credenciais, regras e serviços continua a ser uma decisão do utilizador. Consulte os detalhes e instruções de cada item.

A presença de uma ferramenta não prova que o computador está seguro. A aplicação não faz uma auditoria completa, não gere cofres de senhas e não guarda credenciais para esses programas.

<a id="diagnostico"></a>

## Diagnóstico

Factos e sugestões explicáveis antes de mudar o sistema.

Diagnóstico recolhe factos locais sobre sistema, discos, espaço, bateria e arranque. As sugestões mostram a evidência e podem apontar para um item do catálogo. Abrir o relatório não aplica essas sugestões.

Compare a sugestão com o uso da máquina antes de a selecionar. Aplicações de arranque desconhecidas não são desativadas automaticamente. Informação ausente não significa que o hardware é incapaz; pode significar que uma ferramenta ou serviço não estava disponível.

Entradas de arranque são ficheiros de configuração, não uma lista de aplicações instaladas ou em execução. O diagnóstico distingue comando disponível, executável não encontrado, entrada desativada e disponibilidade por verificar. Ver origem da entrada mostra o ficheiro e o executável referido. Uma referência antiga pode permanecer depois de remover uma aplicação ou importar configurações; o diagnóstico não a executa nem a apaga.

Em simulação, os factos são exemplos. Para avaliar a sua máquina, use o relatório no modo real e mantenha a execução separada da leitura do diagnóstico.

1. Carregue ou atualize o diagnóstico. Clique em Exportar relatório e escolha PDF, CSV ou Excel (.xlsx).
2. Selecione as secções a incluir. Todas começam selecionadas; o resumo indica os registos e avisa quando inclui caminhos de entradas de arranque.
3. Clique em Guardar relatório, escolha a pasta e mantenha a extensão indicada. Cancelar a janela não cria um ficheiro.
4. O ficheiro usa a fotografia do diagnóstico apresentada, com data, versão e modo real/simulado. Atualizar diagnóstico recolhe uma nova fotografia.

CSV usa UTF-8 e ponto e vírgula; Excel guarda texto sem fórmulas executáveis. PDF inclui acentos e paginação; caracteres fora do plano Unicode básico, como emoji, aparecem como [U+1F680]. CSV e XLSX preservam esses caracteres. O relatório contém diagnóstico, não constitui uma cópia de segurança dos seus ficheiros.

<a id="migracao"></a>

## Migração e dados pessoais

Inventário de software, Firefox e checklist de dados.

O importador de inventário reconhece listas produzidas por apt-mark showmanual, snap list e flatpak list e associa entradas ao catálogo quando possível. Reveja as correspondências; uma lista de pacotes não contém os seus dados nem restaura configurações por si.

```sh
apt-mark showmanual > pacotes-apt.txt
snap list > pacotes-snap.txt
flatpak list > pacotes-flatpak.txt
```

1. Na página Migração, detete os perfis Firefox. Feche todas as instâncias do Firefox antes de criar backup ou copiar.
2. Crie um backup do perfil. Os arquivos ficam em ~/MeuUbuntu-backups.
3. Copiar para Flatpak cria um perfil novo Migrado e faz backup do profiles.ini; não substitui um perfil existente.
4. Abra o Firefox Flatpak, escolha o perfil em about:profiles e valide marcadores, extensões e dados.
5. Só depois selecione a remoção do Firefox Snap e reveja a confirmação no plano. O snapd não é removido.

O perfil Firefox pode incluir histórico, cookies e senhas guardadas; proteja o arquivo de backup. O Brave não importa automaticamente um perfil Firefox Snap/Flatpak: use as opções do navegador para sincronizar ou exportar dados.

A checklist cobre chaveiro, navegadores, SSH, GPG, VPN, volumes Docker, projetos e documentos. As marcas são locais à interface. A aplicação não transfere automaticamente estes segredos. Não restaure toda a pasta pessoal por cima das configurações novas sem rever.

<a id="guias-manuais"></a>

## Guias manuais

Instruções copiáveis e progresso pessoal por passo.

Guias manuais aparecem nos detalhes do item, na revisão e no resultado. Alguns são tarefas totalmente manuais; outros complementam a instalação automática. A ação de copiar só coloca o comando ou bloco de configuração na área de transferência. Nunca o executa.

Leia o comando antes de o colar no terminal, especialmente quando altera configurações pessoais. Marque os passos que concluiu para retomar mais tarde. As marcas persistem neste computador e não mudam o estado detetado nem afirmam que uma instalação foi verificada.

Quando uma instrução muda, a marca antiga deixa de ser aproveitada para essa instrução. O progresso pessoal não segue no JSON exportado. Se a área de transferência falhar, a aplicação mostra o erro e pode selecionar/copiar o bloco manualmente.

<a id="permissoes-e-rede"></a>

## Permissões e ligações à rede

O que corre como utilizador, o que precisa de root e quais fontes são usadas.

| Operação | Permissão habitual |
| --- | --- |
| Consultar catálogo, perfis, manual, diagnóstico e estados | Utilizador normal. |
| Preferências GNOME, Git e ficheiros pessoais | Utilizador normal. |
| Flatpak instalado para o utilizador | Utilizador normal; depende da configuração local. |
| APT, repositórios e ajustes administrativos | Auxiliar via polkit/pkexec, com autorização do sistema. |
| Guias e scripts do fornecedor | Dependem do item; scripts automáticos do utilizador não usam root. |

O auxiliar administrativo aceita operações conhecidas e IDs/opções validados pelo catálogo, não comandos arbitrários vindos da interface. Repositórios APT usam fingerprints esperados para as chaves e ficheiros de origem próprios. Isso limita a entrada aceita; continua a haver confiança nos pacotes e fornecedores escolhidos.

Instalações podem ligar ao arquivo Ubuntu, repositórios oficiais de fornecedores, Flathub, serviço de extensões GNOME e fontes declaradas pelos instaladores. O manual e os logos não precisam de rede. Abrir uma ligação de documentação usa o navegador do sistema.

Não há conta, telemetria ou serviço de sincronização do Meu Ubuntu. Gestores de pacotes, aplicações instaladas e fornecedores podem ter os seus próprios comportamentos de rede e privacidade. A simulação não reproduz todas essas condições reais.

<a id="dados-e-privacidade"></a>

## Dados locais e privacidade

Pastas XDG, logs, backups e informação para suporte.

| Local padrão | Conteúdo |
| --- | --- |
| ~/.config/meu-ubuntu/ | prefs.json com escolhas, tema, vista do catálogo e progresso dos guias; profiles/ com receitas guardadas e imported/. |
| ~/.local/state/meu-ubuntu/ | changes.json, runs/, logs/ e backups/ de configurações pessoais antes de scripts. |
| ~/.local/share/meu-ubuntu/ | Dados da aplicação; em simulação, home e raiz simuladas. |
| /var/lib/meu-ubuntu/ e /var/backups/meu-ubuntu/ | Snapshots e backups administrativos protegidos para os fluxos suportados. |
| ~/MeuUbuntu-backups/ | Backups de perfis Firefox e ficheiros relacionados à migração. |
| Dados da WebView | Marcas da checklist de migração guardadas em localStorage. |

XDG_CONFIG_HOME, XDG_STATE_HOME e XDG_DATA_HOME podem mudar estes caminhos. Em simulação, as pastas XDG usam meu-ubuntu-sim e ficam separadas dos dados reais. Backups de ficheiros do sistema são geridos pelo auxiliar com permissões protegidas, não pela interface.

Perfis não guardam segredos. Registos de comandos e backups podem revelar caminhos, nomes ou dados pessoais; leia-os antes de partilhar. Um backup Firefox pode conter dados sensíveis do navegador. Apagar histórico remove a evidência para reverter mudanças antigas.

Sobre fica sempre acessível no rodapé da barra lateral, junto da versão em execução, mesmo com Explorar por área expandido. Permite copiar informações técnicas para suporte: versão, distribuição, arquitetura, GNOME, sessão, modo e ferramentas de build. Esse resumo não inclui as suas escolhas, caminhos pessoais ou logs. Pode sempre ler o resumo antes de o copiar.

<a id="simulacao-e-terminal"></a>

## Simulação e opções de terminal

Conheça os fluxos e consulte a documentação sem abrir a janela.

| Comando | Efeito |
| --- | --- |
| meu-ubuntu | Abre a aplicação no modo real. |
| meu-ubuntu --simular | Abre com sistema e executor simulados; não instala nem altera os ajustes reais. |
| MEU_UBUNTU_SIMULATE=1 meu-ubuntu | Alternativa à flag de simulação. |
| meu-ubuntu --perfil ./perfil.json | Abre uma revisão do perfil; não o executa. |
| meu-ubuntu --version | Mostra a versão do executável. |
| meu-ubuntu --help | Mostra as opções disponíveis. |
| meu-ubuntu --manual | Imprime o manual incluído no executável. |
| meu-ubuntu --licenses | Imprime a GPLv3 do projeto e as licenças dos logos, fonte e biblioteca PDF incluídos. |

A simulação guarda preferências e histórico próprios para permitir testar os fluxos. A interface identifica sempre esse modo. Não valida bibliotecas reais, autenticação polkit, compatibilidade do hardware, extensões em sessão ou acesso aos fornecedores.

--helper é um protocolo interno da aplicação através do pkexec. Não é um modo de administração geral nem uma opção para usar manualmente. A aplicação gráfica recusa iniciar como root.

<a id="problemas-frequentes"></a>

## Problemas frequentes

Respostas para dúvidas de estado, perfis, execução e compatibilidade.

| Sintoma ou dúvida | O que fazer |
| --- | --- |
| Está instalado, mas não está marcado | Clique em Atualizar estado. A deteção marca os instalados automaticamente, exceto exclusões explícitas, reversões pedidas, alternativas já escolhidas e perfis em edição. Pode voltar a marcar os que excluiu em Mais ações do catálogo. |
| Mudei de perfil e as escolhas antigas ficaram | O select filtra a vista. Use Usar perfil… e Substituir na revisão se quiser trocar as escolhas. |
| Não vejo um item | Abra Todo o catálogo, limpe a pesquisa e os filtros e, se necessário, ative Mostrar opções avançadas. |
| O plano pede uma nova revisão | O estado ou a proposta mudou. Gere o plano novamente e confira as diferenças. |
| A aplicação não abre | Confira Ubuntu 26.04 amd64, sessão desktop, bibliotecas de execução e caminho local executável; abra pelo terminal para ver o erro. |
| A janela de senha foi recusada | A execução precisa da autorização do sistema; reveja o erro e volte a tentar quando puder autorizar. |
| Uma instalação terminou, mas requer intervenção | Leia a verificação que falhou; não presuma que o programa ficou utilizável. |
| Posso usar no Fedora? | Ainda não. O backend APT e os ajustes estão preparados para Ubuntu; a portabilidade exige implementação e testes próprios. |
| Posso desfazer tudo? | A reversão cobre ajustes suportados. Pacotes instalados ou removidos não têm rollback automático. |
| Onde obter ajuda? | Consulte o capítulo relevante, copie o resumo técnico em Sobre e acrescente o erro revisto quando reportar o problema. |

Se uma funcionalidade deixou de funcionar ao atualizar, compare um plano pequeno com o que funcionava na versão anterior. Registe o ambiente, a versão e o resultado esperado versus observado para tornar o problema reproduzível.

<a id="pacotes-de-integracao"></a>

## Pacotes de integração

Diagnóstico e preparação revista para uma pasta de projeto local.

Em Pacotes de integração, escolha Instalar pacote… e selecione o ZIP recebido de uma origem que conhece. Reveja nome, versão e fornecedor antes de confirmar. Instalar o pacote acrescenta o separador e as instruções do projeto.

1. Abra o separador instalado e escolha a pasta onde já clonou o projeto.
2. Carregue em Verificar ambiente. Em simulação, os resultados são exemplos; no modo real, as verificações consultam os requisitos locais.
3. Se o pacote inclui Preparar ambiente, reveja as ferramentas necessárias e aplique o plano. Depois reveja a preparação do projeto e confirme. Pacotes só de guia mantêm execução manual dos comandos.
4. A preparação mostra cada etapa e admite pausa/retoma entre operações. Configurações originais, serviços alheios e dados existentes são preservados. Após uma interrupção, reveja novamente para conferir o estado real.
5. Ao atualizar o pacote, a pasta escolhida é preservada. Se o ZIP mudar depois da revisão, escolha-o e reveja-o novamente.

> Use um checkout de desenvolvimento em que confia: Docker e Composer executam código desse projeto. O pacote pode preparar containers isolados, schemas, Redis e dependências. Dumps, credenciais privadas e comandos administrativos continuam manuais. Containers ativos não comprovam que o projeto já funciona com os dados da equipa.

<a id="contribuir-e-licencas"></a>

## Contribuir, licenças e limites

Estado do projeto e como ajudar com evidência útil.

O código do Meu Ubuntu é distribuído sob GNU GPL versão 3 (GPL-3.0-only), incluída em LICENSE. Logos, dependências e a fonte do PDF conservam as suas licenças próprias. Consulte NOTICE.md, Sobre e os avisos do pacote. A candidata inclui um arquivo com o código-fonte correspondente; o repositório público está em https://github.com/mercynet/my-ubuntu. Vulnerabilidades devem ser comunicadas em privado pelo separador Security, opção Report a vulnerability.

Para contribuir, leia docs/desenvolvimento.md e docs/catalogo.md no código-fonte. Descreva problemas com passos, ambiente e resultado, e proponha receitas com fontes oficiais, método de deteção, verificação posterior, dependências e riscos concretos. Não inclua segredos em exemplos ou registos.

A 1.0.0-rc.1 inclui correções de recuperação e reversão verificadas com regressões e revisão adversarial. A atualização de uma instalação anterior, autenticação polkit e os fluxos GNOME/Wayland ainda precisam de testes em desktops reais antes da 1.0 estável.

Fedora, outras arquiteturas, AppImage e um catálogo maior precisam de benefício demonstrado e validação. Nesta fase, melhorar previsibilidade, recuperação e documentação é mais útil do que prometer suporte não testado.

<a id="agentes-de-ia"></a>

## Instalar agentes de IA

Instalar a ferramenta e configurar o acesso são etapas separadas.

Agentes de IA reúne Claude Code, Codex, Gemini CLI, OpenCode, Pi e GitHub Copilot CLI. Escolha apenas os que pretende usar. A instalação usa o pacote oficial do fornecedor com versão fixa, para o seu utilizador, em ~/.local.

1. Abra Agentes de IA no menu e escolha uma ou mais ferramentas.
2. Reveja o plano. A dependência Node.js/npm é apresentada antes da instalação, se necessária.
3. Confirme a instalação do código do fornecedor. O npm usa o registo oficial e verifica os requisitos de versão.
4. Aguarde a verificação do comando com --version. Instalações existentes são preservadas.
5. Abra Como começar depois de instalar no cartão do agente e siga o guia no seu terminal. O início de sessão ocorre dentro do próprio agente.

> O Meu Ubuntu não guarda contas, chaves nem subscrições dos agentes. A instalação da ferramenta não inclui créditos nem garante acesso a modelos; o fornecedor define as condições de acesso.

Se ~/.local/bin ainda não estava no PATH da sua sessão, use o caminho completo indicado no guia ou termine e volte a iniciar sessão no Ubuntu. Não é necessário usar sudo para instalar os pacotes dos agentes.

<a id="flatpak-e-origens"></a>

## Flatpak e origens de instalação

Requisitos, verificação e conservação de aplicações já presentes.

Aplicações com publicação ativa confirmada usam Flatpak para o seu utilizador. Cada cartão mostra o requisito e a origem. Quando falta Flatpak, o plano inclui essa dependência antes da aplicação.

Antes de instalar, a execução confere Flatpak, o remoto do utilizador, origem permitida, verificação de assinaturas e ID da aplicação. Remoto incorreto, desativado, rede indisponível ou aplicação não encontrada impedem a instalação. Não muda silenciosamente para outra origem.

Versões APT e IDs Flatpak anteriores já instalados são reconhecidos e preservados. Esta política aplica-se a novas instalações; não remove programas nem migra dados e perfis automaticamente.

> Flatpak isola a aplicação. Algumas ferramentas precisam de SDKs, permissões de dispositivo ou serviços no host. Consulte o guia da receita. Um editor não verificado no Flathub fica identificado na origem; instaladores do próprio fornecedor também são identificados.
