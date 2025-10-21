package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

type ScanResult struct {
	Port    int
	Open    bool
	Service string
}

type PortScanner struct {
	host    string
	timeout time.Duration
	workers int
}

var commonPorts = map[int]string{
	20:   "FTP Data",
	21:   "FTP Control",
	22:   "SSH",
	23:   "Telnet",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	110:  "POP3",
	143:  "IMAP",
	443:  "HTTPS",
	445:  "SMB",
	3306: "MySQL",
	3389: "RDP",
	5432: "PostgreSQL",
	5900: "VNC",
}

func NewPortScanner(host string, timeout time.Duration, workers int) *PortScanner {
	return &PortScanner {
		host, 
		timeout, 
		workers,
	}
}

func (ps *PortScanner) ScanPort(port int) ScanResult {
	address := fmt.Sprintf("%s:%d", ps.host, port)
	conn, err := net.DialTimeout("tcp", address, ps.timeout)

	result := ScanResult {
		Port: port, 
		Open: false, 
		Service: commonPorts[port], 
	}

	if err == nil {
		result.Open = true 
		conn.Close()
	}

	return result
}

func (ps *PortScanner) ScanRange(startPort, endPort int) []ScanResult {
	// Create channels
	jobs := make(chan int, endPort - startPort+1)
	results := make(chan ScanResult, endPort-startPort+1)

	// Create a waitGroup to wait 
	var wg sync.WaitGroup

	// Start worker goroutines 
	for i := 0; i < ps.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				results <- ps.ScanPort(port)
			}
		}()
	}

	// Send the jobs to workers 
	go func() {
		for port := startPort; port <= endPort; port++ {
			jobs <- port
		}
	}()

	// Close results channel when all workers are done 
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var scanResults []ScanResult
	for result := range results {
		scanResults = append(scanResults, result)
	}

	// Sort the results by port number
	sort.Slice(scanResults, func(i, j int) bool {
		return scanResults[i].Port < scanResults[j].Port
	})

	return scanResults
}

// ScanCommonPorts method scans only 
func (ps *PortScanner) scanCommonPorts() []ScanResult {
	ports := make([]int, 0, len(commonPorts))
	for port := range commonPorts {
		ports = append(ports, port)
	}
	sort.Ints(ports)

	jobs := make(chan int, len(ports))
	results := make(chan ScanResult, len(ports))

	var wg sync.WaitGroup
	// Start the workers 
	for i := 0; i < ps.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				results <- ps.ScanPort(port)
			}
		}()
	}

	// Send jobs 
	go func() {
		for _, port := range ports {
			jobs <- port 
		}
		close(jobs)
	}()

	// Wait for workers to complete and close 
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect the results 
	var scanResults []ScanResult 
	for scanResult := range results {
		scanResults = append(scanResults, scanResult)
	}

	sort.Slice(scanResults, func(i, j int) bool {
		return scanResults[i].Port < scanResults[j].Port
	})

	return scanResults 
}

// ScanSpecificPorts scans a list of specific ports 
func (ps *PortScanner) ScanSpecificPorts(ports []int) []ScanResult {
	jobs := make(chan int, len(ports))
	results := make(chan ScanResult, len(ports))

	var wg sync.WaitGroup

	// Spawn workers 
	for i := 0; i < ps.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				results <- ps.ScanPort(port)
			}
		}()
	}

	// Send the jobs 
	go func() {
		for _, port := range ports {
			jobs <- port
		}
		close(jobs)
	}()

	// Wait and close after completion 
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect the results 
	scanResults := make([]ScanResult, 0, len(ports))
	for scanResult := range results {
		scanResults = append(scanResults, scanResult)
	}
	// Sort scan results
	sort.Slice(scanResults, func(i, j int) bool {
		return scanResults[i].Port < scanResults[j].Port
	})

	return scanResults
}

// PrintResults displays scan results 
func PrintResults(results []ScanResult, showClosed bool) {
	fmt.Println("\n" + "===============================================================")
	fmt.Println("PORT SCAN RESULTS")
	fmt.Println("===============================================================")

	openCount := 0  
	closedCount := 0 

	for _, result := range results {
		if result.Open {
			openCount++
			service := result.Service
			if service == "" {
				service = "Unknown"
			}
			fmt.Printf("Port %d is OPEN - %s\n", result.Port, service)
		} else {
			closedCount++ 
			if showClosed {
				fmt.Printf("Port %d is CLOSED\n", result.Port)
			}
		}
	}

	fmt.Println("===============================================================")
	fmt.Printf("Total Ports Scanned: %d\n", len(results))
	fmt.Printf("Open: %d | Closed: %d\n", openCount, closedCount)
	fmt.Println("===============================================================")
}

func main() {
	// Configuration 
	host := "localhost"
	timeout := 500 * time.Millisecond
	workers := 100 

	scanner := NewPortScanner(host, timeout, workers)

	fmt.Printf("Starting port scan on %s...\n", host)
	fmt.Printf("Workers: %d | Timeout: %v\n", workers, timeout)

	// Scan Common ports  
	fmt.Println("\n--- Scanning Common Ports ---")
	start := time.Now()
	results := scanner.scanCommonPorts()
	elapsed := time.Since(start)
	PrintResults(results, false)
	fmt.Printf("Scan completed in %v\n", elapsed)

	// Usage 2: Scan a range of ports
	start = time.Now()
	results = scanner.ScanRange(1, 1024)
	elapsed = time.Since(start)
	PrintResults(results, false)
	fmt.Printf("Scan completed in %v\n", elapsed)

	// Usage 3: Scan specific ports 
	specificPorts := []int{22, 80, 443, 3000, 5432}
	fmt.Println("\n--- Scanning Specific Ports ---")
	start = time.Now()
	results = scanner.ScanSpecificPorts(specificPorts)
	elapsed = time.Since(start)
	PrintResults(results, true)
	fmt.Printf("Scan completed in %v\n", elapsed)
}
