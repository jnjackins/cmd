package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runLine(line [][]*exec.Cmd) {
	// During parsing, it is convenient to append to the list
	// of pipes from the last to the first, so iterate over them
	// in reverse to restore the correct order.
	for i := len(line) - 1; i >= 0; i-- {
		pipe := line[i]
		for j := 0; j < len(pipe); j++ {
			cmd := pipe[j]
			if fn, ok := builtins[cmd.Path]; ok {
				fn(cmd.Path, cmd.Args)
			} else {
				if err := start(cmd); err != nil {
					log.Print(err)
					break
				} else {
					defer wait(cmd)
				}
			}
		}
	}
}

// Start reports whether it started cmd.
func start(cmd *exec.Cmd) error {
	if filepath.Base(cmd.Path) == cmd.Path {
		if lp, err := exec.LookPath(cmd.Path); err != nil {
			return err
		} else {
			cmd.Path = lp
		}
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func wait(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		log.Print(err)
	}
}

func connect(cmd1, cmd2 *exec.Cmd) {
	stdout, err := cmd1.StdoutPipe()
	if err != nil {
		log.Print(err)
	}
	cmd2.Stdin = stdout
}

func create(path string) io.WriteCloser {
	w, err := os.Create(path)
	if err != nil {
		log.Print(err)
	}
	return w
}

func open(path string) io.ReadCloser {
	r, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	return r
}

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Print(err)
	}
}
