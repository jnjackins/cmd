package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func execute(t *treeNode) int {
	switch t.typ {
	// foo; bar
	case ';':
		execute(t.children[0])
		return execute(t.children[1])

	// foo &
	case '&':
		path, err := os.Executable()
		if err != nil {
			log.Print(err)
			return -1
		}
		cmd := exec.Command(path, "-c", t.children[0].String())
		cmd.Args[0] = filepath.Base(cmd.Args[0] + "(fork)")
		dprintf("running forked command: %#v", cmd.Args)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Print(err)
			return -1
		}
		fmt.Printf("%d\n", cmd.Process.Pid)
		return 0

	// foo && bar
	case AND:
		if status := execute(t.children[0]); status == 0 {
			return execute(t.children[1])
		} else {
			return status
		}

	// foo || bar
	case OR:
		if status := execute(t.children[0]); status != 0 {
			return execute(t.children[1])
		} else {
			return status
		}

	// foo | bar
	case '|':
		r, w, err := os.Pipe()
		if err != nil {
			log.Print(err)
			return -1
		}
		t.children[0].io.stdout = w
		t.children[1].io.stdin = r
		go func() {
			execute(t.children[0])
			w.Close()
		}()
		defer r.Close()
		return execute(t.children[1])

	// foo bar
	case SIMPLE:
		cmd, err := t.mkCmd()
		if err != nil {
			log.Print(err)
			return -1
		}
		dprintf("running simple command: %#v", cmd.Args)
		cmd.Run()
		return exitStatus(cmd)

	default:
		log.Printf("not implemented")
		return -1
	}
}

func exitStatus(cmd *exec.Cmd) int {
	return cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
}

func (t *treeNode) mkCmd() (*exec.Cmd, error) {
	if t.typ != SIMPLE {
		panic("mkCmd: bad node type")
	}
	var args []string
	for _, n := range t.children {
		args = append(args, n.string)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, fmt.Errorf("%s: command not found", args[0])
	}
	cmd := exec.Cmd{
		Path: path,
		Args: args,
	}

	if t.io != nil {
		cmd.Stdin = t.io.stdin
		cmd.Stdout = t.io.stdout
		cmd.Stderr = t.io.stderr
		for fd, path := range t.io.redirs {
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
