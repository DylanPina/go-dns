package client

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

// ConnectUDP sends a single UDP request to the desginated address and returns the raw response bytes.
//
// method, path and headers are just concatenated into the datagram preamble;
// you can customize framing to match your server’s expectations.
func ConnectUDP(port int, method, path string, headers map[string]string, body []byte) ([]byte, error) {
	// Resolve and dial the server address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("resolve UDP addr: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("dial UDP: %w", err)
	}
	defer conn.Close()

	// Application-level request packet
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", method, path)
	for k, v := range headers {
		fmt.Fprintf(&buf, "%s: %s\n", k, v)
	}
	buf.WriteString("\n")
	if len(body) > 0 {
		buf.Write(body)
	}

	// Send the packet
	if _, err := conn.Write(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("write UDP packet: %w", err)
	}

	// Optional: set a read timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Receive response
	resp := make([]byte, 4096)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("read UDP response: %w", err)
	}

	return resp[:n], nil
}

func DNSQuery() {}
