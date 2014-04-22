# LinkCheck

A tool for checking link states.  This can do DNS resolution tests and Ping tests.  Results including elapsed time is displayed.


## Sample Output

```shell
user@host ~/linkcheck: ./bin/linkcheck -p -n 3 -d
addresses for 127.0.0.1: [127.0.0.1]
time: 0.000607 seconds

addresses for 74.125.142.99: [74.125.142.99]
time: 0.000777 seconds

addresses for www.google.com: [74.125.142.147 74.125.142.105 74.125.142.103 74.125.142.99 74.125.142.106 74.125.142.104 2607:f8b0:4001:c03::93]
time: 0.000380 seconds

PING 74.125.142.99 (74.125.142.99): 56 data bytes
64 bytes from 74.125.142.99: icmp_req=1 time=50.216 ms
64 bytes from 74.125.142.99: icmp_req=2 time=51.278 ms
64 bytes from 74.125.142.99: icmp_req=3 time=51.044 ms

--- 74.125.142.99 ping statistics ---
3 packets transmitted, 3 packets received, 0.000% packet loss
round-trip min/avg/max/stddev = 50.216/50.846/51.278/0.000 ms


PING www.google.com (74.125.142.147): 56 data bytes
64 bytes from 74.125.142.147: icmp_req=1 time=52.755 ms
64 bytes from 74.125.142.147: icmp_req=2 time=53.117 ms
64 bytes from 74.125.142.147: icmp_req=3 time=53.309 ms

--- www.google.com ping statistics ---
3 packets transmitted, 3 packets received, 0.000% packet loss
round-trip min/avg/max/stddev = 52.755/53.060/53.309/0.000 ms
```
