package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

//go:generate goyacc -o y.out.go -p "sh" syntax.y
//go:generate rm y.output

var (
	cflag = flag.String("c", "", "Read comands from `string`.")
	debug = flag.Bool("d", false, "Enable debug mode.")
)

func main() {
	flag.Parse()

	if *debug {
		shDebug = 3
	}

	log.SetPrefix("sh: ")
	log.SetFlags(0)

	setupEnv()

	if *cflag != "" {
		setInput(bytes.NewBufferString(*cflag + "\n"))
	} else {
		initPrompt()
	}
	for {
		line, err := getLine()
		if err == io.EOF {
			exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		shParse(&shLex{line: line})
	}
}

func dprintf(format string, args ...interface{}) {
	if *debug {
		log.Printf("debug: "+format, args...)
	}
}

func exit(i int) {
	fixTerminal()
	os.Exit(i)
}
