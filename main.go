package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port    int
	State   bool
	Service string
	Banner  string
	Version string
}

// for Logging banner
var logMutex sync.Mutex

func grabBanner(host string, port int, timeout time.Duration) (string, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	// Sends HTTP HEAD request to respond
	if port == 80 || port == 443 || port == 8080 || port == 8443 {
		fmt.Fprintf(conn, "HEAD / HTTP/1.0\r\n\r\n")
	} else {
		fmt.Fprintf(conn, "\r\n")
	}

	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(banner), nil
}

func identifyService(port int, banner string) (string, string) {
	commonPorts := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		3306:  "MySQL",
		5432:  "PostgreSQL",
		6379:  "Redis",
		8080:  "HTTP-Proxy",
		27017: "MongoDB",
	}

	service := "Unknown"
	if s, exists := commonPorts[port]; exists {
		service = s
	}

	version := "Unknown"

	lowerBanner := strings.ToLower(banner)

	// SSH version detection
	if strings.Contains(lowerBanner, "ssh") {
		service = "SSH"
		parts := strings.Split(banner, " ")
		if len(parts) >= 2 {
			version = parts[1]
		}
	}

	// HTTP server detection
	if strings.Contains(lowerBanner, "http") || strings.Contains(lowerBanner, "apache") ||
		strings.Contains(lowerBanner, "nginx") {
		if port == 443 {
			service = "HTTPS"
		} else {
			service = "HTTP"
		}

		// Try to find server info in format "Server: Apache/2.4.29"
		if strings.Contains(banner, "Server:") {
			parts := strings.Split(banner, "Server:")
			if len(parts) >= 2 {
				version = strings.TrimSpace(parts[1])
			}
		}
	}

	return service, version
}

func scanPort(host string, port int, timeout time.Duration) ScanResult {

	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return ScanResult{Port: port, State: false}
	}

	conn.Close()

	banner, err := grabBanner(host, port, timeout)

	service := "Unknown"
	version := "Unknown"

	if err == nil && banner != "" {
		service, version = identifyService(port, banner)
	}

	return ScanResult{
		Port:    port,
		State:   true,
		Service: service,
		Banner:  banner,
		Version: version,
	}
}

func scanPorts(host string, start, end int, timeout time.Duration) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup

	resultChan := make(chan ScanResult, end-start+1)

	semaphore := make(chan struct{}, 100)

	for port := start; port <= end; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := scanPort(host, p, timeout)
			resultChan <- result
		}(port)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.State {
			results = append(results, result)
		}
	}

	return results
}

func main() {
	// Define CLI flags
	host := flag.String("host", "", "Target Host/IP address")
	startPort := flag.Int("start", 1, "Start port number")
	endPort := flag.Int("end", 65000, "End port number")
	timeoutMs := flag.Int("timeout", 800, "Timeout in milliseconds")

	flag.Parse()

	if *host == "" {
		fmt.Println("Error: host is required")
		fmt.Println("Usage: go run main.go -host <IP> -start <startPort> -end <endPort>")
		return
	}

	timeout := time.Duration(*timeoutMs) * time.Millisecond

	fmt.Printf("Scanning %s from port %d to %d\n", *host, *startPort, *endPort)
	startTime := time.Now()

	results := scanPorts(*host, *startPort, *endPort, timeout)

	elapsed := time.Since(startTime)

	fmt.Printf("\nScan completed in %s\n", elapsed)
	fmt.Printf("Found %d open ports:\n\n", len(results))

	fmt.Println("PORT\tSTATUS\tSERVICE\tBANNER")
	fmt.Println("----\t------\t-------\t------")
	for _, result := range results {
		status := "CLOSED"
		service := "-"
		bannerPreview := "-"

		if result.State {
			status = "OPEN"
			service = result.Service
			if len(result.Banner) > 30 {
				bannerPreview = result.Banner[:30] + "..."
			} else {
				bannerPreview = result.Banner
			}
		}

		fmt.Printf("%d\t%s\t%s\t%s\n",
			result.Port,
			status,
			service,
			bannerPreview)
	}
	fmt.Printf("\nTotal scan time: %s\n", elapsed)

	for {
		fmt.Print("Do you want to run another scan? (y/n): ")
		var again string
		fmt.Scanln(&again)

		if strings.ToLower(again) != "y" {
			return
		}

		fmt.Print("Enter new host/IP: ")
		fmt.Scanln(host)

		fmt.Print("Enter new start port: ")
		fmt.Scanln(startPort)

		fmt.Print("Enter new end port: ")
		fmt.Scanln(endPort)

		fmt.Printf("Scanning %s from port %d to %d\n", *host, *startPort, *endPort)
		startTime = time.Now()
		results = scanPorts(*host, *startPort, *endPort, timeout)
		elapsed = time.Since(startTime)

		fmt.Printf("\nScan completed in %s\n", elapsed)
		fmt.Printf("Found %d open ports:\n\n", len(results))
		fmt.Println("PORT\tSTATUS\tSERVICE\tBANNER")
		fmt.Println("----\t------\t-------\t------")
		for _, result := range results {
			status := "CLOSED"
			service := "-"
			bannerPreview := "-"

			if result.State {
				status = "OPEN"
				service = result.Service
				if len(result.Banner) > 30 {
					bannerPreview = result.Banner[:30] + "..."
				} else {
					bannerPreview = result.Banner
				}
			}

			fmt.Printf("%d\t%s\t%s\t%s\n",
				result.Port,
				status,
				service,
				bannerPreview)
		}
		fmt.Printf("\nTotal scan time: %s\n", elapsed)
	}
}
