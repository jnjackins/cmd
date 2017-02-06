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
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\tPrint link and routing information.")
		fmt.Fprintf(os.Stderr, "%s link\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s link <iface>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s link <iface> up|down\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s link <iface> add|del <addr>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\tConfigure network interfaces.")
		fmt.Fprintf(os.Stderr, "%s route\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s route add <dst> <gw>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s route del <dst>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\tConfigure Internet Protocol routes.")
	}
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Network Interfaces:")
		printIfaces()
		fmt.Println()
		fmt.Println("Routing Tables:")
		printRoutes()
		os.Exit(0)
	}

	switch flag.Arg(0) {
	case "link":
		linkCmd(flag.Args()[1:])
	case "route":
		routeCmd(flag.Args()[1:])
		os.Exit(0)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func linkCmd(args []string) {
	if len(args) == 0 {
		printIfaces()
		os.Exit(0)
	}

	// all other commands require iface
	iface, err := net.InterfaceByName(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 1 {
		if err := printIface(iface); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	cmd := args[1]
	var cmderr error
	switch cmd {
	case "up":
		cmderr = netlink.NetworkLinkUp(iface)
	case "down":
		cmderr = netlink.NetworkLinkDown(iface)
	case "add", "del":
		if len(args) != 3 {
			flag.Usage()
			os.Exit(1)
		}
		ip, ipnet, err := net.ParseCIDR(args[2])
		if err != nil {
			log.Fatal(err)
		}
		f := netlink.NetworkLinkAddIp
		if cmd == "del" {
			f = netlink.NetworkLinkDelIp
		}
		cmderr = f(iface, ip, ipnet)
	default:
		log.Fatalf("invalid verb %s", cmd)
	}
	if cmderr != nil {
		log.Fatalf("%s: %v", cmd, cmderr)
	}
}

func printIfaces() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		if err := printIface(&iface); err != nil {
			log.Fatal(err)
		}
	}
}

func printIface(iface *net.Interface) error {
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
