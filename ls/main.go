package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"sigint.ca/group"
	"sigint.ca/user"
)

var (
	lflag = flag.Bool("l", false, "List in long format")
	pflag = flag.Bool("p", false, "Print only the final path element of each file name")
	sflag = flag.Bool("s", false, "Give size in KB for each entry")
)

func main() {
	elog := log.New(os.Stderr, "ls: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ls [options] [file ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	var status int
	args := flag.Args()
	var noargs bool
	if len(args) == 0 {
		args = []string{"."}
		noargs = true
	}
	for _, path := range args {
		stat, err := os.Stat(path)
		if err != nil {
			elog.Print(err)
			status++
			continue
		}
		if stat.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				elog.Print(err)
				status++
				continue
			}
			fi, err := f.Readdir(0)
			if err != nil {
				elog.Print(err)
				status++
				continue
			}
			f.Close()
			for i := range fi {
				if noargs {
					path = ""
				}
				ls(fi[i], path)
			}
		} else {
			if path != stat.Name() {
				path = filepath.Dir(path)
			} else {
				path = ""
			}
			ls(stat, path)
		}
	}
	os.Exit(status)
}

func ls(info os.FileInfo, path string) {
	if *sflag {
		fmt.Printf("%4d ", info.Size()/1024+1) // +1 is sloppy round-up
	}
	if *lflag {
		modestr := []byte("-rwxrwxrwx")
		mode := info.Mode()
		switch mode & os.ModeType {
		case os.ModeDir:
			modestr[0] = 'd'
		case os.ModeSymlink:
			modestr[0] = 'l'
		}
		bit := os.FileMode(1 << 8)
		for i := range modestr[1:] {
			if mode&bit == 0 {
				modestr[1+i] = '-'
			}
			bit >>= 1
		}
		stat := info.Sys().(*syscall.Stat_t) // TODO non-portable
		uid := strconv.Itoa(int(stat.Uid))
		uname := uid
		u, err := user.LookupId(uid)
		if err == nil {
			uname = u.Username
		}
		gid := strconv.Itoa(int(stat.Gid))
		gname, err := group.Name(gid)
		if err != nil {
			gname = gid
		}
		fmt.Printf("%s %s %s %7d %s ", modestr, uname, gname, info.Size(), info.ModTime().Format("Jan 02 15:04"))
	}

	if path != "" && !*pflag {
		fmt.Println(filepath.Clean(path + "/" + info.Name()))
	} else {
		fmt.Println(info.Name())
	}
}
