# Meu Ubuntu 0.6.0-beta.5 — acesso a Sobre

Data: 2026-10-12. Alvo da aplicação: Ubuntu 26.04 LTS, GNOME, amd64.

## Alteração e resultado

Sobre e a versão do executável ficam no rodapé da barra lateral. O menu tem deslocamento independente, pelo que expandir Explorar por área não oculta esse acesso numa janela de 900 × 600. A página identifica a beta, o autor, a natureza gratuita da aplicação e os limites da licença principal ainda pendente. Mantém o manual, as atribuições e a cópia de informações técnicas.

## Verificação

- `check.log`: documentação gerada, testes Node, TypeScript/Vite, Go vet e todos os testes Go com deteção de concorrência passaram num container Ubuntu 26.04.
- `build.log` e `verify.log`: binário e .deb compilados; versão, checksums, metadados e documentação do pacote conferidos.
- Chrome com assets de produção e ponte de teste: menu expandido em 900 × 600, rodapé visível, deslocamento independente e ausência de overflow horizontal. Versão 0.6.0-beta.5, duas atribuições, abertura do manual e acesso por teclado confirmados.
- Cópia na ponte de teste corresponde exatamente ao resumo apresentado. Recusa da área de transferência mostra erro; falha ao ler Sobre desativa a cópia e permite nova tentativa. Sobre não chamou operações de execução na ponte. Sem erros JavaScript observados.
- Binário Wails real em GTK/WebKit/Xvfb, com diretórios de utilizador isolados e `--simular`: Sobre acessível em 900 × 600; página e versão correta em 1280 × 840; Go 1.26.9/Wails v2.16.0 e sucesso da cópia apresentados. Capturas em `screenshots/`.

## Evidência de uso real anterior

O utilizador relata instalação real bem-sucedida de ferramentas neste PC com uma versão anterior. Leituras em 2026-10-12 confirmam pacote e executável instalados 0.3.2, Ubuntu 26.04.1 LTS, GNOME 50.1 e sessão Wayland. O relato não especifica todos os cenários nem a subversão usada historicamente. O pacote instalado não foi atualizado nesta revisão.

## Limites

Chrome usa uma ponte fictícia; não valida o backend. Xvfb valida a janela Wails em simulação e não equivale a GNOME/Wayland com polkit e instalações reais. O contrato backend de Sobre tem teste próprio de versão, atribuições e ausência de alterações a preferências/histórico. A revisão adversarial em paralelo identificou problemas de recuperação e reversão que precisam de resolução antes da 1.0; estes artefactos não são uma candidata final.
