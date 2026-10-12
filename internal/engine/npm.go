package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"meuubuntu/internal/catalog"
)

func (e *Engine) installNPMUser(ctx context.Context, item *catalog.Item, log func(string)) error {
	if !e.Simulated && os.Geteuid() == 0 {
		return fmt.Errorf("a instalação dos agentes deve correr como utilizador normal, sem root")
	}
	if item == nil || !item.NPM.Valid() {
		return fmt.Errorf("receita npm inválida; nada instalado")
	}
	if !e.R.LookPath("node") || !e.R.LookPath("npm") {
		return fmt.Errorf("Node.js e npm não encontrados; instale a dependência indicada no plano")
	}
	version, err := e.R.Output(ctx, "node", "--version")
	if err != nil {
		return err
	}
	major, err := strconv.Atoi(strings.Split(strings.TrimPrefix(strings.TrimSpace(version), "v"), ".")[0])
	if err != nil || major < 22 {
		return fmt.Errorf("é necessário Node.js 22 ou superior; encontrado %s", strings.TrimSpace(version))
	}
	prefix := filepath.Join(e.Home, ".local")
	args := []string{"install", "--global", "--prefix", prefix, "--registry", "https://registry.npmjs.org", "--engine-strict", "--no-audit", "--no-fund"}
	if item.NPM.IgnoreScripts {
		args = append(args, "--ignore-scripts")
	}
	args = append(args, item.NPM.Package+"@"+item.NPM.Version)
	log("Instalação por utilizador: " + item.NPM.Package + "@" + item.NPM.Version)
	if err := e.R.Stream(ctx, log, "npm", args...); err != nil {
		return err
	}
	bin := filepath.Join(prefix, "bin", item.NPM.Binary)
	output, err := e.R.Output(ctx, bin, "--version")
	if err != nil {
		return fmt.Errorf("a instalação terminou mas o agente não executa: %w", err)
	}
	log("Verificado: " + strings.TrimSpace(output))
	log("Para começar: " + bin + "; inicie a sessão de acesso dentro do próprio agente")
	return nil
}
