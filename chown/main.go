package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"sigint.ca/group"
	"sigint.ca/user"
)

var rflag = flag.Bool("r", false, "Recursively change ownership of all files rooted at the given paths.")

const usage = `usage: chown [-r] owner[:group] file ...
       chown [-r] :group file ...`

func main() {
	elog := log.New(os.Stderr, "chown: ", 0)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	f := strings.SplitN(args[0], ":", 2)
	uname := f[0]
	var gname string
	if len(f) == 2 {
		gname = f[1]
	}

	var u *user.User
	var err error
	if uname != "" {
		u, err = user.Lookup(uname)
		if err != nil {
			elog.Fatalf("error looking up user %#v", uname)
		}
	}

	var g *group.Group
	if gname != "" {
		g, err = group.Lookup(gname)
		if err != nil {
			elog.Fatalf("error looking up group %#v", gname)
		}
	}

	if u == nil && g == nil {
		// no-op
		return
	}

	for _, name := range args[1:] {
		err := chown(name, u, g)
		if err != nil {
			elog.Fatal(err)
		}
	}
}

func currentUid(fi os.FileInfo) int {
	return int(fi.Sys().(*syscall.Stat_t).Uid)
}

func currentGid(fi os.FileInfo) int {
	return int(fi.Sys().(*syscall.Stat_t).Gid)
}

func chown(name string, u *user.User, g *group.Group) error {
	chownFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var uid, gid int
		if u != nil {
			uid = u.Uid
		} else {
			uid = currentUid(fi)
		}
		if g != nil {
			gid = g.Id
		} else {
			gid = currentGid(fi)
		}
		return os.Chown(path, uid, gid)
	}
	if *rflag {
		return filepath.Walk(name, chownFn)
	} else {
		fi, err := os.Stat(name)
		return chownFn(name, fi, err)
	}
}
