package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
)

var (
	fstab   = flag.String("fstab", "/etc/fstab", "Specify path to fstab(5).")
	fstype  = flag.String("type", "", "Specify filesystem type.")
	options = flag.String("options", "", "Specify mount options.")
)

type fsEntry struct {
	devpath string
	mntpt   string
	fstype  string
	options []string
}

func main() {
	log.SetPrefix("mount: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: mount [options] device [mountpoint]")
		fmt.Fprintln(os.Stderr, "       mount")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		if err := printMounts(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if flag.NArg() > 2 {
		flag.Usage()
		os.Exit(1)
	}

	entry, err := readEntry(flag.Arg(0))
	if err != nil {
		// couln't read fstab entry; assume first arg is devpath
		entry = &fsEntry{
			devpath: flag.Arg(0),
		}
	}

	if entry.mntpt == "" && flag.Arg(1) == "" {
		log.Fatal("mountpoint not specified")
	} else if flag.Arg(1) != "" {
		entry.mntpt = flag.Arg(1)
	}

	if *fstype != "" {
		entry.fstype = *fstype
	}
	if entry.fstype == "" {
		log.Fatal("filesystem type not specified")
	}
	if *options != "" {
		entry.options = strings.Split(*options, ",")
	}

	var flags int
	var data []string
	for _, s := range entry.options {
		switch s {
		case "async":
			flags |= syscall.MS_ASYNC
		case "atime":
			flags &^= syscall.MS_NOATIME
		case "noatime":
			flags |= syscall.MS_NOATIME
		case "defaults":
			// rw | suid | dev | exec | nouser | async
			flags = syscall.MS_ASYNC
		case "dev":
			flags &^= syscall.MS_NODEV
		case "nodev":
			flags |= syscall.MS_NODEV
		case "diratime":
			flags &^= syscall.MS_NODIRATIME
		case "nodiratime":
			flags |= syscall.MS_NODIRATIME
		case "dirsync":
			flags |= syscall.MS_DIRSYNC
		case "exec":
			flags &^= syscall.MS_NOEXEC
		case "iversion":
			flags |= syscall.MS_I_VERSION
		case "noiversion":
			flags &^= syscall.MS_I_VERSION
		case "mand":
			flags |= syscall.MS_MANDLOCK
		case "nomand":
			flags &^= syscall.MS_MANDLOCK
		case "relatime":
			flags |= syscall.MS_RELATIME
		case "norelatime":
			flags &^= syscall.MS_RELATIME
		case "strictatime":
			flags |= syscall.MS_STRICTATIME
		case "nostrictatime":
			flags &^= syscall.MS_STRICTATIME
		case "suid":
			flags &^= syscall.MS_NOSUID
		case "nosuid":
			flags |= syscall.MS_NOSUID
		case "remount":
			flags |= syscall.MS_REMOUNT
		case "ro":
			flags |= syscall.MS_RDONLY
		case "rw":
			flags &^= syscall.MS_RDONLY
		case "sync":
			flags |= syscall.MS_SYNC
		case "user":
			flags &^= syscall.MS_NOUSER
		case "nouser":
			flags |= syscall.MS_NOUSER
		default:
			// assume filesystem-specific option
			data = append(data, s)
		}
	}
	err = syscall.Mount(entry.devpath, entry.mntpt, entry.fstype, uintptr(flags), strings.Join(data, ","))
	if err != nil {
		log.Fatal(err)
	}
}

func printMounts() error {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(os.Stdout, f); err != nil {
		return err
	}
	return nil
}

func readEntry(arg string) (*fsEntry, error) {
	f, err := os.Open(*fstab)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var e fsEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		if arg != fields[0] && arg != fields[1] {
			continue
		}
		e.devpath = fields[0]
		e.mntpt = fields[1]
		e.fstype = fields[2]
		e.options = strings.Split(fields[3], ",")
		return &e, nil
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("no entry found")
}
