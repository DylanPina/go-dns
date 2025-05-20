package main

import (
	"flag"
	"os"
	"testing"

	"github.com/DylanPina/go-dns-server/internal/client"
)

var TEST_PORT = flag.Int("port", 2053, "Port that the DNS server to listen on (default: 2053)")

// TestMain is the entry point for the test suite
func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// TestConnect tests the ConnectUDP function
func TestConnect(t *testing.T) {
	respBytes, err := client.ConnectUDP(*TEST_PORT, "GET", "", map[string]string{"X-Foo": "Bar"}, []byte("www.dsp.dev"))
	if err != nil {
		t.Fatalf("ConnectUDP failed: %v\n", err)
	}

	t.Logf("Server replied: %s\n", string(respBytes))
}
