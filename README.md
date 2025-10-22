# Leisure Time Projects
A list of projects I build in my leisure/free time in **Golang**, for study purposes inorder to improve my skills 

### Projects Built So Far 
---
- **Port Scanner**: Scans for open ports in a given address
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
 ### How to run
 ```bash 
  cd port-scanner
  go run main.go
```
----
- **Live Reloader**: Listens for changes in files in a given directory. It works by recursively visiting all files in a given folder, each file's contents are hashed and stored in a hashmap, it recursively reads every file again, hashes the content then compares the hash with the hash in the hashmap, if not them same, then it knows that the file's contents have been modified.

### How to run 
```bash 
cd live-reloader 
go run main.go
```
---
- **Concurrent Task Workers**: A simple implementation of how concurrent workers execute multiple jobs in parallel

### How to run 
```bash 
cd concurrent-task-workers 
go run main.go
```

## Authors
- Ismael Swaleh
