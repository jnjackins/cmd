package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"sigint.ca/user"
	"sigint.ca/user/passwd"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

const (
	authkeys = ".ssh/authorized_keys"
	login    = "/bin/login"
)

var (
	hostkey    = flag.String("key", "/etc/ssh/ssh_host_rsa_key", "Path to the host private key.")
	listenAddr = flag.String("listen", ":22", "Listen address.")
	logPath    = flag.String("log", "/dev/stderr", "Send logs to `file`.")
)

var logWriter io.Writer

func main() {
	flag.Parse()

	logWriter = os.Stderr
	if *logPath != "/dev/stderr" {
		log.Printf("writing logs to %s", *logPath)
		f, err := os.Create(*logPath)
		if err != nil {
			log.Fatal(err)
		}
		logWriter = f
		log.SetOutput(logWriter)
	}

	// set up authentication methods
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			user := c.User()
			log.Printf("attempting to authenticate user %s by public key", user)
			authorized, err := getAuthorizedKeys(user)
			if err != nil {
				return nil, fmt.Errorf("error getting authorized keys %s: %v", user, err)
			}
			if authorized[string(pubKey.Marshal())] {
				return nil, nil
			}
			return nil, fmt.Errorf("unknown public key for %s", user)
		},
		PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			user := c.User()
			log.Printf("attempting to authenticate user %s by password", user)
			e, err := passwd.GetEntry(user)
			if err != nil {
				return nil, fmt.Errorf("error getting passwd entry for %s: %v", user, err)
			}
			if e.Authenticate(string(password)) {
				return nil, nil
			}
			return nil, fmt.Errorf("bad password for %s", user)
		},
	}

	// load host key
	buf, err := ioutil.ReadFile(*hostkey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	private, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	config.AddHostKey(private)

	// listen for connections
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen for connection: %v", err)
	}
	log.Printf("listening on %s", *listenAddr)

	// authenticate and handle connections
	for n := 0; true; n++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection: %v", err)
			continue
		}
		log.Printf("incoming connection from %s", conn.RemoteAddr())

		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Printf("failed to handshake: %v", err)
			continue
		}
		log.Printf("handshake succeeded: [%d]", n)
		go handleSession(sshConn, chans, reqs, n)
	}
}

func handleSession(conn *ssh.ServerConn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, n int) {
	// distinguish between sessions in log messages
	slog := log.New(logWriter, fmt.Sprintf("[%d] ", n), log.LstdFlags)
	if err := handleSession2(conn, chans, reqs, slog); err != nil {
		slog.Printf("handle session: %v", err)
	}
}

func handleSession2(conn *ssh.ServerConn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, slog *log.Logger) error {
	//  we don't do anything with these, for now
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return fmt.Errorf("Could not accept channel: %v", err)
		}
		slog.Printf("accepted channel: %s", newChannel.ChannelType())

		go func(in <-chan *ssh.Request) {
			var env []string
			for req := range in {
				slog.Printf("received request: %s", req.Type)
				switch req.Type {
				case "pty-req":
					req.Reply(true, nil)
				case "env":
					var envMsg struct {
						K, V string
					}
					if err := ssh.Unmarshal(req.Payload, &envMsg); err != nil {
						slog.Printf("unmarshal env: %v", err)
						req.Reply(false, nil)
						break
					}
					envString := envMsg.K + "=" + envMsg.V
					env = append(env, envString)
					req.Reply(true, nil)
				case "shell":
					req.Reply(true, nil)
					if err := shell(channel, conn.User(), env); err != nil {
						slog.Printf("shell: %v", err)
					}
				case "exec":
					var cmdMsg struct {
						Cmd string
					}
					if err := ssh.Unmarshal(req.Payload, &cmdMsg); err != nil {
						slog.Printf("unmarshal cmd: %v", err)
						req.Reply(false, nil)
						break
					}
					req.Reply(true, nil)
					if err := execCmd(channel, cmdMsg.Cmd, conn.User(), env); err != nil {
						slog.Printf("exec: %v", err)
					}
				}
			}
			slog.Println("closing connection")
		}(requests)
	}
	return nil
}

func execCmd(channel ssh.Channel, cmdString, username string, env []string) error {
	defer channel.Close()

	u, err := user.Lookup(username)
	if err != nil {
		return err
	}
	//groups, err := u.GroupIds()
	// if err != nil {
	//	return err
	// }

	entry, err := passwd.GetEntry(username)
	if err != nil {
		return fmt.Errorf("get passwd entry: %v", err)
	}

	cmd := exec.Cmd{
		Path:   entry.Shell,
		Args:   []string{entry.Shell, "-c", cmdString},
		Env:    env,
		Dir:    u.HomeDir,
		Stdout: channel,
		Stderr: channel.Stderr(),
		SysProcAttr: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(u.Uid),
				Gid: uint32(u.Gid),
				//Groups: u.GroupIds(),
			},
		},
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		io.Copy(stdin, channel)
		stdin.Close()
	}()

	cmd.Run()
	code := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if err := sendExitCode(channel, code); err != nil {
		return err
	}
	return nil
}

func shell(channel ssh.Channel, username string, env []string) error {
	defer channel.Close()

	cmd := exec.Command("login", "-f", username)
	cmd.Env = env
	pseudo, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("start login: %v", err)
	}
	defer pseudo.Close()
	go io.Copy(pseudo, channel)
	go io.Copy(channel, pseudo)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait login: %v", err)
	}
	code := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if err := sendExitCode(channel, code); err != nil {
		return err
	}
	return nil
}

func sendExitCode(channel ssh.Channel, code int) error {
	buf := ssh.Marshal(struct{ code uint32 }{uint32(code)})
	_, err := channel.SendRequest("exit-status", false, buf)
	if err != nil {
		return fmt.Errorf("send exit-status failed: %v", err)
	}
	return nil
}

func getAuthorizedKeys(username string) (map[string]bool, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(u.HomeDir, authkeys)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	authorized := make(map[string]bool)
	for len(buf) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(buf)
		if err != nil {
			return nil, err
		}
		authorized[string(pubKey.Marshal())] = true
		buf = rest
	}
	return authorized, nil
}
