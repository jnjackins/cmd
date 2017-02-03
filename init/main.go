package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	sys "golang.org/x/sys/unix"
)

const (
	ntty    = 6
	shell   = "/bin/sh"
	runc    = "/etc/rc"
	dev     = "/dev/"
	console = "/dev/console"
	login   = "/bin/login"
)

var cons *os.File

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {
	signal.Ignore(sys.SIGINT, sys.SIGTERM, sys.SIGHUP)

	f, err := os.OpenFile(console, os.O_RDWR, 0)
	if err != nil {
		panic("no console")
	}
	cons = f

	log.SetPrefix("init: ")
	log.SetOutput(cons)
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	for {
		run()
	}
}

func run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v\n", r)
		}
	}()

	type stage struct {
		name string
		fn   func() error
	}
	for _, s := range []stage{
		// In case of a panic recovery, starting with
		// the shutdown phase ensures proper cleanup.
		{"shutdown", shutdown},
		{"single-user", single},
		{"runcom", runcom},
		{"multi-user", multiple},
	} {
		log.Printf("starting phase %s\n", s.name)
		if err := s.fn(); err != nil {
			log.Printf("stage %s: %v\n", s.name, err)
			return
		}
	}
}

func shutdown() error {
	log.Println("killing tty sessions")
	for _, s := range sessions {
		s.close()
	}

	log.Println("killing all processes")
	for i := 0; i < 5; i++ {
		sys.Kill(-1, sys.SIGKILL)
	}
	for wait(-1) == nil {
	}

	log.Println("closing open file descriptors")
	for i := 0; i < 10; i++ {
		if i != int(cons.Fd()) {
			sys.Close(i)
		}
	}
	return nil
}

func wait(pid int) error {
	_, err := sys.Wait4(pid, nil, 0, nil)
	return err
}

func single() error {
	log.Println("press ctl-d to proceed to multi-user mode")
	return (&exec.Cmd{
		Path:   shell,
		Args:   []string{"-sh", "+m"},
		Stdin:  cons,
		Stdout: cons,
		Stderr: cons,
	}).Run()
}

func runcom() error {
	return (&exec.Cmd{
		Path: shell,
		Args: []string{shell, runc},
	}).Run()
}

func multiple() error {
	// shutdown is initiated by SIGHUP
	c := make(chan os.Signal, 1)
	signal.Notify(c, sys.SIGHUP)

	// start TTY sessions
	sessionWg = sync.WaitGroup{}
	for i := 1; i < ntty; i++ {
		log.Printf("starting session on tty%d\n", i+1)
		sessions[i] = newSession(i)
		go sessions[i].getty()
	}

	// wait for SIGHUP
	<-c
	log.Println("received SIGHUP")
	signal.Ignore(sys.SIGHUP)

	// close all sessions
	for i := 1; i < ntty; i++ {
		go sessions[i].close()
		sessions[i] = nil
	}
	log.Println("waiting for sessions to exit")
	sessionWg.Wait()

	return nil
}

var (
	sessions  [ntty]*session
	sessionWg sync.WaitGroup
)

type session struct {
	tty  string
	quit chan bool
	proc *os.Process
}

func newSession(i int) *session {
	sessionWg.Add(1)
	return &session{
		tty:  fmt.Sprintf("tty%d", i+1),
		quit: make(chan bool, 1),
	}
}

func (s *session) String() string {
	if s == nil {
		return "nil session"
	}
	return s.tty + ".session"
}

func (s *session) getty() {
	for {
		select {
		case <-s.quit:
			log.Printf("%v: exiting\n", s)
			sessionWg.Done()
			return
		default:
			path := dev + s.tty
			f, err := os.OpenFile(path, os.O_RDWR, 0)
			if err != nil {
				log.Printf("%v: open %s: %v\n", s, path, err)
				break
			}
			f.Chown(0, 0)
			f.Chmod(0620)
			cmd := &exec.Cmd{
				Path:   login,
				Args:   []string{"login"},
				Stdin:  f,
				Stdout: f,
				Stderr: f,
				// start a new session and set the controlling TTY
				SysProcAttr: &syscall.SysProcAttr{
					Setsid:  true,
					Setctty: true,
					Ctty:    int(f.Fd()),
				},
			}
			if err := cmd.Start(); err != nil {
				log.Printf("%v: start login: %v\n", s, err)
				f.Close()
				// we're in trouble; return to single-user mode.
				sys.Kill(1, sys.SIGHUP)
				return
			}
			s.proc = cmd.Process
			f.Close()
			if err := cmd.Wait(); err != nil {
				log.Printf("%v: wait login: %v\n", s, err)
			}
		}
	}
}

func (s *session) close() {
	if s == nil {
		return
	}
	if s.quit != nil {
		s.quit <- true
	}
	if s.proc != nil {
		s.proc.Kill()
	}
}
