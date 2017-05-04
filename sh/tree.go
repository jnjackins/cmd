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
	redirs  map[int]redir
	pipeIn  io.Reader
	pipeOut io.Writer
}

type redir struct {
	path   string
	append bool
}

func mkTree(typ int, children ...*treeNode) *treeNode {
	t := &treeNode{
		typ:      typ,
		children: children,
	}
	switch typ {
	case SIMPLE, PAREN, IF, FOR:
		// "item" in yacc grammar
		t.io = &ioSpec{redirs: make(map[int]redir)}
	}
	return t
}

func mkLeaf(typ int, i int, s string) *treeNode {
	return &treeNode{
		typ:    typ,
		int:    i,
		string: s,
	}
}

func (t *treeNode) redirect(fd int, path string, append bool) {
	t.io.redirs[fd] = redir{
		path:   path,
		append: append,
	}
}

func (t *treeNode) String() string {
	var s string
	switch t.typ {
	case ';':
		sep := "; "
		if t.children[0].typ == '&' {
			sep = " "
		}
		s = fmt.Sprintf("%v%s%v", t.children[0], sep, t.children[1])
	case '&':
		s = fmt.Sprintf("%v &", t.children[0])
	case AND:
		s = fmt.Sprintf("%v && %v", t.children[0], t.children[1])
	case OR:
		s = fmt.Sprintf("%v || %v", t.children[0], t.children[1])
	case '|':
		s = fmt.Sprintf("%v |%v", t.children[0], t.children[1])
	case IF:
		s = fmt.Sprintf("if %v; then %v; fi", t.children[0], t.children[1])
	case FOR:
		s = fmt.Sprintf("for %v in %v; do %v; done", t.children[0], t.children[1], t.children[2])
	case PAREN:
		s = fmt.Sprintf("(%v)", t.children[0])
	case SIMPLE:
		s = t.children[0].String()
	case WORDS:
		s = t.children[0].string
		for _, c := range t.children[1:] {
			s += " " + c.string
		}
	case WORD:
		s = t.string
	default:
		panic("bad node type")
	}
	if t.io != nil {
		for fd, redir := range t.io.redirs {
			var c string
			switch fd {
			case 0:
				c = "<"
			case 1:
				c = ">"
			default:
				c = fmt.Sprintf("%d>", fd)
			}
			s += fmt.Sprintf(" %s%s", c, redir.path)
		}
	}
	return s
}
