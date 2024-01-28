package main

import (
	"asynchroza/spawner/src/cli"
	"asynchroza/spawner/src/input"
	"fmt"
)

func main() {
	commands := input.GetCommand()
	err := cli.ParseCommand(commands)

	if err != nil {
		fmt.Println(err)
	}
}
