# Ech.0

A simple TCP port forwarder built in Go. It redirects traffic from a source port to a target port, making local services accessible over a network.

## Features
- Port forwarding between any two TCP ports
- Automatic target port selection if not provided
- Displays available network interfaces for easier LAN access
- Lightweight and fast using Go’s concurrency

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/your-username/go-port-forwarder.git
   cd go-port-forwarder
   ```
2. Build the executable:
   ```sh
   go build -o port-forwarder main.go
   ```
3. Run the forwarder:
   ```sh
   ./port-forwarder -source localhost:3000 -target :8080
   ```

## Usage
Basic port forwarding:
```sh
./port-forwarder -source localhost:3000 -target :8080
```
Automatically find an available port:
```sh
./port-forwarder -source localhost:3000
```
Access over LAN by using your local IP:
```sh
http://192.168.x.x:8080
```

## Development
Run in debug mode:
```sh
go run main.go -source localhost:3000 -target :8080
```
Cross-compile for different OS:
```sh
GOOS=windows GOARCH=amd64 go build -o port-forwarder-windows.exe main.go
GOOS=linux GOARCH=amd64 go build -o port-forwarder-linux main.go
GOOS=darwin GOARCH=amd64 go build -o port-forwarder-mac main.go
```

## Why
- Needed a way to access services running on my pc over the LAN network

## License
This project is open-source under the MIT License.
