package main

import (
	"log"
	"os"
	"strconv"
)

var env map[string]string

// TODO: setup $home
func setupEnv() {
	env = make(map[string]string)
	env["pid"] = strconv.Itoa(os.Getpid())
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
