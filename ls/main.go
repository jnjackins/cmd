package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"time"

	"sigint.ca/group"
	"sigint.ca/text/tabwriter"
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
	dflag = flag.Bool("d", false, "If argument is a directory, list it, not its contents.")
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

	out := bufio.NewWriter(os.Stdout)
	tabw := tabwriter.NewWriter(out, 0, 0, 0, ' ', 0)
	for _, path := range args {
		dir := path
		stat, err := os.Lstat(path)
		if err != nil {
			elog.Print(err)
			status++
			continue
		}
		if stat.IsDir() && !*dflag {
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
			for _, stat := range fi {
				isSym := stat.Mode()&os.ModeType == os.ModeSymlink
				if noargs {
					dir = ""
				}
				ls(tabw, stat, dir, readlink(path+"/"+stat.Name(), isSym))
			}
		} else {
			isSym := stat.Mode()&os.ModeType == os.ModeSymlink
			if path != stat.Name() {
				dir = filepath.Dir(path)
			} else {
				dir = ""
			}
			ls(tabw, stat, dir, readlink(path, isSym))
		}
	}
	tabw.Flush()
	out.Flush()
	os.Exit(status)
}

// returns the destination of the named symbolic link, or an empty string
func readlink(name string, isSym bool) string {
	if !isSym {
		return ""
	}
	contents, err := os.Readlink(name)
	if err != nil {
		return ""
	}
	return contents
}

func ls(w io.Writer, info os.FileInfo, dir, target string) {
	if *sflag {
		fmt.Fprintf(w, "%d \t", info.Size()/1024+1) // +1 is sloppy round-up
	}
	if *lflag {
		// mode string
		fmt.Fprintf(w, "%s \t", modeString(info.Mode()))

		stat := info.Sys().(*syscall.Stat_t)

		// number of links
		fmt.Fprintf(w, "%d \t", stat.Nlink)

		// username
		uid := strconv.Itoa(int(stat.Uid))
		uname := uid
		u, err := user.LookupId(uid)
		if err == nil {
			uname = u.Username
		}
		fmt.Fprintf(w, "%s \t", uname)

		// groupname
		gid := strconv.Itoa(int(stat.Gid))
		gname, err := group.Name(gid)
		if err != nil {
			gname = gid
		}
		fmt.Fprintf(w, "%s \t", gname)

		// major, minor or size
		if info.Mode()&os.ModeDevice > 0 {
			major, minor := devNums(stat.Rdev)
			fmt.Fprintf(w, "%3d, %3d \t", major, minor)
		} else {
			fmt.Fprintf(w, "%d \t", info.Size())
		}

		// modified time
		var modtime string
		year := info.ModTime().Year()
		if year == time.Now().Year() {
			modtime = info.ModTime().Format("Jan _2 15:04")
		} else {
			modtime = info.ModTime().Format("Jan _2  2006")
		}
		fmt.Fprintf(w, "%s ", modtime)
	}

	if dir != "" && !*pflag {
		fmt.Fprint(w, filepath.Clean(dir+"/"+info.Name()))
	} else {
		fmt.Fprint(w, info.Name())
	}

	if target != "" {
		fmt.Fprintf(w, " ￫ %s", target)
	}

	fmt.Fprintln(w)
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
