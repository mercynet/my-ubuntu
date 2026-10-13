package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) beforeClose(ctx context.Context) bool {
	a.operationMu.Lock()
	active := a.setup.Running() || a.eng.PreventClose()
	if !active {
		a.closing = true
	}
	a.operationMu.Unlock()
	if !active {
		return false
	}
	if a.ctx != nil {
		choice, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type: runtime.WarningDialog, Title: "Operação em curso",
			Message: "Há uma execução ativa ou em pausa. Retome-a ou peça cancelamento: o passo em curso termina e o resultado é guardado antes de parar. Depois poderá fechar.",
			Buttons: []string{"Aguardar", "Cancelar e aguardar"}, DefaultButton: "Aguardar", CancelButton: "Aguardar",
		})
		if err == nil && choice == "Cancelar e aguardar" {
			a.Cancel()
		}
	}
	return true
}
