Problem 1: Advanced Network Port Scanner
Security engineers often need to quickly assess the attack surface of unfamiliar systems. When onboarding new infrastructure, auditing third-party services, or validating firewall configurations, a fast and reliable network scanner is essential. While tools like Nmap exist, organizations often need custom scanners tailored to their environments, logging systems, or automation pipelines.
In this challenge, you will build your own lightweight network scanner capable of identifying open ports and inferring running services through banner grabbing. The goal is to demonstrate your ability to work with sockets, concurrency, and basic network protocols.
The goal is to build a concurrent port scanner capable of detecting open ports and performing lightweight service identification.
Requirements
Input: a host or subnet (e.g., scanme.nmap.org or 192.168.1.0/24)
Scan a configurable set of ports
Use concurrency (threads or async) to speed up scanning
Randomize scanning order or add small delays to avoid rate limiting
Attempt banner grabbing to identify services (e.g., SSH, HTTP)
Distinguish between open, closed, and filtered ports
Output:
Port number
Status
Service banner (if detected)
Total scan time
