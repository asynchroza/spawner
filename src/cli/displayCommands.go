package cli

import (
	"asynchroza/spawner/src/parse"
	"fmt"
)

func DisplayCommands() error {
	yaml, err := parse.ParseYaml()

	if err != nil {
		return err
	}

	fmt.Println("\nCommands:")
	fmt.Println("\nrepos <full-path-to-directory> - set path to directory where repos are cloned")
	fmt.Println("\ndev <repo> - spawn vscode development container")
	fmt.Println("\nAvailable repos:")

	for key := range yaml.DevContainers {
		fmt.Printf("\t%s\n", key)
	}

	return nil
}
