package main

import (
	"fmt"
	"io"
)

const (
	typeListSequence = iota
	typeListAnd
	typeListOr

	typeForkAmp = iota
	typeForkParen
)

type node interface {
	String() string
}

type redirecter interface {
	node
	redirect(fd int, path string)
}

type connecter interface {
	node
	setStdin(io.Reader)
	setStdout(io.Writer)
	setStderr(io.Writer)
}

type forkNode struct {
	typ  int
	tree node
}

func (n *forkNode) String() string {
	switch n.typ {
	case typeForkAmp:
		return n.tree.String() + " &"
	case typeForkParen:
		// tree is *parenNode; will print ( and )
		return n.tree.String()
	}
	panic("unreached")
}

type listNode struct {
	typ         int
	left, right node
}

func (n *listNode) String() string {
	var sep string
	switch n.typ {
	case typeListSequence:
		switch n.left.(type) {
		case *forkNode:
			sep = " "
		default:
			sep = "; "
		}
	case typeListAnd:
		sep = " && "
	case typeListOr:
		sep = " || "
	}
	return fmt.Sprintf("%v%s%v", n.left.String(), sep, n.right.String())
}

type pipeNode struct {
	left  connecter
	right connecter
}

func (n *pipeNode) String() string {
	return fmt.Sprintf("%v || %v", n.left, n.right)
}

func (n *pipeNode) setStdin(r io.Reader)  { n.left.setStdin(r) }
func (n *pipeNode) setStdout(w io.Writer) { n.right.setStdout(w) }
func (n *pipeNode) setStderr(w io.Writer) { n.right.setStderr(w) }

type parenNode struct {
	typ  int
	tree node
	ioNode
}

func (n *parenNode) String() string {
	return fmt.Sprintf("(%v)", n.tree) + n.printRedirs()
}

type simpleNode struct {
	args *argNode
	ioNode
}

func (n *simpleNode) String() string {
	return fmt.Sprintf("%v", n.args) + n.printRedirs()
}

type ioNode struct {
	redirs map[int]string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (n *ioNode) printRedirs() string {
	var s string
	for fd, path := range n.redirs {
		var sym rune
		if fd == 0 {
			sym = '<'
		} else {
			sym = '>'
		}
		s += fmt.Sprintf(" %d%c%s", fd, sym, path)
	}
	return s
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

func (n *argNode) String() string {
	s := n.val
	for n = n.next; n != nil; n = n.next {
		s += " " + n.val
	}
	return s
}

func printTree(n node) {
	if n == nil {
		return
	}
	switch t := n.(type) {
	case *argNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		if t.next != nil {
			printTree(t.next)
		}
	case *parenNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		printTree(t.tree)
	case *simpleNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		printTree(t.args)
	case *pipeNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		printTree(t.left)
		printTree(t.right)
	case *listNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		printTree(t.left)
		printTree(t.right)
	case *forkNode:
		fmt.Printf("%T(%[1]p): %[1]v\n", t)
		printTree(t.tree)
	default:
		panic(fmt.Sprintf("unrecognized node type: %T", t))
	}
}
