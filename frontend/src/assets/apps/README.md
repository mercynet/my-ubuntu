Os logos são ficheiros locais: o catálogo nunca contacta um CDN para os carregar.

- Devicon v2.16.0: https://github.com/devicons/devicon/tree/v2.16.0 — licença MIT incluída.
- Simple Icons 16.13.0: https://github.com/simple-icons/simple-icons/tree/16.13.0 — licença CC0 incluída. As cores são as publicadas no respetivo catálogo.

`sources.json` regista a origem, hash, referências da marca e eventuais licenças próprias indicadas pelo upstream. Os nomes e logos identificam os produtos; não indicam associação ou patrocínio. As marcas continuam a pertencer aos seus titulares.

Para atualizar conscientemente este conjunto: `python3 build/scripts/update-app-logos.py`. O script usa versões fixas, valida os SVG e só descarrega os logos mapeados. Não é executado durante o build normal. Aplicações sem logo e preferências do sistema usam os ícones Lucide já existentes.
