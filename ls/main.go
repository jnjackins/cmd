package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"sigint.ca/die"
	"sigint.ca/group"
)

var (
	longflag = flag.Bool("l", false, "List in long format")
	sizeflag = flag.Bool("s", false, "Give size in KB for each entry")
)
var noargs bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [file ...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
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
		die.On(err, "ls: error opening "+dir)
		fi, err := f.Readdir(0)
		die.On(err, "ls: error reading directory "+dir)
		for i := range fi {
			info := fi[i]
			if *sizeflag {
				fmt.Printf("%4d ", info.Size()/1024+1) // +1 is sloppy round-up
			}
			if *longflag {
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
			if noargs == false {
				fmt.Print(dir + "/")
			}
			fmt.Println(info.Name())
		}
	}
}
