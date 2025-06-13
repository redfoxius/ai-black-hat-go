package main

import (
	"fmt"
	"net"
	"time"
)

func scanPort(protocol string, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	hostname := "scanme.nmap.org"

	fmt.Printf("Scanning %s for open ports...\n", hostname)
	for port := 1; port <= 1024; port++ {
		if scanPort("tcp", hostname, port) {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}
