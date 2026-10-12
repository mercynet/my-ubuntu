package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) beforeClose(ctx context.Context) bool {
	if !a.eng.PreventClose() {
		return false
	}
	if a.ctx != nil {
		choice, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type: runtime.WarningDialog, Title: "Operação em curso",
			Message: "Aguarde o fim do passo atual antes de fechar. Pode pedir o cancelamento: o Meu Ubuntu termina o passo em curso e guarda o resultado antes de parar.",
			Buttons: []string{"Aguardar", "Cancelar e aguardar"}, DefaultButton: "Aguardar", CancelButton: "Aguardar",
		})
		if err == nil && choice == "Cancelar e aguardar" {
			a.eng.Cancel()
		}
	}
	return true
}
