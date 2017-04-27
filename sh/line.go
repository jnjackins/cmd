package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/peterh/liner"
)

var (
	inbuf  *bufio.Reader
	lstate *liner.State
)

func initPrompt() {
	lstate = liner.NewLiner()
	lstate.SetCtrlCAborts(false)
	lstate.SetTabCompletionStyle(liner.TabPrints)
	lstate.SetWordCompleter(pathComplete)
}

func fixTerminal() {
	if lstate != nil {
		lstate.Close()
	}
}

func setInput(r io.Reader) {
	inbuf = bufio.NewReader(r)
}

func getLine() ([]byte, error) {
	if inbuf == nil {
		s, err := lstate.Prompt(env["PS1"])
		if err != nil {
			return nil, err
		}
		lstate.AppendHistory(s)
		return []byte(s), nil
	}
	return inbuf.ReadBytes('\n')
}

func pathComplete(line string, pos int) (head string, completions []string, tail string) {
	headpos := strings.LastIndexFunc(line[:pos], unicode.IsSpace)
	headpos++
	head = line[:headpos]
	tailpos := strings.IndexFunc(line[headpos:], unicode.IsSpace)
	if tailpos < 0 {
		tailpos = len(line)
	} else {
		tailpos += headpos
	}
	prefix := line[headpos:tailpos]
	tail = line[tailpos:]

	dir, match := filepath.Split(prefix)

	var f *os.File
	var err error
	if dir == "" {
		f, err = os.Open(".")
	} else {
		f, err = os.Open(dir)
	}
	if err != nil {
		return
	}
	entries, err := f.Readdir(0)
	dprintf("got entries=%v for dir=%s (prefix=%s, match=%s)", entries, dir, prefix, match)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), match) {
			sep := " "
			if entry.IsDir() {
				sep = string(os.PathSeparator)
			}
			path := dir + entry.Name() + sep
			completions = append(completions, path)
		}
	}
	return
}
