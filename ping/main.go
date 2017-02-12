package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	count    = flag.Int("c", 0, "Stop after sending and receiving `count` packets.")
	interval = flag.Duration("i", 1*time.Second, "Wait `interval` between sending each packet.")
	size     = flag.Int("s", 56, "Size of packets in `bytes`.")
)

func main() {
	log.SetPrefix("ping: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s host\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	addrs, err := net.LookupHost(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	addr := addrs[0]
	ip := net.ParseIP(addr)

	var (
		listenNet  string
		listenAddr string
		echoType   icmp.Type
	)
	if ip.To4() != nil {
		listenNet = "udp4"
		listenAddr = "0.0.0.0"
		echoType = ipv4.ICMPTypeEcho
	} else {
		listenNet = "udp6"
		listenAddr = "::"
		echoType = ipv6.ICMPTypeEchoRequest
	}

	conn, err := icmp.ListenPacket(listenNet, listenAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	t := time.NewTicker(*interval)
	for seq := 1; true; seq++ {
		if err := ping(conn, ip, echoType, seq); err != nil {
			log.Print(err)
		}
		if *count == 0 || seq < *count {
			<-t.C
		} else {
			return
		}
	}
}

var readbuf = make([]byte, 1500)

func ping(conn *icmp.PacketConn, ip net.IP, echoType icmp.Type, seq int) error {
	buf, err := (&icmp.Message{
		Type: echoType,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: make([]byte, *size),
		},
	}).Marshal(nil)
	if err != nil {
		return err
	}
	_, err = conn.WriteTo(buf, &net.UDPAddr{IP: ip})
	if err != nil {
		return err
	}
	start := time.Now()
	n, peer, err := conn.ReadFrom(readbuf)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	fmt.Printf("%d bytes from %v: seq=%d time=%v\n", n, peer, seq, elapsed)

	return nil
}
