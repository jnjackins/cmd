package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"strings"

	"path/filepath"

	"golang.org/x/sys/unix"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: kmod load <module> [params]")
	fmt.Fprintln(os.Stderr, "       kmod unload <module>")
	fmt.Fprintln(os.Stderr, "       kmod list")

	os.Exit(1)
}

func main() {
	log.SetPrefix("kmod: ")
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "load":
		var params string
		if len(os.Args) == 4 {
			params = os.Args[3]
		} else if len(os.Args) != 3 {
			usage()
		}
		if err := load(os.Args[2], params); err != nil {
			log.Fatal(err)
		}
	case "unload":
		if len(os.Args) != 3 {
			usage()
		}
		if err := unload(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "list":
		if len(os.Args) != 2 {
			usage()
		}
		if err := list(); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func load(name, params string) error {
	var path string
	fi, err := os.Stat(name)
	if err == nil && fi.Mode().IsRegular() {
		path = name
		name = filepath.Base(name)
		if strings.HasSuffix(name, ".ko") {
			name = name[:len(name)-3]
		}
	} else {
		path, err = getModulePath(name)
		if err != nil {
			return err
		}
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	n := len(buf)
	_, _, errno := unix.Syscall(unix.SYS_INIT_MODULE,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(n),
		uintptr(unsafe.Pointer(&params)))
	if errno == unix.EEXIST {
		return fmt.Errorf("module %s already loaded", name)
	} else if errno != 0 {
		return errno
	}
	return nil
}

func getModulePath(name string) (string, error) {
	modroot := fmt.Sprintf("/lib/modules/%s/", release())
	f, err := os.Open(modroot + "modules.order")
	if err != nil {
		return "", err
	}
	defer f.Close()

	filename := strings.Replace(name, "-", "_", -1) + ".ko"
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		path := scanner.Text()
		if strings.HasSuffix(path, filename) {
			return modroot + path, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("module %s not found", name)
}

func release() string {
	var buf unix.Utsname
	if err := unix.Uname(&buf); err != nil {
		panic(err)
	}
	return intsToString(buf.Release[:])
}

func intsToString(a []int8) string {
	var buf bytes.Buffer
	for _, c := range a {
		if c == 0 {
			break
		}
		buf.WriteByte(byte(c))
	}
	return buf.String()
}

func unload(name string) error {
	namep, err := unix.BytePtrFromString(name)
	if err != nil {
		return err
	}
	flags := uintptr(unix.O_NONBLOCK)
	_, _, errno := unix.Syscall(unix.SYS_DELETE_MODULE,
		uintptr(unsafe.Pointer(namep)),
		flags, 0)
	if errno == unix.ENOENT {
		return fmt.Errorf("no such module %s", name)
	} else if errno != 0 {
		return errno
	}
	return nil
}

func list() error {
	f, err := os.Open("/proc/modules")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("%-19s %8s  Used by\n", "Module", "Size")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		fmt.Printf("%-19s %8s  %s\n", fields[0], fields[1], fields[2])
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
