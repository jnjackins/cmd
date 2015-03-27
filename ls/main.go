package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"sigint.ca/group"
)

var (
	lflag = flag.Bool("l", false, "List in long format")
	pflag = flag.Bool("p", false, "Print only the final path element of each file name")
	sflag = flag.Bool("s", false, "Give size in KB for each entry")
)

var noargs bool

func main() {
	elog := log.New(os.Stderr, "ls: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ls [options] [file ...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
		noargs = true
	}
	for _, dir := range args {
		for dir[len(dir)-1] == '/' {
			dir = dir[:len(dir)-1]
		}
		f, err := os.Open(dir)
		if err != nil {
			elog.Fatal(err)
		}
		fi, err := f.Readdir(0)
		if err != nil {
			elog.Fatal(err)
		}
		for i := range fi {
			info := fi[i]
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
			if noargs == false && !*pflag {
				fmt.Print(dir + "/")
			}
			fmt.Println(info.Name())
		}
	}
}
