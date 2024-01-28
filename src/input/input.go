package input

import "flag"

func GetCommand() []string {
	flag.Parse()
	return flag.Args()
}
