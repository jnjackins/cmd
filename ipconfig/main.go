package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/docker/libcontainer/netlink"
)

func main() {
	log.SetPrefix("ipconfig: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [interface [verb [arg]]]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Verbs:\n")
		fmt.Fprintf(os.Stderr, "\tup|down\n")
		fmt.Fprintf(os.Stderr, "\tadd|remove <cidr>\n")
	}
	flag.Parse()

	if flag.NArg() == 0 {
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}
		for _, iface := range ifaces {
			if err := print(&iface); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}

	iface, err := net.InterfaceByName(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() == 1 {
		if err := print(iface); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	cmd := flag.Arg(1)
	var cmderr error
	switch cmd {
	case "up":
		cmderr = netlink.NetworkLinkUp(iface)
	case "down":
		cmderr = netlink.NetworkLinkDown(iface)
	case "add", "remove":
		if flag.NArg() < 3 {
			flag.Usage()
			os.Exit(1)
		}
		ip, ipnet, err := net.ParseCIDR(flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		f := netlink.NetworkLinkAddIp
		if cmd == "remove" {
			f = netlink.NetworkLinkDelIp
		}
		cmderr = f(iface, ip, ipnet)
	default:
		log.Fatalf("invalid verb %s", flag.Arg(1))
	}
	if cmderr != nil {
		log.Fatalf("%s: %v", cmd, cmderr)
	}
}

func print(iface *net.Interface) error {
	fmt.Printf("%s: flags=%v mtu %d\n", iface.Name, iface.Flags, iface.MTU)
	if iface.HardwareAddr != nil {
		fmt.Printf("\tether %v\n", iface.HardwareAddr)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		fmt.Printf("\taddr %v\n", addr)
	}
	return nil
}
