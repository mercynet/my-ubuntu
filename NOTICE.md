# Licenças e atribuições

O código original do Meu Ubuntu é distribuído sob **GNU GPL versão 3**, identificador `GPL-3.0-only`. O texto completo está em [LICENSE](LICENSE). O arquivo `meu-ubuntu-<versão>-source.tar.gz` contém o código-fonte correspondente aos artefactos entregues, incluindo instruções e scripts de build.

Componentes de terceiros conservam as suas licenças:

- Logos Devicon e Simple Icons: avisos e fontes em `frontend/src/assets/apps/`.
- Biblioteca PDF go-pdf/fpdf: MIT, aviso em `internal/diag/assets/FPDF-LICENSE.txt`.
- Fonte DejaVu Sans: aviso em `internal/diag/assets/DejaVu-LICENSE.txt`; a licença da fonte difere dos avisos de empacotamento Debian incluídos nesse ficheiro.
- Bibliotecas Go, runtime React e Lucide: o build recolhe os avisos das versões resolvidas em `dist/third-party/`, com índice `modules.json`.

O pacote instala estes documentos em `/usr/share/doc/meu-ubuntu/`. `meu-ubuntu --licenses` apresenta a licença do projeto e os avisos dos assets embutidos; Sobre também permite consultá-los. Bibliotecas de sistema GTK/WebKit são fornecidas pelo sistema, sob as suas próprias licenças. Programas instalados pelo catálogo não passam a usar a licença do Meu Ubuntu.
