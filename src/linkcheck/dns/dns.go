package dns

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Source: http://stackoverflow.com/a/15323988
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

func lookupNames(host string) []string {
	names, _ := net.LookupAddr(host)

	return names
}

func lookupIPAddresses(host string) (results []string) {
	addrs, _ := net.LookupIP(host)

	for _, val := range addrs {
		results = append(results, val.String())
	}

	return
}

func TestHost(host string, debug bool) {
	defer wg.Done()

	if debug {
		fmt.Println("testing host: " + host)
	}

	// Start the timer
	var startTime = time.Now()

	// Find IP addresses for the given host string
	ips := lookupIPAddresses(host)

	// Find hostnames for the given host string
	names := lookupNames(host)

	// Stop the timer
	var endTime = time.Now()

	// Calculate elapsed time for the lookups
	var elapsedTime = endTime.Sub(startTime)

	// Combine results into a single array
	var results []string

	// If we successfully resolved a hostname, then the host string was an IP address
	//    and there is no need to merge in the ip addresses since it will be equal to
	//    the host string we tried to look up.
	if len(names) == 0 {
		for _, val := range ips {
			if !stringInSlice(val, results) {
				results = append(results, val)
			}
		}
	}

	// Merge in all hostnames found in the lookup.  This will be empty if the host
	//    string was a hostname to begin with.
	for _, val := range names {
		if !stringInSlice(val, results) {
			var tmpname = strings.TrimRight(val, ".")
			results = append(results, tmpname)
		}
	}

	fmt.Printf("addresses for %s: %s\ntime: %f seconds\n\n", host, results, elapsedTime.Seconds())
}
