package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type cmd struct {
	args   []string
	stdin  string
	stdout string
}

var builtins = map[string]func(args []string){
	"cd":   cdBuiltin,
	"exec": execBuiltin,
}

func run(cmd cmd) {
	if len(cmd.args) == 0 {
		return
	}
	if fn, ok := builtins[cmd.args[0]]; ok {
		fn(cmd.args)
		return
	}
	stdin, stdout := os.Stdin, os.Stdout
	if cmd.stdin != "" {
		f, err := os.Open(cmd.stdin)
		if err != nil {
			log.Print(err)
			return
		}
		stdin = f
	}
	if cmd.stdout != "" {
		f, err := os.Create(cmd.stdout)
		if err != nil {
			log.Print(err)
			return
		}
		stdout = f
	}
	runner := exec.Command(cmd.args[0], cmd.args[1:]...)
	runner.Stdin = stdin
	runner.Stdout = stdout
	runner.Stderr = os.Stderr
	err := runner.Run()
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
