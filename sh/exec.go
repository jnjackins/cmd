package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		cmd, err := t.mkFork()
		if err != nil {
			log.Print(err)
			return -1
		}
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
		args, vars := expandArgs(t)

		if len(args) == 0 {
			// only variable assignments
			for _, v := range vars {
				parts := strings.SplitN(v, "=", 2)
				setEnv(parts[0], parts[1])
			}
			return 0
		}

		if fn, ok := builtins[args[0]]; ok {
			exit, err := fn(args[1:])
			if err != nil {
				log.Printf("%s: %v", args[0], err)
			}
			return exit
		}

		cmd, err := t.mkCmd(args, vars)
		if err != nil {
			log.Print(err)
			return -1
		}
		dprintf("running simple command: %#v", cmd.Args)
		cmd.Run()
		return exitStatus(cmd.ProcessState)

	default:
		log.Printf("not implemented")
		return -1
	}
}

func exitStatus(state *os.ProcessState) int {
	return state.Sys().(syscall.WaitStatus).ExitStatus()
}

func expandArgs(t *treeNode) (args, vars []string) {
	prologue := true
	for _, n := range t.children {
		// variable assignments
		if prologue {
			if i := strings.Index(n.string, "="); i > 0 {
				vars = append(vars, n.string)
				continue
			} else {
				prologue = false
			}
		}

		// regular arguments
		if strings.ContainsAny(n.string, "[?*") {
			matches, err := filepath.Glob(n.string)
			if err == nil {
				args = append(args, matches...)
				continue
			}
		}
		args = append(args, n.string)
	}
	return args, vars
}

func (t *treeNode) mkCmd(args, vars []string) (*exec.Cmd, error) {
	if t.typ != SIMPLE {
		panic("mkCmd: bad node type")
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, fmt.Errorf("%s: command not found", args[0])
	}
	cmd := exec.Cmd{
		Path: path,
		Args: args,
	}
	if len(vars) > 0 {
		cmd.Env = append(os.Environ(), vars...)
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

func (t *treeNode) mkFork() (*exec.Cmd, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, "-c", t.children[0].String())
	cmd.Args[0] = filepath.Base(cmd.Args[0] + "(fork)")
	dprintf("running forked command: %#v", cmd.Args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
