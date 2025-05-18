package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/DylanPina/go-dns-server/internal/dns"
)

func main() {
	port := flag.Int("port", 2053, "Port that the DNS server to listen on (default: 2053)")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	fmt.Println("Server listening on ", udpAddr)

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		header, err := dns.DecodeDNSHeader(receivedData)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode DNS header:", err)
			continue
		}

		response, err := header.Encode()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to encode DNS header:", err)
			continue
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to send response:", err)
		}

		fmt.Printf("Sent %d bytes to %s: %s\n", len(response), source, string(response))
	}
}
