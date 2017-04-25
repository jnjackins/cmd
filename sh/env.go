package main

import (
	"os"
	"strings"
)

var env map[string]string

func setupEnv() {
	environ := os.Environ()
	env = make(map[string]string, len(environ))
	for _, s := range environ {
		split := strings.SplitN(s, "=", 2)
		env[split[0]] = split[1]
	}

	if _, ok := env["PS1"]; !ok {
		env["PS1"] = "$ "
	}

	if *debug {
		for k, v := range env {
			dprintf("env: set %s=%s", k, v)
		}
	}
}
