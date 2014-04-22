package dns

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

func TestHosts(hosts []string, debug bool) {
	for _, host := range hosts {
		_, err := net.ResolveIPAddr("ip4", host) // *IPAddr
		if err != nil {
			fmt.Println("Invalid DNS name or IP address:", host)
			return
		}

		wg.Add(1)
		go TestHost(host, debug)
	}

	wg.Wait()
}

func TestHost(host string, debug bool) {
	defer wg.Done()

	if debug {
		fmt.Println("testing host: " + host)
	}

	var startTime = time.Now()
	addrs, err := net.LookupIP(host)
	var endTime = time.Now()
	var elapsedTime = endTime.Sub(startTime)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("addresses for %s: %v\ntime: %f seconds\n\n", host, addrs, elapsedTime.Seconds())
	}
}
