package cli

import (
	"asynchroza/spawner/src/parse"
	"fmt"

	"github.com/fatih/color"
)

func DisplayCommands() error {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	yaml := parse.GetConfiguration()

	fmt.Println("\nCommands:")

	yellow.Printf("\ninstall-devcontainers")
	fmt.Println(" - install devcontainers cli tool using npm")

	yellow.Printf("\ndev <repo>")
	fmt.Println(" - spawn vscode development container")

	fmt.Printf("\nAvailable repos:\n\n")

	for key := range yaml.DevContainers.Containers {
		green.Printf("\t➡️ %s\n", key)
	}

	return nil
}
