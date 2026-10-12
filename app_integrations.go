package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"meuubuntu/internal/integrations"
)

type integrationReview struct {
	Manifest integrations.Package `json:"manifest"`
	Token    string               `json:"token"`
	Existing string               `json:"existing"`
	path     string
	created  time.Time
}
type integrationPending struct {
	sync.Mutex
	review *integrationReview
}

func (a *App) integrationManager() integrations.Manager {
	return integrations.Manager{Data: a.st.Data, Config: a.st.Config, Simulated: a.simulated}
}
func (a *App) Integrations() []integrations.View { return a.integrationManager().List() }
func (a *App) PreviewIntegrationImport() (*integrationReview, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Instalar pacote de integração", Filters: []runtime.FileFilter{{DisplayName: "Pacote de integração", Pattern: "*.zip"}}})
	if err != nil || path == "" {
		return nil, err
	}
	p, token, err := integrations.ReadArchive(path)
	if err != nil {
		return nil, err
	}
	r := &integrationReview{Manifest: p, Token: token, path: path, created: time.Now()}
	if old, err := a.integrationManager().Load(p.ID); err == nil {
		r.Existing = old.Version
	}
	a.integrationPending.Lock()
	a.integrationPending.review = r
	a.integrationPending.Unlock()
	return r, nil
}
func (a *App) InstallIntegration(token string) error {
	a.integrationPending.Lock()
	defer a.integrationPending.Unlock()
	r := a.integrationPending.review
	if r == nil || token == "" || token != r.Token || time.Since(r.created) > 10*time.Minute {
		return errors.New("reveja o pacote antes de instalar")
	}
	p, hash, err := integrations.ReadArchive(r.path)
	if err != nil {
		return err
	}
	if hash != token {
		return errors.New("o pacote mudou; escolha e reveja novamente")
	}
	if err = a.integrationManager().Install(p); err != nil {
		return err
	}
	a.integrationPending.review = nil
	return nil
}
func (a *App) IntegrationChooseProject(id string) (string, error) {
	if _, err := a.integrationManager().Load(id); err != nil {
		return "", err
	}
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Escolher pasta do projeto"})
	if err != nil || path == "" {
		return "", err
	}
	if err = a.integrationManager().SetProject(id, path); err != nil {
		return "", err
	}
	return path, nil
}
func (a *App) IntegrationDiagnostics(id string) (integrations.Report, error) {
	return a.integrationManager().Diagnose(context.Background(), id)
}
