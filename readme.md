# Advanced Network Port Scanner
The **Advanced Network Port Scanner** is a high-performance, concurrent TCP scanner written in Go.  
It is designed to detect **open ports, running services, banner information, and version details** across any target host.

This scanner supports **full-range scanning (1–65535)** with:
- Banner grabbing
- Service fingerprinting
- Progress tracking
- JSON export

Built as a **robust network vulnerability assessment tool**, it demonstrates Go’s strengths in developing security utilities.  
The scanner rapidly probes ports, identifies the services running behind them, and provides detailed output for analysis — powered by multi-threading and structured result formats.

## Features

- Fast **concurrent TCP port scanning** (1–65535)
- **Banner grabbing** for service fingerprinting
- **Automatic service & version detection**
- **Progress indicator** during scan
- **JSON export support**
- **Configurable timeout & port range**
- **Interactive multi-scan mode**
- **Docker support for easy deployment**

## Installation
Make sure **Go** is installed:
[Go Installation](https://go.dev/dl/)

### Clone the repository:

```bash
https://github.com/OmkarGeedh/Advanced-Network-Port-Scanner.git
cd Advanced-Network-Port-Scanner
go mod tidy
go build .
```
## Usage (without Docker)
### Basic scan of a single host
```bash
go run main.go -host scanme.nmap.org
```
### Scan port range
```bash
go run main.go -host scanme.nmap.org -start 1 -end 1024
```
### Save result to JSON
```
go run main.go -host scanme.nmap.org -format json -output results.json
```
### Or build the binary and run
```bash
go build -o scanner
./scanner -host scanme.nmap.org
```

## 🐳 Running with Docker
### Build Docker Image
In the project root (where main.go is):
```bash
docker build -t advanced-port-scanner .
```
### Run a basic scan for a single host
```bash
docker run --rm advanced-port-scanner -host scanme.nmap.org
```
### Scan a specific port range
```bash
docker run --rm advanced-port-scanner -host scanme.nmap.org -start 1 -end 1024
```
### Save results to a JSON file
```bash
docker run --rm advanced-port-scanner -host scanme.nmap.org -format json -output results.json
```

## ⚙️ Flag Reference

| Flag       | Description |
|------------|-------------|
| `-host`    | Target IP / Domain |
| `-start`   | Start port (default: `1`) |
| `-end`     | End port (default: `65535`) |
| `-timeout` | Timeout in milliseconds (default: `800`) |
| `-format`  | Output format: `text` or `json` |
| `-output`  | JSON export filename |


## Output example
```bash
$go run . -host localhost 
Scanning localhost from port 1 to 65535
Scanning [65359/65535]
Scan completed in 1.325787958s
Found 6 open ports:

PORT	STATUS	SERVICE	BANNER
----	------	-------	------
80	OPEN	HTTP  -	
443	OPEN	HTTPS  Server: nginx/1.18.0	
8080	OPEN	HTTP	HTTP/1.0 404 Not Found
5000	OPEN	upnp	
7000	OPEN	Unknown	
7265	OPEN	Unknown	

Total scan time: 1.325787958s
Do you want to run another scan? (y/n): y
Enter new host/IP: 192.168.1.1
Enter new start port: 1
Enter new end port: 1000
Scanning 192.168.1.1 from port 1 to 1000
Scanning [1000/1000]
Scan completed in 199.411625ms
Found 1 open ports:

PORT	STATUS	SERVICE	BANNER
----	------	-------	------
80	OPEN	HTTP	HTTP/1.0 302 Redirect

Total scan time: 199.411625ms
Do you want to run another scan? (y/n): 

```

## JSON Output
```
[
  {
    "Port": 80,
    "State": true,
    "Service": "HTTP",
    "Banner": "",
    "Version": "Unknown"
  },
  {
    "Port": 443,
    "State": true,
    "Service": "HTTPS",
    "Banner": "Server: nginx/1.18.0",
    "Version": "nginx/1.18.0"
  },
  {
    "Port": 8080,
    "State": true,
    "Service": "HTTP",
    "Banner": "HTTP/1.0 404 Not Found",
    "Version": "Unknown"
  },
  {
    "Port": 5000,
    "State": true,
    "Service": "upnp",
    "Banner": "",
    "Version": "Unknown"
  },
  {
    "Port": 7000,
    "State": true,
    "Service": "Unknown",
    "Banner": "",
    "Version": "Unknown"
  },
  {
    "Port": 7265,
    "State": true,
    "Service": "Unknown",
    "Banner": "",
    "Version": "Unknown"
  }
]
```

