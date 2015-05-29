package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type program struct {
	dirs []*directive
}

type directive struct {
	s       string
	pattern string
	action  string
	re      *regexp.Regexp
	fns     []func(record)
}

const (
	strArg = iota
	numArg
)

func compileProg(r io.Reader) (*program, error) {
	dprintf("compiling program...")
	prog := &program{dirs: make([]*directive, 0, 1)}
	scanner := bufio.NewScanner(r)
	scanner.Split(splitDirective)
	for scanner.Scan() {
		s := scanner.Text()
		dprintf("scanned a directive: %#v", s)
		dir, err := parseDir(s)
		if err != nil {
			log.Printf("error parsing directive: %v", err)
			continue
		}
		prog.dirs = append(prog.dirs, dir)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	dprintf("compiled program with %d directives", len(prog.dirs))
	return prog, nil
}

func splitDirective(data []byte, atEOF bool) (adv int, tok []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	i := bytes.IndexByte(data, '{')
	if !atEOF && i < 0 {
		// Request more data.
		return 0, nil, nil
	} else if i < 0 {
		return 0, nil, fmt.Errorf("bad directive: missing opening brace: %#v", string(data))
	}
	stack := 1
	for i++; i < len(data); i++ {
		switch data[i] {
		case '{':
			stack++
		case '}':
			stack--
		}
		if stack == 0 {
			return i + 1, data[:i+1], nil
		}
	}
	if atEOF {
		return 0, nil, fmt.Errorf("bad directive: missing closing brace: %#v", string(data))
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func (p *program) exec(rec record) {
	for _, d := range p.dirs {
		d.exec(rec)
	}
}

func parseDir(s string) (*directive, error) {
	parts := strings.SplitN(s, "{", 2)
	if len(parts) != 2 {
		panic("couldn't split directive into pattern + action")
	}
	pattern, action := parts[0], parts[1]
	// we split on the opening brace, so remove the closing brace as well
	action = action[:len(action)-1]
	var re *regexp.Regexp
	if pattern == "" {
		re = nil
	} else {
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern: %s", pattern)
		}
	}
	d := &directive{
		pattern: pattern,
		action:  action,
		re:      re,
		fns:     parseAction(action),
	}
	dprintf("compiled a directive with pattern %#v and action %#v", d.pattern, d.action)
	return d, nil
}

func (d *directive) exec(rec record) {
	if d.re == nil || d.re.MatchString(rec.s) {
		dprintf("pattern %#v matches record %#v", d.pattern, rec.s)
		for _, fn := range d.fns {
			fn(rec)
		}
	} else {
		dprintf("pattern %#v does not match record %#v", d.pattern, rec.s)
	}
}
