package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/x6r/venboy/log"
)

func Execute() {
	if err := app().Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func app() *cli.App {
	var uninstall bool

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "uninstall",
				Value:       false,
				Usage:       "Uninstall Vencord",
				Destination: &uninstall,
			},
		},
		Action: func(_ *cli.Context) error {
			resourcesPath := getResourcesPath()
			appPath := resourcesPath + `\app`

			switch uninstall {
			case true:
				if err := os.RemoveAll(appPath); err != nil {
					log.Error("Failed to uninstall. Vencord is already uninstalled.")
					os.Exit(1)
				}
				log.Info("Successfully ", color.RedString("uninstalled"), " Vencord!")
			default:
				patcherPath := getPatcherPath()
				if err := os.Mkdir(appPath, 0755); err != nil {
					log.Error("app folder exists. Looks like your Discord is already modified.")
					os.Exit(1)
				}

				log.Info("Installing...")

				index := fmt.Sprintf("require(\"%s\");\nrequire(\"../app.asar\");",
					strings.ReplaceAll(patcherPath, `\`, `\\`))
				if err := os.WriteFile(appPath+"\\index.js", []byte(index), 0755); err != nil {
					log.Error("Failed to write index.js: ", err)
				}

				pkg := fmt.Sprintf("{\n  \"main\": \"index.js\",\n  \"name\": \"discord\"\n}")
				if err := os.WriteFile(appPath+"\\package.json", []byte(pkg), 0755); err != nil {
					log.Error("Failed to write package.json: ", err)
				}

				log.Info("Successfully ", color.GreenString("installed"), " Vencord!")
			}
			return nil
		},
	}

	return app
}

func getResourcesPath() string {
	appdata, err := os.UserCacheDir()
	if err != nil {
		log.Error("Failed to get %AppData%: ", err)
		os.Exit(1)
	}
	paths, err := os.ReadDir(appdata)
	if err != nil {
		log.Error("Failed to read %AppData%: ", err)
		os.Exit(1)
	}
	var branches []string
	for _, path := range paths {
		name := path.Name()
		if strings.HasPrefix(name, "Discord") &&
			name != "DiscordGames" &&
			name != "DiscordOverlayHost" {
			branches = append(branches, name)
		}
	}

	var branch string
	prompt := &survey.Select{
		Message: "Choose a branch",
		Options: branches,
	}
	if err := survey.AskOne(prompt, &branch); err != nil {
		log.Error("Branch selection error: ", err)
		os.Exit(1)
	}
	branchPath := appdata + `\` + branch

	appPaths, err := os.ReadDir(branchPath)
	if err != nil {
		log.Error("Failed to read ", color.GreenString(branchPath), " direcotry: ", err)
		os.Exit(1)
	}
	var app string
	for _, path := range appPaths {
		name := path.Name()
		if strings.HasPrefix(name, "app-") {
			app = name
		}
	}
	log.Warn("Version: ", strings.ReplaceAll(app, "app-", ""))

	return branchPath + `\` + app + `\resources`
}

func getPatcherPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current working directory: ", err)
		os.Exit(1)
	}
	patcher := pwd + `\dist\patcher.js`
	if !fileExists(patcher) {
		log.Error("patcher.js doesn't exist. Are you running from Vencord folder? Did you forget to build?")
		os.Exit(1)
	}

	return patcher
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
