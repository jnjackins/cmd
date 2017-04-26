package main

import "os"

var builtins = map[string]func([]string) error{
	"cd": cd,
}

func cd(args []string) error {
	if len(args) == 0 {
		return os.Chdir(env["HOME"])
	}
	return os.Chdir(args[0])
}
