package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	aflag = flag.Bool("a", false, "Change only access time.")
	cflag = flag.Bool("c", false, "Don't create files.")
)

func main() {
	log.SetPrefix("touch: ")
	log.SetFlags(0)
	flag.Parse()

	var errc int
	for _, path := range flag.Args() {
		fi, err := os.Stat(path)
		if os.IsNotExist(err) {
			if !*cflag {
				f, err := os.Create(path)
				if err != nil {
					log.Println(err)
					errc++
					continue
				}
				f.Close()
			}
			// quietly succeed if -c is set and the file does not exist
			continue
		} else if err != nil {
			log.Println(err)
			errc++
			continue
		}
		atime := time.Now()
		mtime := atime
		if *aflag {
			if err != nil {
				log.Println(err)
				errc++
				continue
			}
			mtime = fi.ModTime()
		}
		if err := os.Chtimes(path, atime, mtime); err != nil {
			log.Println(err)
			errc++
			continue
		}
	}
	os.Exit(errc)
}
