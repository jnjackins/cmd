package main

import (
	"fmt"
	"io"
)

const (
	typeListSequence = iota
	typeListFork
	typeListAnd
	typeListOr
)

type node interface{}

type redirecter interface {
	redirect(fd int, path string)
}

type connecter interface {
	setStdin(io.Reader)
	setStdout(io.Writer)
	setStderr(io.Writer)
}

type forkNode struct {
	typ  int
	tree node
}

type listNode struct {
	typ         int
	left, right node
}

type pipeNode struct {
	left  connecter
	right connecter
}

func (n *pipeNode) setStdin(r io.Reader)  { n.left.setStdin(r) }
func (n *pipeNode) setStdout(w io.Writer) { n.right.setStdout(w) }
func (n *pipeNode) setStderr(w io.Writer) { n.right.setStderr(w) }

type parenNode struct {
	typ  int
	tree node
	ioNode
}

type simpleNode struct {
	args *argNode
	ioNode
}

type ioNode struct {
	redirs map[int]string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (n *ioNode) redirect(fd int, path string) {
	if n.redirs == nil {
		n.redirs = make(map[int]string)
	}
	n.redirs[fd] = path
}

func (n *ioNode) setStdin(r io.Reader)  { n.stdin = r }
func (n *ioNode) setStdout(w io.Writer) { n.stdout = w }
func (n *ioNode) setStderr(w io.Writer) { n.stderr = w }

type argNode struct {
	val  string
	next *argNode
}

func printTree(n node) {
	if n == nil {
		return
	}
	switch t := n.(type) {
	case *argNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		if t.next != nil {
			printTree(t.next)
		}
	case *parenNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		printTree(t.tree)
	case *simpleNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		printTree(t.args)
	case *pipeNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		printTree(t.left)
		printTree(t.right)
	case *listNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		printTree(t.left)
		printTree(t.right)
	case *forkNode:
		fmt.Printf("%T(%[1]p): %#[1]v\n", t)
		printTree(t.tree)
	default:
		panic(fmt.Sprintf("unrecognized node type: %T", t))
	}
}
