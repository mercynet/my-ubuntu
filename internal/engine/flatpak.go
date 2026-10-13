package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const flathubRepo = "https://dl.flathub.org/repo/"

func (e *Engine) requireFlatpak(ctx context.Context) error {
	if !e.R.LookPath("flatpak") {
		return errors.New("Flatpak não está instalado; instale Flatpak + Flathub no plano antes desta aplicação")
	}
	out, err := e.R.Output(ctx, "flatpak", "--version")
	if err != nil || !strings.HasPrefix(strings.TrimSpace(out), "Flatpak ") {
		return errors.New("Flatpak não respondeu corretamente; confira a instalação antes de continuar")
	}
	return nil
}

func (e *Engine) checkFlatpakRemote(ctx context.Context, name, expected string, required bool) error {
	out, err := e.R.Output(ctx, "flatpak", "remotes", "--user", "--show-disabled", "--columns=name,url,options")
	if err != nil {
		return errors.New("não foi possível validar os remotos Flatpak do utilizador")
	}
	for _, line := range strings.Split(out, "\n") {
		columns := strings.Split(line, "\t")
		if len(columns) < 2 || strings.TrimSpace(columns[0]) != name {
			continue
		}
		actual := strings.TrimRight(strings.TrimSpace(columns[1]), "/")
		allowed := actual == strings.TrimRight(expected, "/")
		if name == "flathub" {
			allowed = allowed || actual == "https://flathub.org/repo"
		}
		if !allowed {
			return fmt.Errorf("o remoto %s aponta para outra origem; confira-o antes de instalar", name)
		}
		if len(columns) > 2 {
			options := strings.FieldsFunc(columns[2], func(r rune) bool { return r == ',' || r == ' ' })
			for _, option := range options {
				if option == "no-gpg-verify" || option == "disabled" {
					return fmt.Errorf("o remoto %s está desativado ou sem verificação de assinaturas; corrija-o antes de instalar", name)
				}
			}
		}
		return nil
	}
	if required {
		return fmt.Errorf("o remoto %s não está configurado para o utilizador; reveja Flatpak + Flathub", name)
	}
	return nil
}

func (e *Engine) checkFlatpakApp(ctx context.Context, app string) error {
	ref, err := e.R.Output(ctx, "flatpak", "remote-info", "--user", "--show-ref", "flathub", app)
	parts := strings.Split(strings.TrimSpace(ref), "/")
	if err != nil || len(parts) != 4 || parts[0] != "app" || parts[1] != app {
		return errors.New("a aplicação não foi confirmada no Flathub; confira ID, rede e disponibilidade antes de tentar novamente")
	}
	return nil
}
