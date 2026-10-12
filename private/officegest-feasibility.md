# Integração privada OfficeGest: viabilidade

Avaliação por leitura do código em `../ogdevel`, em 2026-10-12. Não foram iniciados serviços, modificados ficheiros desse projeto nem consultados dados de tenants. Ainda não há uma extensão implementada.

## Conclusão

É viável distribuir um pacote separado para a equipa que, depois de instalado, faça aparecer um separador OfficeGest. O Meu Ubuntu atual tem catálogo/perfis embutidos e navegação fixa; importar um perfil apenas seleciona receitas existentes. Precisamos de um mecanismo de extensões antes de oferecer instalação de separadores.

Recomendação: manter uma interface de extensões genérica no projeto público e um pacote OfficeGest num repositório privado. Um executável separado pode fornecer operações limitadas e resultados estruturados; o pacote deve declarar versão, compatibilidade e origem verificável. Dados, credenciais e dumps não acompanham o núcleo público. Executar código de uma extensão exige instalação explícita e confiança na sua origem; um manifesto não ganha acesso ao helper root nem aceita comandos arbitrários vindos da interface.

## Obstáculos encontrados no Docker

| Evidência | Consequência para instalação nova |
|---|---|
| `docker-compose.yml`: `shared-network` externa | Precisa de existir antes do arranque; `dock up` não a cria. |
| Worker usa `ogdevel_app:latest`; app não declara essa imagem | Dependência do nome do projeto/pasta e da imagem construída; precisa de normalização. |
| Nomes fixos, como `php-app`, e várias portas publicadas | Podem colidir com outro checkout ou serviços locais. |
| UID/GID configuráveis e bind mounts | O instalador deve conferir identidade/permissões; o entrypoint atual aplica `chmod 777` em pastas, o que merece revisão. |
| Vários serviços sem perfis Compose | Arranque normal inclui também cluster Redis, Elasticsearch/Kibana e correio de teste; faltam modos mínimo/completo claros. |
| Criador de tenant interativo e clonagem de schema | Ainda falta um fluxo reproduzível com dados de demonstração sanitizados e configuração de domínio/certificado. |
| `docker/README.md` e Compose divergem em portas/serviços | A interface deve usar o estado/configuração efetivos, com documentação alinhada. |

São obstáculos demonstrados na configuração, não diagnósticos das falhas concretas dos colegas: ainda não houve reprodução numa instalação limpa.

## Fluxo proposto

1. Escolher pasta existente ou obter o repositório privado, usando a autenticação Git do colega.
2. Verificar Docker/Compose, acesso ao daemon, espaço, memória, portas, UID/GID, rede e ficheiros necessários.
3. Rever as operações e escolher ambiente mínimo ou completo.
4. Preparar configuração local sem sobrepor ficheiros existentes; obter dependências bloqueadas pelas versões do projeto e construir/iniciar os serviços escolhidos.
5. Configurar um tenant de desenvolvimento com dados sanitizados e validar PHP, ligação à base, cache, filas e resposta HTTP.
6. Mostrar estado, iniciar/parar apenas este ambiente, abrir a aplicação e exportar diagnóstico expurgado para suporte.

Cada etapa precisa de verificação e retoma. Volumes existentes são preservados; reset/importação de base são ações separadas com identificação explícita do destino e backup. Não reutilizar os comandos de troubleshooting que param todos os containers ou limpam redes globais.

## Entrega e validação

Primeiro protótipo: Ubuntu 26.04, uma pasta já clonada e diagnóstico sem alterações. Depois, testar preparação completa em VM limpa, segunda execução sem duplicações, portas ocupadas, falta de acesso privado, falha de download e retoma. Uma segunda máquina da equipa deve concluir sem ajuda antes de distribuir o instalador.

O pacote pode permanecer privado para uso interno. A GPLv3 permite modificações internas sem publicação pública; fornecimento a terceiros e plugins que formem uma obra combinada têm requisitos próprios. Não assumir que separar ficheiros muda a licença: manter uma fronteira de programas independentes e confirmar os termos antes de distribuição fora da organização. Referência: [FAQ GNU](https://www.gnu.org/licenses/gpl-faq.en.html#InternalDistribution).

Esta integração fica fora da RC atual. O próximo passo útil é acordar o ambiente mínimo e recolher os erros de duas instalações novas; isso orienta o instalador antes de construir a página.
