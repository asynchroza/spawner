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

	yaml, err := parse.ParseYaml()
	if err != nil {
		return "", err
	}

	sshUrl := yaml.DevContainers[commands[1]]
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
	path := os.Getenv("SPAWN_REPOS_PATH")
	if path == "" {
		return "", errors.New("Directory for storing repos is not defined - run 'spawner help | grep \"set path to directory\"'")
	}

	return path, nil
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
	_, err = exec.Command("git", "clone", sshUrl, fmt.Sprintf("repos/%s", repoName)).Output()

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
