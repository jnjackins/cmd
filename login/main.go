package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"sigint.ca/user/passwd"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	motd    = "/etc/motd"
	timeout = 60 * time.Second
)

var (
	fflag = flag.Bool("f", false, "Don't authenticate.")
)

func main() {
	log.SetPrefix("login: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f] [username]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var username string
	if flag.NArg() == 1 {
		username = flag.Arg(0)
	} else if flag.NArg() == 0 {
		for username == "" {
			s, err := readUser()
			if err == io.EOF {
				os.Exit(1)
			} else if err != nil {
				log.Fatal(err)
			}
			username = s
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
	entry, err := passwd.GetEntry(username)
	if err != nil {
		log.Fatal(err)
	}
	if entry == nil {
		if !*fflag {
			readPw()
		}
		fail()
	}

	if *fflag || authenticate(entry) {
		if err := login(entry); err != nil {
			log.Fatal(err)
		}
	} else {
		fail()
	}
}

func fail() {
	time.Sleep(3 * time.Second)
	fmt.Println("Login incorrect")
	os.Exit(1)
}

func readUser() (string, error) {
	fmt.Print("login: ")
	s, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return s[:len(s)-1], nil
}

func authenticate(u *passwd.Entry) bool {
	hash := u.PasswordHash
	if hash == "" {
		return true
	}
	pw, err := readPw()
	if err != nil {
		log.Fatalf("failed to read password: %v", err)
	}
	return passwd.Authenticate(pw)
}

func readPw() (string, error) {
	done := make(chan bool, 1)
	go func() {
		select {
		case <-done:
		case <-time.After(timeout):
			log.Fatal("timed out")
		}
	}()

	fmt.Print("password: ")
	buf, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	done <- true
	fmt.Println()
	return string(buf), nil
}

func login(u *passwd.Entry) error {
	f, err := os.Open(motd)
	if err == nil {
		io.Copy(os.Stdout, f)
	}
	f.Close()

	shell := exec.Cmd{
		Path: u.Shell,
		Args: []string{"-" + filepath.Base(u.Shell)},
		Env: []string{
			"USER=" + u.Username,
			"HOME=" + u.Homedir,
			"SHELL=" + u.Shell,
			"TERM=vt100",
		},
		Dir: u.Homedir,
		SysProcAttr: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(u.Uid),
				Gid: uint32(u.Gid),
				//Groups: TODO,
			},
		},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		os.Stdin.Chown(u.Uid, u.Gid)
		os.Stdin.Chmod(os.FileMode(0620))
	}

	return shell.Run()
}
