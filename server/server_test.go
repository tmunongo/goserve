package server

import (
	"testing"
	"time"
)

func TestServerInterface(t *testing.T) {
	var _ Server = &TCPServer{}
	var _ Server = &UDPServer{}
}

type mockHandler struct {
	handleCalled bool
	data		[]byte
}

func (m *mockHandler) Handle(data []byte) error {
	m.handleCalled = true
	m.data = data
	return nil
}

func TestServerCreate(t *testing.T) {
	handler := &mockHandler{}
	tcpServer := NewTCPServer("localhost", 8080, handler)

	if tcpServer.config.Host != "localhost" {
		t.Error("Host does not match")
	}

	if tcpServer.config.Port != 8080 {
		t.Error("Port does not match")
	}

	if tcpServer.IsRunning() != false {
		t.Error("Server has not been started. Should not be running")
	}

	udpServer := NewUDPServer("127.0.0.1", 4040, handler)

	if udpServer.config.Host != "127.0.0.1" {
		t.Error("Host does not match")
	}

	if udpServer.config.Port != 4040 {
		t.Error("Port does not match")
	}

	if udpServer.IsRunning() != false {
		t.Error("Server has not been started. Should not be running")
	}
}

func TestServerLifecycle(t *testing.T) {
	handler := &mockHandler{}
	server := NewTCPServer("localhost", 8080, handler)

	if err := server.Start(); err != nil {
		t.Fatalf("Server test failed to start server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// test server running
	if !server.IsRunning() {
		t.Error("Server is not running")
	}

	if err := server.Stop(); err != nil {
		t.Fatalf("Failed to stop the server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if server.IsRunning() {
		t.Error("Server should not be running")
	}
}