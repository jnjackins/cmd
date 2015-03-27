package main

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func isTTY(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
