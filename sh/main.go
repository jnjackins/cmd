// TODO: handle sigint correctly

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

	"github.com/jnjackins/liner"
)

var eflag = flag.Bool("e", true, "enable line editing")
var lflag = flag.String("l", "", "read commands from `file` before reading normal input")

func main() {
	log.SetPrefix("sh: ")
	log.SetFlags(0)
	setupSignals()
	setupEnv()
	flag.Parse()
	if *lflag != "" {
		f, err := os.Open(*lflag)
		if err != nil {
			log.Print(err)
		} else {
			parseDumb(f, false)
		}
	}
	if !*eflag || !isTTY(os.Stdin) {
		parseDumb(os.Stdin, isTTY(os.Stdin))
	} else {
		parse()
	}
}

func parseDumb(r io.Reader, tty bool) {
	in := bufio.NewReader(r)
	for {
		if tty {
			fmt.Print(os.Getenv("prompt"))
		}
		line, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		shParse(&shLex{line: line})
	}
}

func parse() {
	prompt := liner.NewLiner()
	defer prompt.Close()
	for {
		prompt.SetWordCompleter(completer)
		line, err := prompt.Prompt(os.Getenv("prompt"))
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Print(err)
		} else {
			prompt.AppendHistory(line)
			prompt.Stop()
			shParse(&shLex{line: line + "\n"})
			prompt.Start()
		}
	}
}

func completer(line string, pos int) (string, []string, string) {
	runes := []rune(line)
	head := runes[:pos]
	word := ""
	if len(head) > 0 && !unicode.IsSpace(head[len(head)-1]) {
		fields := strings.Fields(string(head))
		word = fields[len(fields)-1]
		head = head[:len(head)-len([]rune(word))]
	}
	completions, _ := filepath.Glob(word + "*")
	if strings.HasPrefix(word, "./") {
		for i := range completions {
			completions[i] = "./" + completions[i]
		}
	}
	return string(head), completions, string(runes[pos:])
}
