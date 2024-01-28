package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var AVAILABLE_SHELLS = map[string]string{
	"zsh":  ".zshrc",
	"bash": ".bashrc",
}

func replaceReposPathInRcFile(rcFilePath string, path string) error {
	content, err := os.ReadFile(rcFilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "export SPAWN_REPOS_PATH=") {
			lines[i] = fmt.Sprintf("export SPAWN_REPOS_PATH=%s", path)
			break
		}
	}

	newContent := strings.Join(lines, "\n")

	err = os.WriteFile(rcFilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func writeReposPathInRcFile(rcFilePath string, path string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo \"export SPAWN_REPOS_PATH=%s\" >> %s", path, rcFilePath))

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func updateRcFile(rcFile string, path string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	rcFilePath := fmt.Sprintf("%s/%s", usr.HomeDir, rcFile)

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

	cmd := exec.Command("/bin/exec", "$SHELL")
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

	var found bool

	for key := range AVAILABLE_SHELLS {
		if strings.Contains(shell, key) {
			rcFile := AVAILABLE_SHELLS[key]
			err := updateRcFile(rcFile, commands[1])

			if err != nil {
				return err
			}

			found = true
			break
		}
	}
	if !found {
		return errors.New("Shell is not supported. Can't set repos path.")
	}

	return reloadShell()
}
