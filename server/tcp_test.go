package server

import (
	"net"
	"testing"
	"time"
)

func TestTCPServer(t *testing.T) {
	handler := &mockHandler{}
	server := NewTCPServer("localhost", 8080, handler)

	// start in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Failed to start the server: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	testData := []byte("Hello, Server!")
	_, err = conn.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to server %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !handler.handleCalled {
		t.Error("Handler was not called")
	}

	if string(handler.data) != string(testData) {
		t.Errorf("Handler received wrong data. Got %s instead of %s", string(handler.data), string(testData))
	}

	// clean up
	conn.Close()
	server.Stop()
}