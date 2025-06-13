package main

import (
	"fmt"
	"net"
	"sort"
	"time"
)

const (
	workerCount = 100
	timeout     = 500 * time.Millisecond
)

type ScanResult struct {
	port int
	open bool
}

func worker(ports <-chan int, results chan<- ScanResult) {
	for port := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", port)
		conn, err := net.DialTimeout("tcp", address, timeout)

		if err != nil {
			results <- ScanResult{port: port, open: false}
			continue
		}

		conn.Close()
		results <- ScanResult{port: port, open: true}
	}
}

func main() {
	ports := make(chan int, workerCount)
	results := make(chan ScanResult)
	var openPorts []int

	// Start workers
	for i := 0; i < workerCount; i++ {
		go worker(ports, results)
	}

	// Send ports to workers
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
		close(ports)
	}()

	// Collect results
	for i := 1; i <= 1024; i++ {
		result := <-results
		if result.open {
			openPorts = append(openPorts, result.port)
		}
	}
	close(results)

	// Sort and print results
	sort.Ints(openPorts)
	fmt.Println("\n=== TCP Port Scanner Results ===")
	fmt.Println("Target: scanme.nmap.org")
	fmt.Println("Scan completed in", timeout)
	fmt.Println("\nOpen ports:")
	for _, port := range openPorts {
		fmt.Printf("  %d/tcp - open\n", port)
	}
	fmt.Println("\nTotal open ports:", len(openPorts))
}
