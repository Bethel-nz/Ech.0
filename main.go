package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	listenPort := flag.String("listen", "", "port to listen on (optional, e.g., :8080)")
	destinationAddr := flag.String("destination", "", "destination address to forward to (e.g., localhost:3000)")
	startPort := flag.Int("start-port", 8000, "start of port range to find available port")
	endPort := flag.Int("end-port", 9000, "end of port range to find available port")
	flag.Parse()

	if *destinationAddr == "" {
		log.Fatal("Please provide a destination address (e.g., -destination localhost:3000)")
	}

	log.Printf("Testing connection to destination %s...\n", *destinationAddr)
	testConn, err := net.DialTimeout("tcp", *destinationAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("⚠️ Cannot connect to destination %s: %v\n", *destinationAddr, err)
	}
	testConn.Close()

	var listener net.Listener
	if *listenPort == "" {
		listener = getAvailableListener(*startPort, *endPort)
	} else {
		listener, err = net.Listen("tcp", "0.0.0.0"+*listenPort)
		if err != nil {
			log.Printf("⚠️ Port %s is in use, finding an available port...\n", *listenPort)
			listener = getAvailableListener(*startPort, *endPort)
		}
	}
	defer listener.Close()

	addr := listener.Addr().String()
	_, port, _ := net.SplitHostPort(addr)

	fmt.Printf("\nPort forwarder started\n")
	fmt.Printf("Listening on port %s\n", port)
	fmt.Printf("Forwarding to %s\n", *destinationAddr)
	fmt.Printf("\nAvailable network interfaces:\n")
	printNetworkInterfaces()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		<-shutdown
		log.Println("\nShutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Println("⚠️ Failed to accept connection:", err)
			continue
		}
		log.Printf("New connection from: %s\n", conn.RemoteAddr())
		go handleConnection(conn, *destinationAddr)
	}

	log.Println("Shutdown complete.")
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
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
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
