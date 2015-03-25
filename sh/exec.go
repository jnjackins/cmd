package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

var builtins = map[string]func(args []string){
	"cd":   cdBuiltin,
	"exec": execBuiltin,
}

func run(args []string) {
	if len(args) == 0 {
		return
	}
	if fn, ok := builtins[args[0]]; ok {
		fn(args)
		return
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Run: %s", err)
	}
}

func cdBuiltin(args []string) {
	var dest string
	if len(args) == 1 {
		dest = os.Getenv("HOME")
	} else {
		dest = args[1]
	}
	err := os.Chdir(dest)
	if err != nil {
		log.Printf("Chdir: %s", err)
	}
}

func execBuiltin(args []string) {
	if len(args) < 2 {
		return
	}
	cmd, err := exec.LookPath(args[1])
	if err != nil {
		log.Printf("LookPath: %s", err)
	}
	err = syscall.Exec(cmd, args[1:], os.Environ())
	if err != nil {
		log.Printf("Exec: %s", err)
	}
}
