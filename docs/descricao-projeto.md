# Meu Ubuntu

O Meu Ubuntu é uma central visual de pós-instalação para Ubuntu 26.04 com GNOME e amd64. Reúne aplicações, desenvolvimento, infraestrutura, segurança, desktop, estilo, extensões e energia num catálogo declarativo. As escolhas são revistas num plano antes de serem executadas, e o resultado é verificado por passo.

Usa Go, Wails e React/TypeScript. Interface, catálogo, perfis, ajuda e logos estão incluídos no executável. Não exige Go ou Node no computador de quem o utiliza. É distribuído como binário e pacote .deb com bibliotecas declaradas.

O select de perfis no topo organiza a exploração sem alterar escolhas. Usar perfil… abre uma revisão para juntar ou substituir; Todas… mostra o catálogo completo. Estado instalado, intenção do perfil e alterações pendentes são informações distintas. Inclui guias manuais, histórico, reversões específicas, diagnóstico e migração orientada do Firefox.

Consulte a documentação conforme o que pretende fazer:

- [Manual completo](manual.md): instalação, utilização, perfis, execução, erros, reversão, permissões, dados e FAQ. Disponível também em Ajuda e manual e por `meu-ubuntu --manual`.
- [Catálogo e perfis](catalogo.md): todas as receitas, opções, fontes, dependências e verificações.
- [Desenvolvimento](desenvolvimento.md): arquitetura, contratos, receitas, testes, segurança, build e lançamento.
- [README](../README.md): início rápido, evolução e evidência das verificações.
- [Pesquisa](pesquisa.md): referências e decisões técnicas anteriores.

O utilizador confirmou uso real da 0.3. Versões posteriores precisam de validar atualização e fluxos alterados no mesmo ambiente. Containers e testes de interface não substituem uma sessão GNOME real. Não há rollback geral de pacotes nem ganhos de desempenho prometidos. Fedora e outras arquiteturas ainda não têm suporte validado.

O código usa GNU GPLv3 (GPL-3.0-only); a candidata inclui a licença, avisos de terceiros e código-fonte correspondente. A publicação pública e a versão 1.0 estável dependem da validação descrita no README.
