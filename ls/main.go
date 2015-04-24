package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"sigint.ca/group"
	"sigint.ca/user"
)

type normalSort []os.FileInfo
type timeSort []os.FileInfo

func (a normalSort) Len() int           { return len(a) }
func (a normalSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a normalSort) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func (a timeSort) Len() int           { return len(a) }
func (a timeSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timeSort) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }

var (
	lflag = flag.Bool("l", false, "List in long format.")
	pflag = flag.Bool("p", false, "Print only the final path element of each file name.")
	rflag = flag.Bool("r", false, "Reverse the order of sort.")
	sflag = flag.Bool("s", false, "Give size in KB for each entry.")
	tflag = flag.Bool("t", false, "Sort by time modified (latest first) instead of by name.")
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

			rev := func(i sort.Interface) sort.Interface { return i }
			if *rflag {
				rev = sort.Reverse
			}
			if *tflag {
				sort.Sort(rev(timeSort(fi)))
			} else {
				sort.Sort(rev(normalSort(fi)))
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight)
			for i := range fi {
				if noargs {
					path = ""
				}
				ls(w, fi[i], path)
			}
			w.Flush()

		} else {
			if path != stat.Name() {
				path = filepath.Dir(path)
			} else {
				path = ""
			}
			ls(os.Stdout, stat, path)
		}
	}
	os.Exit(status)
}

func ls(w io.Writer, info os.FileInfo, path string) {
	if *sflag {
		fmt.Fprintf(w, "%d\t ", info.Size()/1024+1) // +1 is sloppy round-up
	}
	if *lflag {
		fmt.Fprintf(w, "%s\t ", modeString(info.Mode()))

		stat := info.Sys().(*syscall.Stat_t) // TODO: non-portable
		uid := strconv.Itoa(int(stat.Uid))
		uname := uid
		u, err := user.LookupId(uid)
		if err == nil {
			uname = u.Username
		}
		fmt.Fprintf(w, "%s\t ", uname)

		gid := strconv.Itoa(int(stat.Gid))
		gname, err := group.Name(gid)
		if err != nil {
			gname = gid
		}
		fmt.Fprintf(w, "%s\t ", gname)

		// major, minor
		if info.Mode()&os.ModeDevice > 0 {
			major, minor := devNums(stat.Rdev)
			fmt.Fprintf(w, "%3d, %3d\t ", major, minor)
		} else {
			fmt.Fprintf(w, "%d\t ", info.Size())
		}

		// modified time
		var modtime string
		year := info.ModTime().Year()
		if year == time.Now().Year() {
			modtime = info.ModTime().Format("Jan 02 15:04")
		} else {
			modtime = info.ModTime().Format("Jan 02  ") + strconv.Itoa(year)
		}

		fmt.Fprintf(w, "%s ", modtime)
	}

	if path != "" && !*pflag {
		fmt.Fprintln(w, filepath.Clean(path+"/"+info.Name()))
	} else {
		fmt.Fprintln(w, info.Name())
	}
}

func modeString(mode os.FileMode) string {
	modestr := []byte("-rwxrwxrwx")

	// type
	switch mode & os.ModeType {
	case os.ModeDevice:
		if mode&os.ModeCharDevice > 0 {
			modestr[0] = 'c'
		} else {
			modestr[0] = 'b'
		}
	case os.ModeDir:
		modestr[0] = 'd'
	case os.ModeSymlink:
		modestr[0] = 'l'
	case os.ModeNamedPipe:
		modestr[0] = 'p'
	case os.ModeSocket:
		modestr[0] = 's'
	}

	// permissions
	bit := os.FileMode(1 << 8)
	for i := range modestr[1:] {
		if mode&bit == 0 {
			modestr[1+i] = '-'
		}
		bit >>= 1
	}

	// special attributes
	if mode&os.ModeSetuid > 0 {
		if modestr[3] == 'x' {
			modestr[3] = 's'
		} else {
			modestr[3] = 'S'
		}
	}
	if mode&os.ModeSetgid > 0 {
		if modestr[6] == 'x' {
			modestr[6] = 's'
		} else {
			modestr[6] = 'S'
		}
	}
	if mode&os.ModeSticky > 0 {
		if modestr[9] == 'x' {
			modestr[9] = 't'
		} else {
			modestr[9] = 'T'
		}
	}

	return string(modestr)
}
