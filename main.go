package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	targetPort := flag.String("target", "", "target port to forward to (optional, e.g., :8080)")
	sourcePort := flag.String("source", "", "source port to forward from (e.g., localhost:3000)")
	flag.Parse()

	if *sourcePort == "" {
		log.Fatal("Please provide a source port (e.g., -source localhost:3000)")
	}

	log.Printf("Testing connection to source %s...\n", *sourcePort)
	testConn, err := net.DialTimeout("tcp", *sourcePort, 5*time.Second)
	if err != nil {
		log.Fatalf("⚠️ Cannot connect to source %s: %v\n", *sourcePort, err)
	}
	testConn.Close()

	var listener net.Listener
	if *targetPort == "" {
		listener = getAvailableListener(8000, 9000)
	} else {
		listener, err = net.Listen("tcp", "0.0.0.0"+*targetPort)
		if err != nil {
			log.Printf("⚠️ Port %s is in use, finding an available port...\n", *targetPort)
			listener = getAvailableListener(8000, 9000)
		}
	}
	defer listener.Close()

	addr := listener.Addr().String()
	_, port, _ := net.SplitHostPort(addr)

	fmt.Printf("\nPort forwarder started\n")
	fmt.Printf("Forwarding from %s\n", *sourcePort)
	fmt.Printf("To port %s\n", port)
	fmt.Printf("\nAvailable network interfaces:\n")
	printNetworkInterfaces()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("⚠️ Failed to accept connection:", err)
			continue
		}
		log.Printf("New connection from: %s\n", conn.RemoteAddr())
		go handleConnection(conn, *sourcePort)
	}
}

// Finds an available port within a given range
func getAvailableListener(startPort, endPort int) net.Listener {
	for port := startPort; port <= endPort; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err == nil {
			log.Printf(" Found available port: %d\n", port)
			return listener
		}
	}
	log.Fatal(" No available ports found in range")
	return nil
}

// Prints available network interfaces
func printNetworkInterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	fmt.Println("  -----------------------------")
	for _, iface := range interfaces {
		if !isRelevantInterface(iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && !strings.HasPrefix(ipnet.IP.String(), "169.254") {
					fmt.Printf("  %-22s %s\n", iface.Name+":", ipnet.IP.String())
				}
			}
		}
	}
}

// Filters only relevant network interfaces
func isRelevantInterface(name string) bool {
	name = strings.ToLower(name)
	return strings.Contains(name, "ethernet") ||
		strings.Contains(name, "wi-fi") ||
		strings.Contains(name, "local area connection")
}

// Handles TCP connection forwarding
func handleConnection(source net.Conn, targetAddr string) {
	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("⚠️ Failed to connect to target: %v\n", err)
		source.Close()
		return
	}
	defer source.Close()
	defer target.Close()

	go func() { io.Copy(target, source) }()
	io.Copy(source, target)
}
