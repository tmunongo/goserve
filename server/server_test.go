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