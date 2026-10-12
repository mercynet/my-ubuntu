package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"meuubuntu/internal/diag"
	"meuubuntu/internal/safefile"
	"meuubuntu/internal/store"
)

type diagnosticsCache struct {
	sync.Mutex
	report diag.Report
}

func (a *App) Diagnostics() diag.Report {
	rep := diag.Run(context.Background(), a.r, a.simulated)
	rep.SnapshotID, rep.CollectedAt = store.NewID(), time.Now().Format(time.RFC3339)
	a.diagCache.Lock()
	a.diagCache.report = rep
	a.diagCache.Unlock()
	return rep
}

func (a *App) diagnosticExport(snapshot, format string, sections []string) ([]byte, error) {
	a.diagCache.Lock()
	rep := a.diagCache.report
	a.diagCache.Unlock()
	if snapshot == "" || snapshot != rep.SnapshotID {
		return nil, errors.New("o diagnóstico mudou; atualize e reveja os dados antes de exportar")
	}
	s := a.det.Sys
	mode := "real"
	if a.simulated {
		mode = "simulação"
	}
	return diag.Encode(rep, diag.Metadata{Version: version, System: s.Pretty, Architecture: s.Arch, Session: s.Session, Mode: mode}, format, sections)
}

func (a *App) ExportDiagnostics(snapshot, format string, sections []string) (string, error) {
	data, err := a.diagnosticExport(snapshot, format, sections)
	if err != nil {
		return "", err
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{Title: "Exportar diagnóstico", DefaultFilename: "meu-ubuntu-diagnostico-" + time.Now().Format("2006-01-02") + "." + format, Filters: []runtime.FileFilter{{DisplayName: "Relatório " + format, Pattern: "*." + format}}})
	if err != nil || path == "" {
		return "", err
	}
	if !strings.EqualFold(filepath.Ext(path), "."+format) {
		return "", errors.New("guarde o ficheiro com a extensão ." + format + " para confirmar o destino correto")
	}
	if err := safefile.Write(path, data, 0600); err != nil {
		return "", err
	}
	return path, nil
}
