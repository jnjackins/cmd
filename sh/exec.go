package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func execute(n node) {
	switch t := n.(type) {
	case *cmdNode:
		cmd, err := t.mkCmd()
		if err != nil {
			log.Print(err)
			return
		}
		cmd.Run()

	case *pipeNode:
		r, w, err := os.Pipe()
		if err != nil {
			log.Print(err)
			return
		}
		t.left.setStdout(w)
		t.right.setStdin(r)
		go func() {
			execute(t.left)
			w.Close()
		}()
		execute(t.right)
		r.Close()

	default:
		log.Printf("not handled: %T", t)
		printTree(n)
	}
}

func (c *cmdNode) mkCmd() (*exec.Cmd, error) {
	var args []string
	for p := c.args; p != nil; p = p.next {
		args = append(args, p.val)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, fmt.Errorf("%s: command not found", args[0])
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Stdin:  c.stdin,
		Stdout: c.stdout,
		Stderr: c.stderr,
	}

	for fd, path := range c.redirs {
		switch fd {
		case 0:
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			cmd.Stdin = f
		case 1:
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			cmd.Stdout = f
		case 2:
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			cmd.Stderr = f
		default:
			panic("TODO")
		}
	}
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	return &cmd, nil
}
