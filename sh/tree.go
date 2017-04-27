package main

import (
	"fmt"
	"io"
)

type treeNode struct {
	typ      int
	int      int
	string   string
	io       *ioSpec
	children []*treeNode
}

type ioSpec struct {
	redirs map[int]string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func mkTree(typ int, children ...*treeNode) *treeNode {
	return &treeNode{
		typ:      typ,
		children: children,
	}
}

func mkSimple(t *treeNode) *treeNode {
	t.typ = SIMPLE
	t.io = &ioSpec{redirs: make(map[int]string)}
	return t
}

func mkLeaf(typ int, i int, s string) *treeNode {
	return &treeNode{
		typ:    typ,
		int:    i,
		string: s,
	}
}

func (n *treeNode) String() string {
	var s string
	switch n.typ {
	case ';':
		sep := "; "
		if n.children[0].typ == '&' {
			sep = " "
		}
		s = fmt.Sprintf("%v%s%v", n.children[0], sep, n.children[1])
	case '&':
		s = fmt.Sprintf("%v &", n.children[0])
	case AND:
		s = fmt.Sprintf("%v && %v", n.children[0], n.children[1])
	case OR:
		s = fmt.Sprintf("%v || %v", n.children[0], n.children[1])
	case '|':
		s = fmt.Sprintf("%v |%v", n.children[0], n.children[1])
	case PAREN:
		s = fmt.Sprintf("(%v)", n.children[0])
	case SIMPLE:
		s = n.children[0].string
		for _, c := range n.children[1:] {
			s += " " + c.string
		}
	default:
		panic("bad node type")
	}
	if n.io != nil {
		for fd, path := range n.io.redirs {
			var redir string
			switch fd {
			case 0:
				redir = "<"
			case 1:
				redir = ">"
			default:
				redir = fmt.Sprintf("%d>", fd)
			}
			s += fmt.Sprintf(" %s%s", redir, path)
		}
	}
	return s
}
