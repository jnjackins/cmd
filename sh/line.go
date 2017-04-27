package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	lstate.SetCompleter(dumbComplete)
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

// TODO: cache $PATH contents
func dumbComplete(line string) []string {
	fields := strings.Fields(line)
	if len(fields) > 1 {
		return nil
	}
	var prefix string
	if len(fields) == 0 {
		prefix = ""
	} else {
		prefix = fields[0]
	}
	var matches []string
	paths := filepath.SplitList(env["PATH"])
	for _, d := range paths {
		if d == "" {
			d = "."
		}
		fi, err := os.Open(d)
		if err != nil {
			continue
		}
		names, err := fi.Readdirnames(0)
		for _, name := range names {
			if strings.HasPrefix(name, prefix) {
				matches = append(matches, name+" ")
			}
		}
	}
	return matches
}
