# Evidência de QA

Relatórios, logs e capturas revistos pertencem à versão indicada em cada pasta.
`rc.2/` e `rc.3/` preservam provas anteriores; `rc.4/` regista a política Flatpak
e organização do repositório. Consulte [aceitação](../docs/release-1.0.md).

Cada relatório deve identificar revisão/árvore, versão, checksum do executável,
ambiente, passos, resultado e alcance. Use PASS, FAIL, NOT RUN ou INCONCLUSIVE;
uma prova simulada ou em container não equivale a GNOME/Wayland e polkit reais.

Estados XDG, caches, casas de teste, dados de projetos, discos/chaves de VM e
pastas temporárias ficam locais e são ignorados. Reveja dados pessoais antes de
versionar logs ou capturas. Artefactos de distribuição ficam em `dist/`, Actions
ou Releases. Materiais antigos que já entraram no histórico permanecem públicos;
ignorar uma pasta não apaga esse histórico.
