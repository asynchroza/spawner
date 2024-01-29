package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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
	file, err := os.OpenFile(rcFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = fmt.Fprintf(writer, "export SPAWN_REPOS_PATH=%s\n", path)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func getRcFilePath(rcFile string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, rcFile), nil
}

func updateRcFile(rcFilePath string, path string) error {
	content, err := os.ReadFile(rcFilePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "SPAWN_REPOS_PATH") {
		return replaceReposPathInRcFile(rcFilePath, path)
	}

	return writeReposPathInRcFile(rcFilePath, path)
}

func makeBackupOfRcFile(rcFilePath string) error {
	originalFile, err := os.Open(rcFilePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	fmt.Printf("Creating a backup of shell rc file at %s", rcFilePath+".backup\n")

	backupFilePath := rcFilePath + ".backup"
	backupFile, err := os.Create(backupFilePath)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	_, err = io.Copy(backupFile, originalFile)
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
			rcFilePath, err := getRcFilePath(rcFile)
			if err != nil {
				return err
			}

			makeBackupOfRcFile(rcFilePath)
			err = updateRcFile(rcFilePath, commands[1])

			if err != nil {
				return err
			}

			found = true

			fmt.Printf("\n❗ Please, reload your shell or source your shell's rc file by running \"source %s\"❗\n", rcFilePath)
			break
		}
	}
	if !found {
		return errors.New("Shell is not supported. Can't set repos path.")
	}

	return nil
}
