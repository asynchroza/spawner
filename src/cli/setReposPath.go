package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var AVAILABLE_SHELLS = map[string]string{
	"zsh":  ".zshrc",
	"bash": ".bashrc",
}

func replaceReposPathInRcFile(rcFilePath string, path string) error {
	cmd := exec.Command(fmt.Sprintf("awk -F= '/^export SPAWN_REPOS_PATH/ {sub($2, \"%s\"); print}' %s > %s"), path, rcFilePath, rcFilePath)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func writeReposPathInRcFile(rcFilePath string, path string) error {
	cmd := exec.Command("echo", fmt.Sprintf("'export SPAWN_REPOS_PATH=%s'", path))
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func updateRcFile(rcFile string, path string) error {
	rcFilePath := fmt.Sprintf("~/%s", rcFile)

	content, err := os.ReadFile(rcFilePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "SPAWN_REPOS_PATH") {
		return replaceReposPathInRcFile(rcFilePath, path)
	}

	return writeReposPathInRcFile(rcFilePath, path)
}

func reloadShell() error {
	fmt.Println("Reloading shell")

	cmd := exec.Command("exec", "$SHELL")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func SetReposPath(commands []string) error {
	if len(commands) < 2 {
		return errors.New("Path is not provided")
	}

	shell := os.Getenv("SHELL")

	var err error

	if rcFile, exists := AVAILABLE_SHELLS[shell]; exists {
		err = updateRcFile(rcFile, commands[1])
	} else {
		err = errors.New("Shell is not supported. Can't set repos path.")
	}

	if err != nil {
		return err
	}

	return reloadShell()
}
