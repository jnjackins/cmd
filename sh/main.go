//go:generate go tool yacc -p "sh" parse.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/peterh/liner"
)

var (
	eflag = flag.Bool("e", true, "enable line editing")
	lflag = flag.String("l", "", "read commands from `file` before reading normal input")
)

var cooked, raw liner.ModeApplier

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	setupSignals()
	setupEnv()
	loadProfile()
	flag.Parse()
	if *lflag != "" {
		f, err := os.Open(*lflag)
		if err != nil {
			log.Print(err)
		} else {
			parseDumb(f, false)
		}
	}
	args := flag.Args()
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		env["*"] = strings.Join(args[1:], " ")
		updateEnv()
		parseDumb(f, false)
	} else if !*eflag || !isTTY(os.Stdin) {
		parseDumb(os.Stdin, isTTY(os.Stdin))
	} else {
		parse()
	}
}

func loadProfile() {
	path := os.Getenv("home") + "/lib/profile"
	f, err := os.Open(path)
	if err != nil {
		return
	}
	parseDumb(f, false)
	f.Close()
}

func parseDumb(r io.Reader, tty bool) {
	in := bufio.NewScanner(bufio.NewReader(r))
	for in.Scan() {
		if tty {
			fmt.Print(os.Getenv("prompt"))
		}
		line := in.Text() + "\n"
		shParse(&shLex{line: line})
	}
	if err := in.Err(); err != nil {
		log.Print(err)
	}
}

func parse() {
	prompt, err := setupLineEditing()
	if err != nil {
		log.Fatal(err)
	}
	defer prompt.Close()
	for {
		line, err := prompt.Prompt(os.Getenv("prompt"))
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Print(err)
		} else {
			prompt.AppendHistory(line)
			setCooked()
			shParse(&shLex{line: line + "\n"})
			setRaw()
		}
	}
}

func setupLineEditing() (*liner.State, error) {
	var err error
	cooked, err = liner.TerminalMode()
	if err != nil {
		return nil, err
	}
	s := liner.NewLiner()
	s.SetWordCompleter(complete)
	s.SetTabCompletionStyle(liner.TabPrints)
	raw, err = liner.TerminalMode()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func setCooked() {
	if err := cooked.ApplyMode(); err != nil {
		log.Print(err)
	}
}

func setRaw() {
	if err := raw.ApplyMode(); err != nil {
		log.Print(err)
	}
}

func complete(line string, pos int) (string, []string, string) {
	runes := []rune(line)
	head := runes[:pos]
	word := ""
	if len(head) > 0 && !unicode.IsSpace(head[len(head)-1]) {
		fields := strings.Fields(string(head))
		word = fields[len(fields)-1]
		head = head[:len(head)-len([]rune(word))]
	}
	completions, _ := filepath.Glob(word + "*")
	var prefix string
	if strings.HasPrefix(word, "./") {
		prefix = "./"
	}
	for i := range completions {
		completions[i] = prefix + completions[i]
		stat, err := os.Stat(completions[i])
		if err == nil && stat.IsDir() {
			completions[i] += "/"
		}
	}
	return string(head), completions, string(runes[pos:])
}
