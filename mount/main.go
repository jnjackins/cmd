package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"
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
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
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

	var flags uintptr
	for _, s := range entry.options {
		switch s {
		case "ro":
			flags |= unix.MS_RDONLY
		case "rw":
			// rw is default
		case "remount":
			flags |= unix.MS_REMOUNT
		default:
			log.Printf("unrecognized option %q", s)
		}
	}
	if err := unix.Mount(entry.devpath, entry.mntpt, entry.fstype, flags, ""); err != nil {
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

func readEntry(devpath string) (*fsEntry, error) {
	// default is rw if no entry
	e := fsEntry{
		devpath: devpath,
		options: []string{"rw"},
	}

	f, err := os.Open(*fstab)
	if err != nil {
		return &e, err
	}
	defer f.Close()

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
		if fields[0] != devpath {
			continue
		}
		e.mntpt = fields[1]
		e.fstype = fields[2]
		e.options = strings.Split(fields[3], ",")
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &e, nil
}
