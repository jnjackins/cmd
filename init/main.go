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
	"time"

	sys "golang.org/x/sys/unix"
)

const (
	ntty    = 6 // TODO: launch gettys automatically?
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
	type stage struct {
		name string
		fn   func() error
	}
	for _, s := range []stage{
		{"single-user", single},
		{"runcom", runcom},
		{"multi-user", multiple},
		{"shutdown", shutdown},
	} {
		log.Printf("starting phase %s", s.name)
		if err := s.fn(); err != nil {
			log.Printf("stage %s: %v", s.name, err)
			return
		}
	}
}

func shutdown() error {
	log.Println("killing all processes")
	for i := 0; i < 5; i++ {
		sys.Kill(-1, sys.SIGKILL)
	}
	var err error
	for err == nil {
		_, err = sys.Wait4(-1, nil, 0, nil)
	}

	return nil
}

func single() error {
	if _, err := exec.LookPath(shell); err != nil {
		panic("no shell at " + shell)
	}
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
		Path:   shell,
		Args:   []string{shell, runc},
		Stdout: cons,
		Stderr: cons,
	}).Run()
}

func multiple() error {
	// shutdown is initiated by SIGHUP
	c := make(chan os.Signal, 1)
	signal.Notify(c, sys.SIGHUP)

	quit := make(chan bool, 1)
	go reap(quit)

	// start TTY sessions
	sessionWg = sync.WaitGroup{}
	for i := 1; i < ntty; i++ {
		log.Printf("starting session on tty%d", i+1)
		sessions[i] = newSession(i)
		go sessions[i].getty()
	}

	<-c // wait for SIGHUP
	signal.Ignore(sys.SIGHUP)
	log.Println("received SIGHUP")

	// kill reap goroutine
	quit <- true

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
			log.Printf("%v: exiting", s)
			sessionWg.Done()
			return
		default:
			path := dev + s.tty
			f, err := os.OpenFile(path, os.O_RDWR, 0)
			if err != nil {
				log.Printf("%v: open %s: %v", s, path, err)
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
				log.Printf("%v: start login: %v", s, err)
				f.Close()
				// we're in trouble; return to single-user mode.
				sys.Kill(1, sys.SIGHUP)
				return
			}
			s.proc = cmd.Process
			f.Close()
			cmd.Wait()
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

func reap(quit <-chan bool) {
	log.Println("reap: starting")
	ticker := time.NewTicker(10 * time.Second)
	for {
		// reap all zombies
		for {
			pid, err := sys.Wait4(-1, nil, sys.WNOHANG, nil)
			if err != nil {
				log.Printf("reap: wait: %v", err)
				break
			}
			if pid == 0 {
				break
			}
		}

		// quit, or rest a bit
		select {
		case <-quit:
			ticker.Stop()
			log.Println("reap: exiting")
			return
		case <-ticker.C:
			continue
		}
	}
}
