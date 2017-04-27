package main

import (
	"os"
	"strings"
	"sync"
)

var env struct {
	mu   sync.RWMutex
	vars map[string]string
}

func setupEnv() {
	environ := os.Environ()
	env.vars = make(map[string]string, len(environ))
	for _, s := range environ {
		split := strings.SplitN(s, "=", 2)
		env.vars[split[0]] = split[1]
	}

	if _, ok := env.vars["PS1"]; !ok {
		env.vars["PS1"] = "$ "
	}

	if *debug {
		for k, v := range env.vars {
			dprintf("env: set %s=%s", k, v)
		}
	}
}

func setEnv(key, val string) {
	env.mu.Lock()
	defer env.mu.Unlock()

	env.vars[key] = val
}

func getEnv(key string) string {
	env.mu.RLock()
	defer env.mu.RUnlock()

	return env.vars[key]
}

func exportEnv(key string) {
	panic("TODO")
}
