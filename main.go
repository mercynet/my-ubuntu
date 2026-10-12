package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"

	"meuubuntu/internal/catalog"
	"meuubuntu/internal/detect"
	"meuubuntu/internal/diag"
	"meuubuntu/internal/helper"
	"meuubuntu/internal/sys"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed frontend/src/assets/apps/DEVICON-LICENSE.txt frontend/src/assets/apps/SIMPLE-ICONS-LICENSE.txt
var logoLicenses embed.FS

//go:embed docs/manual.md
var manualText string

//go:embed LICENSE
var projectLicense string

var version = "dev"

func main() {
	for _, a := range os.Args[1:] {
		switch a {
		case "--helper":
			os.Exit(runHelper())
		case "--version":
			fmt.Println("meu-ubuntu", version)
			return

		case "--help", "-h":
			fmt.Printf("Meu Ubuntu %s\n\nUso: meu-ubuntu [--simular] [--perfil ficheiro.json]\n\n  --simular        Explorar com executor e dados separados de simulação\n  --perfil JSON    Abrir uma revisão de perfil, sem instalar\n  --version        Mostrar a versão do executável\n  --manual         Imprimir o manual incluído\n  --licenses       Imprimir a licença do projeto e avisos incluídos\n  --help, -h       Mostrar esta ajuda\n\nAbra a aplicação como utilizador normal. As operações administrativas usam pkexec.\n", version)
			return

		case "--manual":
			fmt.Print(manualText)
			return
		case "--licenses":
			fmt.Printf("Meu Ubuntu — GNU GPLv3\n%s\n", projectLicense)
			for _, name := range []string{"DEVICON-LICENSE.txt", "SIMPLE-ICONS-LICENSE.txt"} {
				b, err := logoLicenses.ReadFile("frontend/src/assets/apps/" + name)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("%s\n%s\n", name, b)
			}
			for _, l := range diag.ExportLicenses() {
				fmt.Printf("%s (%s)\n%s\n", l.Name, l.License, l.Text)
			}
			return
		}
	}
	if os.Geteuid() == 0 {
		fmt.Fprintln(os.Stderr, "Não execute o Meu Ubuntu como root. Abra-o como utilizador normal; as operações administrativas pedem autorização ao sistema.")
		os.Exit(1)
	}
	simulated := os.Getenv("MEU_UBUNTU_SIMULATE") == "1"
	for _, a := range os.Args[1:] {
		if a == "--simular" {
			simulated = true
		}
	}
	app, err := NewApp(simulated)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app.pendingPath = profileArg(os.Args[1:])
	title, instanceID := "Meu Ubuntu", "pt.meuubuntu.app"
	if simulated {
		title += " — modo simulação"
		instanceID += ".sim"
	}
	err = wails.Run(&options.App{
		Title:            title,
		Width:            1280,
		Height:           840,
		MinWidth:         900,
		MinHeight:        600,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 27, A: 255},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind:             []interface{}{app},
		// One window only: opening a profile JSON again reviews it in the running app.
		SingleInstanceLock: &options.SingleInstanceLock{UniqueId: instanceID,
			OnSecondInstanceLaunch: func(d options.SecondInstanceData) { app.SecondInstance(d.Args) }},
		Linux: &linux.Options{Icon: icon, ProgramName: "meu-ubuntu", WebviewGpuPolicy: linux.WebviewGpuPolicyOnDemand},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runHelper is the restricted administrative mode started via pkexec.
func runHelper() int {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "--helper só é usado pelo próprio Meu Ubuntu através do pkexec")
		return 2
	}
	uid := os.Getenv("PKEXEC_UID")
	if uid == "" {
		fmt.Fprintln(os.Stderr, "--helper exige pkexec (PKEXEC_UID ausente)")
		return 2
	}
	cat, err := catalog.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	r := sys.Real{}
	system := detect.ReadSystem(r, false)
	env := &helper.Env{R: r, Cat: cat, Codename: system.Codename, Distribution: system.ID, UID: uid, Fetch: helper.HTTPFetch}
	if err := helper.Serve(env, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
