package cli

import (
	"errors"
)

const (
	DEV          = "dev"
	HELP         = "help"
	INSTALL_DEVC = "install-devcontainers"
)

func ParseCommand(commands []string) error {
	if len(commands) == 0 {
		return errors.New("No arguments were provided")
	}

	switch commands[0] {
	case DEV:
		return RunDevContainer(commands)
	case HELP:
		return DisplayCommands()
	case INSTALL_DEVC:
		return InstallDevContainers()
	default:
		return errors.New("Leading argument doesn't match any of the defined commands")
	}
}
