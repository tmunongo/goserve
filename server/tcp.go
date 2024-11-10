package server

import (
	"fmt"
	"net"
	"sync"
)

// implements the Server interface
// for TCP connections
type TCPServer struct {
	config Configuration
	handler Handler

	listener net.Listener
	running bool
	wg sync.WaitGroup
	mu sync.Mutex
}

func NewTCPServer(host string, port int, handler Handler) *TCPServer {
	return &TCPServer{
		config: Configuration{
			Host: host,
			Port: port,
			Type: TCP,
		},
		handler: handler,
	}
}

func (s *TCPServer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to start the TCP sevrer: %w", err)
	}

	s.listener = listener
	s.running = true
	s.mu.Unlock()

	go s.acceptConnections()

	return nil
}

func (s *TCPServer) acceptConnections() {
	for {
		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			return
		}

		listener := s.listener
		s.mu.Unlock()

		conn, err := listener.Accept()
		if err != nil {
			// was it stopped?
			s.mu.Lock()
			if !s.running {
				s.mu.Unlock()
				return
			}

			s.mu.Unlock()

			fmt.Printf("error acceping connection: %v\n", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		if err := s.handler.Handle(buffer[:n]); err != nil {
			fmt.Printf("Error handling data: %v\n", err)
			return
		}
	}
}

func (s *TCPServer) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is not running")
	}

	s.running = false
	if err := s.listener.Close(); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("error closing listener %w", err)
	}
	s.mu.Unlock()

	s.wg.Wait()
	return nil
}

func (s *TCPServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}