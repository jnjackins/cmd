package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/vishvananda/netlink"
)

func routeCmd(args []string) {
	if len(args) == 0 {
		printRoutes()
		os.Exit(0)
	}

	// all other commands require dst
	if len(args) < 2 || len(args) > 3 {
		flag.Usage()
		os.Exit(1)
	}
	var dst *net.IPNet
	if args[1] == "default" {
		_, dst, _ = net.ParseCIDR("0.0.0.0/0")
	} else {
		var err error
		_, dst, err = net.ParseCIDR(args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	switch args[0] {
	case "add":
		var gw net.IP
		if len(args) == 3 {
			gw = net.ParseIP(flag.Arg(3))
			if gw == nil {
				log.Fatalf("failed to parse address: %s", flag.Arg(3))
			}

		}
		err := netlink.RouteAdd(&netlink.Route{Dst: dst, Gw: gw})
		if err != nil {
			log.Fatal(err)
		}
	case "del":
		if flag.NArg() != 3 {
			flag.Usage()
			os.Exit(1)
		}
		err := netlink.RouteDel(&netlink.Route{Dst: dst})
		if err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func printRoutes() {
	links, err := netlink.LinkList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-16s%-16s%-6s %s\n", "Destination", "Gateway", "Flags", "Iface")
	for _, link := range links {
		routes, err := netlink.RouteList(link, netlink.FAMILY_ALL)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range routes {
			printRoute(link, r)
		}
	}
}

func printRoute(link netlink.Link, r netlink.Route) {
	flags := "U"
	dst := "default"
	if r.Dst != nil {
		dst = r.Dst.String()
	}
	gw := "*"
	if r.Gw != nil {
		gw = r.Gw.String()
		flags += "G"
	}
	iface := link.Attrs().Name
	fmt.Printf("%-16s%-16s%-6s %s\n", dst, gw, flags, iface)
}
