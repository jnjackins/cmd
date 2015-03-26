package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

var builtins = map[string]func(path string, args []string){
	"cd":   cdBuiltin,
	"exec": execBuiltin,
}

func cdBuiltin(path string, args []string) {
	var dest string
	if len(args) == 1 {
		dest = os.Getenv("HOME")
	} else {
		dest = args[1]
	}
	err := os.Chdir(dest)
	if err != nil {
		log.Print(err)
	}
}

func execBuiltin(path string, args []string) {
	if len(args) < 2 {
		return
	}
	cmd, err := exec.LookPath(args[1])
	if err != nil {
		log.Print(err)
	}
	err = syscall.Exec(cmd, args[1:], os.Environ())
	if err != nil {
		log.Print(err)
	}
}
