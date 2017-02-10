package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	gflag = flag.Bool("g", false, "Global substitution.")
)

func main() {
	log.SetPrefix("sub: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <pattern> <replace>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	pattern, err := regexp.Compile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	replace := flag.Arg(1)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(sub(scanner.Text(), pattern, replace, *gflag))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func sub(line string, re *regexp.Regexp, replace string, global bool) string {
	if global {
		var final string
		for {
			match := re.FindStringSubmatchIndex(line)
			if match == nil {
				final += line
				break
			}
			result, rsize := sub1(line, match, re, replace)
			line = line[match[1]:]
			final += result[:match[0]+rsize]
		}
		return final
	}
	match := re.FindStringSubmatchIndex(line)
	result, _ := sub1(line, match, re, replace)
	return result
}

func sub1(line string, match []int, re *regexp.Regexp, replace string) (string, int) {
	if len(match) > 2 {
		// there was at least one submatch
		submatches := match[2:]
		for i := 0; i < len(submatches)/2; i++ {
			submatch := []int{submatches[2*i], submatches[2*i+1]}
			submatchIdentifier := `\` + strconv.Itoa(i+1)
			replace = strings.Replace(replace, submatchIdentifier, line[submatch[0]:submatch[1]], 1)
		}
	}
	line = line[:match[0]] + replace + line[match[1]:]
	return line, len(replace)
}
