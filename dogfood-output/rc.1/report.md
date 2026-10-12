# Meu Ubuntu 1.0.0-rc.1: validação

Data: 2026-10-12. Alvo: Ubuntu 26.04 LTS, GNOME, amd64. Candidata local para testes; publicação externa e versão estável não foram feitas. O pacote instalado no computador do utilizador permanece intacto.

## Resultado

Correções de recuperação/reversão, persistência de intervenções, verificação pós-instalação e barreira de fecho normal foram implementadas e protegidas por regressões. A [revisão Astra/high final](review-astra-high.md) não identificou bloqueadores de código remanescentes nos cenários revistos. As descobertas adicionais dessa revisão também foram corrigidas.

PDF, CSV e Excel (.xlsx) exportam a fotografia do diagnóstico apresentado e as secções escolhidas. O projeto usa GPL-3.0-only; a licença, os avisos de terceiros e o código-fonte correspondente acompanham os artefactos.

## Checks executados

- Script de `make check` num container Ubuntu 26.04: documentação gerada, nove testes Node, TypeScript/Vite, Go vet e todos os testes Go com `-race` aprovados.
- Auditoria Go/npm: sem vulnerabilidades reportadas pelas ferramentas nessa execução.
- Regressões dirigidas: reversão parcial, mudança após revisão, snapshot administrativo alterado, journald/reload/reaplicação, Xdebug/configuração, dependências após cancelamento, intervenção persistida, pertença Docker, Git e diretório de templates alterado.
- Exportações com leitores independentes: `csv.reader`, `openpyxl` e `pdftotext`; acentos, grego, emoji, texto sem fórmulas e coerência entre tabelas conferidos. PDF representa emoji como `[U+1F680]`.
- Após o QA visual, ajuste de paginação para manter o título de uma secção com o início do seu primeiro registo; regressões de diagnóstico e ponte repetidas com `-race -count=1` e aprovadas.
- Build e verificação: binário, metadados Debian `1.0.0~rc1`, ordem anterior à versão final, manuais, licença/avisos, arquivo de fontes e checksums aprovados. Avisos recolhidos para 119 dependências resolvidas.
- Código-fonte extraído numa pasta temporária, sem assets compilados nem dependências locais: documentação, frontend e executável recompilados com sucesso. Evidência em `source-build.log`.

O teste adicional de instalação em imagem Ubuntu 26.04 mínima concluiu o processamento APT, incluindo `Setting up meu-ubuntu (1.0.0~rc1)`, sem erros APT/dpkg registados. O script de QA terminou com código 1 num check posterior. A imagem tem `path-exclude=/usr/share/doc/*`, o que exclui os avisos instalados que o script exige; os avisos foram conferidos dentro do `.deb`. Não se declara esse script aprovado nem se confunde esse resultado com instalação/atualização num desktop. A limitação e o log ficam registados; não foi feita outra instalação extensa para repetir o check de documentação.

## Janela nativa

Wails/WebKitGTK em Xvfb, container isolado, utilizador UID 1000 e pastas HOME/XDG temporárias, executor simulado. Não houve receitas privilegiadas no host nem uso de dados pessoais reais.

- Abrir Diagnóstico, escolher e guardar PDF, CSV e XLSX pelo diálogo GTK real.
- Abrir os ficheiros guardados com leitores independentes: CSV e XLSX têm os mesmos 49 registos; PDF apresenta versão e modo de simulação. Ficheiros gravados com permissões 0600. Evidência em [exports/native/validation.json](exports/native/validation.json).
- Cancelar o diálogo e continuar a usar a aplicação.
- Escolher extensão `.csv` ao exportar `.xlsx`, confirmar a substituição de um ficheiro de teste existente: a aplicação recusa a extensão e preserva o conteúdo original. Captura: [native-extension-refused.png](screenshots/native-extension-refused.png).
- Abrir Sobre e Diagnóstico numa janela 900 × 600: versão/Sobre permanecem acessíveis e o conteúdo permite deslocamento.

Este ambiente tem X11 virtual e não valida a sessão GNOME/Wayland, autenticação polkit nem hardware real. As indicações de sistema na aplicação são fixtures da simulação.

O binário final, após o ajuste de paginação, voltou a guardar os três formatos pela janela GTK; os leitores repetiram a validação dos 49 registos e Sobre abriu a 900 × 600. O PDF foi também inspecionado visualmente.

## Interface no navegador

`agent-browser`, sessão própria, assets de produção copiados sem alteração e ponte Wails explicitamente simulada. Janela 900 × 600, tema claro. Conferidos: formatos, exclusão da secção de arranque no pedido, impedir exportação sem secções, erro injetado visível, cancelamento sem erro, libertação do estado ocupado, Sobre/GPLv3, resumo técnico e ausência de deslocamento horizontal. Isto valida interface/pedidos, não escrita real nem o backend.

## Antes da 1.0 estável

Testar atualização da instalação anterior sem perder dados, polkit autorizado/cancelado, fecho e recuperação durante operações numa VM GNOME/Wayland, receitas dependentes de hardware e feedback de utilizadores Ubuntu 26.04. Seguir [o roteiro de teste](../../docs/beta-usabilidade.md). Não atribuir aos checks simulados essas validações reais.

O repositório público e o canal privado de segurança ainda precisam de ser anunciados antes da publicação pública. A integração privada OfficeGest é uma avaliação separada e não acompanha esta release.
