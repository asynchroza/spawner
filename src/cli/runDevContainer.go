package cli

import (
	"asynchroza/spawner/src/parse"
	"errors"
	"fmt"
)

func RunDevContainer(commands []string) error {
	if len(commands) < 2 {
		return errors.New("Container is not provided")
	}

	yaml, err := parse.ParseYaml()
	if err != nil {
		return err
	}

	fmt.Println(yaml)
	return nil
}
