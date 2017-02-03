package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	for _, cmd := range os.Args[1:] {
		if path, err := exec.LookPath(cmd); err == nil {
			fmt.Println(path)
		}
	}
}
