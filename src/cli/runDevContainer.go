package cli

import (
	"asynchroza/spawner/src/parse"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func parseContainerPath(commands []string) (string, error) {
	if len(commands) < 2 {
		return "", errors.New("Container is not provided")
	}

	yaml := parse.GetConfiguration()

	sshUrl := yaml.DevContainers.Containers[commands[1]]
	if len(sshUrl) == 0 {
		return "", errors.New("No path found for container. You may see available options by typing 'spawner help'")
	}

	return sshUrl[0], nil
}

func extractRepoName(sshURL string) (string, error) {
	parts := strings.Split(sshURL, ":")
	if len(parts) < 2 {
		return "", errors.New("Invalid SSH URL")
	}

	repoParts := strings.Split(parts[1], "/")
	if len(repoParts) < 2 {
		return "", errors.New("Invalid repository URL")
	}

	repoName := strings.TrimSuffix(repoParts[1], ".git")

	return repoName, nil
}

func findRepo(baseDir, repoName string) (bool, error) {
	var folderFound bool

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == repoName {
			folderFound = true
			return nil
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return folderFound, nil
}

func getRepoDirectoryPath() (string, error) {
	pathArr := parse.GetConfiguration().DevContainers.ReposPath

	if len(pathArr) == 0 {
		return "", errors.New("Path to repos directory is not provided in configuration")
	}
	return pathArr[0], nil
}

func pullRepoLocallyAndGetName(sshUrl string) (string, error) {
	repoName, err := extractRepoName(sshUrl)

	if err != nil {
		return "", err
	}

	repoDirectory, err := getRepoDirectoryPath()
	if err != nil {
		return "", err
	}

	repoExists, err := findRepo(repoDirectory, repoName)
	if err != nil {
		return "", err
	}

	if repoExists {
		return repoName, nil
	}

	fmt.Println("⬇️ Pulling repo locally")
	cmd := exec.Command("git", "clone", sshUrl, fmt.Sprintf("repos/%s", repoName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	return repoName, nil
}

func spawnContainer(repoName string) error {
	fmt.Println("🚀 Starting development container")

	cmd := exec.Command("devcontainer", "open", fmt.Sprintf("repos/%s", repoName))
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func RunDevContainer(commands []string) error {
	sshUrl, err := parseContainerPath(commands)
	if err != nil {
		return err
	}

	repoName, err := pullRepoLocallyAndGetName(sshUrl)
	if err != nil {
		return err
	}

	err = spawnContainer(repoName)
	if err != nil {
		return err
	}

	return nil
}
