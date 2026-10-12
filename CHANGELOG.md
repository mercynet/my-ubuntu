# Changelog

## 1.0.0-rc.1 — 2026-10-12

- Exportação do diagnóstico em PDF, CSV e Excel (.xlsx), com seleção de secções, data e versão da fotografia apresentada.
- Recuperação de reversões parciais sem repetir chaves restauradas; planos vinculados aos valores atuais e alvos, incluindo snapshots administrativos.
- Nova tentativa de configuração/verificação sem reinstalar pacotes já presentes; falhas de verificação persistem como intervenção no histórico.
- Correção da deteção de grupos Docker, preferências Git ausentes, Xdebug, journald e drivers recomendados.
- Fecho normal protegido durante execução, com cancelamento no limite seguro.
- GPLv3, avisos de terceiros, guias de contribuição/segurança e arquivo do código-fonte correspondente.

Candidata para Ubuntu 26.04 LTS, GNOME e amd64. Testes GNOME/Wayland, polkit, atualização e feedback externo ainda precedem a 1.0 estável. Veja a [revisão adversarial](dogfood-output/rc.1/review-astra-high.md).
