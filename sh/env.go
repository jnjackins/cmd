package main

import (
	"log"
	"os"
)

var env map[string]string

func init() {
	env = make(map[string]string)
}

func updateEnv() {
	for key, val := range env {
		if err := os.Setenv(key, val); err != nil {
			log.Print(err)
		}
		delete(env, key)
	}
}
