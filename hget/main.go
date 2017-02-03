package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	method      = flag.String("method", "GET", "HTTP request method.")
	body        = flag.String("body", "", "Send `data` as the request body.")
	contentType = flag.String("content-type", "", "Content type of body, if any.")
)

func main() {
	log.SetPrefix("hget: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: hget [options] url")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	urlstring := flag.Arg(0)
	if !strings.Contains(urlstring, "://") {
		urlstring = "http://" + urlstring
	}
	url, err := url.Parse(urlstring)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(*method, url.String(), bytes.NewBufferString(*body))
	if err != nil {
		log.Fatal(err)
	}
	if *contentType != "" {
		req.Header.Set("Content-Type", *contentType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
