package server

type Handler interface {
	Handle(data []byte) error
}

type Server interface {
	Start() error
	Stop() error
	IsRunning() bool
}

type ServerType int

const (
	TCP ServerType = iota
	UDP
)

type Configuration struct {
	Host string
	Port int
	Type ServerType
}