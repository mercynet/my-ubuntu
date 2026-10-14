# Meu Ubuntu

[![Check and build](https://github.com/mercynet/my-ubuntu/actions/workflows/check.yml/badge.svg?branch=development)](https://github.com/mercynet/my-ubuntu/actions/workflows/check.yml)
[![License: GPLv3](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

Aplicação desktop para preparar um Ubuntu pessoal: escolher ferramentas e ajustes,
rever um plano, aplicar com autorização do sistema, verificar o resultado e guardar
perfis reutilizáveis. Interface Go/Wails + React/TypeScript, recursos e ajuda offline.

**Versão:** 1.0.0. **Escopo suportado:** Ubuntu 26.04 LTS, GNOME, amd64.
Outras distribuições, arquiteturas e ambientes de hardware específicos não são
promessas desta versão.

![Revisão do plano](docs/screenshots/plano.png)

## O que inclui

- Catálogo de aplicações, ferramentas, ajustes GNOME, extensões e perfis.
- Preferência por Flatpak para aplicações com publicação ativa verificada;
  requisito visível e validação de Flatpak, remoto e ID antes da instalação.
- Plano revisto antes de agir, progresso, pausa/retoma entre operações e histórico.
- Reversão individual de ajustes registados e recuperação de trabalho interrompido.
- Diagnóstico e exportação em PDF, CSV e Excel; ajuda e atribuições offline.
- Pacotes opcionais de integração com diagnóstico e preparação declarativa revista.

Aplicações já presentes por APT ou IDs Flatpak anteriores são preservadas.
Ferramentas CLI, motores e serviços que precisam do host continuam com receitas
próprias. Contas, dumps e configuração de terceiros não são recolhidos automaticamente.

## Instalar e testar

Obtenha o instalador nos [lançamentos](https://github.com/mercynet/my-ubuntu/releases)
quando publicado, ou os artefactos da execução aprovada em
[Actions](https://github.com/mercynet/my-ubuntu/actions/workflows/check.yml).
O instalador reduzido contém .deb, checksum e instruções. Código-fonte e avisos
acompanham a distribuição completa e são disponibilizados separadamente.

Extraia o ZIP, abra um terminal nessa pasta e execute:

```sh
sha256sum -c SHA256SUMS
sudo apt install ./meu-ubuntu_1.0.0_amd64.deb
meu-ubuntu --simular
```

Explore primeiro a simulação. Depois abra `meu-ubuntu` como utilizador normal;
as operações administrativas pedem autorização pelo polkit após revisão.
A simulação não instala ferramentas nem prova comportamento numa sessão real.
Go, Node e npm não são necessários para executar o pacote.

## Desenvolver

No host, instale Docker e Make. A compilação usa um container Ubuntu 26.04 com
toolchain fixo, limitado por padrão a 2 CPUs e 4 GiB:

```sh
make check
make audit
make build
make verify
make sim
```

Os resultados ficam em `dist/`; dependências e frontend compilado são gerados
localmente. Esses ficheiros não são versionados. A versão está em [VERSION](VERSION);
o build recusa divergências nos manifests frontend.
Veja [CONTRIBUTING.md](CONTRIBUTING.md) e [desenvolvimento](docs/desenvolvimento.md).

## Documentação e suporte

| Documento | Conteúdo |
|---|---|
| [Manual](docs/manual.md) | Utilização, perfis, execução, reversão e diagnóstico |
| [Catálogo](docs/catalogo.md) | Receitas, fontes, opções e verificações |
| [Flatpak](docs/flatpak.md) | Política, requisitos, validação e limites |
| [Índice técnico](docs/README.md) | Arquitetura, testes, contexto e investigação |
| [Changelog](CHANGELOG.md) | Alterações por versão |
| [Releases](RELEASING.md) | Versionamento, critérios e publicação |
| [Aceitação da 1.0](docs/release-1.0.md) | Provas concluídas e trabalho restante |
| [Beta](docs/beta-usabilidade.md) | Roteiro para utilizadores de teste |

Reporte bugs nos [issues](https://github.com/mercynet/my-ubuntu/issues), com versão,
ambiente, passos e resultado observado. Vulnerabilidades seguem o
[canal privado](SECURITY.md); remova dados pessoais de anexos e logs.

## Escopo e limites da 1.0.0

Esta versão cobre o uso normal da aplicação no ambiente suportado: seleção,
revisão, execução autorizada, verificação, histórico, perfis, diagnóstico e
reversão dos ajustes suportados. Foi exercitada em dois computadores reais.

Receitas que dependem de hardware, dispositivos, sessão ou permissões externas
podem permanecer manuais ou indisponíveis. O addon OfficeGest é opcional e não
faz parte da promessa de instalação/preparação automática da aplicação; o seu
ambiente completo exige validação própria.

[Contexto histórico](docs/contexto-projeto.md) e [evidência de QA](dogfood-output/README.md)
preservam resultados anteriores sem os apresentar como validação da candidata atual.

## Licença

Código GPL-3.0-only: [LICENSE](LICENSE). Assets e dependências mantêm as licenças
e avisos em [NOTICE.md](NOTICE.md). Não existe telemetria da aplicação.
