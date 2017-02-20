package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetPrefix("daemonize: ")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("need at least 1 arg")
	}
	path, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	devnull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	attr := os.ProcAttr{
		Files: []*os.File{devnull, devnull, devnull},
	}
	if _, err = os.StartProcess(path, os.Args[1:], &attr); err != nil {
		log.Print(err)
	}
}
