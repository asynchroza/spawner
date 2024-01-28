package cli

import (
	"fmt"
	"os/exec"
)

func InstallDevContainers() error {
	fmt.Println("Installing @devcontainers/cli globally")
	cmd := exec.Command("npm", "install", "-g", "@devcontainers/cli")
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
