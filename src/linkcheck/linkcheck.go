package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"linkcheck/dns"
	"linkcheck/ping"
	"os"
	"syscall"
)

type Config struct {
	Global struct {
		Interval int
		Tries    int
	}
	Dns struct {
		Host []string
	}
	Ping struct {
		Packets int
		Host    []string
	}
}

var cfgfile string
var config Config
var debug bool
var pcount int
var chkdns bool
var chkping bool

func testcfgfile() (ok error) {
	var res syscall.Stat_t
	ok = syscall.Stat(cfgfile, &res)

	return
}

func init() {
	flag.StringVar(&cfgfile, "c", "linkcheck.cfg", "configuration file location")
	flag.BoolVar(&debug, "v", false, "enable verbose messages")
	flag.IntVar(&pcount, "n", 0, "packet count for ping test (0 = use value from config file)")
	flag.BoolVar(&chkdns, "d", false, "enable DNS checks")
	flag.BoolVar(&chkping, "p", false, "enable ping checks")
}

func main() {
	flag.Parse()

	// Test cfgfile exists
	res := testcfgfile()
	if res != nil {
		fmt.Println("Configuration file error: " + res.Error())
		os.Exit(1)
	}

	// Parse in config file
	err := gcfg.ReadFileInto(&config, cfgfile)
	if err != nil {
		fmt.Println("Configuration file error: " + err.Error())
		os.Exit(1)
	}

	if chkping && pcount == 0 {
		pcount = config.Ping.Packets
	}

	if chkdns {
		dns.TestHosts(config.Dns.Host, debug)
	}

	if chkping {
		ping.PingIPs(config.Ping.Host, pcount)
	}
}
