package server

import (
	"fmt"
	"net"
	"sync"
)

type UDPServer struct {
	config Configuration
	handler Handler

	listener net.Listener
	running bool
	wg sync.WaitGroup
	mu sync.Mutex
	conn *net.UDPConn
}

func NewUDPServer(host string, port int, handler Handler) *UDPServer {
	return &UDPServer {
		config: Configuration {
			Host: host,
			Port: port,
			Type: UDP,
		},
		handler: handler,
	}
}

func (s *UDPServer) Start() error {
	s.mu.Lock()

	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to start UDP Server: %w", err)
	}

	s.conn = conn
	s.running = true
	s.mu.Unlock()

	s.handlePackets()

	return nil
}

func (s *UDPServer) handlePackets() {
	buffer := make([]byte, 1024)

	for {
		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			return
		}

		conn := s.conn
		s.mu.Unlock()

		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			s.mu.Lock()
			if !s.running {
				s.mu.Unlock()
				return
			}

			s.mu.Unlock()

			fmt.Printf("Error reading UDP packet: %v\n", err)
			continue
		}

		s.wg.Add(1)
		go func (data []byte)  {
			defer s.wg.Done()
			if err := s.handler.Handle(data); err != nil {
				fmt.Printf("Error handling UDP packet: %v\n", err)
			}
		}(append([]byte(nil), buffer[:n]...))
	}
}

func (s *UDPServer) Stop() error {
    s.mu.Lock()
    if !s.running {
        s.mu.Unlock()
        return fmt.Errorf("server is not running")
    }

	// clean up
    s.running = false
    if err := s.conn.Close(); err != nil {
        s.mu.Unlock()
        return fmt.Errorf("error closing connection: %w", err)
    }
    s.mu.Unlock()

    s.wg.Wait()
    return nil
}

func (s *UDPServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}