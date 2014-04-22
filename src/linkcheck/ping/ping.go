// Original Source: https://raw.githubusercontent.com/atomaths/gtug8/master/ping/ping.go
//
// 121.254.177.105
// tcpdump -n icmp and icmp[icmptype] != icmp-echo or icmp[icmptype] != icmp-echoreply
//
// CAP_NET_RAW
//
// http://blog.daum.net/wonho777/5320889
// http://kldp.org/node/35797
// https://groups.google.com/forum/#!searchin/golang-nuts/ping/golang-nuts/yyXqNGIcMzA/mmZ_vQKun9UJ

package ping

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"time"
)

const (
	ICMP_ECHO_REQUEST = 8
	ICMP_ECHO_REPLY   = 0
)

type Stats struct {
	Sent     int
	Received int
	LossPct  float64
	Min      float64
	Max      float64
	Avg      float64
	Stddev   float64
}

func NewPingStats(sent int, received int, etimes []float64) Stats {
	var s Stats
	s.Sent = sent
	s.Received = received
	s.LossPct = (1.0 - float64(received)/float64(sent)) * 100.00

	var total float64
	for i, _ := range etimes {
		if i == 0 {
			s.Max = etimes[i]
			s.Min = etimes[i]
		} else {
			s.Max = math.Max(s.Max, etimes[i])
			s.Min = math.Min(s.Min, etimes[i])
		}

		total += etimes[i]
		s.Avg = total / (float64(i) + 1.00)
	}

	return s
}

// returns a suitable 'ping request' packet, with id & seq and a
// payload length of pktlen
func makePingRequest(id, seq, pktlen int, filler []byte) []byte {
	p := make([]byte, pktlen)
	copy(p[8:], bytes.Repeat(filler, (pktlen-8)/len(filler)+1))

	p[0] = ICMP_ECHO_REQUEST // type
	p[1] = 0                 // code
	p[2] = 0                 // cksum
	p[3] = 0                 // cksum
	p[4] = uint8(id >> 8)    // id
	p[5] = uint8(id & 0xff)  // id
	p[6] = uint8(seq >> 8)   // sequence
	p[7] = uint8(seq & 0xff) // sequence

	// calculate icmp checksum
	cklen := len(p)
	s := uint32(0)
	for i := 0; i < (cklen - 1); i += 2 {
		s += uint32(p[i+1])<<8 | uint32(p[i])
	}
	if cklen&1 == 1 {
		s += uint32(p[cklen-1])
	}
	s = (s >> 16) + (s & 0xffff)
	s = s + (s >> 16)

	// place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero
	p[2] ^= uint8(^s & 0xff)
	p[3] ^= uint8(^s >> 8)

	return p
}

func parsePingReply(p []byte) (id, seq int) {
	id = int(p[4])<<8 | int(p[5])
	seq = int(p[6])<<8 | int(p[7])
	return
}

func elapsedTime(start time.Time) float64 {
	t := float64(time.Since(start).Seconds()) * 1000.0
	return t
}

func PingIPs(hosts []string, pcount int) {
	if pcount == 0 {
		fmt.Println("Ping tests disabled (0 packet ping count specified)")
		return
	}
	for i, host := range hosts {
		raddr, err := net.ResolveIPAddr("ip4", host) // *IPAddr
		if err != nil {
			fmt.Println("Invalid DNS name or IP address:", host)
			return
		}

		sent, received, etimes := PingIP(host, raddr, pcount)
		stats := NewPingStats(sent, received, etimes)

		fmt.Println("\n--- " + host + " ping statistics ---")
		fmt.Printf("%d packets transmitted, %d packets received, %.3f%% packet loss\n", stats.Sent, stats.Received, stats.LossPct)
		fmt.Printf("round-trip min/avg/max/stddev = %.3f/%.3f/%.3f/%.3f ms\n\n\n", stats.Min, stats.Avg, stats.Max, stats.Stddev)

		if i < len(hosts)-1 {
			time.Sleep(2 * time.Second)
		}
	}
}

func PingIP(hostname string, raddr *net.IPAddr, pcount int) (sent int, received int, etimes []float64) {
	// Make the IP connection to the destination host
	ipconn, err := net.DialIP("ip4:icmp", nil, raddr)
	if err != nil {
		fmt.Printf("could not connect to %s: %v\n", raddr.IP, err)
		return
	}

	sendid := os.Getpid() & 0xffff
	sendseq := 1
	pingpktlen := 64

	fmt.Printf("PING %s (%s): %d data bytes\n", hostname, raddr.String(), pingpktlen-8)
	for {
		var check4reply bool
		var etime float64

		// Generate our ICMP packet
		sendpkt := makePingRequest(sendid, sendseq, pingpktlen, []byte("Go Ping"))

		// Start timer
		start := time.Now()

		// Send ICMP packet
		n, err := ipconn.Write(sendpkt)

		// Err out if we couldn't send the full packet size
		if err != nil || n != pingpktlen {
			etime = elapsedTime(start)
			fmt.Printf("0 bytes from %s: icmp_req=%d time=%.3f ms ERR: Network is down\n", raddr.IP, sendseq, etime)
		} else {
			// We will check this flag later to see if we need to wait for a returned packet
			check4reply = true
		}

		// increment our sent counter
		sent++

		// set a 0.5 second timer as our max packet reply wait time
		var deadline = start.Add(500 * time.Millisecond)
		ipconn.SetReadDeadline(deadline)

		// Read response packet (or process timeout if it occurs)
		resp := make([]byte, 1024)
		for {
			if check4reply == false {
				break
			}

			n, _, err := ipconn.ReadFrom(resp)
			if err != nil {
				// Could not read response packet
				etime = elapsedTime(start)
				fmt.Printf("%d bytes from %s: icmp_req=%d time=%.3f ms ERR: %s\n", n, raddr.IP, sendseq, etime, err.Error())
				break
			} else {
				// Response was okay

			}

			if resp[0] != ICMP_ECHO_REPLY {
				// Skip non-ICMP packets
				continue
			}

			rcvid, rcvseq := parsePingReply(resp)
			if rcvid != sendid || rcvseq != sendseq {
				etime = elapsedTime(start)
				fmt.Printf("%d bytes from %s: icmp_req=%d time=%.3f ms ERR: Out of sequence (0x%x,0x%x)\n", n, raddr.IP, sendseq, etime, rcvid, rcvseq)
			} else {
				etime = elapsedTime(start)
				fmt.Printf("%d bytes from %s: icmp_req=%d time=%.3f ms\n", n, raddr.IP, sendseq, etime)
				received++
				etimes = append(etimes, etime)
			}
			break
		}

		sendseq++
		if sendseq > pcount {
			// We've reached the maximum number of packets to send so return
			return
		} else {
			// Sleep 1 second before moving onto the next packet
			time.Sleep(1 * time.Second)
		}
	}
}
