# Leisure Time Projects
A list of projects I build in my leisure/free time, for study purposes inorder to improve my skills 

### Projects Built So Far 
- Port Scanner: Scans for open ports in a given address
```go
  // List of commonPorts
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
```
 #### How to run the port scanner
 ```bash 
  cd port-scanner
  go run main.go
```
