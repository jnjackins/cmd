package main

import (
	"log"
	"os"
	"strconv"

	"sigint.ca/user"
)

var env map[string]string

func setupEnv() {
	env = make(map[string]string)
	env["pid"] = strconv.Itoa(os.Getpid())
	if os.Getenv("HOME") == "" {
		u, err := user.Current()
		if err != nil {
			log.Print(err)
		}
		env["HOME"] = u.HomeDir
	}
	if os.Getenv("prompt") == "" {
		env["prompt"] = "$ "
	}
	updateEnv()
}

func updateEnv() {
	for key, val := range env {
		if err := os.Setenv(key, val); err != nil {
			log.Print(err)
		}
		delete(env, key)
	}
}
